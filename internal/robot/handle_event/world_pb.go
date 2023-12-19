package handle_event

import (
	"context"
	"dawn-server/impl/common/msg"
	"dawn-server/impl/protobuf/world_proto"
	"dawn-server/impl/tool/robot/config"
	"dawn-server/impl/tool/robot/robot_mgr"
	xrerror "dawn-server/impl/xr/lib/error"
	xrlog "dawn-server/impl/xr/lib/log"
	xrpb "dawn-server/impl/xr/lib/pb"
	"sync/atomic"

	"github.com/gogo/protobuf/proto"
)

var GWorldPbFunMgr xrpb.Mgr

func init() {
	GWorldPbFunMgr.Init()
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//常规[0x60000,0x60fff]
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	_ = GWorldPbFunMgr.Register(world_proto.WorldVerifyTokenMsg_CMD, xrpb.NewMessage().SetHandler(WorldVerifyTokenMsgRes).SetNewPBMessage(func() proto.Message { return new(world_proto.WorldVerifyTokenMsgRes) }))
	_ = GWorldPbFunMgr.Register(world_proto.WorldGetUserMsg_CMD, xrpb.NewMessage().SetHandler(WorldGetUserMsgRes).SetNewPBMessage(func() proto.Message { return new(world_proto.WorldGetUserMsgRes) }))
	_ = GWorldPbFunMgr.Register(world_proto.WorldCreatUserMsg_CMD, xrpb.NewMessage().SetHandler(WorldCreatUserMsgRes).SetNewPBMessage(func() proto.Message { return new(world_proto.WorldCreatUserMsgRes) }))
	_ = GWorldPbFunMgr.Register(world_proto.WorldUpdateRegionAreaLevelMsg_CMD, xrpb.NewMessage().SetHandler(WorldUpdateRegionAreaLevelMsgRes).SetNewPBMessage(func() proto.Message { return new(world_proto.WorldUpdateRegionAreaLevelMsgRes) }))
	_ = GWorldPbFunMgr.Register(world_proto.WorldEnterBattleGatewayMsg_CMD, xrpb.NewMessage().SetHandler(WorldEnterBattleGatewayMsgRes).SetNewPBMessage(func() proto.Message { return new(world_proto.WorldEnterBattleGatewayMsgRes) }))
	//_ = GWorldPbFunMgr.Register(room_proto.RoomSetPrimaryWeaponMsg_CMD, xrpb.NewMessage().SetHandler(RoomSetPrimaryWeaponMsgRes).SetNewPBMessage(func() proto.Message { return new(room_proto.RoomSetPrimaryWeaponMsgRes) }))
	_ = GWorldPbFunMgr.Register(world_proto.WorldSetSecondaryWeaponMsg_CMD, xrpb.NewMessage().SetHandler(WorldSetSecondaryWeaponMsgRes).SetNewPBMessage(func() proto.Message { return new(world_proto.WorldSetSecondaryWeaponMsgRes) }))
	//_ = GWorldPbFunMgr.Register(room_proto.RoomSetBTLHeroMsg_CMD, xrpb.NewMessage().SetHandler(RoomSetBTLHeroMsgRes).SetNewPBMessage(func() proto.Message { return new(room_proto.RoomSetBTLHeroMsgRes) }))
	_ = GWorldPbFunMgr.Register(world_proto.WorldGetZoneRegionAreaUserCntMsg_CMD, xrpb.NewMessage().SetHandler(WorldGetZoneRegionAreaUserCntMsgRes).SetNewPBMessage(func() proto.Message { return new(world_proto.WorldGetZoneRegionAreaUserCntMsgRes) }))
	_ = GWorldPbFunMgr.Register(world_proto.WorldNoticeoMsgRes_CMD, xrpb.NewMessage().SetHandler(WorldNoticeMsgRes).SetNewPBMessage(func() proto.Message { return new(world_proto.WorldNoticeoMsgRes) }))
	_ = GWorldPbFunMgr.Register(world_proto.WorldRoomEndMsgRes_CMD, xrpb.NewMessage().SetHandler(WorldRoomEndMsgRes).SetNewPBMessage(func() proto.Message { return new(world_proto.WorldRoomEndMsgRes) }))
	_ = GWorldPbFunMgr.Register(world_proto.WorldStageCreateRegionMsg_CMD, xrpb.NewMessage().SetHandler(WorldStageCreateRegionMsgRes).SetNewPBMessage(func() proto.Message { return new(world_proto.WorldStageCreateRegionMsgRes) }))
	_ = GWorldPbFunMgr.Register(world_proto.WorldTaskCompleteMsg_CMD, xrpb.NewMessage().SetHandler(WorldTaskCompleteMsgRes).SetNewPBMessage(func() proto.Message { return new(world_proto.WorldTaskCompleteMsgRes) }))
	_ = GWorldPbFunMgr.Register(world_proto.WorldStatRoomEndMsg_CMD, xrpb.NewMessage().SetHandler(WorldStatRoomEndMsg).SetNewPBMessage(func() proto.Message { return new(world_proto.WorldStatRoomEndMsg) }))
	_ = GWorldPbFunMgr.Register(world_proto.WorldUserChatMsg_CMD, xrpb.NewMessage().SetHandler(WorldUserChatMsg).SetNewPBMessage(func() proto.Message { return new(world_proto.WorldUserChatMsg) }))
}

func WorldVerifyTokenMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	return nil
}
func WorldCreatUserMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	u := obj.(*robot_mgr.Robot)
	_ = u.Send2World(&world_proto.WorldGetUserMsg{}, world_proto.WorldGetUserMsg_CMD, robot_mgr.GenerateSessionID(), 0, 0)

	return nil
}
func WorldGetUserMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	u := obj.(*robot_mgr.Robot)

	in := protoMessage.(*world_proto.WorldGetUserMsgRes)
	u.UID = in.UID
	u.DBInventory = in.DBInventory

	// 检查创建章节信息
	if config.GRobotCfg.Base.IsBattle {
		// 没有章节信息,创建章节信息. 有章节,进入战网
		if 0 == len(in.DBInventory.Stages) {
			req := &world_proto.WorldStageCreateRegionMsg{
				Region: config.DefaultRegion,
			}
			_ = u.Send2World(req, world_proto.WorldStageCreateRegionMsg_CMD, robot_mgr.GenerateSessionID(), 0, 0)
		} else {
			req := &world_proto.WorldEnterBattleGatewayMsg{
				ZoneID:    u.BattleGatewayInfo.ZoneID,
				ServiceID: u.BattleGatewayInfo.ServiceID,
			}
			_ = u.Send2World(req, world_proto.WorldEnterBattleGatewayMsg_CMD, robot_mgr.GenerateSessionID(), 0, 0)
		}
	} else {
		// 如果机器人不需要战斗 至此已完成登录流程
		u.IsCompleteLogin = true
	}

	return nil
}
func WorldStageCreateRegionMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	u := obj.(*robot_mgr.Robot)

	// 创建章节信息后进入战网
	req := &world_proto.WorldEnterBattleGatewayMsg{
		ZoneID:    u.BattleGatewayInfo.ZoneID,
		ServiceID: u.BattleGatewayInfo.ServiceID,
	}
	_ = u.Send2World(req, world_proto.WorldEnterBattleGatewayMsg_CMD, robot_mgr.GenerateSessionID(), 0, 0)

	return nil
}
func WorldEnterBattleGatewayMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	u := obj.(*robot_mgr.Robot)

	in := protoMessage.(*world_proto.WorldEnterBattleGatewayMsgRes)

	// 是否已连接过战网
	if !u.BattleGatewayTCPClient.IsConn() {
		u.BattleGatewayToken = in.Token
		if err := u.ConnectBGW(OnParsePacketFromBattleGateway, OnEventPacketClientBattleGatewayTCP, OnEventDisConnClientBattleGatewayTCP); err != nil {
			return xrerror.Link
		}
	}

	// 机器人需要战斗 至此完成登录流程
	u.IsCompleteLogin = true

	return nil
}
func WorldSetSecondaryWeaponMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	//	ph := protoHead.(*proto_head.CSProtoHead)

	//u := obj.(*robot_mgr.Robot)
	//
	////获取
	//_ = u.Send(&world_proto.WorldGetUserMsg{}, world_proto.WorldGetUserMsg_CMD, 0, 0, 0)

	return nil
}
func WorldGetZoneRegionAreaUserCntMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	//	ph := protoHead.(*proto_head.CSProtoHead)

	//u := obj.(*robot_mgr.Robot)
	//
	////获取
	//_ = u.Send(&world_proto.WorldGetUserMsg{}, world_proto.WorldGetUserMsg_CMD, 0, 0, 0)

	return nil
}
func WorldNoticeMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	//	ph := protoHead.(*proto_head.CSProtoHead)

	//u := obj.(*robot_mgr.Robot)

	////获取
	//_ = u.Send(&world_proto.WorldGetUserMsg{}, world_proto.WorldGetUserMsg_CMD, 0, 0, 0)

	return nil
}
func WorldRoomEndMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	u := obj.(*robot_mgr.Robot)
	//u.InRoom = false
	atomic.StoreUint32(&u.InRoom, 0)
	u.ResetBattle()
	xrlog.GetInstance().Debugf("robot:%s receive WorldRoomEndMsgRes, robot.InRoom:%v ", u.UID, u.InRoom)
	return nil
}

func WorldTaskCompleteMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	////	ph := protoHead.(*proto_head.CSProtoHead)
	//
	//u := obj.(*robot_mgr.Robot)
	//
	////获取
	//_ = u.Send(&world_proto.WorldGetUserMsg{}, world_proto.WorldGetUserMsg_CMD, 0, 0, 0)

	return nil
}
func WorldUpdateRegionAreaLevelMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	return nil
}
func WorldStatRoomEndMsg(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	return nil
}
func WorldUserChatMsg(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	u := obj.(*robot_mgr.Robot)
	ph := protoHead.(*msg.CSProtoHead)
	in := protoMessage.(*world_proto.WorldEnterBattleGatewayMsgRes)
	xrlog.GetInstance().Debugf("robot:%s receive head:%v, WorldUserChatMsg:%v ", u.UID, ph.String(), in.String())
	return nil
}
