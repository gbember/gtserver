// role.go
package role

type fm interface {
	//数据加载
	LoadData(*RoleData) error
	//上线回调
	HookOnline(*RoleData)
	//下线回调
	HookOffline(*RoleData)
}

var fmList []fm = make([]fm, 0, 100)

func RegisterFM(_fm fm) {
	fmList = append(fmList, _fm)
}

func LoadData(rd *RoleData) error {
	var err error
	for i := 0; i < len(fmList); i++ {
		err = fmList[i].LoadData(rd)
		if err != nil {
			return err
		}
	}
	return nil
}

func HookOnline(rd *RoleData) {
	for i := 0; i < len(fmList); i++ {
		fmList[i].HookOnline(rd)
	}
}

func HookOffline(rd *RoleData) {
	for i := 0; i < len(fmList); i++ {
		fmList[i].HookOffline(rd)
	}
}
