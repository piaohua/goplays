package main

import (
	"goplays/data"
	"goplays/game/handler"
	"goplays/glog"
	"goplays/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//玩家数据请求处理
func (ws *WSConn) HandlerUser(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.CPing:
		arg := msg.(*pb.CPing)
		//glog.Debugf("CPing %#v", arg)
		rsp := handler.Ping(arg)
		ws.Send(rsp)
	case *pb.CNotice:
		arg := msg.(*pb.CNotice)
		glog.Debugf("CNotice %#v", arg)
		rsp := handler.GetNotice(data.NOTICE_TYPE1)
		ws.Send(rsp)
	case *pb.CGetCurrency:
		arg := msg.(*pb.CGetCurrency)
		glog.Debugf("CGetCurrency %#v", arg)
		//响应
		rsp := handler.GetCurrency(arg, ws.User)
		ws.Send(rsp)
	case *pb.CBuy:
		arg := msg.(*pb.CBuy)
		glog.Debugf("CBuy %#v", arg)
		//优化
		rsp, diamond, coin := handler.Buy(arg, ws.User)
		//同步兑换
		ws.addCurrency(diamond, coin, 0, 0, data.LogType18)
		//响应
		ws.Send(rsp)
	case *pb.CShop:
		arg := msg.(*pb.CShop)
		glog.Debugf("CShop %#v", arg)
		//响应
		rsp := handler.Shop(arg, ws.User)
		ws.Send(rsp)
	case *pb.CUserData:
		arg := msg.(*pb.CUserData)
		glog.Debugf("CUserData %#v", arg)
		userid := arg.GetUserid()
		if userid == "" {
			userid = ws.User.GetUserid()
		}
		if userid != ws.User.GetUserid() {
			msg1 := new(pb.GetUserData)
			msg1.Userid = userid
			ws.rolePid.Request(msg1, ctx.Self())
		} else {
			//TODO 添加房间数据返回
			rsp := handler.GetUserDataMsg(arg, ws.User)
			ws.Send(rsp)
		}
	case *pb.GotUserData:
		arg := msg.(*pb.GotUserData)
		glog.Debugf("GotUserData %#v", arg)
		rsp := handler.UserDataMsg(arg)
		ws.Send(rsp)
	default:
		//glog.Errorf("unknown message %v", msg)
		ws.HandlerPay(msg, ctx)
	}
}
