// main.go
package main

import (
	"os"
	"os/signal"
	"runtime"

	"github.com/gbember/gt/console"
	"github.com/gbember/gt/logger"
	"github.com/gbember/gt/module"
	"github.com/gbember/gt/util"
	"github.com/gbember/gtserver/common/name"
	"github.com/gbember/gtserver/config"
	"github.com/gbember/gtserver/db"
	"github.com/gbember/gtserver/gateway"
	"github.com/gbember/gtserver/word"
)

var (
	wait = make(chan struct{})
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	err := config.LoadConfig("setting")
	if err != nil {
		panic(err)
	}
	l, err := logger.StartLog(config.DataSetting.LogDir, config.DataSetting.LogLevel, true)
	if err != nil {
		panic(err)
	}
	logger.Export(l)

	defer util.LogPanicStack()

	logger.Info("开始加载配置....")
	err = config.LoadAllConfig()
	if err != nil {
		logger.Critical("%s", err.Error())
		return
	}
	logger.Info("===>>配置加载成功")

	err = db.Start()
	if err != nil {
		logger.Critical("===>数据库启动错误:%v", err)
	}
	logger.Info("db start...")
	//加载所有名字
	name.InitAllRoleName()

	word.RegisterModule()
	gateway.RegisterModule()
	console.RegisterModule(config.DataSetting.ConsoleAddr, 10, 1024)
	//注册方法运行时间统计命令
	util.RegisterMTCommand()
	module.Init()
	logger.Info("服务器启动成功")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c

	module.Destroy()
	db.Stop()
	logger.Info("服务器退出:%v", sig)
}
