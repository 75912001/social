package robot_mgr

import (
	"context"
	"dawn-server/impl/common/dk"
	"dawn-server/impl/tool/robot/config"
	xrlog "dawn-server/impl/xr/lib/log"
	xrtcp "dawn-server/impl/xr/lib/tcp"
	xrtimer "dawn-server/impl/xr/lib/timer"
	xrutil "dawn-server/impl/xr/lib/util"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var LockOnlineRobot sync.RWMutex
var GRobotMgr RobotMgr

type RobotMap map[*xrtcp.Remote]*Robot

type RobotMgr struct {
	TimerMgr  xrtimer.Mgr
	TimeMgr   xrutil.TimeMgr
	EventChan chan interface{}

	OnlineRobotMap  RobotMap //在线
	OfflineRobotMap RobotMap //离线
	BGSRobotMap     RobotMap //在线

	ReadyRobotChan chan interface{} //准备就绪(即将上线) robot
}

func (p *RobotMgr) AddOnlineRobot(robot *Robot) {
	LockOnlineRobot.Lock()
	defer func() {
		LockOnlineRobot.Unlock()
	}()
	p.OnlineRobotMap[&robot.WorldTCPClient.Remote] = robot
}
func (p *RobotMgr) DelOnlineRobot(robot *Robot) {
	LockOnlineRobot.Lock()
	defer func() {
		LockOnlineRobot.Unlock()
	}()
	delete(p.OnlineRobotMap, &robot.WorldTCPClient.Remote)
}

func (p *RobotMgr) FindOnline(remote *xrtcp.Remote) *Robot {
	LockOnlineRobot.Lock()
	defer func() {
		LockOnlineRobot.Unlock()
	}()
	robot, _ := p.OnlineRobotMap[remote]
	return robot
}

func (p *RobotMgr) AddOfflineRobot(robot *Robot) {
	p.OfflineRobotMap[&robot.WorldTCPClient.Remote] = robot
}

func (p *RobotMgr) DelOfflineRobot(robot *Robot) {
	delete(p.OfflineRobotMap, &robot.WorldTCPClient.Remote)
}

func (p *RobotMgr) FindOffline(client *xrtcp.Client) *Robot {
	robot, _ := p.OfflineRobotMap[&client.Remote]
	return robot
}

func (p *RobotMgr) AddBGSRobot(robot *Robot) {
	LockOnlineRobot.Lock()
	defer func() {
		LockOnlineRobot.Unlock()
	}()
	p.BGSRobotMap[&robot.BattleGatewayTCPClient.Remote] = robot
}
func (p *RobotMgr) DelBGSRobot(robot *Robot) {
	LockOnlineRobot.Lock()
	defer func() {
		LockOnlineRobot.Unlock()
	}()
	delete(p.BGSRobotMap, &robot.BattleGatewayTCPClient.Remote)
}

func (p *RobotMgr) FindBGSOnline(remote *xrtcp.Remote) *Robot {
	LockOnlineRobot.Lock()
	defer func() {
		LockOnlineRobot.Unlock()
	}()
	robot, _ := p.BGSRobotMap[remote]
	return robot
}

func (p *RobotMgr) Init() error {
	p.OnlineRobotMap = make(RobotMap)
	p.OfflineRobotMap = make(RobotMap)
	p.BGSRobotMap = make(RobotMap)
	p.ReadyRobotChan = make(chan interface{}, config.GRobotCfg.Account.TotalNum)
	p.TimeMgr.Update()

	//log
	if err := xrlog.GetInstance().Start(context.TODO(),
		xrlog.NewOptions().
			SetLevel(xrlog.Level(config.GRobotCfg.Base.LogLevel)).
			SetAbsPath(config.GRobotCfg.Base.LogAbsPath).
			SetNamePrefix(fmt.Sprintf("%v", config.GRobotCfg.Account.AccountPre)),
	); err != nil {
		return errors.WithMessagef(err, xrutil.GetCodeLocation(1).String())
	}

	//eventChan
	p.EventChan = make(chan interface{}, config.GRobotCfg.Account.TotalNum*100)
	go func() {
		defer func() {
			if xrutil.IsRelease() {
				if err := recover(); err != nil {
					xrlog.GetInstance().Fatal(dk.GoroutinePanic, err, debug.Stack())
				}
			}
			xrlog.GetInstance().Trace(dk.GoroutineDone)
		}()
		p.handleEvent()
	}()

	secondDuration := time.Millisecond * 100
	millisecondDuration := time.Millisecond * 100
	//timer
	_ = p.TimerMgr.Start(context.TODO(),
		xrtimer.NewOptions().
			SetScanSecondDuration(&secondDuration).
			SetScanMillisecondDuration(&millisecondDuration).
			SetTimerOutChan(p.EventChan),
	)
	return nil
}

func (p *RobotMgr) handleEvent() {
	for v := range p.EventChan {
		p.TimeMgr.Update()
		switch t := v.(type) {
		case *xrtcp.EventDisconnect:
			if !t.Remote.IsConn() {
				continue
			}
			_ = t.Remote.Owner.OnDisconnect(t.Remote)
		case *xrtcp.Packet:
			if !t.Remote.IsConn() {
				continue
			}
			_ = t.Remote.Owner.OnPacket(t)
		case *xrtimer.Second:
			v1, _ := v.(*xrtimer.Second)
			if v1.IsValid() {
				v1.Function(v1.Arg)
			}
		case *xrtimer.Millisecond:
			v1, _ := v.(*xrtimer.Millisecond)
			if v1.IsValid() {
				v1.Function(v1.Arg)
			}
		default:
			xrlog.GetInstance().Fatalf("non-existent event:%v", v)
		}
	}
}
