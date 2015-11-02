// role_client.go
package role

import (
	"sync"
	"time"

	"github.com/gbember/gt/logger"
	"github.com/gbember/gt/timer"
	"github.com/gbember/gt/util"
	"github.com/gbember/gtserver/common/gateway"
	"github.com/gbember/gtserver/common/role"
	"github.com/gbember/gtserver/proto"
)

type roleAgent struct {
	roleID    int32
	recv      chan proto.PMessage //接收route过来的客户端消息
	exitCnt   chan struct{}       //退出控制
	wgExitCnt *sync.WaitGroup     //退出等待控制
	isOK      bool                //是否无错误
	timer     *timer.Timer
	gw        gateway.Gateway
	rd        *role.RoleData
}

func (rc *roleAgent) start() {
	defer util.LogPanicStack()
	defer rc.stop()
	rc.timer = timer.NewTimer()
	rc.rd.Timer = rc.timer

	role.SetRoleClient(rc.rd)
	err := role.LoadData(rc.rd)
	if err != nil {
		logger.Error("加载数据出错:%s", err.Error())
		return
	}
	role.HookOnline(rc.rd)
	rc.isOK = true
	rc.addTimerFun(2*time.Minute, func() { persiteData(rc) })
	rc.loop()
}

func (rc *roleAgent) loop() {
	for {
		select {
		case msg := <-rc.recv:
			if f := role.GetHandle(msg.ID); f != nil {
				f(rc.rd, msg.Msg)
			} else {
				logger.Error("无效消息:%#v", msg)
			}
		case <-rc.timer.T.C:
			rc.timer.Run()
		case <-rc.exitCnt:
			return
		}
	}
}

func (rc *roleAgent) stop() {
	//关闭网关
	rc.gw.Close(0)
	defer rc.wgExitCnt.Done()
	rc.timer.Stop()
	role.HookOffline(rc.rd)
	if rc.isOK {
		//持久化数据
		rc.rd.PersiteData()
	}
}

func (rc *roleAgent) addTimerFun(d time.Duration, fun func()) {
	rc.timer.AddFun(d, fun)
}

func persiteData(rc *roleAgent) {
	rc.addTimerFun(2*time.Minute, func() { persiteData(rc) })
	err := rc.rd.PersiteData()
	if err != nil {
		logger.Error("数据持久化错误:%v", err)
	}
}
