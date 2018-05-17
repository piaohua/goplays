package main

import (
	"time"

	"goplays/glog"
	"goplays/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//玩家桌子常用共有操作请求处理
func (rs *RoleActor) HandlerDesk(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.CloseDesk:
		arg := msg.(*pb.CloseDesk)
		glog.Debugf("CloseDesk %#v", arg)
		//TODO
		//响应
		//rsp := new(pb.ClosedDesk)
		//ctx.Respond(rsp)
	case *pb.LeaveDesk:
		arg := msg.(*pb.LeaveDesk)
		glog.Debugf("LeaveDesk %#v", arg)
		//离开房间
		rs.leaveRoom(arg, ctx)
		//响应
		//rsp := new(pb.LeftDesk)
		//ctx.Respond(rsp)
	case *pb.CChatText:
		arg := msg.(*pb.CChatText)
		glog.Debugf("CChatText %#v", arg)
		if rs.gamePid == nil {
			rsp := new(pb.SChatText)
			rsp.Error = pb.NotInRoom
			rs.Send(rsp)
			return
		}
		rs.gamePid.Request(arg, ctx.Self())
	case *pb.CChatVoice:
		arg := msg.(*pb.CChatVoice)
		glog.Debugf("CChatVoice %#v", arg)
		if rs.gamePid == nil {
			rsp := new(pb.SChatVoice)
			rsp.Error = pb.NotInRoom
			rs.Send(rsp)
			return
		}
		rs.gamePid.Request(arg, ctx.Self())
	case *pb.CHuiYinRoomList:
		arg := msg.(*pb.CHuiYinRoomList)
		glog.Debugf("CHuiYinRoomList %#v", arg)
		rs.roomPid.Request(arg, ctx.Self())
	case *pb.CPk10Record:
		arg := msg.(*pb.CPk10Record)
		glog.Debugf("CPk10Record %#v", arg)
		rs.dbmsPid.Request(arg, ctx.Self())
	case *pb.CHuiYinRecords:
		arg := msg.(*pb.CHuiYinRecords)
		glog.Debugf("CHuiYinRecords %#v", arg)
		arg.Userid = rs.User.GetUserid()
		rs.dbmsPid.Request(arg, ctx.Self())
	case *pb.CHuiYinProfit:
		arg := msg.(*pb.CHuiYinProfit)
		glog.Debugf("CHuiYinProfit %#v", arg)
		arg.Userid = rs.User.GetUserid()
		rs.dbmsPid.Request(arg, ctx.Self())
	case *pb.CGetTrend:
		arg := msg.(*pb.CGetTrend)
		glog.Debugf("CGetTrend %#v", arg)
		//arg.Userid = rs.User.GetUserid()
		if rs.gamePid == nil {
			rsp := new(pb.SGetTrend)
			rsp.Error = pb.NotInRoom
			rs.Send(rsp)
			return
		}
		rs.gamePid.Request(arg, ctx.Self())
	case *pb.CGetOpenResult:
		arg := msg.(*pb.CGetOpenResult)
		glog.Debugf("CGetOpenResult %#v", arg)
		//arg.Userid = rs.User.GetUserid()
		if rs.gamePid == nil {
			rsp := new(pb.SGetOpenResult)
			rsp.Error = pb.NotInRoom
			rs.Send(rsp)
			return
		}
		rs.gamePid.Request(arg, ctx.Self())
	case *pb.CGetLastWins:
		arg := msg.(*pb.CGetLastWins)
		glog.Debugf("CGetLastWins %#v", arg)
		//arg.Userid = rs.User.GetUserid()
		if rs.gamePid == nil {
			rsp := new(pb.SGetLastWins)
			rsp.Error = pb.NotInRoom
			rs.Send(rsp)
			return
		}
		rs.gamePid.Request(arg, ctx.Self())
	case *pb.SetRecord:
		arg := msg.(*pb.SetRecord)
		glog.Debugf("SetRecord %#v", arg)
		rs.status = true
		rs.User.SetRecord(arg.Rtype)
	default:
		//glog.Errorf("unknown message %v", msg)
		rs.HandlerHuiYin(msg, ctx)
	}
}

//离开房间
func (rs *RoleActor) leaveRoom(arg *pb.LeaveDesk, ctx actor.Context) {
	glog.Debugf("leaveRoom: %#v", arg)
	rs.gamePid = nil
	rs.roomPid.Request(arg, ctx.Self())
	rs.hallPid.Request(arg, ctx.Self())
}

//进入房间
func (rs *RoleActor) entryRoom(arg *pb.CHuiYinEnterRoom, ctx actor.Context) bool {
	if rs.entryRoom2(arg, ctx) {
		glog.Debug("enter desk successfully")
		return true
	}
	//进入失败
	glog.Debug("enter desk failed")
	if rs.gamePid != nil {
		//离开消息
		rs.leaveDesk()
		//
		rs.gamePid = nil
	}
	return false
}

//进入房间
func (rs *RoleActor) entryRoom2(arg *pb.CHuiYinEnterRoom, ctx actor.Context) bool {
	if rs.gamePid == nil {
		glog.Errorf("not in the room: %s", rs.User.GetUserid())
		return false
	}
	result4, err4 := json.Marshal(rs.User)
	if err4 != nil {
		glog.Errorf("user Marshal err %v", err4)
		return false
	}
	msg4 := new(pb.EnterDesk)
	//不能用future pid做返回,要用真实pid
	//future和真实Pid不同,如下
	//"127.0.0.1:8004/$G"
	//127.0.0.1:8004/future$G
	msg4.Sender = ctx.Self()
	msg4.Data = result4
	//进入房间
	timeout := 3 * time.Second
	res1, err1 := rs.gamePid.RequestFuture(msg4, timeout).Result()
	if err1 != nil {
		glog.Errorf("entry Room err: %v", err1)
		return false
	}
	if response1, ok := res1.(*pb.EnteredDesk); ok {
		glog.Debugf("response1: %#v", response1)
		if response1.Error != pb.OK {
			//TODO 如果已经下注重连?
			glog.Debugf("enter desk failed: %d", response1.Error)
			return false
		}
	} else {
		glog.Debugf("enter desk failed: %v", res1)
		return false
	}
	//进入房间数据
	rs.gamePid.Request(arg, ctx.Self())
	return true
}

//大厅中匹配可用房间,TODO 暂时通过roomid匹配
func (rs *RoleActor) matchRoom(roomid string, gtype, rtype uint32) *pb.MatchedDesk {
	//匹配可以进入的房间
	msg1 := new(pb.MatchDesk)
	msg1.Roomid = roomid
	msg1.Gtype = gtype
	msg1.Rtype = rtype
	//节点注册名称,TODO 多节点处理
	//msg1.Name = cfg.Section("game.huiyin").Name()
	timeout := 3 * time.Second
	res1, err1 := rs.hallPid.RequestFuture(msg1, timeout).Result()
	if err1 != nil {
		glog.Errorf("matchRoom err: %v", err1)
		return nil
	}
	if response1, ok := res1.(*pb.MatchedDesk); ok {
		glog.Debugf("response1: %#v", response1)
		if response1.Desk == nil {
			glog.Errorf("matchRoom failed: %d", roomid)
			return nil
		}
		return response1
	}
	return nil
}
