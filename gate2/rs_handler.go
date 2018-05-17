package main

import (
	"time"

	"goplays/game/handler"
	"goplays/glog"
	"goplays/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
)

func (rs *RoleActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.OfflineStop:
		glog.Debugf("rs OfflineStop %s", ctx.Self().String())
		//断开连接
		rs.CloseWs()
		//停止,TODO 暂时使用离线方法
		rs.loginElse()
		//关闭
		rs.StopRs()
	case *pb.ServeClose:
		glog.Debugf("rs ServeClose %s", ctx.Self().String())
		arg := new(pb.SLoginOut)
		arg.Rtype = 2 //停服
		rs.Send(arg)
		//断开连接
		rs.CloseWs()
		//停止,TODO 暂时使用离线方法
		rs.loginElse()
		//关闭
		rs.StopRs()
	case *pb.ServeStop:
		glog.Debugf("rs ServeStop %s", ctx.Self().String())
		rs.stop(ctx)
		//响应
		//rsp := new(pb.ServeStarted)
		//ctx.Respond(rsp)
	case *pb.ServeStoped:
	case *pb.ServeStart:
		glog.Debugf("rs ServeStart %s", ctx.Self().String())
		//启动时钟
		go rs.ticker(ctx)
		//响应
		//rsp := new(pb.ServeStarted)
		//ctx.Respond(rsp)
	case *pb.ServeStarted:
	case *pb.Tick:
		rs.ding(ctx)
	case *pb.PayCurrency:
		arg := msg.(*pb.PayCurrency)
		glog.Debugf("PayCurrency %#v", arg)
		//后台或充值同步到game房间
		if rs.gamePid != nil {
			msg2 := handler.Pay2ChangeCurr(arg)
			rs.gamePid.Tell(msg2)
		}
		diamond := arg.Diamond
		coin := arg.Coin
		chip := arg.Chip
		card := arg.Card
		ltype := arg.Type
		rs.addCurrency(diamond, coin, card, chip, ltype)
	case *pb.ChangeCurrency:
		//货币变更
		arg := msg.(*pb.ChangeCurrency)
		diamond := arg.Diamond
		coin := arg.Coin
		chip := arg.Chip
		card := arg.Card
		ltype := arg.Type
		rs.addCurrency(diamond, coin, card, chip, ltype)
	case proto.Message:
		//响应消息
		//rs.Send(msg)
		rs.HandlerLogin(msg, ctx)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

//时钟
func (rs *RoleActor) ticker(ctx actor.Context) {
	tick := time.Tick(time.Second)
	msg := new(pb.Tick)
	for {
		select {
		case <-rs.stopCh:
			glog.Info("rs ticker closed")
			return
		default: //防止阻塞
		}
		select {
		case <-rs.stopCh:
			glog.Info("rs ticker closed")
			return
		case <-tick:
			ctx.Self().Tell(msg)
		}
	}
}

//30秒同步一次
func (rs *RoleActor) ding(ctx actor.Context) {
	rs.timer += 1
	if rs.timer != 30 {
		return
	}
	rs.timer = 0
	if !rs.online {
		return
	}
	//同步数据
	rs.syncUser()
}

//断线
func (rs *RoleActor) stop(ctx actor.Context) {
	glog.Infof("rs stop: %v", ctx.Self().String())
	//已经断开,在别处登录
	if !rs.online {
		return
	}
	//关闭连接
	rs.CloseWs()
	//离开消息
	rs.leaveDesk()
	//回存数据
	rs.syncUser()
	//登出日志
	msg2 := &pb.LogLogout{
		Userid: rs.User.Userid,
		Event:  1, //正常断开
	}
	rs.dbmsPid.Tell(msg2)
	//断开处理
	msg := &pb.Logout{
		Sender: ctx.Self(),
		Userid: rs.User.Userid,
	}
	nodePid.Tell(msg)
	//表示已经断开
	rs.online = false
}

/*
func (rs *RoleActor) addPrize(rtype, ltype, amount int32) {
	switch uint32(rtype) {
	case data.DIAMOND:
		rs.addCurrency(amount, 0, 0, 0, ltype)
	case data.COIN:
		rs.addCurrency(0, amount, 0, 0, ltype)
	case data.CARD:
		rs.addCurrency(0, 0, amount, 0, ltype)
	case data.CHIP:
		rs.addCurrency(0, 0, 0, amount, ltype)
	}
}

//消耗钻石
func (rs *RoleActor) expend(cost uint32, ltype int32) {
	diamond := -1 * int64(cost)
	rs.addCurrency(diamond, 0, 0, 0, ltype)
}
*/

//奖励发放
func (rs *RoleActor) addCurrency(diamond, coin, card, chip int64, ltype int32) {
	if rs.User == nil {
		glog.Errorf("add currency user err: %d", ltype)
		return
	}
	//日志记录
	if diamond < 0 && ((rs.User.GetDiamond() + diamond) < 0) {
		diamond = 0 - rs.User.GetDiamond()
	}
	if chip < 0 && ((rs.User.GetChip() + chip) < 0) {
		chip = 0 - rs.User.GetChip()
	}
	if coin < 0 && ((rs.User.GetCoin() + coin) < 0) {
		coin = 0 - rs.User.GetCoin()
	}
	if card < 0 && ((rs.User.GetCard() + card) < 0) {
		card = 0 - rs.User.GetCard()
	}
	rs.User.AddCurrency(diamond, coin, card, chip)
	//货币变更及时同步
	msg2 := handler.ChangeCurrencyMsg(diamond, coin,
		card, chip, ltype, rs.User.GetUserid())
	rs.rolePid.Tell(msg2)
	//消息
	msg := handler.PushCurrencyMsg(diamond, coin,
		card, chip, ltype)
	rs.Send(msg)
	//TODO 机器人不写日志
	//if rs.User.GetRobot() {
	//	return
	//}
	rs.status = true
	//日志
	//TODO 日志放在dbms中统一写入
	//if diamond != 0 {
	//	msg1 := handler.LogDiamondMsg(diamond, ltype, rs.User)
	//	rs.dbmsPid.Tell(msg1)
	//}
	//if coin != 0 {
	//	msg1 := handler.LogCoinMsg(coin, ltype, rs.User)
	//	rs.dbmsPid.Tell(msg1)
	//}
	//if card != 0 {
	//	msg1 := handler.LogCardMsg(card, ltype, rs.User)
	//	rs.dbmsPid.Tell(msg1)
	//}
	//if chip != 0 {
	//	msg1 := handler.LogChipMsg(chip, ltype, rs.User)
	//	rs.dbmsPid.Tell(msg1)
	//}
}

//同步数据
func (rs *RoleActor) syncUser() {
	if rs.User == nil {
		return
	}
	if rs.rolePid == nil {
		return
	}
	if !rs.status { //有变更才同步
		return
	}
	rs.status = false
	msg := new(pb.SyncUser)
	msg.Userid = rs.User.GetUserid()
	result, err := json.Marshal(rs.User)
	if err != nil {
		glog.Errorf("user %s Marshal err %v", rs.User.GetUserid(), err)
		return
	}
	msg.Data = result
	rs.rolePid.Tell(msg)
}
