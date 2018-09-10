package main

import (
	"goplays/data"
	"goplays/game/algo"
	"goplays/game/handler"
	"goplays/glog"
	"goplays/pb"
	"utils"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func (a *Desk) HandlerLogic(msg interface{}, ctx actor.Context) {
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
		//TODO 同一个房间?
		if arg.Roomid != a.id {
		}
		arg.Roomid = a.id
		//玩家切换房间时离开房间
		nodePid.Tell(arg)
		//响应
		rsp := new(pb.LeftDesk)
		if _, ok := a.players[arg.Userid]; ok {
			errcode := a.Leave(arg.Userid)
			rsp.Error = errcode
		}
		ctx.Respond(rsp)
	case *pb.SyncConfig:
		//更新配置
		arg := msg.(*pb.SyncConfig)
		glog.Debugf("SyncConfig %#v", arg)
		b := make(map[string]data.Game)
		err = json.Unmarshal(arg.Data, &b)
		if err != nil {
			glog.Errorf("syncConfig Unmarshal err %v", err)
			return
		}
		for _, v := range b {
			if a.DeskData.Unique == v.Id {
				deskData := handler.NewDeskData(&v)
				deskData.Rid = a.id
				a.DeskData = deskData
				glog.Infof("SyncConfig v: %#v", v)
				glog.Infof("SyncConfig deskData: %#v", deskData)
				return
			}
		}
	case *pb.PrintDesk:
		//打印牌局状态信息,test
		a.printOver()
	case *pb.EnterDesk:
		arg := msg.(*pb.EnterDesk)
		glog.Debugf("EnterDesk %#v", arg)
		a.enterDesk(arg, ctx)
	case *pb.OfflineDesk:
		arg := msg.(*pb.OfflineDesk)
		glog.Debugf("OfflineDesk %#v", arg)
		//离线消息
		if _, ok := a.players[arg.Userid]; ok {
			a.offline[arg.Userid] = true
		}
	case *pb.ChangeCurrency:
		arg := msg.(*pb.ChangeCurrency)
		//充值或购买同步
		a.changeCurrency(arg)
	case *pb.PushDeskState:
		arg := msg.(*pb.PushDeskState)
		glog.Debugf("PushDeskState %#v", arg)
		a.deakState(arg, ctx)
	case *pb.CChatText:
		arg := msg.(*pb.CChatText)
		glog.Debugf("CChatText %#v", arg)
		userid := a.getRouter(ctx)
		glog.Debugf("CChatText %s", userid)
		//房间消息广播,聊天
		a.broadcast(handler.ChatMsg(0, userid, arg.Content))
	case *pb.CChatVoice:
		arg := msg.(*pb.CChatVoice)
		glog.Debugf("CChatVoice %#v", arg)
		userid := a.getRouter(ctx)
		glog.Debugf("CChatVoice %s", userid)
		//房间消息广播,聊天
		a.broadcast(handler.ChatMsg2(0, userid, arg.Content))
	case *pb.CHuiYinDealer:
		arg := msg.(*pb.CHuiYinDealer)
		glog.Debugf("CHuiYinDealer %#v", arg)
		//userid := a.router[ctx.Sender().String()]
		userid := a.getRouter(ctx)
		glog.Debugf("CHuiYinDealer %s", userid)
		var state uint32 = arg.GetState()
		var num uint32 = arg.GetNum()
		errcode := a.BeDealer(userid, state, num)
		if errcode == pb.OK {
			return
		}
		//响应
		rsp := new(pb.SHuiYinDealer)
		rsp.Error = errcode
		ctx.Respond(rsp)
	case *pb.CHuiYinDealerList:
		arg := msg.(*pb.CHuiYinDealerList)
		glog.Debugf("CHuiYinDealerList %#v", arg)
		//userid := a.getRouter(ctx)
		//glog.Debugf("CHuiYinDealerList %s", userid)
		//上庄列表
		rsp := new(pb.SHuiYinDealerList)
		rsp.List = a.dealerListMsg()
		ctx.Respond(rsp)
	case *pb.CHuiYinRoomBet:
		arg := msg.(*pb.CHuiYinRoomBet)
		glog.Debugf("CHuiYinRoomBet %#v", arg)
		userid := a.getRouter(ctx)
		glog.Debugf("CHuiYinRoomBet %s", userid)
		value := arg.GetValue()
		seatBet := arg.GetSeatbet()
		glog.Debugf("CHuiYinRoomBet %#v", arg)
		errcode := a.ChoiceBet(userid, seatBet, int64(value))
		if errcode == pb.OK {
			return
		}
		//响应
		rsp := new(pb.SHuiYinRoomBet)
		rsp.Error = errcode
		ctx.Respond(rsp)
	case *pb.CHuiYinEnterRoom:
		arg := msg.(*pb.CHuiYinEnterRoom)
		glog.Debugf("CHuiYinEnterRoom %#v", arg)
		glog.Debugf("CHuiYinEnterRoom %s", ctx.Sender().String())
		userid := a.getRouter(ctx)
		glog.Debugf("CHuiYinEnterRoom %s", userid)
		//检测重复进入
		msg1 := a.enterMsg(userid)
		ctx.Respond(msg1)
		//庄信息
		//msg2 := a.pushDealerMsg()
		//ctx.Sender().Tell(msg2)
		//消息
		a.cameinMsg(userid)
	case *pb.CHuiYinLeave:
		arg := msg.(*pb.CHuiYinLeave)
		glog.Debugf("CHuiYinLeave %#v", arg)
		userid := arg.Userid
		if userid == "" {
			userid = a.getRouter(ctx)
		}
		glog.Debugf("CHuiYinLeave %s", userid)
		errcode := a.Leave(userid)
		if errcode != pb.OK {
			return
		}
		glog.Debugf("CHuiYinLeave %s", userid)
		//离开消息
		msg2 := &pb.LeaveDesk{
			Roomid: a.id,
			Userid: userid,
		}
		nodePid.Tell(msg2)
		//离开响应
		ctx.Respond(msg2)
		//响应
		rsp := new(pb.SHuiYinLeave)
		rsp.Error = errcode
		ctx.Respond(rsp)
		delete(a.router, ctx.Sender().String())
	case *pb.CHuiYinRoomRoles:
		arg := msg.(*pb.CHuiYinRoomRoles)
		glog.Debugf("CHuiYinRoomRoles %#v", arg)
		//响应
		a.roleList(ctx)
	case *pb.CHuiYinDeskBetInfo:
		arg := msg.(*pb.CHuiYinDeskBetInfo)
		glog.Debugf("CHuiYinDeskBetInfo %#v", arg)
		//响应
		rsp := a.seatBetInfo(arg.Seat)
		ctx.Respond(rsp)
	case *pb.CHuiYinSit:
		arg := msg.(*pb.CHuiYinSit)
		glog.Debugf("CHuiYinSit %#v", arg)
		userid := a.getRouter(ctx)
		glog.Debugf("CHuiYinSit %s", userid)
		errcode := a.SitDown(userid, arg.Seat, arg.State)
		if errcode == pb.OK {
			return
		}
		rsp := new(pb.SHuiYinSit)
		rsp.Error = errcode
		rsp.State = arg.State
		ctx.Respond(rsp)
	case *pb.CGetTrend:
		arg := msg.(*pb.CGetTrend)
		glog.Debugf("CGetTrend %#v", arg)
		glog.Debugf("CGetTrend %#v", a.Trends)
		rsp := handler.GetHuiYinTrends(a.Trends)
		ctx.Respond(rsp)
	case *pb.CGetOpenResult:
		arg := msg.(*pb.CGetOpenResult)
		glog.Debugf("CGetOpenResult %#v", arg)
		glog.Debugf("CGetOpenResult %#v", a.Trends)
		rsp := handler.GetHuiYinOpenResult(a.Trends)
		ctx.Respond(rsp)
	case *pb.CGetLastWins:
		arg := msg.(*pb.CGetLastWins)
		glog.Debugf("CGetLastWins %#v", arg)
		glog.Debugf("CGetLastWins %#v", a.Winers)
		rsp := handler.GetHuiYinWiners(a.Winers)
		ctx.Respond(rsp)
	case *pb.RobotFake:
		arg := msg.(*pb.RobotFake)
		glog.Debugf("RobotFake %#v", arg)
		if arg.RoomPid != nil && arg.Ltype == a.DeskData.Ltype {
			arg.Roomid = a.id
			_, r := a.realRoles()
			arg.RealNum = r
			switch arg.Type { //1添加,2设置
			case 1:
				arg.FakeNum = uint32(utils.RandIntN(9) + 2)
			case 2:
				if r < 40 {
					arg.FakeNum = uint32(utils.RandIntN(21) + 30)
				}
			}
			arg.RoomPid.Tell(arg)
		}
	case *pb.RobotAllot:
		arg := msg.(*pb.RobotAllot)
		glog.Debugf("RobotAllot %#v", arg)
		if arg.HallPid == nil {
			return
		}
		switch arg.Type {
		case 1:
			msg1 := &pb.RobotMsg{
				EnvBet: arg.EnvBet,
				Roomid: a.id,
				Rtype:  a.DeskData.Rtype,
				Ltype:  a.DeskData.Ltype,
			}
			f, r := a.realRoles()
			if f >= 15 { //TODO 已经有7个机器人
				return
			}
			if r < 20 {
				msg1.Num = uint32(utils.RandIntN(4) + 2)
				arg.HallPid.Tell(msg1)
			} else if r < 50 {
				msg1.Num = uint32(utils.RandIntN(2) + 2)
				arg.HallPid.Tell(msg1)
			}
		case 2:
			f, r := a.realRoles()
			if f >= 15 {
				return //TODO 已经有7个机器人
			}
			glog.Debugf("f %d, r %d", f, r)
			//TODO 存在真实玩家且有下注
			//if r != 0 && a.realRoleBet() {
			//	msg1 := &pb.RobotMsg{
			//		EnvBet: arg.EnvBet,
			//		Roomid: a.id,
			//		Rtype:  a.DeskData.Rtype,
			//		Ltype:  a.DeskData.Ltype,
			//		Num:    1,
			//	}
			//	arg.HallPid.Tell(msg1)
			//}
			//TODO 暂时不限制
			msg1 := &pb.RobotMsg{
				EnvBet: arg.EnvBet,
				Roomid: a.id,
				Rtype:  a.DeskData.Rtype,
				Ltype:  a.DeskData.Ltype,
				Num:    15 - f,
			}
			arg.HallPid.Tell(msg1)
		}
	case *pb.SPing:
		for k, v := range a.players {
			if !v.GetRobot() {
				continue
			}
			if p, ok := a.pids[k]; ok && p != nil {
				p.Tell(msg)
			}
		}
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

//状态更新
func (a *Desk) deakState(arg *pb.PushDeskState, ctx actor.Context) {
	a.nexttime = arg.Nexttime
	glog.Debugf("deakState %d, %d", a.state, arg.State)
	if a.state == arg.State {
		return
	}
	a.state = arg.State
	glog.Debugf("deakState %d, %d", a.state, len(a.pids))
	timer := a.nexttime - utils.Timestamp()
	if timer < 0 {
		timer = 0
	}
	glog.Debugf("deakState %d, %d", a.nexttime, timer)
	msg1 := &pb.SHuiYinDeskState{
		State:    a.state,
		Nexttime: timer,
	}
	//广播
	a.broadcast(msg1)
	switch arg.State {
	case data.STATE_BET: //开始下注
		a.gameStart()
	case data.STATE_SEAL: //封盘
		if arg.Expect != "" {
			//封盘重启
			a.setExpect(arg)
		}
		a.gameSeal()
	case data.STATE_OVER: //结算
		a.setExpect(arg)
		a.gameOver()
	}
}

//设置期号
func (a *Desk) setExpect(arg *pb.PushDeskState) {
	a.HuiYinDeskData.expect = arg.Expect
	a.HuiYinDeskData.opencode = arg.Opencode
	a.HuiYinDeskData.opentime = arg.Opentime
	a.HuiYinDeskData.opentimestamp = arg.Opentimestamp
	//封盘时重启服务设置上期期号和点数
	switch arg.State {
	case data.STATE_SEAL:
		if a.dealCard() {
			for k, v := range a.handCards {
				//点数大小
				a.power[k] = algo.Point(v)
			}
			a.gameInit()
		}
	}
}

//获取路由
func (a *Desk) getRouter(ctx actor.Context) string {
	glog.Debugf("getRouter %s", ctx.Sender().String())
	return a.router[ctx.Sender().String()]
}

//进入
func (a *Desk) enterDesk(arg *pb.EnterDesk, ctx actor.Context) {
	rsp := new(pb.EnteredDesk)
	user := new(data.User)
	err2 := json.Unmarshal(arg.Data, user)
	if err2 != nil {
		glog.Errorf("user Unmarshal err %v", err2)
		rsp.Error = pb.RoomNotExist
		ctx.Respond(rsp)
		return
	}
	errcode := a.Enter(user)
	if errcode != pb.OK {
		glog.Errorf("entry Desk err: %d", errcode)
		rsp.Error = errcode
		ctx.Respond(rsp)
		return
	}
	rsp.Roomid = a.id
	rsp.Rtype = a.DeskData.Rtype
	rsp.Gtype = a.DeskData.Gtype
	rsp.Userid = user.GetUserid()
	ctx.Respond(rsp)
	//加入游戏
	a.pids[user.Userid] = arg.Sender
	//设置路由
	a.router[arg.Sender.String()] = user.GetUserid()
	//进入消息
	msg3 := new(pb.JoinDesk)
	msg3.Roomid = a.id
	msg3.Rtype = a.DeskData.Rtype
	msg3.Gtype = a.DeskData.Gtype
	msg3.Userid = user.Userid
	msg3.Sender = arg.Sender
	nodePid.Request(msg3, ctx.Self())
}

//进入消息
func (a *Desk) cameinMsg(userid string) {
	msg2 := new(pb.SHuiYinCamein)
	user := a.getPlayer(userid)
	if user == nil {
		glog.Debugf("cameinMsg err %s", userid)
		return
	}
	msg2.Userdata = handler.PackUserData(user)
	a.broadcast(msg2)
}

//更新货币
func (a *Desk) changeCurrency(arg *pb.ChangeCurrency) {
	p := a.getPlayer(arg.Userid)
	if p == nil {
		glog.Debugf("changeCurrency err %s", arg.Userid)
		return
	}
	p.AddCurrency(arg.Diamond, arg.Coin, arg.Card, arg.Chip)
}
