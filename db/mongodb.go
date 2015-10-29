package db

import (
	"strconv"
	"sync"
	"time"

	"github.com/gbember/gtserver/config"
	"labix.org/v2/mgo"
)

type db_tabler interface {
	init() error
	destroy()
}

const (
	DEFAULT_MGO_TIMEOUT = 300
)

var (
	session   *mgo.Session
	database  *mgo.Database
	tableList = make([]db_tabler, 0, 50)
	tmut      sync.Mutex
)

//连接数据库
func Start() error {
	addr := config.DataSetting.DBAddrIP + ":" + strconv.Itoa(config.DataSetting.DBAddrPort)
	sess, err := mgo.Dial(addr)
	if err != nil {
		return err
	}

	sess.SetBatch(1000)
	sess.SetSocketTimeout(DEFAULT_MGO_TIMEOUT * time.Second)
	sess.SetCursorTimeout(0)

	session = sess
	database = sess.DB(config.DataSetting.DBName)
	err = database.Login(config.DataSetting.DBUser, config.DataSetting.DBPassword)
	if err != nil {
		return err
	}
	sess.SetMode(mgo.Strong, true)

	tmut.Lock()
	for _, table := range tableList {
		err := table.init()
		if err != nil {
			return err
		}
	}
	tmut.Unlock()

	return nil
}

func Stop() {
	tmut.Lock()
	for _, table := range tableList {
		table.destroy()
	}
	tmut.Unlock()
	session.Close()
}

func register(table db_tabler) {
	tmut.Lock()
	tableList = append(tableList, table)
	tmut.Unlock()
}

func c(collection string) (*mgo.Session, *mgo.Collection) {
	ms := session.Copy()
	c := ms.DB(config.DataSetting.DBName).C(collection)
	return ms, c
}
