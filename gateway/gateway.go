// gateway.go
package gateway

import (
	"strconv"

	"github.com/gbember/gt/logger"
	"github.com/gbember/gt/module"
	"github.com/gbember/gt/network"
	"github.com/gbember/gt/network/msg"
	"github.com/gbember/gt/util/concurrent"
	"github.com/gbember/gtserver/config"
)

type gateway struct {
	server     *network.TCPServer
	addr       string                    //监听地址
	maxConnNum int                       //最高在线数
	headLen    int                       //消息长度字节数
	msgIDLen   int                       //消息ID字节数
	maxDataLen int                       //消息最大字节数
	agants     *concurrent.ConcurrentMap //所有在线网关
}

var globalGW *gateway

func RegisterModule() {
	gw := new(gateway)
	gw.addr = ":" + strconv.Itoa(config.DataSetting.GatewayPort)
	gw.maxConnNum = config.DataSetting.MaxOnlineNum
	gw.headLen = config.DataSetting.HeadLen
	gw.msgIDLen = config.DataSetting.MsgIDLen
	gw.maxDataLen = config.DataSetting.MaxDataLen
	gw.agants = concurrent.NewConcurrentMap(config.DataSetting.MaxOnlineNum)
	module.Register(gw)
	globalGW = gw
}

//注册agent
func registerAgent(agent *gateway_agent) {
	if v := globalGW.agants.Replace(agent.roleID, agent); v != nil {
		//踢下线
		v.(*gateway_agent).Close(3)
	}
}

//注销agent
func unRegisterAgent(agent *gateway_agent) {
	globalGW.agants.Delete(agent.roleID)
}

func (gw *gateway) OnInit() {
	msgParser, err := msg.NewMsgParserProtobuf(gw.headLen, gw.msgIDLen, gw.maxDataLen)
	if err != nil {
		panic(err)
	}
	server, err := network.StartTCPServer(gw.addr, gw.maxConnNum, msgParser, NewAgent)
	if err != nil {
		panic(err)
	}
	gw.server = server
	logger.Info("gateway start...")
}

func (gw *gateway) OnDestroy() {
	if gw.server != nil {
		gw.server.Close()
	}
}

func (gw *gateway) Run(chan bool) {
}
