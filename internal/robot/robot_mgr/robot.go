package robot_mgr

import (
	"context"
	"dawn-server/impl/common"
	"dawn-server/impl/common/msg"
	"dawn-server/impl/protobuf/battlegateway_proto"
	"dawn-server/impl/protobuf/common_proto"
	"dawn-server/impl/protobuf/room_proto"
	"dawn-server/impl/protobuf/world_proto"
	"dawn-server/impl/service/login/battle_gateway_mgr"
	"dawn-server/impl/service/login/login_msg"
	"dawn-server/impl/tool/robot/config"
	xrerror "dawn-server/impl/xr/lib/error"
	xrhttp "dawn-server/impl/xr/lib/http"
	xrlog "dawn-server/impl/xr/lib/log"
	xrpb_func "dawn-server/impl/xr/lib/pb"
	xrtcp "dawn-server/impl/xr/lib/tcp"
	xrtimer "dawn-server/impl/xr/lib/timer"
	xrutil "dawn-server/impl/xr/lib/util"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gogo/protobuf/proto"
)

type Robot struct {
	Account   string
	Name      string
	UID       uint64
	ExistUser bool   //是否有创角
	ServerId  uint32 // 所属服务器ID

	WorldTCPClient xrtcp.Client
	WorldIP        string
	WorldPort      uint16
	WorldToken     string

	BattleGatewayTCPClient xrtcp.Client
	BattleGatewayInfo      *battle_gateway_mgr.BattleGateway // 战斗网关信息.
	BattleGatewayToken     string

	IsCompleteLogin bool // 是否完成登录流程（包括http登录，连接world服、battle gateway服，并拉取用户数据）

	IsPause bool // 是否暂停

	InRoom            uint32               // 是否在房间中 需要用atomic原子操作 所以用了uint32类型
	InBattle          bool                 // 战斗中
	FrameMillisecond  *xrtimer.Millisecond // 发送帧数据定时器
	ServerFrameID     uint32               // 服务器返回的帧ID
	LastSendFrameID   uint32               // 上次发送给服务器的帧ID
	FrameHealthUpload uint32               // 服务器返回的帧同步健康值

	tickerStopChan        chan struct{} // 消息定时器关闭channel
	prepareMessageIDSlice []uint32      // 待发送的消息ID
	SendMessageTimeMap    sync.Map      // 用户发送各个消息的开始时间
	LastResultId          uint32        // 上个请求结果

	DBInventory *common_proto.DBInventory // 用户DB数据
}

func (p *Robot) Init() {
	p.tickerStopChan = make(chan struct{}, 1)
}

func (p *Robot) String() string {
	return fmt.Sprintf("client account:%v, name:%v uid:%v worldID:%v, worldPort:%v", p.Account, p.Name, p.UID, p.WorldIP, p.WorldPort)
}

func GenerateSessionID() uint32 {
	return uint32(xrutil.RandomInt(10000000, 99999999))
}

func (p *Robot) Send2World(pb proto.Message, messageID uint32, sessionID uint32, resultID uint32, frameID uint32) error {
	// 更新发出的请求数
	GProfiler.ChanSendNum <- 1
	GProfiler.ChanCycleSendNum <- 1

	// 如果有请求唯一ID 记录开始时间
	if sessionID > 0 {
		p.SendMessageTimeMap.Store(sessionID, time.Now())
	}

	xrlog.GetInstance().Debugf("robot:%s Send2World messageID:%#x sessionID:%d ", p.Name, messageID, sessionID)

	return p.WorldTCPClient.Remote.Send(&xrpb_func.UnserializedPacket{
		Header: &msg.CSProtoHead{
			MessageID: messageID,
			SessionID: sessionID,
			ResultID:  resultID,
			FrameID:   frameID,
		},
		Message: pb,
	})
}

func (p *Robot) Send2BGS(pb proto.Message, messageID uint32, sessionID uint32, resultID uint32, frameID uint32) error {
	// 更新发出的请求数
	GProfiler.ChanSendNum <- 1
	GProfiler.ChanCycleSendNum <- 1

	// 如果有请求唯一ID 记录开始时间
	if sessionID > 0 {
		p.SendMessageTimeMap.Store(sessionID, time.Now())
	}

	xrlog.GetInstance().Debugf("robot:%s Send2BGS messageID:%#x sessionID:%d frameID:%d", p.Name, messageID, sessionID, frameID)

	return p.BattleGatewayTCPClient.Remote.Send(&xrpb_func.UnserializedPacket{
		Header: &msg.CSProtoHead{
			MessageID: messageID,
			SessionID: sessionID,
			ResultID:  resultID,
			FrameID:   frameID,
		},
		Message: pb,
	})
}

func (p *Robot) Online(OnParsePacket xrtcp.OnUnmarshalPacket, OnPacket xrtcp.OnPacket, OnDisconnect xrtcp.OnDisconnect) error {
	// http登录
	if err := p.httpLogin(); err != nil {
		return err
	}

	// 连接world服并验证token
	if err := p.ConnectWorld(OnParsePacket, OnPacket, OnDisconnect); err != nil {
		return err
	}

	// 登录后请求的消息存在依赖关系：
	// 1. 在创角world_proto.WorldCreatUserMsg_CMD返回的消息处理中请求world_proto.WorldGetUserMsg_CMD
	// 2. 在world_proto.WorldGetUserMsg_CMD返回的消息处理中 检查并创建章节信息，再进入战网

	// 是否已创角
	if !p.ExistUser {
		err := GMessageFuncMap.Handler(world_proto.WorldCreatUserMsg_CMD, p)
		if err != nil {
			xrlog.GetInstance().Warnf("robot:%s sendMessage error:%v", p.Name, err)
		}
	} else {
		err := GMessageFuncMap.Handler(world_proto.WorldGetUserMsg_CMD, p)
		if err != nil {
			xrlog.GetInstance().Warnf("robot:%s sendMessage error:%v", p.Name, err)
		}
	}

	// 定时请求
	p.StartMessageTicker()

	return nil
}

func (p *Robot) httpLogin() error {
	res := &login_msg.LoginJsonRes{}
	lj := &login_msg.LoginJson{
		Account:         p.Account,
		Verify:          common.GenLoginVerify(p.Account, "dawn-taptap"),
		OpID:            9001,
		ClientDevice:    "ClientDevice-robot-" + p.Account,
		AdID:            9002,
		ChannelID:       9003,
		MultiscreenType: "pc",
		Model:           "centos-robto",
		Country:         "中国",
		Region:          "东亚",
		Timezone:        "+8",
		Language:        "简体中文",
		YzDeviceID:      "YzDeviceID-robot-" + p.Account,
		DeviceType:      "device_type",
		DeviceOS:        "device_os",
		OaID:            "oaid",
		GameVersion:     "beta01-202302080539",
	}

	// 随机一个login service
	addr := config.GRobotCfg.Base.LoginAddr
	slice := strings.Split(addr, ":")
	if port, err := strconv.Atoi(slice[1]); err != nil {
		return errors.WithMessagef(err, xrutil.GetCodeLocation(1).String())
	} else {
		if data, err := json.Marshal(lj); err != nil {
			return errors.WithMessagef(err, xrutil.GetCodeLocation(1).String())
		} else {
			if result, err := xrhttp.Post(slice[0], uint16(port), "/login", data); err != nil {
				return errors.WithMessagef(err, xrutil.GetCodeLocation(1).String())
			} else {
				// {"ip":"10.18.32.72","port":5071,"token":"caf35c895e8f88f49cab510360544d3a","version":"Beta.3.0","existuser":1,"errorcode":0,"battlegateway_slice":[{"zone_id":10005,"service_id":1,"tcp_addrs":":0","kcp_addrs":":0","version":""},{"zone_id":10007,"service_id":1,"tcp_addrs":"10.18.32.72:8071","kcp_addrs":"10.18.32.72:8071","version":"Beta.3.0"},{"zone_id":10004,"service_id":1,"tcp_addrs":"10.18.32.71:8041","kcp_addrs":"10.18.32.71:8041","version":"Beta.3.0.1"}]}
				xrlog.GetInstance().Debugf("login response:%s", result)

				if err = json.Unmarshal(result, res); err != nil {
					return errors.WithMessagef(err, xrutil.GetCodeLocation(1).String())
				}
				if res.ErrorCode != 0 {
					// 无可用的world服
					//if res.ErrorCode == xrerror.NonExistent.Code {
					if xrerror.NonExistent.Code > 0 {
						GProfiler.ChanConnectFailedNum <- 1
					}
					return fmt.Errorf("[%+v xrhttp.Post res error code::%v]", xrutil.GetCodeLocation(1), res.ErrorCode)
				}
				if len(res.BattleGateways) <= 0 {
					return fmt.Errorf("[%+v xrhttp.Post res BattleGateways is 0]", xrutil.GetCodeLocation(1))
				}

				p.ExistUser = res.ExistUser == 1

				// 从符合版本号的多个战网中随机取一个
				var battleGatewaySlice []*battle_gateway_mgr.BattleGateway
				for _, v := range res.BattleGateways {
					if len(config.GRobotCfg.Base.BattleVersion) > 0 && v.Version != config.GRobotCfg.Base.BattleVersion {
						continue
					}

					battleGatewaySlice = append(battleGatewaySlice, v)
				}
				if 0 == len(battleGatewaySlice) && config.GRobotCfg.Base.IsBattle {
					GProfiler.ChanConnectFailedNum <- 1
					return fmt.Errorf("[%+v err:%v]", xrutil.GetCodeLocation(1), fmt.Errorf("no-has BGS %v", config.GRobotCfg.Base.BattleVersion))
				}

				randomIndex := xrutil.RandomInt(0, len(battleGatewaySlice)-1)
				p.BattleGatewayInfo = battleGatewaySlice[randomIndex]
			}
		}
	}
	p.WorldIP = res.Ip
	p.WorldPort = res.Port
	p.WorldToken = res.Token

	return nil
}

func (p *Robot) ConnectWorld(OnParsePacket xrtcp.OnUnmarshalPacket, OnPacket xrtcp.OnPacket, OnDisconnect xrtcp.OnDisconnect) error {
	address := fmt.Sprintf("%v:%v", p.WorldIP, p.WorldPort)
	if err := p.WorldTCPClient.Connect(context.TODO(),
		xrtcp.NewClientOptions().
			SetAddress(address).
			SetEventChan(GRobotMgr.EventChan).
			SetSendChanCapacity(10000).
			SetPacket(&msg.SCPacket{}).
			SetOnUnmarshalPacket(OnParsePacket).
			SetOnPacket(OnPacket).
			SetOnDisconnect(OnDisconnect),
	); err != nil {
		GProfiler.ChanConnectFailedNum <- 1
		xrlog.GetInstance().Errorf("机器人:%s 连接world:%s 失败, err:%v", p.Name, address, err)
		return errors.WithMessagef(err, xrutil.GetCodeLocation(1).String())
	}

	req := &world_proto.WorldVerifyTokenMsg{
		Token:         p.WorldToken,
		BattleVersion: config.BattleVersion,
	}
	_ = p.Send2World(req, world_proto.WorldVerifyTokenMsg_CMD, GenerateSessionID(), 0, 0)

	return nil
}

func (p *Robot) ConnectBGW(OnParsePacket xrtcp.OnUnmarshalPacket, OnPacket xrtcp.OnPacket, OnDisconnect xrtcp.OnDisconnect) error {
	if err := p.BattleGatewayTCPClient.Connect(context.TODO(),
		xrtcp.NewClientOptions().
			SetAddress(p.BattleGatewayInfo.TCPAddrs).
			SetEventChan(GRobotMgr.EventChan).
			SetSendChanCapacity(1000).
			SetPacket(&msg.SCPacket{}).
			SetOnUnmarshalPacket(OnParsePacket).
			SetOnPacket(OnPacket).
			SetOnDisconnect(OnDisconnect),
	); err != nil {
		GProfiler.ChanConnectFailedNum <- 1
		xrlog.GetInstance().Errorf("机器人:%s 连接battle gateway:%s 失败, err:%v", p.Name, p.BattleGatewayInfo.TCPAddrs, err)
		return errors.WithMessagef(err, xrutil.GetCodeLocation(1).String())
	}

	xrlog.GetInstance().Errorf("机器人:%s 连接battle gateway:%s 成功", p.Name, p.BattleGatewayInfo.TCPAddrs)

	req := &battlegateway_proto.BattleGatewayTCPVerifyTokenMsg{
		Token: p.BattleGatewayToken,
		UID:   p.UID,
	}
	_ = p.Send2BGS(req, battlegateway_proto.BattleGatewayTCPVerifyTokenMsg_CMD, GenerateSessionID(), 0, 0)

	GRobotMgr.AddBGSRobot(p)

	return nil
}

// Offline 离线
func (p *Robot) Offline() error {
	p.WorldIP = ""
	p.WorldPort = 0
	p.WorldToken = ""
	p.BattleGatewayInfo = nil
	p.prepareMessageIDSlice = p.prepareMessageIDSlice[0:0]

	p.SendMessageTimeMap.Range(func(k, v interface{}) bool {
		p.SendMessageTimeMap.Delete(k)
		return true
	})

	p.UID = 0
	p.Account = ""
	p.Name = ""
	p.ExistUser = false
	p.ServerId = 0
	p.InRoom = 0
	p.InBattle = false
	p.IsCompleteLogin = false
	p.IsPause = false
	p.ServerFrameID = 0
	p.DBInventory = nil

	if p.FrameMillisecond != nil {
		xrtimer.DelMillisecond(p.FrameMillisecond)
	}

	p.tickerStopChan <- struct{}{}

	// 更新在线用户数
	GProfiler.ChanUserOffline <- p.Account

	_ = p.WorldTCPClient.ActiveDisconnect()

	if p.BattleGatewayTCPClient.IsConn() {
		_ = p.BattleGatewayTCPClient.ActiveDisconnect()
	}

	// 是否放进离线机器人中
	GRobotMgr.AddOfflineRobot(p)

	return nil
}

func (p *Robot) ResetBattle() {
	p.InBattle = false
	p.FrameMillisecond = nil
	p.ServerFrameID = 0
	p.LastSendFrameID = 0
	p.FrameHealthUpload = 0
}

func (p *Robot) StartMessageTicker() {
	ticker := time.NewTicker(time.Duration(config.GRobotCfg.Base.MessageInterval) * time.Second)
	go func(ticker *time.Ticker) {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				p.sendMessage()
			case <-p.tickerStopChan:
				return
			}
		}
	}(ticker)
}

func (p *Robot) sendMessage() {
	// 是否已完成登录流程
	if !p.IsCompleteLogin {
		xrlog.GetInstance().Errorf("robot:%s have not complete login", p.Name)
		return
	}

	// 是否暂停
	if p.IsPause {
		xrlog.GetInstance().Errorf("robot:%s is pause", p.Name)
		return
	}

	// 生成待发送的消息
	p.createPrepareMessage()

	// 取得待发送的消息
	messageID := p.popMessage()

	// 根据当前机器人状态检查发送的消息
	if !p.isMessageInvalid(messageID) {
		xrlog.GetInstance().Debugf("robot:%s messageID:%#x invalid, prepareMessageIDSlice:%v",
			p.Name, messageID, p.prepareMessageIDSlice)

		// 消息放回队列首位 等待下次ticker检查
		prepareMessageIDSliceCopy := make([]uint32, len(p.prepareMessageIDSlice)+1)
		prepareMessageIDSliceCopy[0] = messageID
		copy(prepareMessageIDSliceCopy[1:], p.prepareMessageIDSlice)
		p.prepareMessageIDSlice = prepareMessageIDSliceCopy

		return
	}

	// 发送消息
	err := GMessageFuncMap.Handler(messageID, p)
	if err != nil {
		xrlog.GetInstance().Debugf("robot:%s sendMessage error:%v", p.Name, err)
	}
}

func (p *Robot) isMessageInvalid(messageID uint32) bool {
	switch messageID {
	// 创建房间时 必须不在房间内
	case battlegateway_proto.BattleGatewayCreateRoomMsg_CMD:
		p.ResetBattle() // 测试战斗校验 未收到响应 先重置战斗 TODO
		//if 1 == atomic.LoadUint32(&p.InRoom) {
		//	return false
		//}
	default:
	}

	return true
}

func (p *Robot) setPrepareMessageByRobotStatus() {
	if 1 == atomic.LoadUint32(&p.InRoom) {
		xrlog.GetInstance().Warnf("robot:%s in room, append prepareMessageIDSlice with RoomEndMsg_CMD", p.Name)
		p.prepareMessageIDSlice = append(p.prepareMessageIDSlice, room_proto.RoomEndMsg_CMD)
	}
}

func (p *Robot) createPrepareMessage() {
	//// 如果当前行为组某次请求有错误码 清空后续请求 重新随机
	//if 0 != atomic.LoadUint32(&p.LastResultId) {
	//	p.prepareMessageIDSlice = p.prepareMessageIDSlice[0:0]
	//	atomic.StoreUint32(&p.LastResultId, 0)
	//	xrlog.GetInstance().Warnf("robot:%s lastResultID != 0, truncate prepareMessageIDSlice", p.Name)
	//	p.setPrepareMessageByRobotStatus()
	//}

	// 有消息未发送完
	if len(p.prepareMessageIDSlice) > 0 {
		return
	}

	// 按权重随机请求
	randAction := config.GetKeyByWeight(config.GActionWeightAll)

	idList, ok := config.GActionMessageIDAll[randAction]
	if !ok {
		xrlog.GetInstance().Errorf("robot:%s action:%s not found", p.Name, randAction)
		return
	}

	if 0 == len(idList) {
		xrlog.GetInstance().Errorf("robot:%s action:%s message id list empty", p.Name, randAction)
		return
	}

	p.prepareMessageIDSlice = append(p.prepareMessageIDSlice, idList...)

	xrlog.GetInstance().Debugf("robot:%s randAction:%s prepareMessageIDSlice:%v", p.Name, randAction, p.prepareMessageIDSlice)
}

func (p *Robot) popMessage() uint32 {
	messageID := p.prepareMessageIDSlice[0]
	p.prepareMessageIDSlice = p.prepareMessageIDSlice[1:]

	return messageID
}
