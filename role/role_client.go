// role_client.go
package role

import (
	"sync"
	"time"

	"github.com/gbember/gt/logger"
	"github.com/gbember/gt/util"
	"github.com/gbember/gtserver/common/gateway"
	"github.com/gbember/gtserver/proto"
	"github.com/gbember/gtserver/types"
)

type roleClient struct {
	roleID        int32
	recv          chan proto.PMessage //接收route过来的客户端消息
	exitCnt       chan struct{}       //退出控制
	wgExitCnt     *sync.WaitGroup     //退出等待控制
	gw            gateway.Gateway     //网关
	isOK          bool                //是否无错误
	persistTicker *time.Ticker        //持久化Ticker
	hourTimer     *time.Timer         //整点Timer
	dataRoleInfo  *types.RoleInfo
}

func (rc *roleClient) start() {
	defer util.LogPanicStack()
	defer rc.stop()
	setRoleClient(rc)
	err := rc.load_data()
	if err != nil {
		logger.Error("加载数据出错:%s", err.Error())
		return
	}
	rc.isOK = true
	rc.persistTicker = time.NewTicker(time.Minute * 2)
	now := time.Now()
	h, m, s := now.Clock()
	sec := time.Hour*time.Duration(24-h) - time.Minute*time.Duration(m) - time.Second*time.Duration(s)
	rc.hourTimer = time.NewTimer(time.Duration(sec) * time.Second)
	rc.loop()
}

func (rc *roleClient) loop() {
	for {
		select {
		case msg := <-rc.recv:
			if f := get_handle(msg.ID); f != nil {
				f(rc, msg.Msg)
			} else {
				logger.Error("无效消息:%#v", msg)
			}
		case <-rc.persistTicker.C:
			rc.persist_data()
		case <-rc.hourTimer.C:
			rc.hook_integral_hour()
		case <-rc.exitCnt:
			return
		}
	}
}

func (rc *roleClient) stop() {
	//关闭网关
	rc.gw.Close(0)
	defer rc.wgExitCnt.Done()
	if rc.persistTicker != nil {
		rc.persistTicker.Stop()
	}
	if rc.hourTimer != nil {
		rc.hourTimer.Stop()
	}
	if rc.isOK {
		//持久化数据
		rc.persist_data()
	}
}

//加载数据
func (rc *roleClient) load_data() error {
	return nil
}

//上线操作(所有数据都已经加载)
func (rc *roleClient) role_online() {

}

//下线
func (rc *roleClient) role_offline() {

}

//持久化
func (rc *roleClient) persist_data() {

}

//整点回调
func (rc *roleClient) hook_integral_hour() {
	now := time.Now()
	h, m, s := now.Clock()
	sec := time.Hour*time.Duration(24-h) - time.Minute*time.Duration(m) - time.Second*time.Duration(s)
	rc.hourTimer = time.NewTimer(time.Duration(sec) * time.Second)

	switch h {
	case 0:
	default:
	}
}
