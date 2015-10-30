// role_all_name.go
package name

import (
	"github.com/gbember/gt/logger"
	"github.com/gbember/gt/util/concurrent"
	"github.com/gbember/gtserver/db"
)

var allNames *concurrent.ConcurrentMap

//加载所有名字
func InitAllRoleName() {
	allNames = concurrent.NewConcurrentMap(100000)
	logger.Info("===========>>开启缓存所有名字")
	allRoleNames, err := db.FindAllName()
	if err != nil {
		panic(err)
	}
	for _, RN := range allRoleNames {
		allNames.Put(RN.RoleName, RN.RoleID)
	}
	logger.Info("===========>>缓存所有名字成功")
}

//注册名字 返回是否注册成功
func RegisterRoleName(roleID int32, roleName string) bool {
	return allNames.PutNotExist(roleName, roleID)
}
