// role.go
package role

import (
	"time"

	"github.com/gbember/gtserver/proto"

	"github.com/gbember/gtserver/common/role"
)

type _fm_role struct{}

func init() {
	role.RegisterFM(new(_fm_role))
}

func (*_fm_role) LoadData(rd *role.RoleData) error {
	return nil
}
func (*_fm_role) HookOnline(rd *role.RoleData) {
	bi := rd.BaseInfo
	rd.GW.Send(&proto.Sc_role_info{
		RoleID:        bi.RoleID,
		RoleName:      bi.RoleName,
		Sex:           bi.Sex,
		Level:         bi.Level,
		Exp:           bi.Exp,
		VipLevel:      bi.VipLevel,
		Gold:          bi.Gold,
		Coin:          bi.Coin,
		LastLogoutSec: bi.LastLogoutSec,
	})
}
func (*_fm_role) HookOffline(rd *role.RoleData) {
	rd.LastLogoutSec = int32(time.Now().Unix())
}
