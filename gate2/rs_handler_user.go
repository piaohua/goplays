package main

import (
	"goplays/data"
	"goplays/game/handler"
	"goplays/glog"
	"goplays/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//玩家数据请求处理
func (rs *RoleActor) HandlerUser(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.CPing:
		arg := msg.(*pb.CPing)
		//glog.Debugf("CPing %#v", arg)
		rsp := handler.Ping(arg)
		rs.Send(rsp)
	case *pb.CNotice:
		arg := msg.(*pb.CNotice)
		glog.Debugf("CNotice %#v", arg)
		rsp := handler.GetNotice(data.NOTICE_TYPE1)
		rs.Send(rsp)
	case *pb.CGetCurrency:
		arg := msg.(*pb.CGetCurrency)
		glog.Debugf("CGetCurrency %#v", arg)
		//响应
		rsp := handler.GetCurrency(arg, rs.User)
		rs.Send(rsp)
	case *pb.CBuy:
		arg := msg.(*pb.CBuy)
		glog.Debugf("CBuy %#v", arg)
		//优化
		rsp, diamond, coin := handler.Buy(arg, rs.User)
		//同步兑换
		rs.addCurrency(diamond, coin, 0, 0, data.LogType18)
		//响应
		rs.Send(rsp)
	case *pb.CShop:
		arg := msg.(*pb.CShop)
		glog.Debugf("CShop %#v", arg)
		//响应
		rsp := handler.Shop(arg, rs.User)
		rs.Send(rsp)
	case *pb.CUserData:
		arg := msg.(*pb.CUserData)
		glog.Debugf("CUserData %#v", arg)
		userid := arg.GetUserid()
		if userid == "" {
			userid = rs.User.GetUserid()
		}
		if userid != rs.User.GetUserid() {
			msg1 := new(pb.GetUserData)
			msg1.Userid = userid
			rs.rolePid.Request(msg1, ctx.Self())
		} else {
			//TODO 添加房间数据返回
			rsp := handler.GetUserDataMsg(arg, rs.User)
			rs.Send(rsp)
		}
	case *pb.GotUserData:
		arg := msg.(*pb.GotUserData)
		glog.Debugf("GotUserData %#v", arg)
		rsp := handler.UserDataMsg(arg)
		rs.Send(rsp)
	default:
		//glog.Errorf("unknown message %v", msg)
		rs.HandlerPay(msg, ctx)
	}
}
