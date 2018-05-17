package main

import (
	"time"

	"goplays/game/handler"
	"goplays/glog"
	"goplays/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func (a *GateActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.Connected:
		//连接成功
		arg := msg.(*pb.Connected)
		glog.Infof("Connected %s", arg.Name)
	case *pb.Disconnected:
		//成功断开
		arg := msg.(*pb.Disconnected)
		glog.Infof("Disconnected %s", arg.Name)
	case *pb.LoginElse:
		//在别的节点登录
		arg := msg.(*pb.LoginElse)
		glog.Debugf("LoginElse %#v", arg)
		userid := arg.GetUserid()
		if p, ok := a.roles[userid]; ok {
			p.Tell(arg)
		}
		glog.Debugf("LoginElse userid: %s", userid)
		//移除
		delete(a.roles, userid)
	case *pb.SetLogin:
		arg := msg.(*pb.SetLogin)
		glog.Debugf("SetLogin %#v", arg)
		set := &pb.SetLogined{
			DbmsPid: a.dbmsPid,
			RoomPid: a.roomPid,
			RolePid: a.rolePid,
			HallPid: a.hallPid,
		}
		arg.Sender.Tell(set)
		glog.Infof("SetLogin %s", arg.Sender.String())
	case *pb.LoginGate:
		//登录成功
		arg := msg.(*pb.LoginGate)
		glog.Debugf("LoginGate %#v", arg)
		userid := arg.GetUserid()
		//断开当前节点旧连接
		if p, ok := a.roles[userid]; ok {
			msg1 := new(pb.LoginElse)
			msg1.Userid = userid
			p.Tell(msg1)
		} else {
			//断开其它节点连接
			//TODO 数据一致性,防止数据覆盖
			msg1 := new(pb.LoginHall)
			msg1.Userid = userid
			msg1.NodeName = a.Name
			a.hallPid.Tell(msg1)
		}
		//添加
		a.roles[userid] = arg.Sender
		glog.Debugf("LoginGate userid: %s", userid)
		//响应登录
		rsp := new(pb.LoginedGate)
		rsp.Message = ctx.Self().String()
		ctx.Respond(rsp)
	case *pb.Logout:
		//登出成功
		arg := msg.(*pb.Logout)
		glog.Debugf("Logout %#v", arg)
		userid := arg.GetUserid()
		glog.Debugf("Logout userid: %s", userid)
		//glog.Infof("hallPid %s", a.hallPid.String())
		//glog.Infof("rolePid %s", a.rolePid.String())
		//glog.Infof("roomPid %s", a.roomPid.String())
		//移除
		delete(a.roles, userid)
		a.hallPid.Tell(arg)
		a.rolePid.Tell(arg)
		a.roomPid.Tell(arg)
	case *pb.ServeStop:
		//关闭服务
		a.handlerStop(ctx)
		//响应登录
		rsp := new(pb.ServeStoped)
		ctx.Respond(rsp)
	case *pb.ServeStart:
		//初始化建立连接
		a.start(ctx)
		//响应
		//rsp := new(pb.ServeStarted)
		//ctx.Respond(rsp)
	case *pb.Tick:
		a.ding(ctx)
	case *pb.SyncConfig:
		//同步配置
		arg := msg.(*pb.SyncConfig)
		glog.Debugf("SyncConfig %#v", arg)
		handler.SyncConfig(arg)
	case *pb.PayCurrency:
		//后台或充值同步到game房间
		arg := msg.(*pb.PayCurrency)
		glog.Debugf("PayCurrency %#v", arg)
		userid := arg.Userid
		if v, ok := a.roles[userid]; ok {
			v.Tell(msg)
		} else {
			//离线
			a.rolePid.Tell(msg)
		}
	case *pb.ChangeCurrency:
		arg := msg.(*pb.ChangeCurrency)
		glog.Debugf("ChangeCurrency %#v", arg)
		//货币变更同步
		userid := arg.Userid
		if v, ok := a.roles[userid]; ok {
			v.Tell(msg)
		} else {
			//离线
			a.rolePid.Tell(msg)
		}
	case *pb.WxpayCallback:
		arg := msg.(*pb.WxpayCallback)
		glog.Debugf("WxpayCallback %#v", arg)
		if !handler.WxpayVerify(arg) {
			return
		}
		a.rolePid.Tell(arg)
	case *pb.WxpayGoods:
		arg := msg.(*pb.WxpayGoods)
		glog.Debugf("WxpayGoods: %v", arg)
		userid := arg.Userid
		if v, ok := a.roles[userid]; ok {
			v.Tell(arg)
		} else {
			glog.Errorf("WxpayGoods: %v", arg)
		}
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

//启动服务
func (a *GateActor) start(ctx actor.Context) {
	glog.Infof("gate start: %v", ctx.Self().String())
	//dbms
	bind := cfg.Section("dbms").Key("bind").Value()
	name := cfg.Section("cookie").Key("name").Value()
	room := cfg.Section("cookie").Key("room").Value()
	role := cfg.Section("cookie").Key("role").Value()
	a.dbmsPid = actor.NewPID(bind, name)
	a.roomPid = actor.NewPID(bind, room)
	a.rolePid = actor.NewPID(bind, role)
	//hall
	bind = cfg.Section("hall").Key("bind").Value()
	a.hallPid = actor.NewPID(bind, name)
	glog.Infof("a.hallPid: %s", a.hallPid.String())
	connect := &pb.Connect{
		Name: a.Name,
	}
	a.dbmsPid.Request(connect, ctx.Self())
	a.hallPid.Request(connect, ctx.Self())
	glog.Infof("a.dbmsPid: %s", a.dbmsPid.String())
	glog.Infof("a.roomPid: %s", a.roomPid.String())
	glog.Infof("a.rolePid: %s", a.rolePid.String())
	//启动
	go a.ticker(ctx)
}

//时钟
func (a *GateActor) ticker(ctx actor.Context) {
	tick := time.Tick(30 * time.Second)
	msg := new(pb.Tick)
	for {
		select {
		case <-a.stopCh:
			glog.Info("gate ticker closed")
			return
		default: //防止阻塞
		}
		select {
		case <-a.stopCh:
			glog.Info("gate ticker closed")
			return
		case <-tick:
			ctx.Self().Tell(msg)
		}
	}
}

//钟声
func (a *GateActor) ding(ctx actor.Context) {
	//glog.Debugf("ding: %v", ctx.Self().String())
	//TODO 下线离线玩家
}

//关闭时钟
func (a *GateActor) closeTick() {
	select {
	case <-a.stopCh:
		return
	default:
		//停止发送消息
		close(a.stopCh)
	}
}

func (a *GateActor) handlerStop(ctx actor.Context) {
	glog.Debugf("handlerStop: %s", a.Name)
	//关闭
	a.closeTick()
	//关闭消息
	msg1 := new(pb.ServeStop)
	for k, v := range a.offline {
		glog.Debugf("Stop offline role: %s", k)
		v.Tell(msg1)
	}
	//关闭消息
	msg2 := new(pb.ServeClose)
	for k, v := range a.roles {
		glog.Debugf("Stop role: %s", k)
		v.Tell(msg1)
		v.Tell(msg2)
	}
	//延迟
	<-time.After(3 * time.Second)
	//断开处理
	msg := &pb.Disconnect{
		Name: a.Name,
	}
	if a.dbmsPid != nil {
		a.dbmsPid.Request(msg, ctx.Self())
	}
	if a.hallPid != nil {
		a.hallPid.Request(msg, ctx.Self())
	}
	//延迟
	<-time.After(2 * time.Second)
}
