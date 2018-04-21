/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2018-01-22 17:06:12
 * Filename      : recv.go
 * Description   : 机器人
 * *******************************************************/
package main

import (
	"goplays/data"
	"goplays/glog"
	"goplays/pb"
	"utils"
)

func (r *Robot) receive(msg interface{}) {
	switch msg.(type) {
	case *pb.SRegist:
		r.recvRegist(msg.(*pb.SRegist))
	case *pb.SLogin:
		r.recvLogin(msg.(*pb.SLogin))
	case *pb.SUserData:
		r.recvdata(msg.(*pb.SUserData))
	case *pb.SPushCurrency:
		r.recvPushCurrency(msg.(*pb.SPushCurrency))
	case *pb.SHuiYinGames:
		r.recvGames(msg.(*pb.SHuiYinGames))
	case *pb.SHuiYinRoomList:
		r.recvRoomList(msg.(*pb.SHuiYinRoomList))
	case *pb.SHuiYinLeave:
		r.recvLeave(msg.(*pb.SHuiYinLeave))
	case *pb.SHuiYinEnterRoom:
		r.recvComein(msg.(*pb.SHuiYinEnterRoom))
	case *pb.SHuiYinCamein:
		r.recvCamein(msg.(*pb.SHuiYinCamein))
	case *pb.SHuiYinSit:
		r.recvSitDown(msg.(*pb.SHuiYinSit))
	case *pb.SHuiYinRoomBet:
		r.recvBet(msg.(*pb.SHuiYinRoomBet))
	case *pb.SHuiYinDeskState:
		r.recvGamestate(msg.(*pb.SHuiYinDeskState))
	case *pb.SHuiYinGameover:
		r.recvGameover(msg.(*pb.SHuiYinGameover))
	case *pb.SPing:
		r.recvPing(msg.(*pb.SPing))
		//glog.Debugf("pong : %#v", msg)
	case *pb.SHuiYinPushDealer:
		r.recvDealer(msg.(*pb.SHuiYinPushDealer))
	default:
		glog.Errorf("unknow message: %#v", msg)
	}
}

//' 接收到服务器登录返回
func (r *Robot) recvRegist(stoc *pb.SRegist) {
	var errcode = stoc.GetError()
	switch errcode {
	case pb.OK:
		Logined(r.data.Phone, r.ltype) //登录成功
		r.regist = true                //注册成功
		r.data.Userid = stoc.GetUserid()
		glog.Infof("regist successful -> %s", r.data.Userid)
		r.SendUserData() // 获取玩家数据
		return
	case pb.PhoneRegisted:
		glog.Infof("phone registed -> %s", r.data.Phone)
		r.SendLogin() //尝试登录
		return
	default:
		glog.Infof("regist err -> %d", errcode)
	}
	//重新尝试登录
	//go ReLogined(r.roomid, r.data.Phone, r.code, r.rtype, r.envBet)
	r.Close()
}

//.

//' 接收到服务器登录返回
func (r *Robot) recvLogin(stoc *pb.SLogin) {
	var errcode = stoc.GetError()
	switch errcode {
	case pb.OK:
		Logined(r.data.Phone, r.ltype) //登录成功
		r.data.Userid = stoc.GetUserid()
		glog.Infof("login successful -> %s", r.data.Userid)
		r.SendUserData() // 获取玩家数据
		return
	default:
		glog.Infof("login err -> %d", errcode)
	}
	r.Close()
}

//.

//' 接收到玩家数据
func (r *Robot) recvdata(stoc *pb.SUserData) {
	var errcode = stoc.GetError()
	if errcode != pb.OK {
		glog.Infof("get data err -> %d", errcode)
		r.Close() //断开
		return
	}
	userdata := stoc.GetData()
	// 设置数据
	r.data.Userid = userdata.GetUserid()     // 用户id
	r.data.Nickname = userdata.GetNickname() // 用户昵称
	r.data.Sex = userdata.GetSex()           // 用户性别,男1 女2 非男非女3
	r.data.Coin = userdata.GetCoin()         // 金币
	r.data.Diamond = userdata.GetDiamond()   // 钻石
	r.data.Chip = userdata.GetChip()
	r.data.Card = userdata.GetCard()
	//chip 单位为分
	if r.data.Chip < 650000 {
		//TODO 自动充值
		r.AddCurrency()
		r.SendStandup()
		return
	}
	//获取游戏列表
	r.SendGames()
	//获取房间列表
	r.SendRoomList()
}

//更新金币
func (r *Robot) recvPushCurrency(stoc *pb.SPushCurrency) {
	currencyData := stoc.GetData()
	r.data.Coin += currencyData.GetCoin()
	r.data.Card += currencyData.GetCard()
	r.data.Chip += currencyData.GetChip()
	r.data.Diamond += currencyData.GetDiamond()
	if r.data.Chip < 650000 {
		//TODO 自动充值
		r.SendStandup()
	}
}

//游戏列表
func (r *Robot) recvGames(stoc *pb.SHuiYinGames) {
	//list := stoc.GetList()
	//glog.Debugf("game list %#v, len %d", list, len(list))
}

//房间列表
func (r *Robot) recvRoomList(stoc *pb.SHuiYinRoomList) {
	list := stoc.GetList()
	//glog.Debugf("room list %#v, len %d", list, len(list))
	for _, v := range list {
		roomid := v.GetInfo().Roomid
		RegistRoom(roomid, r.ltype)
	}
	rbet.SetRoom(list)
	if r.roomid != "" {
		r.SendEntryRoom(r.roomid)
	} else {
		r.Close() //下线
		//for _, v := range list {
		//	if v.GetInfo().Roomid != "" &&
		//		r.data.Chip > int64(v.GetInfo().Chip) &&
		//		r.data.Vip >= v.GetInfo().Vip {
		//		glog.Debugf("room id %s", v.GetInfo().Roomid)
		//		r.SendEntryRoom(v.GetInfo().Roomid)
		//		break
		//	}
		//}
	}
}

//游戏
func (r *Robot) recvPing(stoc *pb.SPing) {
	//TODO 暂时用这个协议控制机器人下注,添加新协议替换
	if stoc.GetTime() == 10 {
		glog.Debugf("ping %s", r.data.Userid)
		//设置100%全部下完
		rbet.SetPercent(100)
		if r.bits > 0 {
			//r.SendRoomBet()
			r.SendRoomBet4()
		}
	}
}

//.

//' 离开房间
func (r *Robot) recvLeave(stoc *pb.SHuiYinLeave) {
	if stoc.GetUserid() == r.data.Userid {
		r.Close() //下线
	}
}

//.

//' 进入房间
func (r *Robot) recvComein(stoc *pb.SHuiYinEnterRoom) {
	var errcode = stoc.GetError()
	switch errcode {
	case pb.OK:
		roominfo := stoc.GetRoominfo()
		r.gtype = roominfo.Info.Gtype
		r.rtype = roominfo.Info.Rtype
		r.roomid = roominfo.Info.Roomid
		userinfo := stoc.GetUserinfo()
		for _, v := range userinfo {
			//只返回坐下玩家
			if v.Data.Userid == r.data.Userid {
				glog.Debugf("comein user info -> %s", v.Data.Userid)
				r.seat = v.Seat
				break
			}
		}
		glog.Debugf("comein desk info -> %#v", roominfo)
		//坐下
		r.SendSitDown()
		//下注
		glog.Debugf("comein desk state -> %d", roominfo.State)
		switch roominfo.State {
		case data.STATE_BET:
			//r.SendRoomBet()
			go func() {
				//延迟下注
				utils.Sleep(5)
				r.SendRoomBet4()
			}()
		}
	default:
		glog.Infof("comein err -> %d", errcode)
		r.Close() //进入出错,关闭
	}
}

//进入房间
func (r *Robot) recvCamein(stoc *pb.SHuiYinCamein) {
	if stoc.GetUserdata().GetUserid() == r.data.Userid {
		EnterRoom(r.data.Phone, r.roomid, r.ltype)
	}
}

//.

//' huiyin

//坐下
func (r *Robot) recvSitDown(stoc *pb.SHuiYinSit) {
	var errcode = stoc.GetError()
	var seat uint32 = stoc.GetSeat()
	var userid string = stoc.GetUserid()
	if userid != r.data.Userid {
		return
	}
	switch errcode {
	case pb.OK:
		r.seat = seat //坐下位置
	default:
		if r.sits > 4 { //尝试次数过多
			r.SendStandup()
		} else {
			r.SendSitDown()
		}
	}
}

//下注
func (r *Robot) recvBet(stoc *pb.SHuiYinRoomBet) {
	var errcode = stoc.GetError()
	var userid string = stoc.GetUserid()
	if userid != r.data.Userid {
		//TODO 有人下注时才下注
		return
	}
	glog.Debugf("bet userid %s, errcode %d", userid, errcode)
	switch errcode {
	case pb.OK:
		if userid == r.data.Userid {
			//下注成功检测限制
			if rbet.SetBet(r.roomid, int64(stoc.GetValue())) {
				//达到限制条件停止下注
				return
			}
		}
		//if r.bits > 0 && r.bitNum > 0 {
		if r.bits > 0 {
			//r.SendRoomBet()
			r.SendRoomBet4()
		}
	default:
		r.SendStandup()
	}
}

//状态更新
func (r *Robot) recvGamestate(stoc *pb.SHuiYinDeskState) {
	//glog.Debugf("game state : %#v", stoc)
	var state uint32 = stoc.GetState()
	switch state {
	case data.STATE_READY:
		r.SendStandup()
	case data.STATE_BET:
		//随机下注次数
		//r.bits = uint32(utils.RandInt32N(20) + 1)
		//r.bitNum = uint32(utils.RandInt32N(7) * 500)
		//r.SendRoomBet() //下注
		if !r.setBits() {
			return
		}
		r.bitNum = 0
		rbet.SetState()
		//r.SendRoomBet4() //下注
		go func() {
			//延迟下注
			utils.Sleep(8)
			r.SendRoomBet4()
		}()
	case data.STATE_SEAL:
	case data.STATE_OVER:
		rbet.Reset()  //重置
		r.betSeat = 0 //重置
	default:
		r.SendStandup()
	}
}

//庄家位置
func (r *Robot) recvDealer(stoc *pb.SHuiYinPushDealer) {
	r.dealerSeat = stoc.GetSeat()
	r.setBetSeat()
}

//下注位置
func (r *Robot) setBetSeat() {
	//TODO 选择多个自己下注位置
	var a1 = []uint32{1, 2, 3, 4, 5}
	if r.dealerSeat != 0 {
		for k, v := range a1 {
			//庄家位置过滤掉
			if v == r.dealerSeat {
				a1 = append(a1[:k], a1[k+1:]...)
				break
			}
		}
	}
	r.betSeat = a1[utils.RandIntN(len(a1))] //随机
}

//结束
func (r *Robot) recvGameover(stoc *pb.SHuiYinGameover) {
	r.round++
	if r.round >= 5 { //打10局下线
		r.SendStandup()
		return
	}
	switch r.ltype {
	case 1:
		if !pkBetTime() {
			r.SendStandup()
		}
	case 2:
		if !ftBetTime() {
			r.SendStandup()
		}
	default:
		r.SendStandup()
	}
}

//.

// vim: set foldmethod=marker foldmarker=//',//.:
