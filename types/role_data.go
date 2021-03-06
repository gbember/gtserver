// role_data.go
package types

type RoleData struct {
	Base *RoleInfo
}

type RoleInfo struct {
	RoleID        int32
	RoleName      string
	Sex           int8  //0 女 1 男
	Level         uint8 //等级
	Exp           int32 //经验值
	VipLevel      int8  //vip 等级
	Gold          int32 //元宝
	Coin          int32 //铜币
	LastLogoutSec int32 //上一次下线时间 为0表示在线
}

type RoleName struct {
	RoleID   int32
	RoleName string
}
