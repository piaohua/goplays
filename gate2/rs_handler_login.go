package main

import (
	"goplays/data"
	"goplays/game/config"
	"goplays/glog"
	"goplays/pb"
	"utils"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
)

//玩家数据请求处理
func (rs *RoleActor) HandlerLogin(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.CRegist:
		//注册消息
		arg := msg.(*pb.CRegist)
		glog.Debugf("CRegist %#v", arg)
		//重复登录
		stoc := new(pb.SRegist)
		stoc.Error = pb.RepeatLogin
		rs.Send(stoc)
	case *pb.CLogin:
		//登录消息
		arg := msg.(*pb.CLogin)
		glog.Debugf("CLogin %#v", arg)
		//重复登录
		stoc := new(pb.SLogin)
		stoc.Error = pb.RepeatLogin
		rs.Send(stoc)
	case *pb.CWxLogin:
		//登录消息
		arg := msg.(*pb.CWxLogin)
		glog.Debugf("CWxLogin %#v", arg)
		//重复登录
		stoc := new(pb.SWxLogin)
		stoc.Error = pb.RepeatLogin
		rs.Send(stoc)
	case *pb.CResetPwd:
		//重置密码消息
		arg := msg.(*pb.CResetPwd)
		glog.Debugf("CResetPwd %#v", arg)
		//重复登录
		stoc := new(pb.SResetPwd)
		stoc.Error = pb.RepeatLogin
		rs.Send(stoc)
	case *pb.CTourist:
		//登录消息
		arg := msg.(*pb.CTourist)
		glog.Debugf("CTourist %#v", arg)
		//重复登录
		stoc := new(pb.STourist)
		stoc.Error = pb.RepeatLogin
		rs.Send(stoc)
	case *pb.LoginSuccess:
		//登录成功处理
		arg := msg.(*pb.LoginSuccess)
		glog.Debugf("LoginSuccess %#v", arg)
		rs.logined(arg, ctx)
	case *pb.LoginElse:
		rs.loginElse() //别处登录
	case proto.Message:
		//glog.Errorf("unknown message %v", msg)
		if rs.User == nil {
			glog.Errorf("user empty message %v", msg)
			return
		}
		rs.HandlerUser(msg, ctx)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

//别处登录
func (rs *RoleActor) loginElse() {
	arg := new(pb.SLoginOut)
	glog.Debugf("SLoginOut %s", rs.User.Userid)
	arg.Rtype = 1 //别处登录
	rs.Send(arg)
	//已经断开
	if !rs.online {
		return
	}
	//断开连接
	rs.CloseWs()
	//离开消息
	rs.leaveDesk()
	//同步数据
	rs.syncUser()
	//登出日志
	msg3 := &pb.LogLogout{
		Userid: rs.User.Userid,
		Event:  4, //别处登录
	}
	rs.dbmsPid.Tell(msg3)
	//表示已经断开
	rs.online = false
}

//离开游戏处理
func (rs *RoleActor) leaveDesk() {
	if rs.gamePid == nil {
		return
	}
	//站起
	msg1 := new(pb.CHuiYinSit)
	msg1.State = false
	rs.gamePid.Tell(msg1)
	//离线
	msg2 := new(pb.OfflineDesk)
	if rs.User != nil {
		msg2.Userid = rs.User.GetUserid()
	}
	rs.gamePid.Tell(msg2)
	//下线
	msg3 := new(pb.CHuiYinLeave)
	if rs.User != nil {
		msg3.Userid = rs.User.GetUserid()
	}
	rs.gamePid.Tell(msg3)
}

//登录成功处理
func (rs *RoleActor) logined(arg *pb.LoginSuccess, ctx actor.Context) {
	rs.CloseWs()
	rs.wsPid = arg.WsPid
	//头像
	rs.setHeadImag(arg.IsRegist, ctx)
	//日志
	rs.loginedLog(arg)
	//登录成功
	rs.online = true
}

//默认头像
func (rs *RoleActor) setHeadImag(isRegist bool, ctx actor.Context) {
	if !isRegist {
		return
	}
	if rs.User == nil {
		return
	}
	if rs.User.GetPhoto() != "" {
		return
	}
	if len(HeadImagList) == 0 {
		return
	}
	head := cfg.Section("domain").Key("headimag").Value()
	if head == "" {
		return
	}
	i := utils.RandIntN(len(HeadImagList))
	rs.User.Photo = head + "/" + HeadImagList[i].Photo
}

//登录成功日志处理
func (rs *RoleActor) loginedLog(arg *pb.LoginSuccess) {
	rs.User.LoginIp = arg.Ip
	rs.User.LoginTime = utils.BsonNow()
	if arg.IsRegist {
		//注册ip
		rs.User.RegistIp = arg.Ip
		if !rs.User.IsTourist() {
			//注册奖励发放
			var diamond int64 = int64(config.GetEnv(data.ENV1))
			var coin int64 = int64(config.GetEnv(data.ENV2))
			var chip int64 = int64(config.GetEnv(data.ENV3))
			var card int64 = int64(config.GetEnv(data.ENV4))
			rs.addCurrency(diamond, coin, card, chip, data.LogType1)
			//注册日志
			msg1 := &pb.LogRegist{
				Userid:   rs.User.Userid,
				Nickname: rs.User.Nickname,
				Ip:       arg.Ip,
			}
			rs.dbmsPid.Tell(msg1)
		}
	}
	//登录日志
	msg2 := &pb.LogLogin{
		Userid: rs.User.Userid,
		Ip:     arg.Ip,
	}
	rs.dbmsPid.Tell(msg2)
}

//发送消息
func (rs *RoleActor) Send(msg interface{}) {
	//glog.Debugf("Send %#v", msg)
	if rs.stopCh == nil {
		glog.Errorf("rs msg channel closed %v", msg)
		return
	}
	if rs.wsPid == nil {
		glog.Errorf("ws pid stoped %v", msg)
		return
	}
	//glog.Debugf("send message %s", rs.wsPid.String())
	select {
	case <-rs.stopCh:
		return
	default:
	}
	select {
	case <-rs.stopCh:
		return
	default:
		//glog.Debugf("send message %#v", msg)
		rs.wsPid.Tell(msg)
	}
}

//关闭
func (rs *RoleActor) StopRs() {
	select {
	case <-rs.stopCh:
		return
	default:
		//停止发送消息
		close(rs.stopCh)
	}
	//停止
	rs.pid.Stop()
}

//关闭连接
func (rs *RoleActor) CloseWs() {
	if rs.wsPid == nil {
		return
	}
	glog.Debugf("CloseWs pid : %s", rs.wsPid.String())
	msg1 := new(pb.ServeStop)
	//关闭连接
	rs.wsPid.Tell(msg1)
	//断开
	rs.wsPid = nil
}
