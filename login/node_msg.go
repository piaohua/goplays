package main

import (
	"fmt"
	"time"

	"goplays/glog"
	"goplays/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func (a *LoginActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.ServeStop:
		//关闭服务
		a.handlerStop(ctx)
		//响应登录
		rsp := new(pb.ServeStoped)
		ctx.Respond(rsp)
	case *pb.ServeStart:
		a.start(ctx)
		//响应
		//rsp := new(pb.ServeStarted)
		//ctx.Respond(rsp)
	case *pb.Tick:
		a.ding(ctx)
	case *pb.WxpayCallback:
		arg := msg.(*pb.WxpayCallback)
		glog.Debugf("WxpayCallback: %v", arg)
		a.hallPid.Tell(arg)
	case *pb.SmscodeRegist:
		arg := msg.(*pb.SmscodeRegist)
		glog.Debugf("SmscodeRegist %v", arg)
		//验证限制
		if v, ok := a.smsphone[arg.Phone]; ok && v > 10 {
			glog.Debugf("SmscodeRegist Phone %s, v %d", arg.Phone, v)
			return
		}
		if v, ok := a.smstimes[arg.Ipaddr]; ok && v > 10 {
			glog.Debugf("SmscodeRegist Ipaddr %s, v %d", arg.Ipaddr, v)
			return
		}
		a.smsphone[arg.Phone] += 1  //次数
		a.smstimes[arg.Ipaddr] += 1 //3分钟
		a.rolePid.Tell(arg)
	case *pb.WebRequest:
		arg := msg.(*pb.WebRequest)
		glog.Debugf("WebRequest %#v", arg)
		timeout := 5 * time.Second
		res1, err1 := a.hallPid.RequestFuture(arg, timeout).Result()
		if err1 != nil {
			rsp := new(pb.WebResponse)
			rsp.ErrMsg = fmt.Sprintf("hall request err1 %v", err1)
			ctx.Respond(rsp)
			return
		}
		ctx.Respond(res1)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

//启动服务
func (a *LoginActor) start(ctx actor.Context) {
	glog.Infof("desk start: %v", ctx.Self().String())
	//dbms
	bind := cfg.Section("dbms").Key("bind").Value()
	name := cfg.Section("cookie").Key("name").Value()
	role := cfg.Section("cookie").Key("role").Value()
	a.dbmsPid = actor.NewPID(bind, name)
	a.rolePid = actor.NewPID(bind, role)
	glog.Infof("a.dbmsPid: %s", a.dbmsPid.String())
	glog.Infof("a.rolePid: %s", a.rolePid.String())
	//hall
	bind = cfg.Section("hall").Key("bind").Value()
	a.hallPid = actor.NewPID(bind, name)
	glog.Infof("a.hallPid: %s", a.hallPid.String())
	connect := &pb.Connect{
		Name: a.Name,
	}
	a.hallPid.Request(connect, ctx.Self())
	//启动
	go a.ticker(ctx)
}

//时钟
func (a *LoginActor) ticker(ctx actor.Context) {
	tick := time.Tick(30 * time.Second)
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
func (a *LoginActor) ding(ctx actor.Context) {
	//glog.Debugf("ding: %v", ctx.Self().String())
	//TODO 同步节点
	//a.hallPid.Request(msg, ctx.Self())
	//TODO 路由方式
	//err := cfg.Section("gate.node"+version).NewKey("host", "x:x")
	switch a.timer {
	case 2: //1分钟
		a.smsExpire()
		a.timer = 0
	default:
		a.timer += 1
	}
}

//过期检测
func (a *LoginActor) smsExpire() {
	for k, v := range a.smsphone {
		if v >= 1 {
			a.smsphone[k] -= 1
			continue
		}
		delete(a.smsphone, k)
	}
	for k, v := range a.smstimes {
		if v >= 1 {
			a.smstimes[k] -= 1
			continue
		}
		delete(a.smstimes, k)
	}
}

//关闭时钟
func (a *LoginActor) closeTick() {
	select {
	case <-a.stopCh:
		return
	default:
		//停止发送消息
		close(a.stopCh)
	}
}

func (a *LoginActor) handlerStop(ctx actor.Context) {
	glog.Debugf("handlerStop: %s", a.Name)
	//关闭
	a.closeTick()
	//TODO 关闭消息
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
	<-time.After(2 * time.Second)
}
