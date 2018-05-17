package main

import (
	"time"

	"goplays/data"
	"goplays/glog"
	"goplays/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
)

//玩家请求处理
func (ws *WSConn) HandlerHuiYin(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.HuiYinOpenedTime:
		arg := msg.(*pb.HuiYinOpenedTime)
		glog.Debugf("HuiYinOpenedTime %#v", arg)
		rsp := new(pb.SHuiYinGames)
		rsp.List = arg.List
		rsp.Error = arg.Error
		ws.Send(rsp)
	case *pb.CHuiYinGames:
		arg := msg.(*pb.CHuiYinGames)
		glog.Debugf("CHuiYinGames %#v", arg)
		msg1 := new(pb.HuiYinOpenTime)
		msg1.Sender = ctx.Self()
		//TODO 优化映射关系
		msg1.Name = cfg.Section("game.huiyin1").Name()
		ws.hallPid.Tell(msg1)
	case *pb.CHuiYinEnterRoom:
		arg := msg.(*pb.CHuiYinEnterRoom)
		glog.Debugf("CHuiYinEnterRoom %#v", arg)
		ws.huiyinEnter(arg, ctx)
	case *pb.CHuiYinDealer:
		arg := msg.(*pb.CHuiYinDealer)
		glog.Debugf("CHuiYinDealer %#v", arg)
		if ws.gamePid == nil {
			rsp := new(pb.SHuiYinDealer)
			rsp.Error = pb.NotInRoom
			ws.Send(rsp)
			return
		}
		var num int64 = int64(arg.GetNum())
		if ws.User.GetChip() < num {
			rsp := new(pb.SHuiYinDealer)
			rsp.Error = pb.NotEnoughCoin
			ws.Send(rsp)
			return
		}
		ws.gamePid.Request(arg, ctx.Self())
	case *pb.CHuiYinDealerList:
		arg := msg.(*pb.CHuiYinDealerList)
		glog.Debugf("CHuiYinDealerList %#v", arg)
		if ws.gamePid == nil {
			rsp := new(pb.SHuiYinDealerList)
			rsp.Error = pb.NotInRoom
			ws.Send(rsp)
			return
		}
		ws.gamePid.Request(arg, ctx.Self())
	case *pb.CHuiYinRoomRoles:
		arg := msg.(*pb.CHuiYinRoomRoles)
		glog.Debugf("CHuiYinRoomRoles %#v", arg)
		if ws.gamePid == nil {
			rsp := new(pb.SHuiYinRoomRoles)
			rsp.Error = pb.NotInRoom
			ws.Send(rsp)
			return
		}
		ws.gamePid.Request(arg, ctx.Self())
	case *pb.CHuiYinRoomBet:
		arg := msg.(*pb.CHuiYinRoomBet)
		glog.Debugf("CHuiYinRoomBet %#v", arg)
		ws.huiyinBet(arg, ctx)
	case *pb.CHuiYinLeave:
		arg := msg.(*pb.CHuiYinLeave)
		glog.Debugf("CHuiYinLeave %#v", arg)
		if ws.gamePid == nil {
			rsp := new(pb.SHuiYinLeave)
			rsp.Error = pb.NotInRoom
			ws.Send(rsp)
			return
		}
		arg.Userid = ws.User.GetUserid()
		ws.gamePid.Request(arg, ctx.Self())
	case *pb.CHuiYinSit:
		arg := msg.(*pb.CHuiYinSit)
		glog.Debugf("CHuiYinSit %#v", arg)
		if ws.gamePid == nil {
			rsp := new(pb.SHuiYinSit)
			rsp.Error = pb.NotInRoom
			rsp.State = arg.State
			ws.Send(rsp)
			return
		}
		if ws.User.IsTourist() {
			rsp := new(pb.SHuiYinSit)
			rsp.Error = pb.TouristInoperable
			ws.Send(rsp)
			return
		}
		ws.gamePid.Request(arg, ctx.Self())
	case proto.Message:
		//响应消息
		ws.Send(msg)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

//下注
func (ws *WSConn) huiyinBet(arg *pb.CHuiYinRoomBet, ctx actor.Context) {
	if ws.User.IsTourist() {
		rsp := new(pb.SHuiYinRoomBet)
		rsp.Error = pb.TouristInoperable
		ws.Send(rsp)
		return
	}
	if ws.gamePid == nil {
		rsp := new(pb.SHuiYinRoomBet)
		rsp.Error = pb.NotInRoom
		ws.Send(rsp)
		return
	}
	value := arg.GetValue()
	seatBet := arg.GetSeatbet()
	if !(seatBet >= data.SEAT1 && seatBet <= data.SEAT5) {
		rsp := new(pb.SHuiYinRoomBet)
		rsp.Error = pb.OperateError
		ws.Send(rsp)
		return
	}
	if value <= 0 {
		rsp := new(pb.SHuiYinRoomBet)
		rsp.Error = pb.OperateError
		ws.Send(rsp)
		return
	}
	if ws.User.GetChip() < int64(value) {
		rsp := new(pb.SHuiYinRoomBet)
		rsp.Error = pb.NotEnoughCoin
		ws.Send(rsp)
		return
	}
	ws.gamePid.Request(arg, ctx.Self())
}

//进入房间
func (ws *WSConn) huiyinEnter(arg *pb.CHuiYinEnterRoom, ctx actor.Context) {
	if !ws.huiyinLeave(ctx) {
		return
	}
	if ws.gamePid == nil {
		//匹配可以进入的房间
		response1 := ws.matchRoom(arg.Roomid, arg.Gtype, arg.Rtype)
		if response1 != nil {
			ws.gamePid = response1.Desk
		}
	}
	glog.Debugf("gamePid: %s", ws.gamePid.String())
	if ws.gamePid == nil {
		stoc := new(pb.SHuiYinEnterRoom)
		stoc.Error = pb.RoomNotExist
		ws.Send(stoc)
		return
	}
	//进入房间
	if !ws.entryRoom(arg, ctx) {
		stoc := new(pb.SHuiYinEnterRoom)
		stoc.Error = pb.RoomNotExist
		ws.Send(stoc)
		return
	}
}

//离开房间
func (ws *WSConn) huiyinLeave(ctx actor.Context) bool {
	if ws.gamePid == nil {
		return true
	}
	msg1 := new(pb.LeaveDesk)
	msg1.Userid = ws.User.GetUserid()
	//离开房间
	timeout := 3 * time.Second
	res1, err1 := ws.gamePid.RequestFuture(msg1, timeout).Result()
	if err1 != nil {
		glog.Errorf("entry Room err: %v", err1)
		stoc := new(pb.SHuiYinEnterRoom)
		stoc.Error = pb.RoomNotExist
		ws.Send(stoc)
		return false
	}
	if response1, ok := res1.(*pb.LeftDesk); ok {
		glog.Debugf("response1: %#v", response1)
		if response1.Error != pb.OK {
			//离开失败
			glog.Debugf("leave desk failed: %d", response1.Error)
			stoc := new(pb.SHuiYinEnterRoom)
			stoc.Error = response1.Error
			ws.Send(stoc)
			return false
		}
	} else {
		glog.Debugf("enter desk failed: %v", res1)
		stoc := new(pb.SHuiYinEnterRoom)
		stoc.Error = pb.RoomNotExist
		ws.Send(stoc)
		return false
	}
	ws.gamePid = nil
	return true
}
