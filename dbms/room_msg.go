package main

import (
	"time"

	"goplays/game/handler"
	"goplays/glog"
	"goplays/pb"
	"utils"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func (a *RoomActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.Connected:
		arg := msg.(*pb.Connected)
		glog.Infof("Connected %s", arg.Name)
	case *pb.Disconnected:
		arg := msg.(*pb.Disconnected)
		glog.Infof("Disconnected %s", arg.Name)
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
	case *pb.GenDesk:
		arg := msg.(*pb.GenDesk)
		glog.Debugf("GenDesk: %v", arg)
		a.genDesk(arg, ctx)
	case *pb.AddDesk:
		arg := msg.(*pb.AddDesk)
		glog.Debugf("AddDesk: %v", arg)
		a.addDesk(arg, ctx)
	case *pb.JoinDesk:
		arg := msg.(*pb.JoinDesk)
		glog.Debugf("JoinDesk %#v", arg)
		//房间数据变更
		if _, ok := a.router[arg.Userid]; !ok {
			a.router[arg.Userid] = arg.Roomid
			a.count[arg.Roomid] += 1
		}
		//响应
		//rsp := new(pb.JoinedDesk)
		//ctx.Respond(rsp)
	case *pb.LeaveDesk:
		arg := msg.(*pb.LeaveDesk)
		glog.Debugf("LeaveDesk %#v", arg)
		//移除
		if _, ok := a.router[arg.Userid]; ok {
			delete(a.router, arg.Userid)
			if n, ok := a.count[arg.Roomid]; ok && n > 0 {
				a.count[arg.Roomid] = n - 1
			}
		}
		//响应
		//rsp := new(pb.LeftDesk)
		//ctx.Respond(rsp)
	case *pb.Logout:
		arg := msg.(*pb.Logout)
		glog.Debugf("Logout %#v", arg)
		//TODO 暂时不处理
	case *pb.CloseDesk:
		arg := msg.(*pb.CloseDesk)
		glog.Debugf("CloseDesk %#v", arg)
		//TODO 私人房间
		//移除
		glog.Debugf("CloseDesk router %#v", a.router)
		glog.Debugf("CloseDesk count %#v", a.count)
		glog.Debugf("CloseDesk rules %#v", a.rules)
		delete(a.count, arg.Roomid)
		delete(a.codes, arg.Code)
		if v, ok := a.rooms[arg.Roomid]; ok {
			delete(a.rules, v.Unique)
		}
		delete(a.rooms, arg.Roomid)
		glog.Debugf("CloseDesk %d", len(a.rooms))
		//响应
		//rsp := new(pb.ClosedDesk)
		//ctx.Respond(rsp)
	case *pb.CHuiYinRoomList:
		arg := msg.(*pb.CHuiYinRoomList)
		glog.Debugf("CHuiYinRoomList %#v", arg)
		rsp := new(pb.SHuiYinRoomList)
		for k, v := range a.rooms {
			if v.Ltype != arg.Ltype {
				continue
			}
			l := new(pb.HuiYinRoom)
			l.Info = handler.PackGameInfo(v)
			l.Num = a.count[k]
			rsp.List = append(rsp.List, l)
		}
		ctx.Respond(rsp)
	case *pb.RobotFake:
		arg := msg.(*pb.RobotFake)
		glog.Debugf("RobotFake %#v", arg)
		switch arg.Type {
		case 1:
			a.count[arg.Roomid] += arg.FakeNum
		case 2:
			a.count[arg.Roomid] = arg.RealNum + arg.FakeNum
		}
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

//启动服务
func (a *RoomActor) start(ctx actor.Context) {
	glog.Infof("room start: %v", ctx.Self().String())
	//初始化建立连接
	bind := cfg.Section("hall").Key("bind").Value()
	name := cfg.Section("cookie").Key("name").Value()
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
func (a *RoomActor) ticker(ctx actor.Context) {
	tick := time.Tick(30 * time.Second)
	msg := new(pb.Tick)
	for {
		select {
		case <-a.stopCh:
			glog.Info("room ticker closed")
			return
		default: //防止阻塞
		}
		select {
		case <-a.stopCh:
			glog.Info("room ticker closed")
			return
		case <-tick:
			ctx.Self().Tell(msg)
		}
	}
}

//钟声
func (a *RoomActor) ding(ctx actor.Context) {
	//glog.Debugf("ding: %v", ctx.Self().String())
	//TODO
}

//关闭时钟
func (a *RoomActor) closeTick() {
	select {
	case <-a.stopCh:
		return
	default:
		//停止发送消息
		close(a.stopCh)
	}
}

func (a *RoomActor) handlerStop(ctx actor.Context) {
	glog.Debugf("handlerStop: %s", a.Name)
	//关闭
	a.closeTick()
	//回存数据
	if a.uniqueid != nil {
		a.uniqueid.Save()
	}
	for k, _ := range a.rooms {
		glog.Debugf("Stop room: %s", k)
		//TODO
		//v.Save()
	}
}

//生成一个牌桌邀请码,全列表中唯一
func (a *RoomActor) GenCode() (s string) {
	s = utils.RandStr(6)
	//是否已经存在
	if _, ok := a.codes[s]; ok {
		return a.GenCode() //重复尝试,TODO:一定次数后放弃尝试
	}
	return
}

//生成一个牌桌邀请码,全列表中唯一
func (a *RoomActor) GenCodeFree() (s string) {
	s = utils.RandStr(7) //区别于私人房间
	//是否已经存在
	if _, ok := a.codes[s]; ok {
		return a.GenCode() //重复尝试,TODO:一定次数后放弃尝试
	}
	return
}

//生成房间ID
func (a *RoomActor) genDesk(arg *pb.GenDesk, ctx actor.Context) {
	glog.Debugf("genDesk Rtype: %d", arg.Rtype)
	rsp := new(pb.GenedDesk)
	rsp.Roomid = a.uniqueid.GenID()
	//TODO
	//百人
	//rsp.Code = a.GenCodeFree()
	//私人
	//rsp.Code = a.GenCode()
	//响应
	ctx.Respond(rsp)
}

//添加房间
func (a *RoomActor) addDesk(arg *pb.AddDesk, ctx actor.Context) {
	glog.Debugf("addDesk Rtype: %d, Roomid: %s", arg.Rtype, arg.Roomid)
	rsp := new(pb.AddedDesk)
	deskData := handler.Data2Desk(arg.Data)
	if deskData == nil {
		glog.Errorf("addDesk err Rtype: %d, Roomid: %s", arg.Rtype, arg.Roomid)
		rsp.Error = pb.Failed
		ctx.Respond(rsp)
		return
	}
	//已经存在
	if _, ok := a.rules[deskData.Unique]; ok {
		glog.Errorf("addDesk err Rtype: %d, Roomid: %s, Unique: %s",
			arg.Rtype, arg.Roomid, deskData.Unique)
		rsp.Error = pb.Failed
		ctx.Respond(rsp)
		return
	}
	//添加房间
	a.rooms[arg.Roomid] = deskData
	//TODO 私人房间
	//a.codes[arg.Code] = deskData.Rid
	a.rules[deskData.Unique] = deskData.Rid
	ctx.Respond(rsp)
}
