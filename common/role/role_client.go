package role

import (
	"github.com/gbember/gt/util/concurrent"
	"github.com/gbember/gtserver/proto"
)

var (
	allRoles       = concurrent.NewConcurrentMap(1000)
	resgisterFuncs = make(map[uint16]handleFunc)
)

type handleFunc func(*RoleData, proto.Messager)

func GetRoleClient(roleID int32) *RoleData {
	rd := allRoles.Get(roleID)
	if rd != nil {
		return rd.(*RoleData)
	}
	return nil
}
func SetRoleClient(rd *RoleData) {
	allRoles.Put(rd.RoleID, rd)
}

//注册处理函数
func RegisterHandle(msgID uint16, f handleFunc) {
	resgisterFuncs[msgID] = f
}

//查找处理函数(没有返回nil)
func GetHandle(msgID uint16) handleFunc {
	if f, ok := resgisterFuncs[msgID]; ok {
		return f
	}
	return nil
}
