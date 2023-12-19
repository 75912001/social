// Package config 逻辑配置
package config

import (
	"dawn-server/impl/xr/lib/util"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

const DefaultRegion = 1001 // 初始章节ID
const DefaultArea = 1101   // 初始区域ID
const DefaultLevel = 1001  // 初始关卡ID

const BattleVersion = 9000

// GFrameData 模拟器战斗房间玩家帧数据
var GFrameData = []byte{'1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1'}

var GActionWeightAll = make(map[string]int)         // 行为组及权重
var GActionMessageIDAll = make(map[string][]uint32) // 行为组所有消息ID

var GRobotCfg *RobotCfg

type RobotCfg struct {
	Account struct {
		AccountPre   string `yaml:"accountPre"`
		AccountBegin uint32 `yaml:"accountBegin"`
		TotalNum     uint32 `yaml:"totalNum"`
		OnlineNum    uint32 `yaml:"onlineNum"`
	}
	Base struct {
		LogAbsPath      string `yaml:"logAbsPath"`
		LogLevel        uint32 `yaml:"logLevel"`
		LoginAddr       string `yaml:"loginAddr"`
		BattleVersion   string `yaml:"battleVersion"`
		IsBattle        bool   `yaml:"isBattle"`
		CheckInterval   uint32 `yaml:"checkInterval"`
		MessageInterval uint32 `yaml:"messageInterval"`
	}
	Action []struct {
		Name     string `yaml:"name"`
		Desc     string `yaml:"desc"`
		Required bool   `yaml:"required"`
		Weight   int    `yaml:"weight"`
		Message  []struct {
			Id   uint32 `yaml:"id"`
			Name string `yaml:"name"`
			Desc string `yaml:"desc"`
		}
	}
}

func ParseCfg(pathFile string) (*RobotCfg, error) {
	robotCfg := &RobotCfg{}

	yamlFile, err := ioutil.ReadFile(pathFile)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, robotCfg)
	if err != nil {
		return nil, err
	}

	// 解析行为权重及行为消息列表
	messageIdList := make([]uint32, 0)
	for _, v := range robotCfg.Action {
		name := v.Name
		if 0 == len(name) {
			continue
		}
		GActionWeightAll[name] = v.Weight

		for _, messageInfo := range v.Message {
			messageIdList = append(messageIdList, messageInfo.Id)
		}

		messageIDCopy := make([]uint32, len(messageIdList))
		copy(messageIDCopy, messageIdList)
		messageIdList = messageIdList[0:0]

		GActionMessageIDAll[name] = messageIDCopy
	}

	return robotCfg, nil
}

func GetKeyByWeight(ActionWeight map[string]int) string {
	weightNumAll := 0
	for _, weight := range ActionWeight {
		weightNumAll += weight
	}

	var getKey string
	for moduleName, weight := range ActionWeight {
		if 0 == weight {
			continue
		}
		randNum := util.RandomInt(1, weightNumAll)
		if randNum <= weight {
			getKey = moduleName
			break
		}
		weightNumAll -= weight
	}

	return getKey
}
