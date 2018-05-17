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
func (rs *RoleActor) HandlerHuiYin(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.HuiYinOpenedTime:
		arg := msg.(*pb.HuiYinOpenedTime)
		glog.Debugf("HuiYinOpenedTime %#v", arg)
		rsp := new(pb.SHuiYinGames)
		rsp.List = arg.List
		rsp.Error = arg.Error
		rs.Send(rsp)
	case *pb.CHuiYinGames:
		arg := msg.(*pb.CHuiYinGames)
		glog.Debugf("CHuiYinGames %#v", arg)
		msg1 := new(pb.HuiYinOpenTime)
		msg1.Sender = ctx.Self()
		//TODO 优化映射关系
		msg1.Name = cfg.Section("game.huiyin1").Name()
		rs.hallPid.Tell(msg1)
	case *pb.CHuiYinEnterRoom:
		arg := msg.(*pb.CHuiYinEnterRoom)
		glog.Debugf("CHuiYinEnterRoom %#v", arg)
		rs.huiyinEnter(arg, ctx)
	case *pb.CHuiYinDealer:
		arg := msg.(*pb.CHuiYinDealer)
		glog.Debugf("CHuiYinDealer %#v", arg)
		if rs.gamePid == nil {
			rsp := new(pb.SHuiYinDealer)
			rsp.Error = pb.NotInRoom
			rs.Send(rsp)
			return
		}
		//var num int64 = int64(arg.GetNum())
		//TODO 暂时全部带上庄
		var num int64 = int64(rs.User.GetChip())
		if rs.User.GetChip() < num {
			rsp := new(pb.SHuiYinDealer)
			rsp.Error = pb.NotEnoughCoin
			rs.Send(rsp)
			return
		}
		rs.gamePid.Request(arg, ctx.Self())
	case *pb.CHuiYinDealerList:
		arg := msg.(*pb.CHuiYinDealerList)
		glog.Debugf("CHuiYinDealerList %#v", arg)
		if rs.gamePid == nil {
			rsp := new(pb.SHuiYinDealerList)
			rsp.Error = pb.NotInRoom
			rs.Send(rsp)
			return
		}
		rs.gamePid.Request(arg, ctx.Self())
	case *pb.CHuiYinRoomRoles:
		arg := msg.(*pb.CHuiYinRoomRoles)
		glog.Debugf("CHuiYinRoomRoles %#v", arg)
		if rs.gamePid == nil {
			rsp := new(pb.SHuiYinRoomRoles)
			rsp.Error = pb.NotInRoom
			rs.Send(rsp)
			return
		}
		rs.gamePid.Request(arg, ctx.Self())
	case *pb.CHuiYinDeskBetInfo:
		arg := msg.(*pb.CHuiYinDeskBetInfo)
		glog.Debugf("CHuiYinDeskBetInfo %#v", arg)
		if rs.gamePid == nil {
			rsp := new(pb.SHuiYinDeskBetInfo)
			rsp.Error = pb.NotInRoom
			rs.Send(rsp)
			return
		}
		rs.gamePid.Request(arg, ctx.Self())
	case *pb.CHuiYinRoomBet:
		arg := msg.(*pb.CHuiYinRoomBet)
		glog.Debugf("CHuiYinRoomBet %#v", arg)
		//下注金额转换为分
		arg.Value = arg.Value * 100
		rs.huiyinBet(arg, ctx)
	case *pb.CHuiYinLeave:
		arg := msg.(*pb.CHuiYinLeave)
		glog.Debugf("CHuiYinLeave %#v", arg)
		if rs.gamePid == nil {
			rsp := new(pb.SHuiYinLeave)
			rsp.Error = pb.NotInRoom
			rs.Send(rsp)
			return
		}
		arg.Userid = rs.User.GetUserid()
		rs.gamePid.Request(arg, ctx.Self())
	case *pb.CHuiYinSit:
		arg := msg.(*pb.CHuiYinSit)
		glog.Debugf("CHuiYinSit %#v", arg)
		if rs.gamePid == nil {
			rsp := new(pb.SHuiYinSit)
			rsp.Error = pb.NotInRoom
			rsp.State = arg.State
			rs.Send(rsp)
			return
		}
		if rs.User.IsTourist() {
			rsp := new(pb.SHuiYinSit)
			rsp.Error = pb.TouristInoperable
			rs.Send(rsp)
			return
		}
		rs.gamePid.Request(arg, ctx.Self())
	case proto.Message:
		//响应消息
		rs.Send(msg)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

//下注
func (rs *RoleActor) huiyinBet(arg *pb.CHuiYinRoomBet, ctx actor.Context) {
	if rs.User.IsTourist() {
		rsp := new(pb.SHuiYinRoomBet)
		rsp.Error = pb.TouristInoperable
		rs.Send(rsp)
		return
	}
	if rs.gamePid == nil {
		rsp := new(pb.SHuiYinRoomBet)
		rsp.Error = pb.NotInRoom
		rs.Send(rsp)
		return
	}
	value := arg.GetValue()
	seatBet := arg.GetSeatbet()
	if !(seatBet >= data.SEAT1 && seatBet <= data.SEAT5) {
		rsp := new(pb.SHuiYinRoomBet)
		rsp.Error = pb.OperateError
		rs.Send(rsp)
		return
	}
	if value <= 0 {
		rsp := new(pb.SHuiYinRoomBet)
		rsp.Error = pb.OperateError
		rs.Send(rsp)
		return
	}
	if rs.User.GetChip() < int64(value) {
		rsp := new(pb.SHuiYinRoomBet)
		rsp.Error = pb.NotEnoughCoin
		rs.Send(rsp)
		return
	}
	rs.gamePid.Request(arg, ctx.Self())
}

//进入房间
func (rs *RoleActor) huiyinEnter(arg *pb.CHuiYinEnterRoom, ctx actor.Context) {
	if !rs.huiyinLeave(arg.Roomid, ctx) {
		glog.Errorf("huiyinEnter failed: %v", arg)
		return
	}
	if rs.gamePid == nil {
		//匹配可以进入的房间
		response1 := rs.matchRoom(arg.Roomid, arg.Gtype, arg.Rtype)
		if response1 != nil {
			rs.gamePid = response1.Desk
		}
	}
	if rs.gamePid == nil {
		glog.Errorf("huiyinEnter failed: %v", arg)
		stoc := new(pb.SHuiYinEnterRoom)
		stoc.Error = pb.RoomNotExist
		rs.Send(stoc)
		return
	}
	glog.Debugf("gamePid: %s", rs.gamePid.String())
	//进入房间
	if !rs.entryRoom(arg, ctx) {
		glog.Errorf("huiyinEnter failed: %v", arg)
		stoc := new(pb.SHuiYinEnterRoom)
		stoc.Error = pb.RoomNotExist
		rs.Send(stoc)
		return
	}
}

//离开房间
func (rs *RoleActor) huiyinLeave(roomid string, ctx actor.Context) bool {
	if rs.gamePid == nil {
		return true
	}
	msg1 := new(pb.LeaveDesk)
	msg1.Userid = rs.User.GetUserid()
	msg1.Roomid = roomid
	//离开房间
	timeout := 3 * time.Second
	res1, err1 := rs.gamePid.RequestFuture(msg1, timeout).Result()
	if err1 != nil {
		glog.Errorf("entry Room err: %v", err1)
		stoc := new(pb.SHuiYinEnterRoom)
		stoc.Error = pb.RoomNotExist
		rs.Send(stoc)
		return false
	}
	if response1, ok := res1.(*pb.LeftDesk); ok {
		glog.Debugf("response1: %#v", response1)
		if response1.Error != pb.OK {
			//离开失败
			glog.Errorf("leave desk failed: %d", response1.Error)
			stoc := new(pb.SHuiYinEnterRoom)
			stoc.Error = response1.Error
			rs.Send(stoc)
			return false
		}
	} else {
		glog.Errorf("enter desk failed: %v", res1)
		stoc := new(pb.SHuiYinEnterRoom)
		stoc.Error = pb.RoomNotExist
		rs.Send(stoc)
		return false
	}
	rs.gamePid = nil
	return true
}
