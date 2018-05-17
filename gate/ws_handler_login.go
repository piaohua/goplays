package main

import (
	"time"

	"goplays/data"
	"goplays/game/config"
	"goplays/game/login"
	"goplays/glog"
	"goplays/pb"
	"utils"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//玩家数据请求处理
func (ws *WSConn) HandlerLogin(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.CRegist:
		//注册消息
		arg := msg.(*pb.CRegist)
		glog.Debugf("CRegist %#v", arg)
		ws.regist(arg, ctx)
	case *pb.CLogin:
		//登录消息
		arg := msg.(*pb.CLogin)
		glog.Debugf("CLogin %#v", arg)
		ws.login(arg, ctx)
	case *pb.CWxLogin:
		//登录消息
		arg := msg.(*pb.CWxLogin)
		glog.Debugf("CWxLogin %#v", arg)
		ws.wxlogin(arg, ctx)
	case *pb.CResetPwd:
		//重置密码消息
		arg := msg.(*pb.CResetPwd)
		glog.Debugf("CResetPwd %#v", arg)
		ws.resetPwd(arg, ctx)
	case *pb.CTourist:
		//登录消息
		arg := msg.(*pb.CTourist)
		glog.Debugf("CTourist %#v", arg)
		ws.touristLogin(arg, ctx)
	default:
		//glog.Errorf("unknown message %v", msg)
		if ws.User == nil {
			glog.Errorf("user empty message %v", msg)
			return
		}
		ws.HandlerUser(msg, ctx)
	}
}

func (ws *WSConn) resetPwd(arg *pb.CResetPwd, ctx actor.Context) {
	stoc := login.RestPwdCheck(arg)
	if stoc.Error != pb.OK {
		ws.Send(stoc)
		return
	}
	//已经登录
	if ws.online {
		stoc.Error = pb.RepeatLogin
		ws.Send(stoc)
		return
	}
	//重置
	res1 := ws.reqRole(arg, ctx)
	var response1 *pb.SResetPwd
	var ok bool
	if response1, ok = res1.(*pb.SResetPwd); ok {
		if response1.Error != pb.OK {
			glog.Errorf("CRegist fail %d", response1.Error)
			stoc.Error = response1.Error
			ws.Send(stoc)
			return
		}
	} else {
		glog.Error("CResetPwd fail")
		stoc.Error = pb.ResetPwdFaild
		ws.Send(stoc)
		return
	}
	userid := response1.GetUserid()
	glog.Debugf("CResetPwd successfully %s", userid)
	if !ws.logining(userid, ctx) {
		glog.Debugf("CResetPwd failed %s", ctx.Self())
		stoc.Error = pb.ResetPwdFaild
		ws.Send(stoc)
		return
	}
	stoc.Userid = userid
	glog.Debugf("CResetPwd successfully %s", userid)
	ws.Send(stoc)
	//成功后处理
	ws.logined(true, ctx)
}

func (ws *WSConn) regist(arg *pb.CRegist, ctx actor.Context) {
	stoc := login.RegistCheck(arg)
	if stoc.Error != pb.OK {
		ws.Send(stoc)
		return
	}
	//重复登录
	if ws.online {
		stoc.Error = pb.RepeatLogin
		ws.Send(stoc)
		return
	}
	msg1 := new(pb.RoleRegist)
	msg1.Phone = arg.GetPhone()
	msg1.Nickname = arg.GetNickname()
	msg1.Password = arg.GetPassword()
	msg1.Smscode = arg.GetSmscode()
	//注册
	res1 := ws.reqRole(msg1, ctx)
	var response1 *pb.RoleRegisted
	var ok bool
	if response1, ok = res1.(*pb.RoleRegisted); ok {
		if response1.Error != pb.OK {
			glog.Errorf("CRegist fail %d", response1.Error)
			stoc.Error = response1.Error
			ws.Send(stoc)
			return
		}
	} else {
		glog.Error("CRegist fail")
		stoc.Error = pb.RegistError
		ws.Send(stoc)
		return
	}
	userid := response1.GetUserid()
	glog.Debugf("regist successfully %s", userid)
	if !ws.logining(userid, ctx) {
		glog.Debugf("regist failed %s", ctx.Self())
		stoc.Error = pb.RegistError
		ws.Send(stoc)
		return
	}
	stoc.Userid = userid
	glog.Debugf("regist successfully %s", userid)
	ws.Send(stoc)
	//成功后处理
	ws.logined(true, ctx)
}

func (ws *WSConn) login(arg *pb.CLogin, ctx actor.Context) {
	//检测参数
	stoc := login.LoginCheck(arg)
	if stoc.Error != pb.OK {
		ws.Send(stoc)
		return
	}
	//重复登录
	if ws.online {
		stoc.Error = pb.RepeatLogin
		ws.Send(stoc)
		return
	}
	msg1 := new(pb.RoleLogin)
	msg1.Phone = arg.GetPhone()
	msg1.Password = arg.GetPassword()
	//登录
	res1 := ws.reqRole(msg1, ctx)
	var response1 *pb.RoleLogined
	var ok bool
	if response1, ok = res1.(*pb.RoleLogined); ok {
		if response1.Error != pb.OK {
			glog.Errorf("CLogin fail %d", response1.Error)
			stoc.Error = response1.Error
			ws.Send(stoc)
			return
		}
	} else {
		glog.Error("CLogin fail")
		stoc.Error = pb.LoginError
		ws.Send(stoc)
		return
	}
	userid := response1.GetUserid()
	glog.Debugf("login successfully %s", userid)
	if !ws.logining(userid, ctx) {
		glog.Debugf("login failed %s", ctx.Self())
		stoc.Error = pb.LoginError
		ws.Send(stoc)
		return
	}
	stoc.Userid = userid
	glog.Debugf("login successfully %s", userid)
	ws.Send(stoc)
	//成功后处理
	ws.logined(false, ctx)
}

//游客登录
func (ws *WSConn) touristLogin(arg *pb.CTourist, ctx actor.Context) {
	//检测参数
	key := cfg.Section("gate").Key("tourist").Value()
	stoc := login.TouristLoginCheck(arg, key)
	if stoc.Error != pb.OK {
		ws.Send(stoc)
		return
	}
	//重复登录
	if ws.online {
		stoc.Error = pb.RepeatLogin
		ws.Send(stoc)
		return
	}
	msg1 := new(pb.TouristLogin)
	msg1.Account = arg.GetAccount()
	msg1.Password = arg.GetPassword()
	msg1.Registip = ws.GetIPAddr()
	//登录
	res1 := ws.reqRole(msg1, ctx)
	var response1 *pb.TouristLogined
	var ok bool
	if response1, ok = res1.(*pb.TouristLogined); ok {
		if response1.Error != pb.OK {
			glog.Errorf("CTourist fail %d", response1.Error)
			stoc.Error = response1.Error
			ws.Send(stoc)
			return
		}
	} else {
		glog.Error("CTourist fail")
		stoc.Error = pb.LoginError
		ws.Send(stoc)
		return
	}
	userid := response1.GetUserid()
	glog.Debugf("tourist login successfully %s", userid)
	if !ws.logining(userid, ctx) {
		glog.Debugf("login failed %s", ctx.Self())
		stoc.Error = pb.LoginError
		ws.Send(stoc)
		return
	}
	stoc.Userid = userid
	glog.Debugf("tourist login successfully %s", userid)
	ws.Send(stoc)
	//成功后处理
	ws.logined(response1.IsRegist, ctx)
}

//微信
func (ws *WSConn) wxlogin(arg *pb.CWxLogin, ctx actor.Context) {
	stoc, wxdata := login.WxLoginCheck(arg)
	if stoc.Error != pb.OK {
		ws.Send(stoc)
		return
	}
	//重复登录
	if ws.online {
		stoc.Error = pb.RepeatLogin
		ws.Send(stoc)
		return
	}
	msg1 := new(pb.WxLogin)
	msg1.Wxuid = wxdata.OpenId
	msg1.Nickname = wxdata.Nickname
	msg1.Photo = wxdata.HeadImagUrl
	msg1.Sex = uint32(wxdata.Sex)
	//登录
	res1 := ws.reqRole(msg1, ctx)
	var response1 *pb.WxLogined
	var ok bool
	if response1, ok = res1.(*pb.WxLogined); ok {
		if response1.Error != pb.OK {
			glog.Errorf("CWxLogin fail %d", response1.Error)
			stoc.Error = response1.Error
			ws.Send(stoc)
			return
		}
	} else {
		glog.Error("CWxLogin fail")
		stoc.Error = pb.GetWechatUserInfoFail
		ws.Send(stoc)
		return
	}
	userid := response1.GetUserid()
	glog.Debugf("weixin login successfully %s", userid)
	if !ws.logining(userid, ctx) {
		stoc.Error = pb.GetWechatUserInfoFail
		ws.Send(stoc)
		return
	}
	stoc.Userid = userid
	glog.Debugf("weixin login successfully %s", userid)
	ws.Send(stoc)
	//成功后处理
	ws.logined(response1.IsRegist, ctx)
}

//登录节点
func (ws *WSConn) loginGate(userid string, ctx actor.Context) bool {
	//成功后登录网关
	msg2 := new(pb.LoginGate)
	msg2.Sender = ctx.Self()
	msg2.Userid = userid
	timeout := 3 * time.Second
	res2, err2 := nodePid.RequestFuture(msg2, timeout).Result()
	if err2 != nil {
		glog.Errorf("LoginGate err: %v", err2)
		return false
	}
	if response2, ok := res2.(*pb.LoginedGate); ok {
		glog.Debugf("response2: %#v", response2)
	}
	return true
}

//登录成功数据处理
func (ws *WSConn) loginUser(userid string, ctx actor.Context) bool {
	msg4 := new(pb.Login)
	msg4.Userid = userid
	res4 := ws.reqRole(msg4, ctx)
	var response4 *pb.Logined
	var ok bool
	if response4, ok = res4.(*pb.Logined); !ok {
		glog.Debugf("loginUser failed %s", userid)
		return false
	}
	//glog.Debugf("response4: %#v", response4)
	//数据
	ws.User = new(data.User)
	err2 := json.Unmarshal(response4.Data, ws.User)
	if err2 != nil {
		glog.Errorf("user Unmarshal err %v", err2)
		return false
	}
	if ws.User.GetUserid() == "" {
		glog.Error("CLogin fail")
		return false
	}
	return true
}

//登录成功数据处理
func (ws *WSConn) reqRole(msg interface{}, ctx actor.Context) interface{} {
	glog.Debugf("reqRole msg %#v", msg)
	timeout := 3 * time.Second
	res1, err1 := ws.rolePid.RequestFuture(msg, timeout).Result()
	if err1 != nil {
		glog.Errorf("reqRole err: %v, msg %#v", err1, msg)
		return nil
	}
	return res1
}

//登录流程处理
func (ws *WSConn) logining(userid string, ctx actor.Context) bool {
	if userid == "" {
		glog.Debugf("logining failed %s", userid)
		return false
	}
	//TODO 在节点中spawn一个玩家进程
	if !ws.loginGate(userid, ctx) {
		glog.Debugf("logining loginGate failed %s", userid)
		return false
	}
	if !ws.loginUser(userid, ctx) {
		glog.Debugf("logining loginUser failed %s", userid)
		return false
	}
	return true
}

//登录成功处理
func (ws *WSConn) logined(isRegist bool, ctx actor.Context) {
	//头像
	ws.setHeadImag(isRegist, ctx)
	//日志
	ws.loginedLog(isRegist)
	//登录成功
	ws.online = true
	//成功
	ctx.SetReceiveTimeout(0) //login Successfully, timeout off
	//启动时钟
	go ws.ticker(ctx)
}

//默认头像
func (ws *WSConn) setHeadImag(isRegist bool, ctx actor.Context) {
	if !isRegist {
		return
	}
	if ws.User == nil {
		return
	}
	if ws.User.GetPhoto() != "" {
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
	ws.User.Photo = head + "/" + HeadImagList[i].Photo
}

//登录成功日志处理
func (ws *WSConn) loginedLog(isRegist bool) {
	ws.User.LoginIp = ws.GetIPAddr()
	if isRegist {
		//注册ip
		ws.User.RegistIp = ws.GetIPAddr()
		if !ws.User.IsTourist() {
			//注册奖励发放
			var diamond int64 = int64(config.GetEnv(data.ENV1))
			var coin int64 = int64(config.GetEnv(data.ENV2))
			var chip int64 = int64(config.GetEnv(data.ENV3))
			var card int64 = int64(config.GetEnv(data.ENV4))
			ws.addCurrency(diamond, coin, card, chip, data.LogType1)
			//注册日志
			msg1 := &pb.LogRegist{
				Userid:   ws.User.Userid,
				Nickname: ws.User.Nickname,
				Ip:       ws.GetIPAddr(),
			}
			ws.dbmsPid.Tell(msg1)
		}
	}
	//登录日志
	msg2 := &pb.LogLogin{
		Userid: ws.User.Userid,
		Ip:     ws.GetIPAddr(),
	}
	ws.dbmsPid.Tell(msg2)
}
