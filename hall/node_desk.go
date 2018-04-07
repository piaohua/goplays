package main

import (
	"goplays/data"
	"goplays/game/config"
	"goplays/game/handler"
	"goplays/glog"
	"goplays/pb"
	"strings"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//玩家桌子请求处理
func (a *HallActor) HandlerDesk(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.MatchDesk:
		arg := msg.(*pb.MatchDesk)
		glog.Debugf("MatchDesk: %v", arg)
		//按规则匹配房间
		a.matchDesk(arg, ctx)
	case *pb.JoinDesk:
		arg := msg.(*pb.JoinDesk)
		glog.Debugf("JoinDesk %#v", arg)
		//房间数据变更
		if _, ok := a.router[arg.Userid]; !ok {
			a.rnums[arg.Roomid] += 1
		}
		a.router[arg.Userid] = arg.Roomid
		//响应
		//rsp := new(pb.JoinedDesk)
		//ctx.Respond(rsp)
	case *pb.AddDesk:
		arg := msg.(*pb.AddDesk)
		glog.Debugf("AddDesk %#v", arg)
		//房间数据变更
		a.desks[arg.Roomid] = arg.Desk
		a.rtype[arg.Roomid] = arg.Rtype
		a.gtype[arg.Roomid] = arg.Gtype
		//TODO 添加桌子匹配规则
		//响应
		//rsp := new(pb.AddedDesk)
		//ctx.Respond(rsp)
		glog.Debugf("AddDesk %d", len(a.rtype))
		glog.Debugf("AddDesk %d", len(a.gtype))
		glog.Debugf("AddDesk %d", len(a.rnums))
		glog.Debugf("AddDesk %d", len(a.desks))
	case *pb.LeaveDesk:
		arg := msg.(*pb.LeaveDesk)
		glog.Debugf("LeaveDesk %#v", arg)
		//移除
		if _, ok := a.router[arg.Userid]; ok {
			delete(a.router, arg.Userid)
			if n, ok := a.rnums[arg.Roomid]; ok && n > 0 {
				a.rnums[arg.Roomid] = n - 1
			}
		}
		glog.Debugf("LeaveDesk %d", len(a.rtype))
		glog.Debugf("LeaveDesk %d", len(a.gtype))
		glog.Debugf("LeaveDesk %d", len(a.rnums))
		glog.Debugf("LeaveDesk %d", len(a.desks))
		//响应
		//rsp := new(pb.LeftDesk)
		//ctx.Respond(rsp)
	case *pb.CloseDesk:
		arg := msg.(*pb.CloseDesk)
		glog.Debugf("CloseDesk %#v", arg)
		//移除
		delete(a.rtype, arg.Roomid)
		delete(a.gtype, arg.Roomid)
		delete(a.rnums, arg.Roomid)
		delete(a.desks, arg.Roomid)
		glog.Debugf("CloseDesk %d", len(a.rtype))
		glog.Debugf("CloseDesk %d", len(a.gtype))
		glog.Debugf("CloseDesk %d", len(a.rnums))
		glog.Debugf("CloseDesk %d", len(a.desks))
		//响应
		//rsp := new(pb.ClosedDesk)
		//ctx.Respond(rsp)
	case *pb.SyncConfig:
		//启动后同步配置
		arg := msg.(*pb.SyncConfig)
		glog.Debugf("SyncConfig %#v", arg)
		handler.SyncConfig(arg)
	case *pb.WebRequest:
		arg := msg.(*pb.WebRequest)
		glog.Debugf("WebRequest %#v", arg)
		rsp := new(pb.WebResponse)
		rsp.Code = arg.Code
		a.HandlerWeb(arg, rsp, ctx)
		ctx.Respond(rsp)
	case *pb.HuiYinOpenTime:
		//TODO 优化
		arg := msg.(*pb.HuiYinOpenTime)
		glog.Debugf("HuiYinOpenTime %#v", arg)
		for k, v := range a.serve {
			if strings.Contains(k, arg.Name) {
				v.Tell(arg)
				return
			}
		}
		rsp := new(pb.SHuiYinGames)
		rsp.Error = pb.Failed
		arg.Sender.Tell(rsp)
	case *pb.RobotFake:
		arg := msg.(*pb.RobotFake)
		glog.Debugf("RobotFake %#v", arg)
		var num int32 = config.GetEnv(data.ENV12)
		if num == 0 { //关闭状态
			return
		}
		roomName := cfg.Section("room").Name()
		if v, ok := a.serve[roomName]; ok {
			arg.RoomPid = v
			for _, v2 := range a.desks {
				v2.Tell(arg)
			}
		}
	case *pb.RobotMsg:
		arg := msg.(*pb.RobotMsg)
		glog.Debugf("RobotMsg %#v", arg)
		robotName := cfg.Section("robot").Name()
		if v, ok := a.serve[robotName]; ok {
			v.Tell(arg)
		}
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

//匹配房间
func (a *HallActor) matchDesk(arg *pb.MatchDesk, ctx actor.Context) {
	rsp := new(pb.MatchedDesk)
	//存在id直接查找
	var roomid string
	//TODO 添加人数匹配
	if arg.Roomid != "" {
		roomid = arg.Roomid
	} else if arg.Gtype != 0 && arg.Rtype != 0 {
		for k, v := range a.rtype {
			if v == arg.Rtype && a.gtype[k] == arg.Gtype {
				roomid = k
				break
			}
		}
	} else if arg.Gtype != 0 {
		for k, v := range a.gtype {
			if v == arg.Gtype {
				roomid = k
				break
			}
		}
	} else if arg.Rtype != 0 {
		for k, v := range a.rtype {
			if v == arg.Rtype {
				roomid = k
				break
			}
		}
	}
	glog.Debugf("MatchDesk roomid : %s", roomid)
	if v, ok := a.desks[roomid]; ok {
		rsp.Desk = v
	}
	//响应
	ctx.Respond(rsp)
}
