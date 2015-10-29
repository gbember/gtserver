// gateway.go
package common

import (
	"github.com/gbember/gt/util/concurrent"
	"github.com/gbember/gtserver/config"
	"github.com/gbember/gtserver/proto"
)

type Gateway interface {
	Send(msg proto.Messager)
	RealSend(msg proto.Messager)
	Close()
}

var (
	gateways *concurrent.ConcurrentMap
)

func Init() {
	gateways = concurrent.NewConcurrentMap(config.DataSetting.MaxOnlineNum)
}

func GetGW(roleID int32) Gateway {
	if v := gateways.Get(roleID); v != nil {
		return v.(Gateway)
	}
	return nil
}

func SetGW(roleID int32, gw Gateway) {
	gateways.Put(roleID, gw)
}
func DeleteGW(roleID int32) {
	gateways.Delete(roleID)
}
