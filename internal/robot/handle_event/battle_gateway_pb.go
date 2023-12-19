package handle_event

import (
	"context"
	"dawn-server/impl/common/msg"
	"dawn-server/impl/protobuf/battlegateway_proto"
	"dawn-server/impl/protobuf/room_proto"
	"dawn-server/impl/tool/robot/robot_mgr"
	xrerror "dawn-server/impl/xr/lib/error"
	xrlog "dawn-server/impl/xr/lib/log"
	xrpb "dawn-server/impl/xr/lib/pb"
	"sync/atomic"

	"github.com/gogo/protobuf/proto"
)

var GBattleGatewayPbFunMgr xrpb.Mgr

func init() {
	GBattleGatewayPbFunMgr.Init()

	_ = GBattleGatewayPbFunMgr.Register(battlegateway_proto.BattleGatewayTCPVerifyTokenMsg_CMD, xrpb.NewMessage().SetHandler(BattleGatewayTCPVerifyTokenMsgRes).SetNewPBMessage(func() proto.Message { return new(battlegateway_proto.BattleGatewayTCPVerifyTokenMsgRes) }))
	_ = GBattleGatewayPbFunMgr.Register(battlegateway_proto.BattleGatewayKCPVerifyTokenMsg_CMD, xrpb.NewMessage().SetHandler(BattleGatewayKCPVerifyTokenMsgRes).SetNewPBMessage(func() proto.Message { return new(battlegateway_proto.BattleGatewayKCPVerifyTokenMsgRes) }))
	_ = GBattleGatewayPbFunMgr.Register(battlegateway_proto.BattleGatewayGetRoomListMsg_CMD, xrpb.NewMessage().SetHandler(BattleGatewayGetRoomListMsgRes).SetNewPBMessage(func() proto.Message { return new(battlegateway_proto.BattleGatewayGetRoomListMsgRes) }))
	_ = GBattleGatewayPbFunMgr.Register(battlegateway_proto.BattleGatewayTCPHeartBeatMsgRes_CMD, xrpb.NewMessage().SetHandler(BattleGatewayTCPHeartBeatMsgRes).SetNewPBMessage(func() proto.Message { return new(battlegateway_proto.BattleGatewayTCPHeartBeatMsgRes) }))
	_ = GBattleGatewayPbFunMgr.Register(battlegateway_proto.BattleGatewayKCPHeartBeatMsgRes_CMD, xrpb.NewMessage().SetHandler(BattleGatewayKCPHeartBeatMsgRes).SetNewPBMessage(func() proto.Message { return new(battlegateway_proto.BattleGatewayKCPHeartBeatMsgRes) }))
	_ = GBattleGatewayPbFunMgr.Register(battlegateway_proto.BattleGatewayTestMsgRes_CMD, xrpb.NewMessage().SetHandler(BattleGatewayTestMsgRes).SetNewPBMessage(func() proto.Message { return new(battlegateway_proto.BattleGatewayTestMsgRes) }))

	_ = GBattleGatewayPbFunMgr.Register(room_proto.RoomCreateRoomMsg_CMD, xrpb.NewMessage().SetHandler(RoomCreateRoomMsgRes).SetNewPBMessage(func() proto.Message { return new(room_proto.RoomCreateRoomMsgRes) }))
	_ = GBattleGatewayPbFunMgr.Register(room_proto.RoomJoinRoomMsg_CMD, xrpb.NewMessage().SetHandler(RoomJoinRoomMsgRes).SetNewPBMessage(func() proto.Message { return new(room_proto.RoomJoinRoomMsgRes) }))
	_ = GBattleGatewayPbFunMgr.Register(room_proto.RoomJoinBattleRoomMsg_CMD, xrpb.NewMessage().SetHandler(RoomJoinBattleRoomMsgRes).SetNewPBMessage(func() proto.Message { return new(room_proto.RoomJoinBattleRoomMsgRes) }))
	_ = GBattleGatewayPbFunMgr.Register(room_proto.RoomExitRoomMsg_CMD, xrpb.NewMessage().SetHandler(RoomExitRoomMsgRes).SetNewPBMessage(func() proto.Message { return new(room_proto.RoomExitRoomMsgRes) }))
	_ = GBattleGatewayPbFunMgr.Register(room_proto.RoomKickUserMsg_CMD, xrpb.NewMessage().SetHandler(RoomKickUserMsgRes).SetNewPBMessage(func() proto.Message { return new(room_proto.RoomKickUserMsgRes) }))
	//在world pb中处理了 GPbFunMgr.Register(room_proto.RoomSetPrimaryWeaponMsg_CMD,  xrpb.NewEntity().SetHandler(RoomSetPrimaryWeaponMsgRes).SetNewMessage(func()  proto.Message { return new(room_proto.RoomSetPrimaryWeaponMsgRes)}))
	//在world pb中处理了 GPbFunMgr.Register(room_proto.RoomSetBTLHeroMsg_CMD,  xrpb.NewEntity().SetHandler(RoomSetBTLHeroMsgRes).SetNewMessage(func()  proto.Message { return new(room_proto.RoomSetBTLHeroMsgRes)}))
	_ = GBattleGatewayPbFunMgr.Register(room_proto.RoomGetRoomUserDetailInformationMsg_CMD, xrpb.NewMessage().SetHandler(RoomGetRoomUserDetailInformationMsgRes).SetNewPBMessage(func() proto.Message { return new(room_proto.RoomGetRoomUserDetailInformationMsgRes) }))
	_ = GBattleGatewayPbFunMgr.Register(room_proto.RoomUserReadyMsg_CMD, xrpb.NewMessage().SetHandler(RoomUserReadyMsgRes).SetNewPBMessage(func() proto.Message { return new(room_proto.RoomUserReadyMsgRes) }))
	_ = GBattleGatewayPbFunMgr.Register(room_proto.RoomModifyRoomTypeMsg_CMD, xrpb.NewMessage().SetHandler(RoomModifyRoomTypeMsgRes).SetNewPBMessage(func() proto.Message { return new(room_proto.RoomModifyRoomTypeMsgRes) }))
	_ = GBattleGatewayPbFunMgr.Register(room_proto.RoomStartMsgRes_CMD, xrpb.NewMessage().SetHandler(RoomStartMsgRes).SetNewPBMessage(func() proto.Message { return new(room_proto.RoomStartMsgRes) }))
	_ = GBattleGatewayPbFunMgr.Register(room_proto.RoomBattleStartMsgRes_CMD, xrpb.NewMessage().SetHandler(RoomBattleStartMsgRes).SetNewPBMessage(func() proto.Message { return new(room_proto.RoomBattleStartMsgRes) }))
	_ = GBattleGatewayPbFunMgr.Register(room_proto.RoomChooseLevelMsg_CMD, xrpb.NewMessage().SetHandler(RoomChooseLevelMsgRes).SetNewPBMessage(func() proto.Message { return new(room_proto.RoomChooseLevelMsgRes) }))
	_ = GBattleGatewayPbFunMgr.Register(room_proto.RoomTimeMsg_CMD, xrpb.NewMessage().SetHandler(RoomTimeMsgRes).SetNewPBMessage(func() proto.Message { return new(room_proto.RoomTimeMsgRes) }))
	_ = GBattleGatewayPbFunMgr.Register(room_proto.RoomGetFrameDataMsg_CMD, xrpb.NewMessage().SetHandler(RoomGetFrameDataMsgRes).SetNewPBMessage(func() proto.Message { return new(room_proto.RoomGetFrameDataMsgRes) }))
	_ = GBattleGatewayPbFunMgr.Register(room_proto.RoomSyncUserDataMsg_CMD, xrpb.NewMessage().SetHandler(RoomSyncUserDataMsgRes).SetNewPBMessage(func() proto.Message { return new(room_proto.RoomSyncUserDataMsgRes) }))
	_ = GBattleGatewayPbFunMgr.Register(room_proto.RoomFrameDataMsg_CMD, xrpb.NewMessage().SetHandler(RoomFrameDataMsgRes).SetNewPBMessage(func() proto.Message { return new(room_proto.RoomFrameDataMsgRes) }))
	_ = GBattleGatewayPbFunMgr.Register(room_proto.RoomVerifyHashMsg_CMD, xrpb.NewMessage().SetHandler(RoomVerifyHashMsgRes).SetNewPBMessage(func() proto.Message { return new(room_proto.RoomVerifyHashMsgRes) }))
}

func BattleGatewayTCPVerifyTokenMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	return nil
}
func BattleGatewayKCPVerifyTokenMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	return nil
}
func BattleGatewayGetRoomListMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	//ph := protoHead.(*proto_head.CSProtoHead)

	//robot := obj.(*robot_mgr.Robot)
	//in := (*protoMessage).(*world_proto.WorldGetRoomListMsgRes)
	//robot_mgr.GRobotMgr.Log.Trace(in.String())
	//
	//if robot.IsInRoom() { // 在房间, 离开房间
	//	reqExitRoot := &room_proto.RoomExitRoomMsg{}
	//	_ = msg.C2S(&robot.WorldTCP.Remote, reqExitRoot, room_proto.RoomExitRoomMsg_CMD, 0, 0, 0)
	//
	//	robot.InRoom = false
	//	robot.InBattle = false
	//}
	//{ // 加入新的房间
	//	for _, v := range in.RoomList {
	//		if v.UserMaxCnt <= v.UserCnt {
	//			continue
	//		}
	//		//加入房间
	//		req := &world_proto.WorldJoinRoomMsg{
	//			RoomID: v.RoomID,
	//		}
	//		_ = msg.C2S(&robot.WorldTCP.Remote, req, world_proto.WorldJoinRoomMsg_CMD, 0, 0, 0)
	//		break
	//	}
	//}

	return nil
}
func BattleGatewayTCPHeartBeatMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	robot := obj.(*robot_mgr.Robot)

	ph := protoHead.(*msg.CSProtoHead)

	// 收到服务端主动心跳 回复心跳
	_ = robot.Send2BGS(&battlegateway_proto.BattleGatewayTCPHeartBeatMsg{}, battlegateway_proto.BattleGatewayTCPHeartBeatMsg_CMD, ph.SessionID, 0, 0)
	return nil
}
func BattleGatewayKCPHeartBeatMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	robot := obj.(*robot_mgr.Robot)

	ph := protoHead.(*msg.CSProtoHead)

	// 收到服务端主动心跳 回复心跳
	_ = robot.Send2BGS(&battlegateway_proto.BattleGatewayKCPHeartBeatMsg{}, battlegateway_proto.BattleGatewayKCPHeartBeatMsg_CMD, ph.SessionID, 0, 0)
	return nil
}
func BattleGatewayTestMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	robot := obj.(*robot_mgr.Robot)
	in := (protoMessage).(*battlegateway_proto.BattleGatewayTestMsgRes)
	xrlog.GetInstance().Debugf("robot:%d receive BattleGatewayTestMsgRes, in:%v ", robot.UID, in.String())
	return nil
}
func RoomCreateRoomMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	robot := obj.(*robot_mgr.Robot)
	//robot.InRoom = true
	atomic.StoreUint32(&robot.InRoom, 1)
	return nil
}
func RoomJoinRoomMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	robot := obj.(*robot_mgr.Robot)
	//robot.InRoom = true
	atomic.StoreUint32(&robot.InRoom, 1)
	return nil
}
func RoomJoinBattleRoomMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	return nil
}
func RoomExitRoomMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	robot := obj.(*robot_mgr.Robot)
	//robot.InRoom = false
	atomic.StoreUint32(&robot.InRoom, 0)
	robot.ResetBattle()
	xrlog.GetInstance().Debugf("robot:%d receive RoomExitRoomMsgRes, robot.InRoom:%v ", robot.UID, robot.InRoom)

	return nil
}

func RoomKickUserMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	return nil
}
func RoomSetPrimaryWeaponMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	return nil
}
func RoomSetBTLHeroMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	return nil
}
func RoomGetRoomUserDetailInformationMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	return nil
}
func RoomUserReadyMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	return nil
}
func RoomModifyRoomTypeMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	return nil
}
func RoomStartMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	return nil
}
func RoomBattleStartMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	return nil
}
func RoomChooseLevelMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	return nil
}
func RoomTimeMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	return nil
}
func RoomGetFrameDataMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	return nil
}
func RoomSyncUserDataMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	return nil
}
func RoomFrameDataMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	ph := protoHead.(*msg.CSProtoHead)

	in := (protoMessage).(*room_proto.RoomFrameDataMsgRes)

	// TODO 返回的健康值处理
	lastHealthUpload := uint32(0)
	for _, frameData := range in.RoomFrame {
		lastHealthUpload = uint32(frameData.HealthUpload)
	}

	// 记录服务器返回的帧ID
	robot := obj.(*robot_mgr.Robot)
	atomic.StoreUint32(&robot.ServerFrameID, ph.FrameID)
	atomic.StoreUint32(&robot.FrameHealthUpload, lastHealthUpload)

	xrlog.GetInstance().Debugf("robot:%s receive RoomFrameDataMsgRes save frameID:%d", robot.Name, ph.FrameID)

	return nil
}

func RoomVerifyHashMsgRes(ctx context.Context, protoHead xrpb.IHeader, protoMessage proto.Message, obj interface{}) *xrerror.Error {
	return nil
}
