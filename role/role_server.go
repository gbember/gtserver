// role_server.go
package role

import (
	"sync"

	"github.com/gbember/gt/util/concurrent"
	"github.com/gbember/gtserver/common/gateway"
	"github.com/gbember/gtserver/proto"
	"github.com/gbember/gtserver/types"
)

var (
	allRoles       = concurrent.NewConcurrentMap(1000)
	resgisterFuncs = make(map[uint16]handleFunc)
)

type handleFunc func(*roleClient, proto.Messager)

func Start(r *types.RoleInfo, recv chan proto.PMessage, exitCnt chan struct{}, wgExitCnt *sync.WaitGroup, gw gateway.Gateway) {
	rc := new(roleClient)
	roleData := new(types.RoleData)
	roleData.Base = r
	rc.RoleData = roleData
	rc.recv = recv
	rc.exitCnt = exitCnt
	rc.gw = gw
	rc.wgExitCnt = wgExitCnt
	go rc.start()
}

func GetRoleClient(roleID int32) *roleClient {
	role := allRoles.Get(roleID)
	if role != nil {
		return role.(*roleClient)
	}
	return nil
}
func setRoleClient(rc *roleClient) {
	allRoles.Put(rc.RoleID, rc)
}

//注册处理函数
func register_handle(msgID uint16, f handleFunc) {
	resgisterFuncs[msgID] = f
}

//查找处理函数(没有返回nil)
func get_handle(msgID uint16) handleFunc {
	if f, ok := resgisterFuncs[msgID]; ok {
		return f
	}
	return nil
}
