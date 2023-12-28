package gate

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	libactor "social/lib/actor"
	libconsts "social/lib/consts"
	liberror "social/lib/error"
	libetcd "social/lib/etcd"
	liblog "social/lib/log"
	libpb "social/lib/pb"
	libruntime "social/lib/runtime"
	libtime "social/lib/time"
	libtimer "social/lib/timer"
	libutil "social/lib/util"
	pkgcommon "social/pkg/common"
	pkgmsg "social/pkg/msg"
	protogate "social/pkg/proto/gate"
	"sync"
)

func (p *UserMgr) SpawnUser(key string, stream grpc.ServerStream) *User {
	p.lock.Lock()
	defer p.lock.Unlock()
	user := &User{
		key:    key,
		Stream: stream,
	}
	p.actorMgr.SpawnActor(context.Background(), key, libactor.NewOptions().WithDefaultHandler(user.OnMailBox))
	p.utilMgr.Add(key, user)
	return user
}

func (p *UserMgr) DeleteUser(key string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.actorMgr.DeleteActor(context.Background(), key)
	p.utilMgr.Del(key)
}

func (p *UserMgr) Find(key string) *User {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.utilMgr.Find(key)
}

func (p *UserMgr) FindByStream(stream grpc.ServerStream) *User {
	p.lock.Lock()
	defer p.lock.Unlock()

	// 自定义查找条件函数
	return p.utilMgr.FindOneWithCondition(
		func(key string, user *User) bool {
			return user.Stream == stream
		},
	)
}

func NewUserMgr() *UserMgr {
	userMgr := &UserMgr{
		utilMgr:  libutil.NewMgr[string, *User](),
		actorMgr: libactor.NewMgr[string](),
	}
	// todo menglingchao 修改channel的数量
	userMgr.busChannel = make(chan interface{}, 10000)
	go func() {
		defer func() {
			// 主事件channel报错 不recover
			app.LogMgr.Fatalf(libconsts.GoroutineDone)
		}()
		userMgr.OnBus()
	}()
	_ = app.userPBFunMgr.Register(protogate.GateRegisterReqCMD,
		libpb.NewMessage().SetHandler(GateRegisterReq).SetNewPBMessage(func() proto.Message { return new(protogate.GateRegisterReq) }))
	_ = app.userPBFunMgr.Register(protogate.GateLogoutReqCMD,
		libpb.NewMessage().SetHandler(GateLogoutReq).SetNewPBMessage(func() proto.Message { return new(protogate.GateLogoutReq) }))

	return userMgr
}

type UserMgr struct {
	actor      *libactor.Normal[string]
	actorMgr   *libactor.Mgr[string] // e.g.:1.lp.1
	utilMgr    *libutil.Mgr[string, *User]
	lock       sync.RWMutex
	busChannel chan interface{} //总线 channel
}

func (p *UserMgr) OnBus() {
	for {
		select {
		case v := <-p.busChannel:
			var err error
			switch t := v.(type) {
			case *busUserMgrAddUser:

				if t.IsValid() {
					t.Function(t.Arg)
				}
			case *libtimer.Millisecond:
				if t.IsValid() {
					t.Function(t.Arg)
				}
			case *libetcd.KV:
				err = p.EtcdMgr.Handler(t.Key, t.Value)
			default:
				if p.Options.defaultHandler == nil {
					p.LogMgr.Fatalf("non-existent event:%v %v", v, t)
				} else {
					err = p.Options.defaultHandler(v)
				}
			}
			if err != nil {
				liblog.PrintErr(v, err)
			}
			if libutil.IsDebug() {
				dt := libtime.NowTime().Sub(p.TimeMgr.Time).Milliseconds()
				if dt > 50 {
					p.LogMgr.Warnf("cost time50: %v Millisecond with event type:%T", dt, v)
				} else if dt > 20 {
					p.LogMgr.Warnf("cost time20: %v Millisecond with event type:%T", dt, v)
				} else if dt > 10 {
					p.LogMgr.Warnf("cost time10: %v Millisecond with event type:%T", dt, v)
				}
			}
		}
	}
}

func GateRegisterReq(protoHead libpb.IHeader, protoMessage proto.Message, obj interface{}) *liberror.Error {
	ph := protoHead.(*pkgmsg.Header)
	in := protoMessage.(*protogate.GateRegisterReq)
	stream := obj.(*grpc.ServerStream)
	app.LogMgr.Trace(ph, in, stream)
	//todo 从redis中验证token...

	// 处理 RegisterReq
	key := pkgcommon.GenerateServiceKey(in.ServiceKey.ZoneID, in.ServiceKey.ServiceName, in.ServiceKey.ServiceID)
	user := p.userMgr.Find(key)
	if user != nil { //注册过,返回错误码
		err = send2User(stream, protogate.GateRegisterResCMD, liberror.Duplicate.Code, &protogate.GateRegisterRes{})
		return errors.WithMessagef(err, "user already registered %v", libruntime.Location())
	}
	user = p.userMgr.FindByStream(stream)
	if user != nil { //注册过,返回错误码
		err = send2User(stream, protogate.GateRegisterResCMD, liberror.Duplicate.Code, &protogate.GateRegisterRes{})
		return errors.WithMessagef(err, "user already registered %v", libruntime.Location())
	}
	//注册
	user = p.userMgr.SpawnUser(key, stream)
	//注册-成功
	err = send2User(stream, protogate.GateRegisterResCMD, 0, &protogate.GateRegisterRes{})
	if err != nil {
		return errors.WithMessagef(err, "send2User %v", libruntime.Location())
	}
	return nil
}

func GateLogoutReq(protoHead libpb.IHeader, protoMessage proto.Message, obj interface{}) *liberror.Error {
	return nil
}
