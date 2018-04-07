package main

import (
	"goplays/glog"
	"goplays/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//玩家登录请求处理
func (a *HallActor) HandlerLogin(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.LoginHall:
		arg := msg.(*pb.LoginHall)
		glog.Debugf("LoginHall: %v", arg)
		userid := arg.GetUserid()
		nodeName := arg.GetNodeName()
		//断开旧连接
		if name, ok := a.roles[userid]; ok {
			//已经登录且不同节点
			if p, ok := a.serve[name]; ok && name != nodeName {
				msg1 := new(pb.LoginElse)
				msg1.Userid = userid
				p.Tell(msg1)
			}
		} else {
			//增加
			a.count[nodeName] += 1
		}
		//添加
		a.roles[userid] = nodeName
		//响应登录
		rsp := new(pb.LoginedHall)
		ctx.Respond(rsp)
	case *pb.Logout:
		//登出成功
		arg := msg.(*pb.Logout)
		glog.Debugf("Logout: %v", arg)
		//减少
		userid := arg.GetUserid()
		nodeName := a.roles[userid]
		a.count[nodeName] -= 1
		//移除
		delete(a.roles, userid)
	default:
		//glog.Errorf("unknown message %v", msg)
		a.HandlerPay(msg, ctx)
	}
}
