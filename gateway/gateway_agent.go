// gateway_agent.go
package gateway

import (
	"encoding/binary"
	"net"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gbember/gt/logger"
	"github.com/gbember/gt/network"
	"github.com/gbember/gt/network/msg"
	"github.com/gbember/gt/util"
	cgw "github.com/gbember/gtserver/common/gateway"
	"github.com/gbember/gtserver/common/name"
	"github.com/gbember/gtserver/config"
	"github.com/gbember/gtserver/db"
	"github.com/gbember/gtserver/proto"
	"github.com/gbember/gtserver/role"
	"github.com/gbember/gtserver/types"
)

const (
	//3秒
	_MAX_MSG_SPEED_NS = int64(3 * time.Second)
	//消息在时间段内的最大个数
	_MAX_MSG_SPEED_NUM = 30
	//速度超过几次之后踢下线
	_LAGER_SPEED_MAX_NUM = 3
)

var (
	//名字匹配(字母数字汉字组合)
	nameRe = regexp.MustCompile("^[0-9a-zA-Z\u4e00-\u9fa5]*$")
)

type gateway_agent struct {
	roleID            int32               //角色ID
	conn              net.Conn            //连接
	msgParser         msg.MsgParser       //消息解析器
	isCreated         bool                //是否创建角色
	receiveMsgNum     int                 //接收消息个数(N个之后赋值0)
	lastCheckSpeedSec int64               //最后一次检查消息速度的时间
	lagerSpeedNum     int                 //超过速度次数
	pk                *proto.Packet       //协议encode缓存
	closeMut          sync.Mutex          //关闭锁
	isClosed          bool                //是否关闭
	exitCnt           chan struct{}       //退出控制
	wgExitCnt         sync.WaitGroup      //等待发送和role的goroutine退出
	roleRecv          chan proto.PMessage //role 接收消息
	recv              chan proto.Messager
	rpk               *proto.Packet
}

func NewAgent(conn net.Conn, msgParser msg.MsgParser) network.TCPAgent {
	agent := new(gateway_agent)
	agent.conn = conn
	agent.msgParser = msgParser
	agent.lastCheckSpeedSec = time.Now().Unix()
	agent.pk = proto.NewWriter()
	agent.exitCnt = make(chan struct{})
	logger.Debug("socket:%v", conn.RemoteAddr())
	return agent
}

func (agent *gateway_agent) Run() {
	defer util.LogPanicStack()
	for {
		agent.conn.SetReadDeadline(time.Now().Add(120 * time.Second))
		dataBytes, err := agent.msgParser.Read(agent.conn)
		if err != nil {
			if !agent.isClosed {
				if ne, ok := err.(net.Error); ok {
					logger.Debug("网关消息接收错误:%v", ne)
				}
			}
			return
		}
		if agent.check_msg_speed_lager() {
			logger.Error("发送消息太快")
			agent.Close(4)
			return
		}
		id, msg, err := proto.DecodeProto(dataBytes)
		if err != nil {
			logger.Debug("网关协议解析错误:%s", err.Error())
			return
		}
		logger.Debug("receive msg:%d  %#v", id, msg)
		agent.dispatcher(id, msg)
	}
}

//检查消息速度  返回true表示消息速度太快
func (agent *gateway_agent) check_msg_speed_lager() bool {
	agent.receiveMsgNum++
	if agent.receiveMsgNum >= config.DataSetting.MsgMaxSpeedNum {
		lastCheckSpeedSec := time.Now().Unix()
		if lastCheckSpeedSec-agent.lastCheckSpeedSec <= config.DataSetting.MsgMaxSpeedSec*int64(time.Second) {
			agent.lagerSpeedNum++
			if agent.lagerSpeedNum > config.DataSetting.MsgMaxSpeedLager {
				return true
			}
		}
		agent.lastCheckSpeedSec = lastCheckSpeedSec
		agent.receiveMsgNum = 0
	}
	return false
}

//登陆
func (agent *gateway_agent) handle_login(msg proto.Messager) {
	if agent.roleID != 0 {
		send_msg(agent.conn, &proto.Sc_account_login{Result: 1}, agent.pk)
		return
	}
	msgLogin := msg.(*proto.Cs_account_login)
	//服务器ID
	if msgLogin.ServerID != config.DataSetting.ServerID {
		send_msg(agent.conn, &proto.Sc_account_login{Result: 2}, agent.pk)
		return
	}
	//超时
	if msgLogin.UnixTime+10 < int32(time.Now().Second()) {
		send_msg(agent.conn, &proto.Sc_account_login{Result: 3}, agent.pk)
		return
	}

	Tiket := util.MD5(msgLogin.AccountName + strconv.Itoa(int(msgLogin.UserID)) +
		strconv.Itoa(int(msgLogin.UnixTime)) + config.DataSetting.LoginAuthKey)
	if msgLogin.Ticket != Tiket {
		send_msg(agent.conn, &proto.Sc_account_login{Result: 4}, agent.pk)
		return
	}
	agent.roleID = msgLogin.UserID
	//注册agent
	registerAgent(agent)
	role, _ := db.FindByRoleID(msgLogin.UserID)
	if role.RoleID == 0 {
		send_msg(agent.conn, &proto.Sc_account_login{IsCreate: false}, agent.pk)
	} else {
		send_msg(agent.conn, &proto.Sc_account_login{IsCreate: true}, agent.pk)
		agent.isCreated = true
		agent.enter_game(role)
	}
}

//创建角色
func (agent *gateway_agent) handle_create(msg proto.Messager) {
	if agent.roleID == 0 {
		send_msg(agent.conn, &proto.Sc_account_login{Result: 1}, agent.pk)
		return
	}
	if agent.isCreated {
		send_msg(agent.conn, &proto.Sc_account_create{Result: 2}, agent.pk)
		return
	}
	msgCreate := msg.(*proto.Cs_account_create)
	if msgCreate.Sex != 0 && msgCreate.Sex != 1 {
		send_msg(agent.conn, &proto.Sc_account_create{Result: 3}, agent.pk)
		return
	}
	roleName := strings.TrimSpace(msgCreate.RoleName)
	if l := len(roleName); l < 6 || l > 18 {
		send_msg(agent.conn, &proto.Sc_account_create{Result: 4}, agent.pk)
		return
	}
	if !nameRe.MatchString(roleName) {
		send_msg(agent.conn, &proto.Sc_account_create{Result: 5}, agent.pk)
		return
	}
	//注册使用名字
	if !name.RegisterRoleName(agent.roleID, roleName) {
		send_msg(agent.conn, &proto.Sc_account_create{Result: 6}, agent.pk)
		return
	}
	role := &types.RoleInfo{RoleID: agent.roleID, RoleName: roleName, Sex: msgCreate.Sex, Level: 1}
	db.InsertRoleInfo(role)
	send_msg(agent.conn, &proto.Sc_account_create{}, agent.pk)
	agent.isCreated = true
	agent.enter_game(role)
}

//心跳包
func (agent *gateway_agent) handle_heart(proto.Messager) {
	sec := int32(time.Now().Unix())
	send_msg(agent.conn, &proto.Sc_account_heart{UnixTime: sec}, agent.pk)
}

//注销
func (agent *gateway_agent) handle_logout(proto.Messager) {
	agent.Close(0)
}

//进入游戏
func (agent *gateway_agent) enter_game(r *types.RoleInfo) {
	agent.recv = make(chan proto.Messager, 200)
	agent.rpk = proto.NewWriter()
	go agent.send_loop()
	cgw.SetGW(agent.roleID, agent)
	roleRecv := make(chan proto.PMessage, 1)
	agent.roleRecv = roleRecv
	agent.wgExitCnt.Add(1)
	role.Start(r, roleRecv, agent.exitCnt, &agent.wgExitCnt, agent)
}

//消息分发
func (agent *gateway_agent) dispatcher(msgID uint16, msg proto.Messager) {
	id := msgID / 100
	if id == 100 {
		switch msgID {
		case proto.CS_ACCOUNT_HEART:
			agent.handle_heart(msg)
		case proto.CS_ACCOUNT_LOGIN:
			agent.handle_login(msg)
		case proto.CS_ACCOUNT_CREATE:
			agent.handle_create(msg)
		default:
			//分发消息到role处理
			agent.roleRecv <- proto.PMessage{ID: msgID, Msg: msg}
		}
	} else {
		//分发消息到role处理
		agent.roleRecv <- proto.PMessage{ID: msgID, Msg: msg}
	}
}

func (agent *gateway_agent) send_loop() {
	defer util.LogPanicStack()
	for {
		select {
		case msg := <-agent.recv:
			logger.Debug("发送消息:%v", msg)
			send_msg(agent.conn, msg, agent.rpk)
		case <-agent.exitCnt:
			return
		}
	}
}

func (agent *gateway_agent) Close(reason int8) {
	if !agent.isClosed {
		agent.closeMut.Lock()
		if !agent.isClosed {
			agent.isClosed = true
			agent.closeMut.Unlock()
			cgw.DeleteGW(agent.roleID)
			//不是异地登陆
			if reason != 3 {
				unRegisterAgent(agent)
			}
			//发送关闭消息
			if reason != 0 {
				send_msg(agent.conn, &proto.Sc_account_kick{Reason: reason}, agent.pk)
			}
			agent.conn.Close()
			close(agent.exitCnt)
			agent.wgExitCnt.Wait()
			logger.Debug("gw exit====:%d", agent.roleID)
		} else {
			agent.closeMut.Unlock()
		}
	}
}

//发送消息到发送goroutine
func (agent *gateway_agent) Send(msg proto.Messager) {
	agent.recv <- msg
}

//直接发送消息
func (agent *gateway_agent) RealSend(msg proto.Messager) {
	send_msg(agent.conn, msg, agent.pk)
}

//发送proto消息
func send_msg(conn net.Conn, msg proto.Messager, pk *proto.Packet) {
	logger.Debug("发送消息:%#v", msg)
	pk.SeekTo(2)
	data := proto.EncodeProtoPacket(msg, pk)
	bs := data[0:2]
	binary.BigEndian.PutUint16(bs, uint16(len(data)-2))
	conn.Write(data)
}
