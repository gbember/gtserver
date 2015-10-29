// db_role_info
package db

import (
	"time"

	"github.com/gbember/gt/logger"
	"github.com/gbember/gt/util/concurrent"
	"github.com/gbember/gtserver/types"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

const (
	_DB_ROLE_INFO = "db_role_info"
)

type db_table_role_info struct {
	session    *mgo.Session
	collection *mgo.Collection
	cmap       *concurrent.ConcurrentMap
}

var (
	db_role_info = &db_table_role_info{}
)

func init() {
	register(db_role_info)
}

func (this *db_table_role_info) init() error {
	session, collection := c(_DB_ROLE_INFO)
	this.collection = collection
	this.session = session
	this.cmap = concurrent.NewConcurrentMap(1000)
	err := this.ensureIndexs("roleid", "roleid_", true, true)
	if err != nil {
		return err
	}
	err = this.ensureIndexs("rolename", "rolename_", true, true)
	if err != nil {
		return err
	}
	err = this.ensureIndexs("level", "level_", false, false)
	if err != nil {
		return err
	}
	go this.loop_dump()
	return nil
}

func (this *db_table_role_info) destroy() {
	for this.cmap.Lenght() != 0 {
		time.Sleep(time.Second)
	}
	time.Sleep(time.Second * 2)
	this.session.Close()
}

func (this *db_table_role_info) loop_dump() {
	for {
		key, value := this.cmap.DeleteAnyOneWait()
		_, err := this.collection.Upsert(bson.M{"roleid": key}, value)
		if err != nil {
			logger.Error("持久化错误:%s", err.Error())
			this.cmap.PutAndSignal(key, value)
		}
	}
}

func (this *db_table_role_info) ensureIndexs(key string, name string, unique bool, dropDups bool) error {
	return this.collection.EnsureIndex(mgo.Index{Key: []string{key}, Name: name, Unique: unique, DropDups: dropDups})
}

func (this *db_table_role_info) insert(role *types.RoleInfo) {
	this.cmap.PutAndSignal(role.RoleID, role)
}

func (this *db_table_role_info) update(role *types.RoleInfo) {
	this.cmap.PutAndSignal(role.RoleID, role)
}

//根据roleID先找缓存没有找到再找数据库
func (this *db_table_role_info) findByRoleID(roleID int32) (*types.RoleInfo, error) {
	if role := this.cmap.Get(roleID); role != nil {
		return role.(*types.RoleInfo), nil
	}
	role := &types.RoleInfo{}
	err := this.collection.Find(bson.M{"roleid": roleID}).One(role)
	return role, err
}

//根据名字直接找数据库
func (this *db_table_role_info) findByName(roleName string) (*types.RoleInfo, error) {
	role := &types.RoleInfo{}
	err := this.collection.Find(bson.M{"rolename": roleName}).One(role)
	return role, err
}

func (this *db_table_role_info) findAllName() ([]*types.RoleName, error) {
	v := make([]*types.RoleName, 0, 100000)
	err := this.collection.Find(nil).Batch(1000).Select(map[string]int{"roleid": 1, "rolename": 1}).All(&v)
	return v, err
}

//根据等级范围返回一定数量
func (this *db_table_role_info) findByLevel(minLevel, maxLevel uint8, limitNum int) ([]*types.RoleInfo, error) {
	var result []*types.RoleInfo
	err := this.collection.Find(bson.M{"level": bson.M{"$gte": minLevel, "$lte": maxLevel}}).Limit(limitNum).All(result)
	return result, err
}

func InsertRoleInfo(role *types.RoleInfo) {
	db_role_info.insert(role)
}
func UpdateRoleInfo(role *types.RoleInfo) {
	db_role_info.update(role)
}
func FindByRoleID(roleID int32) (*types.RoleInfo, error) {
	return db_role_info.findByRoleID(roleID)
}
func FindByName(roleName string) (*types.RoleInfo, error) {
	return db_role_info.findByName(roleName)
}
func FindByLevel(minLevel, maxLevel uint8, limitNum int) ([]*types.RoleInfo, error) {
	return db_role_info.findByLevel(minLevel, maxLevel, limitNum)
}
func FindAllName() ([]*types.RoleName, error) {
	return db_role_info.findAllName()
}
