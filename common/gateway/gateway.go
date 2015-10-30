// gateway.go
package gateway

import (
	"github.com/gbember/gt/util/concurrent"
	"github.com/gbember/gtserver/proto"
)

type Gateway interface {
	Send(msg proto.Messager)
	RealSend(msg proto.Messager)
	Close(int8)
}

var (
	gateways *concurrent.ConcurrentMap = concurrent.NewConcurrentMap(10000)
)

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
