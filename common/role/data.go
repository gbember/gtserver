// data.go
package role

import (
	"time"

	"github.com/gbember/gt/timer"
	"github.com/gbember/gtserver/common/gateway"
	"github.com/gbember/gtserver/db"
	"github.com/gbember/gtserver/types"
)

type RoleData struct {
	RoleID        int32
	GW            gateway.Gateway
	Timer         *timer.Timer
	BaseInfo      *types.RoleInfo
	LastLogoutSec int32
}

func (rd *RoleData) PersiteData() error {
	db.UpdateRoleInfo(rd.BaseInfo)
	return nil
}

func (rd *RoleData) AddTimerFun(d time.Duration, fun func()) {
	rd.Timer.AddFun(d, fun)
}
