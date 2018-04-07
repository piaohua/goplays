package main

import (
	"os"
	"os/signal"
	"runtime"
	"time"

	"goplays/game/config"
	"goplays/glog"

	//_ "net/http/pprof"

	jsoniter "github.com/json-iterator/go"
	ini "gopkg.in/ini.v1"
)

var (
	cfg *ini.File
	sec *ini.Section
	err error

	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	defer glog.Flush()
	//go func() {
	//	http.ListenAndServe("localhost:6060", nil)
	//}()
	//日志定义
	glog.Init()
	//加载配置
	cfg, err = ini.Load("conf.ini")
	if err != nil {
		panic(err)
	}
	cfg.BlockMode = false //只读
	bind := cfg.Section("hall").Key("bind").Value()
	name := cfg.Section("cookie").Key("name").Value()
	NewRemote(bind, name)
	//初始化
	config.Init2Game()
	signalListen()
	//关闭服务
	Stop()
	//延迟等待
	<-time.After(10 * time.Second) //延迟关闭
}

func signalListen() {
	c := make(chan os.Signal)
	//signal.Notify(c)
	signal.Notify(c, os.Interrupt, os.Kill) //监听SIGINT和SIGKILL信号
	//signal.Stop(c)
	for {
		s := <-c
		glog.Error("get signal:", s)
		return
	}
}
