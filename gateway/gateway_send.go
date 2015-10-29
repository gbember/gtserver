// gateway_send.go
package gateway

import (
	"encoding/binary"
	"net"

	"github.com/gbember/gt/logger"
	"github.com/gbember/gt/util"
	"github.com/gbember/gtserver/common"
	"github.com/gbember/gtserver/proto"
)

type Send struct {
	roleID  int32
	conn    net.Conn
	exitCnt chan struct{}
	recv    chan proto.Messager
	pk      *proto.Packet //协议encode缓存
}

func startSend(roleID int32, conn net.Conn, exitCnt chan struct{}) *Send {
	s := new(Send)
	s.roleID = roleID
	s.conn = conn
	s.exitCnt = exitCnt
	s.recv = make(chan proto.Messager, 200)
	s.pk = proto.NewWriter()
	go s.run()
	common.SetGW(roleID, s)
	return s
}

func (s *Send) run() {
	defer util.LogPanicStack()
	defer common.DeleteGW(s.roleID)
	for {
		select {
		case msg := <-s.recv:
			logger.Debug("发送消息:%v", msg)
			send_msg(s.conn, msg, s.pk)
		case <-s.exitCnt:
			return
		}
	}
}

//关闭网关
func (s *Send) Close() {
	s.conn.Close()
}

//发送消息到发送goroutine
func (s *Send) Send(msg proto.Messager) {
	s.recv <- msg
}

//直接发送消息
func (s *Send) RealSend(msg proto.Messager) {
	send_msg(s.conn, msg, s.pk)
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
