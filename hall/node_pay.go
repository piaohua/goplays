package main

import (
	"goplays/glog"
	"goplays/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//玩家支付请求处理
func (a *HallActor) HandlerPay(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	//case *pb.WxpayCallback:
	//	arg := msg.(*pb.WxpayCallback)
	//	glog.Debugf("WxpayCallback: %v", arg)
	//	//支付回调
	//	for k, v := range a.serve {
	//		//TODO 优化
	//		if strings.Contains(k, "gate.node") {
	//			v.Tell(arg)
	//			//TODO 优化
	//			//选择一个验证即可,
	//			//暂时不知道哪个节点的订单
	//			break
	//		}
	//	}
	//case *pb.WxpayGoods:
	//	arg := msg.(*pb.WxpayGoods)
	//	//发货
	//	glog.Debugf("WxpayGoods: %v", arg)
	//	userid := arg.Userid
	//	gate := a.roles[userid]
	//	if v, ok := a.serve[gate]; ok {
	//		v.Tell(arg)
	//	} else {
	//		glog.Errorf("WxpayGoods: %v", arg)
	//	}
	case *pb.ChangeCurrency:
		//货币变更同步
		arg := msg.(*pb.ChangeCurrency)
		glog.Infof("ChangeCurrency %#v", arg)
		//消息
		a.msg2role(arg.Userid, arg)
	case *pb.PayCurrency:
		//充值同步到game房间
		arg := msg.(*pb.PayCurrency)
		glog.Debugf("PayCurrency %#v", arg)
		//消息
		a.msg2role(arg.Userid, arg)
	default:
		//glog.Errorf("unknown message %v", msg)
		a.HandlerDesk(msg, ctx)
	}
}
