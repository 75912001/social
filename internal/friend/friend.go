package friend

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"net"
	"runtime"
	"runtime/debug"
	"social/lib/actor"
	libconsts "social/lib/consts"
	libmongodb "social/lib/mongodb"
	libruntime "social/lib/runtime"
	libutil "social/lib/util"
	pkgmdb "social/pkg/mdb"
	pkgmdbfriend "social/pkg/mdb/friend"
	protofriend "social/pkg/proto/friend"
	pkgserver "social/pkg/server"
)

var (
	app *Friend
)

func NewFriend(normal *pkgserver.Normal) *Friend {
	app = &Friend{
		Normal:     normal,
		mongodbMgr: new(libmongodb.Mgr),
	}
	app.bus.Normal = normal
	normal.Options.
		WithDefaultHandler(app.bus.OnEventBus).
		WithEtcdHandler(app.bus.OnEventEtcd).
		WithTimerEachSecond(&pkgserver.NormalTimerSecond{
			OnTimerFun: app.OnTimerEachSecondFun,
			Arg:        app,
		}).
		WithTimerEachDay(&pkgserver.NormalTimerSecond{
			OnTimerFun: app.OnTimerEachDayFun,
			Arg:        app,
		}).WithSubBench(&app.subBenchMgr)
	app.gateMgr.actorMgr = actor.NewMgr[string]()
	return app
}

type Friend struct {
	*pkgserver.Normal
	subBenchMgr subBenchMgr
	bus         Bus
	gateMgr     UserMgr
	mongodbMgr  *libmongodb.Mgr
}

func (p *Friend) String() string {
	return pkgserver.NameFriend
}

func (p *Friend) OnStart(ctx context.Context) (err error) {
	// 定时器-可用负载
	timerAvailableLoadExpireTimestamp = p.TimeMgr.ShadowTimeSecond()
	// 连接mongodb数据库
	if err := p.mongodbMgr.Connect(context.TODO(),
		libmongodb.NewOptions().
			WithAddrs(p.subBenchMgr.ZoneMongoDB.Addrs).
			WithUserName(p.subBenchMgr.ZoneMongoDB.User).
			WithPW(p.subBenchMgr.ZoneMongoDB.Password).
			WithDBName(p.subBenchMgr.ZoneMongoDB.DBName).
			WithMaxPoolSize(p.subBenchMgr.ZoneMongoDB.MaxPoolSize).
			WithMinPoolSize(p.subBenchMgr.ZoneMongoDB.MinPoolSize).
			WithTimeoutDuration(p.subBenchMgr.ZoneMongoDB.TimeoutDuration).
			WithMaxConnIdleTime(p.subBenchMgr.ZoneMongoDB.MaxConnIdleTime).
			WithMaxConnecting(p.subBenchMgr.ZoneMongoDB.MaxConnecting),
	); err != nil {
		return errors.WithMessagef(err, "%v %v", p.subBenchMgr.ZoneMongoDB, libruntime.Location())
	} else {
		p.mongodbMgr.SwitchedDatabase(pkgmdb.GenDBName(p.Normal.ZoneID, *p.subBenchMgr.ZoneMongoDB.DBName))
		pkgmdbfriend.Collection = p.mongodbMgr.SwitchedCollection(pkgmdbfriend.CollectionName)
	}

	// TODO ... 链接 redis
	// ...
	go func() { //启动grpc服务
		defer func() {
			if libutil.IsRelease() {
				if err := recover(); err != nil {
					p.LogMgr.Fatalf(libconsts.GoroutinePanic, err, debug.Stack())
				}
			}
			p.LogMgr.Fatalf(libconsts.GoroutineDone)
		}()
		addr := fmt.Sprintf("%v:%v", p.BenchMgr.Server.IP, p.BenchMgr.Server.Port)
		listen, err := net.Listen("tcp", addr)
		if err != nil {
			p.LogMgr.Fatalf("Failed to listen: %v", err)
		}
		p.GrpcServer = grpc.NewServer(grpc.MaxRecvMsgSize(1024 * 1024 * 1024)) // todo menglingchao 设置接受大小
		protofriend.RegisterFriendServiceServer(p.GrpcServer, &APIServer{})
		p.LogMgr.Tracef("%v is running on %v", p.String(), addr)
		if err = p.GrpcServer.Serve(listen); err != nil {
			p.LogMgr.Fatalf("Failed to serve: %v", err)
		}
	}()
	runtime.GC()
	return nil
}

func (p *Friend) OnPreStop(_ context.Context) (err error) {
	p.LogMgr.Warn("OnPreStop stop")
	{ // 关机前处理...业务逻辑
		p.LogMgr.Warn("grpc Service stop")
	}
	return nil
}
