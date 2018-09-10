package main

import (
	"time"

	"goplays/glog"
	"goplays/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func (a *DeskActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.ServeStop:
		//关闭服务
		a.handlerStop(ctx)
		//响应登录
		rsp := new(pb.ServeStoped)
		ctx.Respond(rsp)
	case *pb.ServeStart:
		a.start(ctx)
		//a.start2(ctx)
		//响应
		//rsp := new(pb.ServeStarted)
		//ctx.Respond(rsp)
	case *pb.Tick:
		a.ding(ctx)
	default:
		//glog.Errorf("unknown message %v", msg)
		a.HandlerMsg(msg, ctx)
	}
}

//test
func (a *DeskActor) start2(ctx actor.Context) {
	msg1 := &pb.CloseDesk{
		Roomid: "1727",
	}
	msg2 := &pb.CloseDesk{
		Roomid: "1728",
	}
	msg3 := &pb.CloseDesk{
		Roomid: "1729",
	}
	msg4 := &pb.CloseDesk{
		Roomid: "1730",
	}
	msg5 := &pb.CloseDesk{
		Roomid: "1731",
	}
	a.roomPid.Tell(msg1)
	a.roomPid.Tell(msg2)
	a.roomPid.Tell(msg3)
	a.roomPid.Tell(msg4)
	a.roomPid.Tell(msg5)
	a.hallPid.Tell(msg1)
	a.hallPid.Tell(msg2)
	a.hallPid.Tell(msg3)
	a.hallPid.Tell(msg4)
	a.hallPid.Tell(msg5)
}

//启动服务
func (a *DeskActor) start(ctx actor.Context) {
	glog.Infof("desk start: %v", ctx.Self().String())
	//dbms
	bind := cfg.Section("dbms").Key("bind").Value()
	name := cfg.Section("cookie").Key("name").Value()
	room := cfg.Section("cookie").Key("room").Value()
	a.dbmsPid = actor.NewPID(bind, name)
	a.roomPid = actor.NewPID(bind, room)
	glog.Infof("a.dbmsPid: %s", a.dbmsPid.String())
	glog.Infof("a.roomPid: %s", a.roomPid.String())
	//hall
	bind = cfg.Section("hall").Key("bind").Value()
	a.hallPid = actor.NewPID(bind, name)
	glog.Infof("a.hallPid: %s", a.hallPid.String())
	connect := &pb.Connect{
		Name: a.Name,
	}
	a.hallPid.Request(connect, ctx.Self())
	//主动同步配置
	msg2 := new(pb.GetConfig)
	a.dbmsPid.Request(msg2, ctx.Self())
	//启动
	//go a.ticker(ctx)
	go a.initTicker(ctx)
}

//时钟
func (a *DeskActor) ticker(ctx actor.Context) {
	tick := time.Tick(time.Second)
	msg := new(pb.Tick)
	for {
		select {
		case <-a.stopCh:
			glog.Info("desk ticker closed")
			return
		default: //防止阻塞
		}
		select {
		case <-a.stopCh:
			glog.Info("desk ticker closed")
			return
		case <-tick:
			ctx.Self().Tell(msg)
		}
	}
}

//钟声
func (a *DeskActor) ding(ctx actor.Context) {
	//glog.Debugf("ding: %v", ctx.Self().String())
	//TODO 优化超时造成的延迟,单独routine处理
	a.handing(ctx)
}

//关闭时钟
func (a *DeskActor) closeTick() {
	select {
	case <-a.stopCh:
		return
	default:
		//停止发送消息
		close(a.stopCh)
	}
}

func (a *DeskActor) handlerStop(ctx actor.Context) {
	glog.Debugf("handlerStop: %s", a.Name)
	//关闭
	a.closeTick()
	//关闭消息
	msg1 := new(pb.ServeStop)
	for k, v := range a.desks {
		//关闭房间消息
		msg2 := new(pb.CloseDesk)
		msg2.Roomid = k
		//TODO 如果节点直接挂掉?
		//TODO 添加类型,桌子中关闭
		a.roomPid.Request(msg2, ctx.Self())
		a.hallPid.Request(msg2, ctx.Self())
		//关闭房间服务
		glog.Infof("Stop desk: %s", k)
		v.Request(msg1, ctx.Self())
		//停掉服务
		//v.Stop()
	}
	//延迟
	<-time.After(5 * time.Second)
	//断开处理
	msg := &pb.Disconnect{
		Name: a.Name,
	}
	if a.dbmsPid != nil {
		a.dbmsPid.Tell(msg)
	}
	if a.hallPid != nil {
		a.hallPid.Tell(msg)
	}
	//延迟
	<-time.After(3 * time.Second)
	for _, v := range a.desks {
		//停掉服务
		v.Stop()
	}
}
