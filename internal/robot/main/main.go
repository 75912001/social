package main

import (
	"dawn-server/impl/common/dk"
	"dawn-server/impl/tool/robot/config"
	"dawn-server/impl/tool/robot/handle_event"
	"dawn-server/impl/tool/robot/robot_mgr"
	xrlog "dawn-server/impl/xr/lib/log"
	xrutil "dawn-server/impl/xr/lib/util"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"runtime/debug"
	"syscall"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var err error
	var currentPath string
	currentPath, err = xrutil.GetCurrentPath()
	if err != nil {
		xrlog.PrintfErr("GetCurrentPath err:%v", err)
	}
	pathFile := path.Join(currentPath, "robot.yaml")

	config.GRobotCfg, err = config.ParseCfg(pathFile)
	if err != nil {
		xrlog.PrintfErr("ParseCfg err:%v", err)
		return
	}

	if err = robot_mgr.GRobotMgr.Init(); err != nil {
		xrlog.PrintfErr("robotMgr init err:%v", err)
		return
	}

	// 监控
	go robot_mgr.GProfiler.Watch()

	// 将机器人放入离线列表中
	for i := config.GRobotCfg.Account.AccountBegin; i < config.GRobotCfg.Account.AccountBegin+config.GRobotCfg.Account.TotalNum; i++ {
		robot := &robot_mgr.Robot{
			Account: fmt.Sprintf("%v%v", config.GRobotCfg.Account.AccountPre, i),
			Name:    fmt.Sprintf("%v%v", config.GRobotCfg.Account.AccountPre, i),
		}
		robot.Init()
		robot_mgr.GRobotMgr.AddOfflineRobot(robot)
	}

	// 记录机器人总数
	robot_mgr.GProfiler.TotalUserNum = int(config.GRobotCfg.Account.TotalNum)

	// 启动机器人
	for i := 0; i <= 8; i++ {
		go func() {
			defer func() {
				if xrutil.IsRelease() {
					if err := recover(); err != nil {
						xrlog.GetInstance().Fatalf(dk.GoroutinePanic, err, debug.Stack())
					}
				}
				xrlog.GetInstance().Fatal(dk.GoroutineDone)
			}()

			for v := range robot_mgr.GRobotMgr.ReadyRobotChan {
				switch v.(type) {
				case *robot_mgr.Robot:
					robot, _ := v.(*robot_mgr.Robot)
					if err = robot.Online(
						handle_event.OnParsePacketFromWorld,
						handle_event.OnEventPacketClientWorld,
						handle_event.OnEventDisConnClient,
					); err != nil {
						xrlog.GetInstance().Warnf("%v robot Online err. err:%v", robot.String(), err)
					} else {
						xrlog.GetInstance().Tracef("robot Online. %v", robot.String())
						robot_mgr.GRobotMgr.AddOnlineRobot(robot)
					}
				default:
					xrlog.GetInstance().Fatalf("%v", v)
					return
				}
			}
		}()
	}

	time.Sleep(time.Second * 1)
	robot_mgr.GRobotMgr.TimerMgr.AddSecond(offline2readyTimeout, nil, robot_mgr.GRobotMgr.TimeMgr.Second+int64(config.GRobotCfg.Base.CheckInterval))

	// 退出服务
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-sigChan

	xrlog.GetInstance().Warn("robot quit")
}

// 离线 -> 准备就绪
func offline2readyTimeout(_ interface{}) {
	robot_mgr.GProfiler.RunTicks++
	robot_mgr.GProfiler.ShowLog()

	for k, v := range robot_mgr.GRobotMgr.OfflineRobotMap {
		// 是否达到在线人数上限
		if robot_mgr.GProfiler.CurrentUserNum >= int(config.GRobotCfg.Account.OnlineNum) {
			break
		}

		// 更新在线用户数
		robot_mgr.GProfiler.AddCurrentUserNum(1)

		// 更新已登录用户总数
		robot_mgr.GProfiler.AddCurrentUserTotalNum(1)

		robot_mgr.GRobotMgr.ReadyRobotChan <- v
		delete(robot_mgr.GRobotMgr.OfflineRobotMap, k)
	}

	robot_mgr.GRobotMgr.TimerMgr.AddSecond(offline2readyTimeout, nil, robot_mgr.GRobotMgr.TimeMgr.Second+int64(config.GRobotCfg.Base.CheckInterval))
}
