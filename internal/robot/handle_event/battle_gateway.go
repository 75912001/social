package handle_event

import (
	"dawn-server/impl/common/msg"
	"dawn-server/impl/tool/robot/robot_mgr"
	xrerror "dawn-server/impl/xr/lib/error"
	xrlog "dawn-server/impl/xr/lib/log"
	xrpb "dawn-server/impl/xr/lib/pb"
	xrtcp "dawn-server/impl/xr/lib/tcp"
	xrutil "dawn-server/impl/xr/lib/util"
	"sync/atomic"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

func OnEventDisConnClientBattleGatewayTCP(remote *xrtcp.Remote) error {
	xrlog.GetInstance().Tracef("remote:%p", remote)
	return nil
}

func OnParsePacketFromBattleGateway(remote *xrtcp.Remote, data []byte) (parsePacket *xrtcp.Packet, err error) {
	ph := &msg.CSProtoHead{}
	ph.Unpack(data)

	var pbFunHandle *xrpb.Message
	var pb proto.Message
	if 0 == ph.ResultID {
		pbFunHandle = GBattleGatewayPbFunMgr.Find(ph.MessageID)
		if pbFunHandle == nil {
			return nil, errors.WithMessagef(xrerror.MessageIDNonExistent, xrutil.GetCodeLocation(1).String())
		}
		pb, err = pbFunHandle.Unmarshal(data[msg.GCSProtoHeadLength:])
		if err != nil {
			return nil, errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
		}
	}

	parsePacket = &xrtcp.Packet{
		Remote:  remote,
		Header:  ph,
		Message: pb,
		Entity:  pbFunHandle,
	}
	return parsePacket, nil
}

func OnEventPacketClientBattleGatewayTCP(parsePacket *xrtcp.Packet) error {
	//处理客户端消息
	robot := robot_mgr.GRobotMgr.FindBGSOnline(parsePacket.Remote)
	if robot == nil {
		return errors.WithMessagef(xrerror.Packet, "%+v client err. remote:%p", xrutil.GetCodeLocation(1), parsePacket.Remote)
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
				xrlog.GetInstance().Warnf("packet from battleGateway , costMs:%d, ph:%v, robot:%v", costMs, ph.String(), robot.String())
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
