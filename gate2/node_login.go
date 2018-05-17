package main

import (
	"goplays/data"
	"goplays/glog"
	"goplays/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//玩家数据请求处理
func (a *GateActor) HandlerLogin(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.LoginElse:
		//在别的节点登录
		arg := msg.(*pb.LoginElse)
		glog.Debugf("LoginElse %#v", arg)
		userid := arg.GetUserid()
		if p, ok := a.roles[userid]; ok {
			p.Tell(arg)
			//直接关闭
			p.Tell(new(pb.ServeClose))
		}
		glog.Debugf("LoginElse userid: %s", userid)
		//移除
		delete(a.roles, userid)
		delete(a.offline, userid)
		delete(a.offtime, userid)
	case *pb.SetLogin:
		arg := msg.(*pb.SetLogin)
		//glog.Debugf("SetLogin %#v", arg)
		rsp := &pb.SetLogined{
			RolePid: a.rolePid,
		}
		//arg.Sender.Tell(rsp)
		glog.Infof("SetLogin %s", arg.Sender.String())
		ctx.Respond(rsp)
	case *pb.Logout:
		//登出成功
		arg := msg.(*pb.Logout)
		glog.Debugf("Logout %#v", arg)
		userid := arg.GetUserid()
		glog.Debugf("Logout userid: %s", userid)
		//离线 arg.Sender == a.roles[userid]
		a.offline[userid] = arg.Sender
		a.offtime[userid] = 10 //缓存5分钟
		//移除
		delete(a.roles, userid)
		//不能离开,因为缓存了数据,离开会出现数据不同步
		//a.hallPid.Tell(arg)
		//a.rolePid.Tell(arg)
		//a.roomPid.Tell(arg)
	case *pb.SelectGate:
		//登录成功
		arg := msg.(*pb.SelectGate)
		glog.Debugf("SelectGate %#v", arg)
		a.selectRole(arg, ctx)
	case *pb.Login2Gate:
		//登录成功
		arg := msg.(*pb.Login2Gate)
		glog.Debugf("Login2Gate %#v", arg)
		a.spawnRole(arg, ctx)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

//下线离线玩家
func (a *GateActor) offlineStop(ctx actor.Context) {
	for k, v := range a.offline {
		if a.offtime[k] <= 0 {
			v.Tell(new(pb.OfflineStop))
			delete(a.offline, k)
			delete(a.offtime, k)
			//正式下线消息
			arg := new(pb.Logout)
			arg.Userid = k
			arg.Sender = v
			a.hallPid.Tell(arg)
			a.rolePid.Tell(arg)
			a.roomPid.Tell(arg)
			continue
		}
		a.offtime[k] -= 1
	}
}

//断开其它节点连接
func (a *GateActor) logoutOther(userid string, ctx actor.Context) {
	//TODO 数据一致性,防止数据覆盖
	msg1 := new(pb.LoginHall)
	msg1.Userid = userid
	msg1.NodeName = a.Name
	a.hallPid.Tell(msg1)
}

//登录成功查询
func (a *GateActor) selectRole(arg *pb.SelectGate, ctx actor.Context) {
	userid := arg.GetUserid()
	glog.Debugf("SelectGate userid: %s", userid)
	//断开其它节点连接
	a.logoutOther(userid, ctx)
	//在线表查询
	if p, ok := a.roles[userid]; ok {
		//断开当前节点旧连接
		msg1 := new(pb.LoginElse)
		msg1.Userid = userid
		p.Tell(msg1)
		glog.Debugf("SelectGate userid: %s", userid)
		glog.Debugf("SelectGate userid: %v", p.String())
		//响应登录
		rsp := new(pb.SelectedGate)
		rsp.Role = p
		ctx.Respond(rsp)
		return
	}
	//离线表查找
	if p, ok := a.offline[userid]; ok {
		//切换到在线表
		a.roles[userid] = p
		delete(a.offline, userid)
		delete(a.offtime, userid)
		glog.Debugf("SelectGate offline userid: %s", userid)
		glog.Debugf("SelectGate offline userid: %v", p.String())
		//响应登录
		rsp := new(pb.SelectedGate)
		rsp.Role = p
		ctx.Respond(rsp)
		return
	}
	//响应登录,不存在
	rsp := new(pb.SelectedGate)
	rsp.Error = pb.Failed
	ctx.Respond(rsp)
}

//新玩家
func (a *GateActor) spawnRole(arg *pb.Login2Gate, ctx actor.Context) {
	userid := arg.GetUserid()
	rsp := new(pb.Logined2Gate)
	if _, ok := a.roles[userid]; ok {
		//算失败处理
		rsp.Error = pb.Failed
		ctx.Respond(rsp)
		return
	} else if _, ok := a.offline[userid]; ok {
		//算失败处理
		rsp.Error = pb.Failed
		ctx.Respond(rsp)
		return
	}
	user := new(data.User)
	err2 := json.Unmarshal(arg.Data, user)
	if err2 != nil {
		glog.Errorf("user Unmarshal err %v", err2)
		rsp.Error = pb.Failed
		ctx.Respond(rsp)
		return
	}
	if user.GetUserid() == "" {
		glog.Error("CLogin fail")
		rsp.Error = pb.Failed
		ctx.Respond(rsp)
		return
	}
	//新玩家
	newRole := NewRole(user)
	newRole.dbmsPid = a.dbmsPid
	newRole.roomPid = a.roomPid
	newRole.rolePid = a.rolePid
	newRole.hallPid = a.hallPid
	rolePid := newRole.initRs()
	newRole.pid = rolePid
	rolePid.Tell(new(pb.ServeStart))
	rsp.Role = rolePid
	ctx.Respond(rsp)
	//添加
	a.roles[userid] = rolePid
}
