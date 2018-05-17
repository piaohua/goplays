package main

import (
	"time"

	"goplays/data"
	"goplays/game/handler"
	"goplays/glog"
	"goplays/pb"
	"utils"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func (a *RoleActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.Connected:
		arg := msg.(*pb.Connected)
		glog.Infof("Connected %s", arg.Name)
	case *pb.Disconnected:
		arg := msg.(*pb.Disconnected)
		glog.Infof("Disconnected %s", arg.Name)
	case *pb.ServeStop:
		//关闭服务
		a.handlerStop(ctx)
		//响应登录
		rsp := new(pb.ServeStoped)
		ctx.Respond(rsp)
	case *pb.ServeStart:
		a.start(ctx)
		//响应
		//rsp := new(pb.ServeStarted)
		//ctx.Respond(rsp)
	case *pb.Tick:
		a.ding(ctx)
	case *pb.SyncUser:
		arg := msg.(*pb.SyncUser)
		a.syncUser(arg, ctx)
	case *pb.ChangeCurrency:
		arg := msg.(*pb.ChangeCurrency)
		//glog.Debugf("ChangeCurrency %#v", arg)
		//更新货币
		a.syncCurrency(arg.Diamond, arg.Coin, arg.Card,
			arg.Chip, arg.Type, arg.Userid)
	case *pb.PayCurrency:
		arg := msg.(*pb.PayCurrency)
		glog.Debugf("PayCurrency %#v", arg)
		//后台或充值同步到game房间
		a.syncCurrency(arg.Diamond, arg.Coin, arg.Card,
			arg.Chip, arg.Type, arg.Userid)
	case *pb.Login:
		//登录成功
		arg := msg.(*pb.Login)
		glog.Debugf("login : %#v", arg)
		a.logined(arg, ctx)
	case *pb.Logout:
		//登出成功
		arg := msg.(*pb.Logout)
		a.logouted(arg, ctx)
	case *pb.RoleRegist:
		arg := msg.(*pb.RoleRegist)
		glog.Debugf("RoleRegist %#v", arg)
		a.regist(arg, ctx)
	case *pb.RoleLogin:
		arg := msg.(*pb.RoleLogin)
		glog.Debugf("RoleLogin %#v", arg)
		a.loginByPhone(arg, ctx)
	case *pb.TouristLogin:
		arg := msg.(*pb.TouristLogin)
		glog.Debugf("TouristLogin %#v", arg)
		a.loginByTourist(arg, ctx)
	case *pb.WxLogin:
		arg := msg.(*pb.WxLogin)
		glog.Debugf("WxLogin %#v", arg)
		a.loginByWx(arg, ctx)
	case *pb.GetUserData:
		arg := msg.(*pb.GetUserData)
		user := a.getUserById(arg.Userid)
		rsp := handler.GetUserData(user)
		ctx.Respond(rsp)
	case *pb.ApplePay:
		arg := msg.(*pb.ApplePay)
		rsp := handler.AppleVerify(arg)
		ctx.Respond(rsp)
	case *pb.WxpayCallback:
		arg := msg.(*pb.WxpayCallback)
		a.payHandler(arg)
	case *pb.SmscodeRegist:
		arg := msg.(*pb.SmscodeRegist)
		glog.Debugf("SmscodeRegist %#v", arg)
		a.smsbao(arg, ctx)
	case *pb.CResetPwd:
		//重置密码消息
		arg := msg.(*pb.CResetPwd)
		glog.Debugf("CResetPwd %#v", arg)
		a.resetPwd(arg, ctx)
	case *pb.GetNumber:
		//后台请求
		arg := msg.(*pb.GetNumber)
		glog.Debugf("GetNumber %#v", arg)
		rsp := new(pb.GotNumber)
		for k, v := range a.roles {
			if v.GetRobot() {
				rsp.Robot = append(rsp.Robot, k)
			} else {
				rsp.Role = append(rsp.Role, k)
			}
		}
		ctx.Respond(rsp)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

//启动服务
func (a *RoleActor) start(ctx actor.Context) {
	glog.Infof("role start: %v", ctx.Self().String())
	//初始化建立连接
	bind := cfg.Section("hall").Key("bind").Value()
	name := cfg.Section("cookie").Key("name").Value()
	a.hallPid = actor.NewPID(bind, name)
	glog.Infof("a.hallPid: %s", a.hallPid.String())
	connect := &pb.Connect{
		Name: a.Name,
	}
	a.hallPid.Request(connect, ctx.Self())
	//注册更新机器人
	phone := cfg.Section("robot").Key("phone").Value()
	passwd := cfg.Section("robot").Key("passwd").Value()
	//data.RegistRobots(phone, passwd, a.uniqueid)
	//新机器人
	head := cfg.Section("domain").Key("headimag").Value()
	//data.RegistRobots2(head, phone, passwd, a.uniqueid)
	data.RegistRobots3(head, phone, passwd, a.uniqueid)
	//启动
	go a.ticker(ctx)
}

//时钟
func (a *RoleActor) ticker(ctx actor.Context) {
	tick := time.Tick(30 * time.Second)
	msg := new(pb.Tick)
	for {
		select {
		case <-a.stopCh:
			glog.Info("role ticker closed")
			return
		default: //防止阻塞
		}
		select {
		case <-a.stopCh:
			glog.Info("role ticker closed")
			return
		case <-tick:
			ctx.Self().Tell(msg)
		}
	}
}

//钟声
func (a *RoleActor) ding(ctx actor.Context) {
	//glog.Debugf("ding: %v", ctx.Self().String())
	//glog.Debugf("timer: %d", a.timer)
	switch a.timer {
	case 4: //2分钟
		a.saveUser()
		a.timer = 0
		a.smsExpire()
		a.touristIP()
	case 2: //1分钟
		a.smsExpire()
		a.timer += 1
	default:
		a.timer += 1
	}
}

//关闭时钟
func (a *RoleActor) closeTick() {
	select {
	case <-a.stopCh:
		return
	default:
		//停止发送消息
		close(a.stopCh)
	}
}

func (a *RoleActor) handlerStop(ctx actor.Context) {
	glog.Debugf("handlerStop: %s", a.Name)
	//关闭
	a.closeTick()
	//回存数据
	if a.uniqueid != nil {
		a.uniqueid.Save()
	}
	for k, v := range a.offline {
		glog.Debugf("Stop offline: %s", k)
		v.Save()
	}
	for k, v := range a.roles {
		glog.Debugf("Stop online: %s", k)
		v.Save()
	}
}

//定期离线数据清理,移除,存储
func (a *RoleActor) saveUser() {
	glog.Debugf("saveUser caches %#v", a.caches)
	glog.Debugf("saveUser %d, %d", len(a.offline), len(a.roles))
	//离线表
	for k, v := range a.offline {
		//TODO 优化缓存策略
		if a.states[k] {
			v.Save()
			delete(a.states, k)
		}
		glog.Debugf("saveUser offline %s, %d", k, v.GetChip())
		if a.caches[k] < 0 {
			if v.Save() {
				a.delUserMap(v)
				delete(a.caches, k)
				//移除离线表
				delete(a.offline, k)
			} else {
				glog.Errorf("saveUser offline failed %s", k)
			}
		} else {
			a.caches[k] -= 1
		}
	}
	//在线表
	for k, v := range a.roles {
		//TODO 优化缓存策略
		if !a.states[k] {
			continue
		}
		glog.Debugf("saveUser roles %s, %d", k, v.GetChip())
		v.Save()
		delete(a.states, k)
	}
}

//立即更新数据库
func (a *RoleActor) saveUserQuickly(userid string) {
	user := a.getUser(userid)
	if user == nil {
		glog.Errorf("saveUser Quickly failed %s", userid)
		return
	}
	user.Save()
}

//在线表中查找,不存在时离线表中获取
func (a *RoleActor) getUser(userid string) *data.User {
	if user, ok := a.roles[userid]; ok {
		return user
	}
	if user, ok := a.offline[userid]; ok {
		return user
	}
	return nil
}

//在线表中查找,不存在时离线表中获取,不在离线表从数据库中加载
func (a *RoleActor) getUserById(userid string) *data.User {
	user := a.getUser(userid)
	if user != nil {
		return user
	}
	newUser := new(data.User)
	newUser.GetById(userid) //数据库中取
	if newUser.Userid == "" {
		glog.Debugf("getUserById failed %s", userid)
		return nil
	}
	a.loadingUser(newUser)
	return newUser
}

//在线表中查找
func (a *RoleActor) getUserByTourist(account string) *data.User {
	if v, ok := a.players[account]; ok {
		return a.getUserById(v)
	}
	user := new(data.User)
	user.Tourist = account
	user.GetByTourist() //数据库中取
	if user.Userid == "" {
		glog.Debugf("getUserByTourist failed %s", account)
		return nil
	}
	a.loadingUser(user)
	return user
}

//在线表中查找
func (a *RoleActor) getUserByPhone(account string) *data.User {
	if v, ok := a.players[account]; ok {
		return a.getUserById(v)
	}
	user := new(data.User)
	user.Phone = account
	user.GetByPhone() //数据库中取
	if user.Userid == "" {
		glog.Debugf("getUserByPhone failed %s", account)
		return nil
	}
	a.loadingUser(user)
	return user
}

//在线表中查找
func (a *RoleActor) getUserByWx(account string) *data.User {
	if v, ok := a.players[account]; ok {
		return a.getUserById(v)
	}
	user := new(data.User)
	user.Wxuid = account
	user.GetByWechat() //数据库中取
	if user.GetUserid() == "" {
		glog.Debugf("getUserByWx failed %s", account)
		return nil
	}
	a.loadingUser(user)
	return user
}

//加载
func (a *RoleActor) loadingUser(user *data.User) {
	a.offline[user.GetUserid()] = user
	//映射
	a.setUserMap(user)
	a.caches[user.GetUserid()] += 1
}

//添加映射
func (a *RoleActor) setUserMap(user *data.User) {
	if user.GetWxuid() != "" {
		a.players[user.GetWxuid()] = user.GetUserid()
		glog.Debugf("setUserMap %s = %s", user.GetWxuid(), user.GetUserid())
	} else if user.GetPhone() != "" {
		a.players[user.GetPhone()] = user.GetUserid()
		glog.Debugf("setUserMap %s = %s", user.GetPhone(), user.GetUserid())
	} else if user.GetTourist() != "" {
		a.players[user.GetTourist()] = user.GetUserid()
		glog.Debugf("setUserMap %s = %s", user.GetTourist(), user.GetUserid())
	} else {
		glog.Errorf("user mapping err %s", user.GetUserid())
	}
}

//移除映射
func (a *RoleActor) delUserMap(user *data.User) {
	if user.GetWxuid() != "" {
		delete(a.players, user.GetWxuid())
		glog.Debugf("delUserMap %s", user.GetWxuid())
	} else if user.GetPhone() != "" {
		delete(a.players, user.GetPhone())
		glog.Debugf("delUserMap %s", user.GetPhone())
	} else if user.GetTourist() != "" {
		delete(a.players, user.GetTourist())
		glog.Debugf("delUserMap %s", user.GetTourist())
	} else {
		glog.Errorf("user mapping err %s", user.GetUserid())
	}
}

//生成一个验证码,唯一
func (a *RoleActor) GenCode() (s string) {
	s = utils.RandStr(6)
	//是否已经存在
	if _, ok := a.smscode[s]; ok {
		return a.GenCode() //重复尝试,TODO:一定次数后放弃尝试
	}
	return
}

//在线同步数据
func (a *RoleActor) syncUser(arg *pb.SyncUser, ctx actor.Context) {
	glog.Debugf("SyncUser %#v", arg.Userid)
	user := a.getUserById(arg.Userid)
	if user == nil {
		glog.Errorf("syncUser user err %s", arg.Userid)
		return
	}
	err := json.Unmarshal(arg.Data, user)
	if err != nil {
		glog.Errorf("userid %s Unmarshal err %v", arg.Userid, err)
		return
	}
	glog.Debugf("sync user successful %s", arg.Userid)
	a.states[arg.Userid] = true
}

//货币变更
func (a *RoleActor) syncCurrency(diamond, coin, card, chip int64,
	ltype int32, userid string) {
	//日志记录
	user := a.getUserById(userid)
	if user == nil {
		glog.Errorf("syncCurrency err userid %s, type %d, chip %d",
			userid, ltype, chip)
		return
	}
	if chip < 0 && ((chip + user.GetChip()) < 0) {
		chip = 0 - user.GetChip()
	}
	if diamond < 0 && ((diamond + user.GetDiamond()) < 0) {
		diamond = 0 - user.GetDiamond()
	}
	if coin < 0 && ((coin + user.GetCoin()) < 0) {
		coin = 0 - user.GetCoin()
	}
	if card < 0 && ((card + user.GetCard()) < 0) {
		card = 0 - user.GetCard()
	}
	//更新操作
	user.AddCurrency(diamond, coin, card, chip)
	//更新状态
	a.states[userid] = true
	//暂时实时写入
	user.UpdateCurrency()
	//TODO 机器人不写日志
	//if user.GetRobot() {
	//	return
	//}
	//日志记录
	if diamond != 0 {
		msg1 := handler.LogDiamondMsg(diamond, ltype, user)
		nodePid.Tell(msg1)
	}
	if coin != 0 {
		msg1 := handler.LogCoinMsg(coin, ltype, user)
		nodePid.Tell(msg1)
	}
	if card != 0 {
		msg1 := handler.LogCardMsg(card, ltype, user)
		nodePid.Tell(msg1)
	}
	if chip != 0 {
		msg1 := handler.LogChipMsg(chip, ltype, user)
		nodePid.Tell(msg1)
	}
}
