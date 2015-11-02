// role_server.go
package role

import (
	"sync"

	"github.com/gbember/gtserver/common/gateway"
	"github.com/gbember/gtserver/common/role"
	"github.com/gbember/gtserver/proto"
	"github.com/gbember/gtserver/types"
)

func Start(r *types.RoleInfo, recv chan proto.PMessage, exitCnt chan struct{}, wgExitCnt *sync.WaitGroup, gw gateway.Gateway) {
	rc := new(roleAgent)
	roleData := new(role.RoleData)
	roleData.BaseInfo = r
	roleData.LastLogoutSec = r.LastLogoutSec
	roleData.GW = gw
	r.LastLogoutSec = 0
	rc.rd = roleData
	rc.recv = recv
	rc.exitCnt = exitCnt
	rc.gw = gw
	rc.wgExitCnt = wgExitCnt
	go rc.start()
}
