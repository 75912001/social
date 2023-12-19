package robot_mgr

import (
	"dawn-server/impl/protobuf/battlegateway_proto"
	"dawn-server/impl/protobuf/common_proto"
	"dawn-server/impl/protobuf/room_proto"
	"dawn-server/impl/protobuf/world_proto"
	"dawn-server/impl/tool/robot/config"
	xrerror "dawn-server/impl/xr/lib/error"
	xrlog "dawn-server/impl/xr/lib/log"
	xrutil "dawn-server/impl/xr/lib/util"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"sync/atomic"
	"time"
)

type messageFunc func(robot *Robot) error
type messageFuncMap map[uint32]messageFunc // 发送消息ID对应的方法

var GMessageFuncMap = make(messageFuncMap)

func init() {
	// 特殊请求
	GMessageFuncMap[0x0] = Offline
	GMessageFuncMap[0x1] = Disconnect
	GMessageFuncMap[0x2] = Invalid

	// world
	GMessageFuncMap[world_proto.WorldVerifyTokenMsg_CMD] = WorldVerifyTokenMsg
	GMessageFuncMap[world_proto.WorldGetUserMsg_CMD] = WorldGetUserMsg
	GMessageFuncMap[world_proto.WorldCreatUserMsg_CMD] = WorldCreatUserMsg
	GMessageFuncMap[world_proto.WorldStageCreateRegionMsg_CMD] = WorldStageCreateRegionMsg
	GMessageFuncMap[world_proto.WorldEnterBattleGatewayMsg_CMD] = WorldEnterBattleGatewayMsg
	GMessageFuncMap[world_proto.WorldSetPrimaryWeaponMsg_CMD] = WorldSetPrimaryWeaponMsg
	GMessageFuncMap[world_proto.WorldSetSecondaryWeaponMsg_CMD] = WorldSetSecondaryWeaponMsg
	GMessageFuncMap[world_proto.WorldTaskCompleteMsg_CMD] = WorldTaskCompleteMsg
	GMessageFuncMap[world_proto.WorldStatRoomEndMsg_CMD] = WorldStatRoomEndMsg
	GMessageFuncMap[world_proto.WorldUserChatMsg_CMD] = WorldUserChatMsg

	// battle gateway
	GMessageFuncMap[battlegateway_proto.BattleGatewayTCPVerifyTokenMsg_CMD] = BattleGatewayTCPVerifyTokenMsg
	GMessageFuncMap[battlegateway_proto.BattleGatewayKCPVerifyTokenMsg_CMD] = BattleGatewayKCPVerifyTokenMsg
	GMessageFuncMap[battlegateway_proto.BattleGatewayTCPHeartBeatMsg_CMD] = BattleGatewayTCPHeartBeatMsg
	GMessageFuncMap[battlegateway_proto.BattleGatewayKCPHeartBeatMsg_CMD] = BattleGatewayKCPHeartBeatMsg
	GMessageFuncMap[battlegateway_proto.BattleGatewayRoomChooseLevelMsg_CMD] = BattleGatewayRoomChooseLevelMsg
	GMessageFuncMap[battlegateway_proto.BattleGatewayGetRoomListMsg_CMD] = BattleGatewayGetRoomListMsg
	GMessageFuncMap[battlegateway_proto.BattleGatewayCreateRoomMsg_CMD] = BattleGatewayCreateRoomMsg
	GMessageFuncMap[battlegateway_proto.BattleGatewayJoinRoomMsg_CMD] = BattleGatewayJoinRoomMsg
	GMessageFuncMap[battlegateway_proto.BattleGatewayJoinBattleRoomMsg_CMD] = BattleGatewayJoinBattleRoomMsg
	GMessageFuncMap[battlegateway_proto.BattleGatewayTestMsg_CMD] = BattleGatewayTestMsg

	GMessageFuncMap[room_proto.RoomExitRoomMsg_CMD] = RoomExitRoomMsg
	GMessageFuncMap[room_proto.RoomGetRoomUserDetailInformationMsg_CMD] = RoomGetRoomUserDetailInformationMsg
	GMessageFuncMap[room_proto.RoomUserReadyMsg_CMD] = RoomUserReadyMsg
	GMessageFuncMap[room_proto.RoomStartMsg_CMD] = RoomStartMsg
	GMessageFuncMap[room_proto.RoomBattleStartMsg_CMD] = RoomBattleStartMsg
	GMessageFuncMap[room_proto.RoomFrameDataMsg_CMD] = RoomFrameDataMsg
	GMessageFuncMap[room_proto.RoomEndMsg_CMD] = RoomEndMsg

}

func (mf messageFuncMap) Handler(messageID uint32, robot *Robot) error {
	handler, ok := mf[messageID]
	if !ok {
		return errors.WithMessagef(xrerror.MessageIDNonExistent, xrutil.GetCodeLocation(1).String())
	}
	err := handler(robot)
	if err != nil {
		return errors.WithMessagef(err, xrutil.GetCodeLocation(1).String())
	}

	return nil
}

//************************************** 特殊请求 ******************************************************

// Offline 下线
func Offline(robot *Robot) error {
	xrlog.GetInstance().Debugf("robot:%s Offline", robot.Account)

	return robot.Offline()
}

// Disconnect 断开连接
func Disconnect(robot *Robot) error {
	xrlog.GetInstance().Debugf("robot:%s Disconnect", robot.Account)

	_ = robot.WorldTCPClient.ActiveDisconnect()

	if robot.BattleGatewayTCPClient.IsConn() {
		_ = robot.BattleGatewayTCPClient.ActiveDisconnect()
	}

	return nil
}

// Invalid 非法请求
func Invalid(robot *Robot) error {
	xrlog.GetInstance().Debugf("robot:%s send Invalid message", robot.Account)

	_ = robot.Send2World(nil, 0x0, GenerateSessionID(), 0, 0)
	_ = robot.Send2BGS(nil, 0x0, GenerateSessionID(), 0, 0)

	return nil
}

//************************************** world ******************************************************

func WorldVerifyTokenMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec WorldVerifyTokenMsg")

	req := &world_proto.WorldVerifyTokenMsg{
		Token:         robot.WorldToken,
		BattleVersion: config.BattleVersion,
	}
	_ = robot.Send2World(req, world_proto.WorldVerifyTokenMsg_CMD, GenerateSessionID(), 0, 0)

	return nil
}
func WorldGetUserMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec WorldGetUserMsg")

	_ = robot.Send2World(&world_proto.WorldGetUserMsg{}, world_proto.WorldGetUserMsg_CMD, GenerateSessionID(), 0, 0)

	return nil
}
func WorldCreatUserMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec WorldCreatUserMsg")

	req := &world_proto.WorldCreatUserMsg{
		Name:        robot.Name,
		CharacterID: uint32(xrutil.RandomInt(1, 2)),
		Head:        uint32(xrutil.RandomInt(1, 3)),
	}
	_ = robot.Send2World(req, world_proto.WorldCreatUserMsg_CMD, GenerateSessionID(), 0, 0)

	return nil
}
func WorldStageCreateRegionMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec WorldStageCreateRegionMsg")

	req := &world_proto.WorldStageCreateRegionMsg{
		Region: config.DefaultRegion,
	}
	_ = robot.Send2World(req, world_proto.WorldStageCreateRegionMsg_CMD, GenerateSessionID(), 0, 0)

	return nil
}
func WorldEnterBattleGatewayMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec WorldEnterBattleGatewayMsg")

	req := &world_proto.WorldEnterBattleGatewayMsg{
		ZoneID:    robot.BattleGatewayInfo.ZoneID,
		ServiceID: robot.BattleGatewayInfo.ServiceID,
	}
	_ = robot.Send2World(req, world_proto.WorldEnterBattleGatewayMsg_CMD, GenerateSessionID(), 0, 0)

	return nil
}
func WorldSetPrimaryWeaponMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec WorldSetPrimaryWeaponMsg")

	if robot.DBInventory == nil {
		return fmt.Errorf("[%+v err:DBInventory is nil]", xrutil.GetCodeLocation(1))
	}

	// 忽略已设置的主、副武器 从当前武器背包中取一个
	var setUUID uint64
	for uuid, _ := range robot.DBInventory.WeaponMap {
		if uuid == robot.DBInventory.PrimaryUUID || uuid == robot.DBInventory.SecondaryUUID {
			continue
		}
		setUUID = uuid
	}

	req := &world_proto.WorldSetPrimaryWeaponMsg{
		UUID: setUUID,
	}
	_ = robot.Send2World(req, world_proto.WorldSetPrimaryWeaponMsg_CMD, GenerateSessionID(), 0, 0)

	return nil
}
func WorldSetSecondaryWeaponMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec WorldSetSecondaryWeaponMsg")

	if robot.DBInventory == nil {
		return fmt.Errorf("[%+v err:DBInventory is nil]", xrutil.GetCodeLocation(1))
	}

	// 忽略已设置的主、副武器 从当前武器背包中取一个
	var setUUID uint64
	for uuid, _ := range robot.DBInventory.WeaponMap {
		if uuid == robot.DBInventory.PrimaryUUID || uuid == robot.DBInventory.SecondaryUUID {
			continue
		}
		setUUID = uuid
	}

	req := &world_proto.WorldSetSecondaryWeaponMsg{
		UUID: setUUID,
	}
	_ = robot.Send2World(req, world_proto.WorldSetSecondaryWeaponMsg_CMD, GenerateSessionID(), 0, 0)
	return nil
}
func WorldTaskCompleteMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec WorldTaskCompleteMsg")

	// 当前任务中随机取一个 TODO
	taskID := uint32(11)

	req := &world_proto.WorldTaskCompleteMsg{
		TaskID: taskID,
	}
	_ = robot.Send2World(req, world_proto.WorldTaskCompleteMsg_CMD, GenerateSessionID(), 0, 0)

	return nil
}

func WorldStatRoomEndMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec WorldStatRoomEndMsg")

	bid, _ := uuid.NewRandom()
	req := &world_proto.WorldStatRoomEndMsg{
		BattleUUID:      bid.String(),
		PrimaryWeaponID: 111,
		HeroID:          222,
		RegionAreaLevelMission: &common_proto.RegionAreaLevelMission{
			Region:  1001,
			Area:    1101,
			Level:   1001,
			Mission: []uint32{1, 2, 3},
		},
		Victory:        1,
		DurationSecond: 600,
		MissionsReport: []*common_proto.MissionsReport{
			{MissionsID: 1, Complete: 1},
			{MissionsID: 2, Complete: 1},
			{MissionsID: 3, Complete: 1},
		},
		UserReport: []*common_proto.BattleRoomUserReport{
			{
				UID:                      robot.UID,
				ParticipationNumerator:   1,
				ParticipationDenominator: 1,
				Performance:              1000,
				DeathCount:               10,
				Extraction:               20,
				Kills:                    30,
				Shots:                    40,
				Accidentals:              50,
				Accuracy:                 60,
				EXP:                      1000,
				TimeSec:                  600,
				SupportIds:               []uint32{1001, 1002},
				Wheel:                    1,
				ImageQuality:             2,
				SideType:                 1,
				SameScreen:               3,
				SupportSuccessCnt:        3,
				SupportFailedCnt:         1,
				SupportCancelCnt:         2,
			},
		},
	}

	_ = robot.Send2World(req, world_proto.WorldStatRoomEndMsg_CMD, GenerateSessionID(), 0, 0)
	return nil
}

func WorldUserChatMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec WorldUserChatMsg")

	req := &world_proto.WorldUserChatMsg{
		FromUID:  robot.UID,
		FromName: robot.Name,
		ToUID:    8000000006,
		Content:  "this is a test chat message",
	}

	_ = robot.Send2World(req, world_proto.WorldUserChatMsg_CMD, GenerateSessionID(), 0, 0)
	return nil
}

//*********************************** battle gateway **********************************************

func BattleGatewayTCPVerifyTokenMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec BattleGatewayTCPVerifyTokenMsg")

	req := &battlegateway_proto.BattleGatewayTCPVerifyTokenMsg{
		Token: robot.BattleGatewayToken,
		UID:   robot.UID,
	}
	_ = robot.Send2BGS(req, battlegateway_proto.BattleGatewayTCPVerifyTokenMsg_CMD, 0, 0, 0)

	return nil
}
func BattleGatewayKCPVerifyTokenMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec BattleGatewayKCPVerifyTokenMsg")

	req := &battlegateway_proto.BattleGatewayKCPVerifyTokenMsg{
		Token: robot.BattleGatewayToken,
		UID:   robot.UID,
	}
	_ = robot.Send2BGS(req, battlegateway_proto.BattleGatewayKCPVerifyTokenMsg_CMD, 0, 0, 0)

	return nil
}
func BattleGatewayTCPHeartBeatMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec BattleGatewayTCPHeartBeatMsg")
	// 此处不发送 已在battle_gateway_pb.go中收到服务端心跳时发送客户端心跳
	return nil
}
func BattleGatewayKCPHeartBeatMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec BattleGatewayKCPHeartBeatMsg")
	// 此处不发送 已在battle_gateway_pb.go中收到服务端心跳时发送客户端心跳
	return nil
}
func BattleGatewayRoomChooseLevelMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec BattleGatewayRoomChooseMissionMsg")

	req := &battlegateway_proto.BattleGatewayRoomChooseLevelMsg{
		Region: config.DefaultRegion,
		Area:   config.DefaultArea,
		Level:  config.DefaultLevel,
	}
	_ = robot.Send2BGS(req, battlegateway_proto.BattleGatewayRoomChooseLevelMsg_CMD, GenerateSessionID(), 0, 0)

	return nil
}
func BattleGatewayGetRoomListMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec BattleGatewayGetRoomListMsg")

	req := &battlegateway_proto.BattleGatewayGetRoomListMsg{
		RoomType:      uint32(common_proto.ROOM_TYPE_RT_ALL),
		BattleVersion: config.BattleVersion,
	}
	_ = robot.Send2BGS(req, battlegateway_proto.BattleGatewayGetRoomListMsg_CMD, GenerateSessionID(), 0, 0)

	return nil
}
func BattleGatewayCreateRoomMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec BattleGatewayCreateRoomMsg")

	req := &battlegateway_proto.BattleGatewayCreateRoomMsg{
		RoomType: uint32(common_proto.ROOM_TYPE_RT_ALL),
	}
	_ = robot.Send2BGS(req, battlegateway_proto.BattleGatewayCreateRoomMsg_CMD, GenerateSessionID(), 0, 0)

	return nil
}
func BattleGatewayJoinRoomMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec BattleGatewayJoinRoomMsg")

	// TODO 根据房间列表返回的数据来设置RoomServiceID和RoomID
	req := &battlegateway_proto.BattleGatewayJoinRoomMsg{
		RoomServiceID: 1,
		RoomID:        "1",
	}
	_ = robot.Send2BGS(req, battlegateway_proto.BattleGatewayJoinRoomMsg_CMD, GenerateSessionID(), 0, 0)

	return nil
}
func BattleGatewayJoinBattleRoomMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec BattleGatewayJoinBattleRoomMsg")

	// TODO 根据房间列表返回的数据来设置RoomServiceID和RoomID
	req := &battlegateway_proto.BattleGatewayJoinBattleRoomMsg{
		RoomServiceID: 1,
		RoomID:        "1",
	}
	_ = robot.Send2BGS(req, battlegateway_proto.BattleGatewayJoinBattleRoomMsg_CMD, GenerateSessionID(), 0, 0)

	return nil
}
func BattleGatewayTestMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec BattleGatewayTestMsg")

	req := &battlegateway_proto.BattleGatewayTestMsg{
		Timestamp: uint64(time.Now().Unix()),
	}
	_ = robot.Send2BGS(req, battlegateway_proto.BattleGatewayTestMsg_CMD, GenerateSessionID(), 0, 0)
	return nil
}
func RoomExitRoomMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec RoomExitRoomMsg")

	_ = robot.Send2BGS(&room_proto.RoomExitRoomMsg{}, room_proto.RoomExitRoomMsg_CMD, GenerateSessionID(), 0, 0)

	return nil
}
func RoomGetRoomUserDetailInformationMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec RoomGetRoomUserDetailInformationMsg")

	_ = robot.Send2BGS(&room_proto.RoomGetRoomUserDetailInformationMsg{}, room_proto.RoomGetRoomUserDetailInformationMsg_CMD, GenerateSessionID(), 0, 0)

	return nil
}
func RoomUserReadyMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec RoomUserReadyMsg")

	req := &room_proto.RoomUserReadyMsg{
		Ready: 1,
	}
	_ = robot.Send2BGS(req, room_proto.RoomUserReadyMsg_CMD, GenerateSessionID(), 0, 0)

	return nil
}
func RoomStartMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec RoomStartMsg")

	_ = robot.Send2BGS(&room_proto.RoomStartMsg{}, room_proto.RoomStartMsg_CMD, GenerateSessionID(), 0, 0)

	return nil
}

func RoomBattleStartMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec RoomBattleStartMsg")

	_ = robot.Send2BGS(&room_proto.RoomBattleStartMsg{LoadPercent: 100}, room_proto.RoomBattleStartMsg_CMD, GenerateSessionID(), 0, 0)

	return nil
}

func RoomFrameDataMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec RoomFrameDataMsg")

	// 模拟客户端帧同步数据 暂停机器人的定时器请求
	robot.IsPause = true

	// 发送1000次 间隔50ms
	for i := 0; i < 100; i++ {
		req := &room_proto.RoomFrameDataMsg{
			Data: config.GFrameData,
		}

		serverFrameId := atomic.LoadUint32(&robot.ServerFrameID)
		lastSendFrameID := serverFrameId + 1

		if lastSendFrameID <= robot.LastSendFrameID {
			lastSendFrameID = robot.LastSendFrameID + 1
		}
		robot.LastSendFrameID = lastSendFrameID
		_ = robot.Send2BGS(req, room_proto.RoomFrameDataMsg_CMD, 0, 0, lastSendFrameID)
		time.Sleep(time.Millisecond * 50)
	}

	robot.IsPause = false

	return nil
}

func RoomEndMsg(robot *Robot) error {
	xrlog.GetInstance().Debug("exec RoomEndMsg")

	req := &room_proto.RoomEndMsg{
		Victory:        1,
		DurationSecond: 500,
		MissionsReport: []*common_proto.MissionsReport{
			{MissionsID: config.DefaultLevel, Complete: 1},
		},
		UserReport: []*common_proto.BattleRoomUserReport{
			{UID: robot.UID, Kills: 10, Shots: 10, EXP: 1000},
		},
	}
	_ = robot.Send2BGS(req, room_proto.RoomEndMsg_CMD, GenerateSessionID(), 0, 0)

	return nil
}
