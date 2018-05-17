package main

import (
	"time"

	"goplays/game/handler"
	"goplays/glog"
	"goplays/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
)

func (ws *WSConn) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.SetLogined:
		//设置连接,还未登录
		arg := msg.(*pb.SetLogined)
		//glog.Debugf("SetLogined %#v", arg)
		ws.dbmsPid = arg.DbmsPid
		ws.roomPid = arg.RoomPid
		ws.rolePid = arg.RolePid
		ws.hallPid = arg.HallPid
		//glog.Infof("SetLogined %v", arg.HallPid)
	case *pb.ServeClose:
		arg := new(pb.SLoginOut)
		arg.Rtype = 2 //停服
		ws.Send(arg)
		//断开连接
		ws.Close()
	case *pb.ServeStop:
		ws.stop(ctx)
		//响应
		//rsp := new(pb.ServeStarted)
		//ctx.Respond(rsp)
	case *pb.ServeStoped:
	case *pb.ServeStart:
		ws.start(ctx)
		//响应
		//rsp := new(pb.ServeStarted)
		//ctx.Respond(rsp)
	case *pb.ServeStarted:
	case *pb.Tick:
		ws.ding(ctx)
	case *pb.LoginElse:
		ws.loginElse() //别处登录
	case *pb.PayCurrency:
		arg := msg.(*pb.PayCurrency)
		glog.Debugf("PayCurrency %#v", arg)
		//后台或充值同步到game房间
		if ws.gamePid != nil {
			msg2 := handler.Pay2ChangeCurr(arg)
			ws.gamePid.Tell(msg2)
		} else {
			diamond := arg.Diamond
			coin := arg.Coin
			chip := arg.Chip
			card := arg.Card
			ltype := arg.Type
			ws.addCurrency(diamond, coin, card, chip, ltype)
		}
	case *pb.ChangeCurrency:
		//货币变更
		arg := msg.(*pb.ChangeCurrency)
		diamond := arg.Diamond
		coin := arg.Coin
		chip := arg.Chip
		card := arg.Card
		ltype := arg.Type
		ws.addCurrency(diamond, coin, card, chip, ltype)
	case proto.Message:
		//响应消息
		//ws.Send(msg)
		ws.HandlerLogin(msg, ctx)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

func (ws *WSConn) start(ctx actor.Context) {
	glog.Infof("ws start: %v", ctx.Self().String())
	ctx.SetReceiveTimeout(waitForLogin) //login timeout set
	set := &pb.SetLogin{
		Sender: ctx.Self(),
	}
	nodePid.Tell(set)
}

//时钟
func (ws *WSConn) ticker(ctx actor.Context) {
	tick := time.Tick(time.Second)
	msg := new(pb.Tick)
	for {
		select {
		case <-ws.stopCh:
			glog.Info("ws ticker closed")
			return
		default: //防止阻塞
		}
		select {
		case <-ws.stopCh:
			glog.Info("ws ticker closed")
			return
		case <-tick:
			if ws.pid != nil {
				ws.pid.Tell(msg)
			}
		}
	}
}

//30秒同步一次
func (ws *WSConn) ding(ctx actor.Context) {
	ws.timer += 1
	if ws.timer != 30 {
		return
	}
	ws.timer = 0
	if !ws.online {
		//断开连接
		ws.Close()
		return
	}
	//同步数据
	ws.syncUser()
}

func (ws *WSConn) stop(ctx actor.Context) {
	glog.Infof("ws stop: %v", ctx.Self().String())
	//已经断开,在别处登录
	if !ws.online {
		//直接关闭
		glog.Infof("ws  %s stoped", ws.pid.String())
		//ws.pid.Stop()
		return
	}
	//离开消息
	ws.leaveDesk()
	if ws.User == nil {
		//直接关闭
		glog.Infof("ws  %s stoped", ws.pid.String())
		//ws.pid.Stop()
		return
	}
	//回存数据
	ws.syncUser()
	//登出日志
	if ws.dbmsPid != nil {
		msg2 := &pb.LogLogout{
			Userid: ws.User.Userid,
			Event:  1, //正常断开
		}
		ws.dbmsPid.Tell(msg2)
	}
	//断开处理
	msg := &pb.Logout{
		Sender: ctx.Self(),
		Userid: ws.User.Userid,
	}
	nodePid.Tell(msg)
	//表示已经断开
	ws.online = false
}

func (ws *WSConn) loginElse() {
	arg := new(pb.SLoginOut)
	glog.Debugf("SLoginOut %s", ws.User.Userid)
	arg.Rtype = 1 //别处登录
	ws.Send(arg)
	//已经断开
	if !ws.online {
		return
	}
	//离开消息
	ws.leaveDesk()
	//同步数据
	ws.syncUser()
	//登出日志
	msg3 := &pb.LogLogout{
		Userid: ws.User.Userid,
		Event:  4, //别处登录
	}
	ws.dbmsPid.Tell(msg3)
	//表示已经断开
	ws.online = false
	//断开连接
	ws.Close()
}

//离开游戏处理
func (ws *WSConn) leaveDesk() {
	if ws.gamePid == nil {
		return
	}
	glog.Debugf("leaveDesk %s", ws.gamePid.String())
	//站起
	msg1 := new(pb.CHuiYinSit)
	msg1.State = false
	ws.gamePid.Tell(msg1)
	//离线
	msg2 := new(pb.OfflineDesk)
	if ws.User != nil {
		msg2.Userid = ws.User.GetUserid()
	}
	ws.gamePid.Tell(msg2)
	//下线
	msg3 := new(pb.CHuiYinLeave)
	if ws.User != nil {
		msg3.Userid = ws.User.GetUserid()
	}
	ws.gamePid.Tell(msg3)
}

/*
func (ws *WSConn) addPrize(rtype, ltype, amount int32) {
	switch uint32(rtype) {
	case data.DIAMOND:
		ws.addCurrency(amount, 0, 0, 0, ltype)
	case data.COIN:
		ws.addCurrency(0, amount, 0, 0, ltype)
	case data.CARD:
		ws.addCurrency(0, 0, amount, 0, ltype)
	case data.CHIP:
		ws.addCurrency(0, 0, 0, amount, ltype)
	}
}

//消耗钻石
func (ws *WSConn) expend(cost uint32, ltype int32) {
	diamond := -1 * int64(cost)
	ws.addCurrency(diamond, 0, 0, 0, ltype)
}
*/

//奖励发放
func (ws *WSConn) addCurrency(diamond, coin, card, chip int64, ltype int32) {
	if ws.User == nil {
		glog.Errorf("add currency user err: %d", ltype)
		return
	}
	ws.User.AddCurrency(diamond, coin, card, chip)
	//货币变更及时同步
	msg2 := handler.ChangeCurrencyMsg(diamond, coin,
		card, chip, ltype, ws.User.GetUserid())
	ws.rolePid.Tell(msg2)
	//消息
	msg := handler.PushCurrencyMsg(diamond, coin,
		card, chip, ltype)
	ws.Send(msg)
	//机器人不写日志
	if ws.User.GetRobot() {
		return
	}
	ws.status = true
	//日志
	if diamond != 0 {
		msg1 := handler.LogDiamondMsg(diamond, ltype, ws.User)
		ws.dbmsPid.Tell(msg1)
	}
	if coin != 0 {
		msg1 := handler.LogCoinMsg(coin, ltype, ws.User)
		ws.dbmsPid.Tell(msg1)
	}
	if card != 0 {
		msg1 := handler.LogCardMsg(card, ltype, ws.User)
		ws.dbmsPid.Tell(msg1)
	}
	if chip != 0 {
		msg1 := handler.LogChipMsg(chip, ltype, ws.User)
		ws.dbmsPid.Tell(msg1)
	}
}

//同步数据
func (ws *WSConn) syncUser() {
	if ws.User == nil {
		return
	}
	if ws.rolePid == nil {
		return
	}
	if !ws.status { //有变更才同步
		return
	}
	ws.status = false
	msg := new(pb.SyncUser)
	msg.Userid = ws.User.GetUserid()
	result, err := json.Marshal(ws.User)
	if err != nil {
		glog.Errorf("user %s Marshal err %v", ws.User.GetUserid(), err)
		return
	}
	msg.Data = result
	ws.rolePid.Tell(msg)
}
