package handle_event

import (
	"dawn-server/impl/common/msg"
	"dawn-server/impl/tool/robot/robot_mgr"
	xrerror "dawn-server/impl/xr/lib/error"
	xrlog "dawn-server/impl/xr/lib/log"
	xrtcp "dawn-server/impl/xr/lib/tcp"
	xrutil "dawn-server/impl/xr/lib/util"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

func OnEventDisConnClient(remote *xrtcp.Remote) error {
	xrlog.GetInstance().Tracef("remote:%p", remote)
	return nil
}

func OnEventPacketClientWorld(parsePacket *xrtcp.Packet) error {
	//处理客户端消息
	robot := robot_mgr.GRobotMgr.FindOnline(parsePacket.Remote)
	if robot == nil {
		return errors.WithMessagef(xrerror.Packet, "%+v client err. remote:%p", xrutil.GetCodeLocation(1).String(), parsePacket.Remote)
	}

	robot_mgr.GProfiler.ChanCycleResponseNum <- 1

	ph := parsePacket.Header.(*msg.CSProtoHead)

	xrlog.GetInstance().Debugf("robot:%s receive message, ph:%v", robot.Name, ph.String())

	// 记录请求耗时信息
	if 0 != ph.SessionID {
		startTimeVal, ok := robot.SendMessageTimeMap.Load(ph.SessionID)
		if ok {
			startTime := startTimeVal.(time.Time)
			costMs := (int)(time.Since(startTime).Milliseconds())
			robot_mgr.GProfiler.ChanCostTime <- costMs
			robot.SendMessageTimeMap.Delete(ph.SessionID)

			if costMs > 100 {
				xrlog.GetInstance().Warnf("packet from world , costMs:%d, ph:%v, robot:%v", costMs, ph.String(), robot.String())
			}
		}
	}

	if ph.ResultID != 0 {
		xrlog.GetInstance().Warnf("proto head result != 0 ph:%v robot:%v", ph.String(), robot.String())
		robot_mgr.GProfiler.ChanErrorNum <- 1
		atomic.StoreUint32(&robot.LastResultId, ph.ResultID)
		return nil
	}

	ret := parsePacket.Entity.Handler(nil, ph, parsePacket.Message, robot)
	if ret != nil {
		xrlog.GetInstance().Warnf("%v, ret:%#x", ph.String(), ret)
		robot_mgr.GProfiler.ChanErrorNum <- 1
	} else {
		robot_mgr.GProfiler.ChanSuccessNum <- 1
		robot_mgr.GProfiler.ChanCycleSuccessSendNum <- 1
	}

	return nil
}

func OnParsePacketFromWorld(remote *xrtcp.Remote, data []byte) (parsePacket *xrtcp.Packet, err error) {
	h := &msg.CSProtoHead{}
	h.Unpack(data)

	entity := GWorldPbFunMgr.Find(h.MessageID)
	if entity == nil {
		return nil, errors.WithMessage(xrerror.MessageIDNonExistent, fmt.Sprintf("h:%s %v", h.String(), xrutil.GetCodeLocation(1).String()))
	}
	message, err := entity.Unmarshal(data[msg.GCSProtoHeadLength:])
	if err != nil {
		return nil, errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}
	parsePacket = &xrtcp.Packet{
		Remote:  remote,
		Header:  h,
		Message: message,
		Entity:  entity,
	}
	return parsePacket, nil
}
