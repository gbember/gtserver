// data.go
package role

import (
	"github.com/gbember/gtserver/db"
	"github.com/gbember/gtserver/types"
)

type RoleData struct {
	RoleID   int32
	BaseInfo *types.RoleInfo
}

func (rd *RoleData) PersiteData() error {
	db.UpdateRoleInfo(rd.BaseInfo)
	return nil
}
