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
	"github.com/gbember/gtserver/config"
	"github.com/gbember/gtserver/db"
	"github.com/gbember/gtserver/gateway"
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
	err = logger.StartLog(config.DataSetting.LogDir, config.DataSetting.LogLevel)
	if err != nil {
		panic(err)
	}

	defer util.LogPanicStack()

	logger.Info("开始加载配置....")
	err = config.LoadAllConfig()
	if err != nil {
		logger.Critical("%s", err.Error())
		return
	}
	logger.Info("===>>配置加载成功")

	console.RegisterModule(config.DataSetting.ConsoleAddr, 10, 1024)
	gateway.RegisterModule()
	module.Init()

	err = db.Start()
	if err != nil {
		logger.Critical("===>数据库启动错误:%v", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c

	module.Destroy()
	db.Stop()
	logger.Info("服务器退出:%v", sig)
}
