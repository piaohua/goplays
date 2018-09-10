/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2018-01-13 10:56:33
 * Filename      : desk.go
 * Description   : 玩牌逻辑
 * *******************************************************/
package main

import (
	"fmt"
	"math"
	"utils"

	"goplays/data"
	"goplays/game/handler"
	"goplays/glog"
	"goplays/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//// external function

//新建一张牌桌
func NewDesk(deskData *data.DeskData) *Desk {
	desk := &Desk{
		id:      deskData.Rid,
		players: make(map[string]*data.User),
		//-------------
		pids:    make(map[string]*actor.PID),
		router:  make(map[string]string),
		offline: make(map[string]bool),
		stopCh:  make(chan struct{}),
		Name:    deskData.Rid,
	}
	desk.DeskData = deskData
	desk.HuiYinDeskData = new(HuiYinDeskData)
	desk.dealers = make([]map[string]int64, 0)
	desk.seats = make(map[string]uint32)
	desk.Trends = make([]*data.Trend, 0)
	desk.Winers = make([]*data.Winer, 0)
	//初始化
	desk.gameInit()
	return desk
}

//初始化
func (t *Desk) gameInit() {
	//重置
	t.setLast()
	//
	t.betNum = 0
	t.roleBets = make(map[string]int64)
	t.seatBets = make(map[uint32]int64)
	t.handCards = make(map[uint32][]uint32) //手牌
	t.power = make(map[uint32]uint32)
	t.seatPower = make([]uint32, 0)
	t.multiple = make(map[uint32]int64) //倍数

	t.seatRoleBets = make(map[uint32]map[string]int64)   //玩家下注结果seat:userid: value
	t.seatRoleCost = make(map[uint32]map[string]int64)   //个人本金结果seat:userid: value
	t.seatRoleWins = make(map[uint32]map[string]int64)   //个人输赢结果seat:userid: value
	t.seatRoleRefund = make(map[uint32]map[string]int64) //退款seat:userid: value

	t.dealerSeatWins = make(map[uint32]int64) //庄家在每个位置上的输赢
	t.roleWins = make(map[string]int64)       //每个玩家总赢利所得userid: value
	t.roleLose = make(map[string]int64)       //每个玩家输掉金额userid: value
	t.roleProfits = make(map[string]int64)    //每个玩家总输赢userid: value
	t.roleCost = make(map[string]int64)       //成本返还seat: value
	t.roleRefund = make(map[string]int64)     //玩家退款userid: value
	t.seatRefund = make(map[uint32]int64)     //位置退款总额seat: value
	t.seatWins = make(map[uint32]int64)       //位置(1-5)输赢总量

	t.result = make([]data.TrendResult, 0)
	//抽佣数量
	t.HuiYinDeskData.feeNum = 0
	t.HuiYinDeskData.robotNum = 0
	t.HuiYinDeskData.playerNum = 0
	t.fees = make(map[string]int64)
	t.seatFees = make(map[uint32]int64)
	//玩家抽佣明细seat:userid: value
	t.seatRoleFees = make(map[uint32]map[string]int64)
	//玩家下注结果(抽佣后)seat:userid: value
	t.seatRoleFeeBets = make(map[uint32]map[string]int64)

}

//重置
func (t *Desk) setLast() {
	if t.HuiYinDeskData.expect == "" {
		return
	}
	t.HuiYinDeskData.lastexpect = t.HuiYinDeskData.expect
	t.HuiYinDeskData.lastopencode = t.HuiYinDeskData.opencode
	t.HuiYinDeskData.lastopentime = t.HuiYinDeskData.opentime
	t.HuiYinDeskData.lastopentimestamp = t.HuiYinDeskData.opentimestamp
	t.HuiYinDeskData.lastpower = make([]uint32, 0)
	//
	var i uint32
	for i = 1; i <= data.SEAT5; i++ {
		t.lastpower = append(t.lastpower, t.power[i])
	}
	//
	t.HuiYinDeskData.expect = ""
	t.HuiYinDeskData.opencode = ""
	t.HuiYinDeskData.opentime = ""
	t.HuiYinDeskData.opentimestamp = 0
}

// 打印
func (t *Desk) printOver() {
	glog.Debugf("game over data -> %#v", t.DeskData)
	glog.Debugf("game over players -> %d", len(t.players))
	glog.Debugf("game over seats -> %#v", t.seats)
	glog.Debugf("game over router -> %#v", t.router)
	glog.Debugf("game over dealer -> %s", t.dealer)
	glog.Debugf("game over dealerSeat -> %d", t.dealerSeat)
	glog.Debugf("game over dealerNum -> %d", t.dealerNum)
	glog.Debugf("game over carry -> %d", t.carry)
	glog.Debugf("game over betNum -> %d", t.betNum)
	glog.Debugf("game over dealers -> %#v", t.dealers)
	glog.Debugf("game over roleBets -> %#v", t.roleBets)
	glog.Debugf("game over seatBets -> %#v", t.seatBets)
	glog.Debugf("game over handCards -> %+x", t.handCards)
	glog.Debugf("game over power -> %#v", t.power)
	glog.Debugf("game over seatPower -> %#v", t.seatPower)
	glog.Debugf("game over multiple -> %#v", t.multiple)
	glog.Debugf("game over seatRoleBets -> %#v", t.seatRoleBets)
	glog.Debugf("game over seatRoleCost -> %#v", t.seatRoleCost)
	glog.Debugf("game over seatRoleWins -> %#v", t.seatRoleWins)
	glog.Debugf("game over seatRoleRefund -> %#v", t.seatRoleRefund)
	glog.Debugf("game over dealerSeatWins -> %#v", t.dealerSeatWins)
	glog.Debugf("game over seatWins -> %#v", t.seatWins)
	glog.Debugf("game over roleWins -> %#v", t.roleWins)
	glog.Debugf("game over roleLose -> %#v", t.roleLose)
	glog.Debugf("game over roleProfits -> %#v", t.roleProfits)
	glog.Debugf("game over roleCost -> %#v", t.roleCost)
	glog.Debugf("game over roleRefund -> %#v", t.roleRefund)
	glog.Debugf("game over seatRefund -> %#v", t.seatRefund)
	glog.Debugf("game over result -> %#v", t.result)
}

//房间消息广播
func (t *Desk) broadcast(msg interface{}) {
	//glog.Debugf("enter desk %#v, %#v", t.router, t.pids)
	for _, p := range t.pids {
		if p == nil {
			continue
		}
		//glog.Debugf("broadcast %s, msg %#v", p.String(), msg)
		p.Tell(msg)
	}
}

//房间消息广播(除seat外)
func (t *Desk) broadcast_(userid string, msg interface{}) {
	for i, p := range t.pids {
		if p == nil {
			continue
		}
		if i != userid {
			//glog.Debugf("broadcast %s, msg %#v", p.String(), msg)
			p.Tell(msg)
		}
	}
}

//获取
func (t *Desk) getPid(userid string) *actor.PID {
	if v, ok := t.pids[userid]; ok && v != nil {
		return v
	}
	//panic(fmt.Sprintf("getPlayer error:%s", userid))
	return nil
}

//获取
func (t *Desk) getPlayer(userid string) *data.User {
	if v, ok := t.players[userid]; ok && v != nil {
		return v
	}
	//panic(fmt.Sprintf("getPlayer error:%s", userid))
	return nil
}

//获取手牌
func (t *Desk) getHandCards(seat uint32) []uint32 {
	if v, ok := t.handCards[seat]; ok && v != nil {
		return v
	}
	//return []uint32{}
	panic(fmt.Sprintf("getHandCards error:%d", seat))
}

//收益
func (a *Desk) sendChip(userid string, num int64, rtype int32) {
	if num == 0 {
		return
	}
	p := a.getPlayer(userid)
	if p == nil {
		glog.Errorf("sendChip err userid %s num %d, rtype %d", userid, num, rtype)
		a.syncChip2(userid, num, rtype)
		return
	}
	p.AddChip(num)
	a.syncChip(userid, num, rtype)
}

//同步收益
func (a *Desk) syncChip(userid string, chip int64, ltype int32) {
	//货币变更及时同步
	msg2 := &pb.ChangeCurrency{
		Userid: userid,
		Chip:   chip,
		Type:   ltype,
	}
	if v, ok := a.pids[userid]; ok && v != nil {
		v.Tell(msg2)
	} else {
		a.syncChip2(userid, chip, ltype)
	}
}

//离线同步数据
func (a *Desk) syncChip2(userid string, chip int64, ltype int32) {
	glog.Infof("syncChip2 userid %s, chip %d, ltype %d", userid, chip, ltype)
	msg3 := &pb.OfflineCurrency{
		Userid: userid,
		Chip:   chip,
		Type:   ltype,
	}
	//通过大厅通知其它节点
	nodePid.Tell(msg3)
}

//离开位置
func (t *Desk) leaveSeat(userid string, st bool, p *data.User) {
	if seat, ok := t.seats[userid]; ok {
		delete(t.seats, userid)
		//广播消息
		msg := handler.SitMsg(userid, seat, st, p) //离开坐位
		t.broadcast(msg)
	}
}

//' 庄家操作

//0下庄 1上庄 2补庄
func (t *Desk) addBeDealer(userid string,
	st uint32, num int64, p *data.User) {
	if userid == t.dealer && st == data.DEALER_BU { //庄家补庄
		t.carry += num
	} else {
		//if t.state != data.STATE_BET && t.dealer == "" &&
		//	st == data.DEALER_BU { //系统做庄
		//	//成为庄家
		//	t.dealer = userid
		//	t.carry = num
		//	//t.dealerNum = 0
		//	t.setDealerSeat()
		//	//t.leaveSeat(userid, false, p)
		//} else if st == data.DEALER_UP {
		//	m := map[string]int64{userid: num}
		//	t.dealers = append(t.dealers, m)
		//}
		if st == data.DEALER_BU {
			var has bool
			for i, m := range t.dealers {
				if _, ok := m[userid]; ok {
					m[userid] += num
					t.dealers[i] = m
					has = true
					break
				}
			}
			if !has {
				m := map[string]int64{userid: num}
				t.dealers = append(t.dealers, m)
			}
		} else {
			m := map[string]int64{userid: num}
			t.dealers = append(t.dealers, m)
		}
	}
	//TODO 暂时不扣除，真实上庄才扣除
	//t.sendChip(userid, (-1 * num), data.LogType7)
	msg := handler.BeDealerMsg(st, num, t.dealer, userid, p.GetNickname(), t.dealerDown)
	t.broadcast(msg)
}

//下庄
func (t *Desk) delBeDealer(userid string, p *data.User) {
	var num int64
	if t.dealer == userid && t.dealer != "" {
		num = t.carry
		t.sendChip(userid, num, data.LogType8)
		t.dealer = ""
		t.carry = 0
		t.dealerSeat = 0
		t.dealerDown = false
	}
	//t.dealerNum = 0
	msg := handler.BeDealerMsg(0, num, t.dealer, userid, p.GetNickname(), t.dealerDown)
	t.broadcast(msg)
}

//离开房间返还上庄列表
func (t *Desk) leaveBeDealer(userid string) {
	for {
		had := true
		for i, m := range t.dealers {
			//if num, ok := m[userid]; ok {
			if _, ok := m[userid]; ok {
				//TODO 不返还,因为修改为排队时不扣除
				//t.sendChip(userid, num, data.LogType8)
				delete(m, userid)
				t.dealers = append(t.dealers[:i], t.dealers[i+1:]...)
				had = false
				break
			}
		}
		if had {
			break
		}
	}
}

//成为庄家
func (t *Desk) beComeDealer() {
	if !t.DeskData.Deal { //非上庄房间
		return
	}
	if t.dealer != "" { //已经有庄家
		//每次开始选择一个位置
		t.setDealerSeat()
		return
	}
	i, userid, num := t.findBeDealer2()
	if userid == "" || num < int64(t.DeskData.Carry) {
		//glog.Errorf("beComeDealer failed %s, %d, %d",
		//	userid, num, t.DeskData.Carry)
		return
	}
	p := t.players[userid]
	if p == nil || t.offline[userid] { //掉线
		glog.Errorf("beComeDealer failed %s, %d, %d",
			userid, num, t.DeskData.Carry)
		return
	}
	//上庄成功扣除
	t.sendChip(userid, (-1 * num), data.LogType7)
	//成为庄家
	t.dealer = userid
	t.carry = num
	//t.dealerNum = 0
	t.dealers = append(t.dealers[:i], t.dealers[i+1:]...)
	t.setDealerSeat()
	//庄家随机一个位置
	msg := handler.BeDealerMsg(1, num, t.dealer, userid, p.GetNickname(), t.dealerDown)
	t.broadcast(msg)
}

//选择庄家位置
func (t *Desk) setDealerSeat() {
	switch t.DeskData.Gtype {
	case data.GAME_NIU:
		t.dealerSeat = uint32(utils.RandIntN(5)) + 1
	case data.GAME_SAN:
		t.dealerSeat = uint32(utils.RandIntN(5)) + 1
	case data.GAME_JIU:
		t.dealerSeat = uint32(utils.RandIntN(5)) + 1
	}
	if t.dealerSeat == 0 {
		t.dealerSeat = data.SEAT1
	}
	msg := t.pushBeDealerMsg()
	t.broadcast(msg)
}

//携带最大的优先做庄
func (t *Desk) findBeDealer() (int, string, int64) {
	var index int
	var userid string
	var maxNum int64
	for i, m := range t.dealers {
		for k, v := range m {
			if v > maxNum {
				index = i
				userid = k
				maxNum = v
			}
		}
	}
	return index, userid, maxNum
}

//携带最大的优先做庄
func (t *Desk) findBeDealer2() (int, string, int64) {
	var index int
	var userid string
	var maxNum int64
	for i, m := range t.dealers {
		for k, _ := range m {
			p := t.players[k]
			if p == nil {
				continue
			}
			//TODO 自动下庄金额不足玩家
			//全部资金上庄
			v := p.GetChip()
			if v < int64(t.DeskData.Carry) {
				continue
			}
			if v > maxNum {
				index = i
				userid = k
				maxNum = v
			}
		}
	}
	return index, userid, maxNum
}

//不足做庄,或者检测是否有人上庄
func (t *Desk) checkDealer() {
	glog.Debugf("offline info %#v", t.offline)
	glog.Debugf("dealer info %s, %d, %d, %#v",
		t.dealer, t.carry, t.dealerSeat, t.dealers)
	if !t.DeskData.Deal {
		return
	}
	if t.dealer != "" {
		//离线或者不足
		if t.offline[t.dealer] ||
			t.carry <= int64(t.DeskData.Down) ||
			t.carry >= int64(t.DeskData.Top) ||
			t.dealerDown {
			p := t.players[t.dealer]
			t.delBeDealer(t.dealer, p)
			return
		}
		//t.dealerNum += 1
	}
}

//是否已经是庄家或者已经申请上庄
func (t *Desk) alreadyDealer(userid string) bool {
	if t.dealer == userid {
		return true
	}
	for _, m := range t.dealers {
		if _, ok := m[userid]; ok {
			return true
		}
	}
	return false
}

//.

//玩家离开
func (t *Desk) SitDown(userid string, seat uint32, st bool) pb.ErrCode {
	if p, ok := t.players[userid]; ok {
		if p.GetChip() < int64(t.DeskData.Sit) {
			return pb.SitNotEnough
		}
	} else {
		return pb.NotInRoom
	}
	if _, ok := t.seats[userid]; ok && st { //坐下
		return pb.SitDownFailed
	}
	if _, ok := t.seats[userid]; !ok && !st { //站起
		return pb.StandUpFailed
	}
	//if userid == t.dealer { //庄家不能坐
	//	return pb.DealerSitFailed
	//}
	if st {
		for _, s := range t.seats {
			if s == seat {
				return pb.SitDownFailed
			}
		}
		t.seats[userid] = seat
	} else {
		delete(t.seats, userid)
	}
	//glog.Infof("SitDown -> %s, %d, %v", userid, seat, st)
	//广播消息
	p := t.players[userid]
	msg := handler.SitMsg(userid, seat, st, p)
	t.broadcast(msg)
	return pb.OK
}

//没人上庄时都可以选择上庄,可以多次上庄，已经上庄的人可以补庄
//0下庄 1上庄 2补庄
func (t *Desk) BeDealer(userid string, st, num uint32) pb.ErrCode {
	if !t.DeskData.Deal {
		return pb.NotDealerRoom
	}
	if _, ok := t.players[userid]; !ok {
		return pb.NotInRoom
	}
	p := t.players[userid]
	//TODO 暂时全部带上庄
	num = uint32(p.GetChip())
	//上庄限制
	if num < t.DeskData.Carry && st == data.DEALER_UP {
		return pb.BeDealerNotEnough
	}
	//下注和封盘时不能下庄, 结算状态下庄
	if (t.state == data.STATE_BET || t.state == data.STATE_SEAL) &&
		userid == t.dealer && st == data.DEALER_DOWN {
		//结束后庄家下庄
		t.dealerDown = true
		p := t.players[userid]
		msg := handler.BeDealerMsg(0, int64(num), t.dealer, userid, p.GetNickname(), t.dealerDown)
		t.broadcast(msg)
		//return pb.GameStartedCannotLeave
		return pb.OK
	}
	//已经上庄,暂时不能重复上
	if st == data.DEALER_UP && t.alreadyDealer(userid) {
		return pb.BeDealerAlready
	}
	if st == data.DEALER_UP || st == data.DEALER_BU {
		if p.GetChip() < int64(num) {
			return pb.NotEnoughCoin
		}
		t.addBeDealer(userid, st, int64(num), p)
	} else {
		t.delBeDealer(userid, p)
	}
	return pb.OK
}

//下注
func (t *Desk) ChoiceBet(userid string, seatBet uint32, num int64) pb.ErrCode {
	if userid == t.dealer { //庄家不用下注
		return pb.BetDealerFailed
	}
	if t.state != data.STATE_BET {
		return pb.GameNotStart
	}
	//有庄家时不能下注庄家位置
	if seatBet == t.dealerSeat && t.dealer != "" {
		return pb.BetSeatWrong
	}
	p := t.getPlayer(userid)
	if p == nil {
		return pb.NotInRoom
	}
	if num <= 0 {
		return pb.OperateError
	}
	if p.GetChip() < num {
		return pb.NotEnoughCoin
	}
	//TODO 限制优化
	//下注不能大于庄家携带1/4
	//if t.dealer != "" && ((t.betNum + num) > (t.carry / 4)) {
	//	return pb.BetTopLimit //下注限制
	//}
	//chip := p.GetChip()             //剩余金额
	//betsCount := t.roleBets[userid] //已经下注额
	//本轮下注不能超过1/4
	//if (num + betsCount) > ((chip + betsCount) / 4) {
	//	return pb.BetTopLimit //下注限制
	//}
	if t.betCheck(num) {
		return pb.BetTopLimit //下注限制
	}
	t.roleBets[userid] += num  //个人总下注额
	t.seatBets[seatBet] += num //当前位置总下注额
	//位置详细记录
	var seatNum int64
	if m, ok := t.seatRoleBets[seatBet]; ok {
		m[userid] += num
		seatNum = m[userid]
		t.seatRoleBets[seatBet] = m
	} else {
		m := make(map[string]int64)
		m[userid] += num
		seatNum = m[userid]
		t.seatRoleBets[seatBet] = m
	}
	t.betNum += num //当局总下注额
	t.sendChip(userid, (-1 * num), data.LogType5)
	msg := handler.RoomBetMsg(t.seats[userid], seatBet, uint32(num),
		t.seatBets[seatBet], seatNum, userid)
	t.broadcast(msg)
	return pb.OK
}

//下注限制,有庄时只对庄家抽佣
func (t *Desk) betCheck(num int64) bool {
	if t.DeskData.Rtype != data.ROOM_TYPE1 { //免佣
		return false
	}
	if t.dealer == "" { //无庄
		return false
	}
	v := t.betNum + num
	n := int64(math.Ceil((float64(t.DeskData.Cost) / float64(100)) * float64(v)))
	if t.carry > (v + n) {
		return false
	}
	return true
}

//进入限制
func (t *Desk) enterCheck(p *data.User) pb.ErrCode {
	//if p.GetVip() < t.DeskData.Vip {
	//	return pb.VipTooLow
	//}
	if p.GetChip() < int64(t.DeskData.Chip) {
		return pb.ChipNotEnough
	}
	if _, ok := t.players[p.GetUserid()]; !ok {
		if uint32(len(t.players)) >= t.DeskData.Count {
			return pb.RoomFull //人数已满
		}
	}
	return pb.OK
}

//进入
func (t *Desk) Enter(p *data.User) pb.ErrCode {
	errcode := t.enterCheck(p)
	if errcode != pb.OK {
		return errcode
	}
	//TODO 已经在房间防止数据覆盖
	//if _, ok := t.players[p.GetUserid()]; !ok {
	//	t.players[p.GetUserid()] = p
	//}
	t.players[p.GetUserid()] = p
	delete(t.offline, p.GetUserid())
	return pb.OK
}

//玩家离开,下注也可以离开
func (t *Desk) Leave(userid string) pb.ErrCode {
	if _, ok := t.players[userid]; !ok {
		return pb.NotInRoom
	}
	//庄家下注时不能离开
	//if t.state != data.STATE_OVER && userid == t.dealer {
	//	return pb.GameStartedCannotLeave
	//}
	p := t.players[userid]
	if t.dealer == userid && userid != "" && t.state == data.STATE_OVER {
		t.delBeDealer(userid, p)
	}
	t.leaveBeDealer(userid) //清空上庄列表
	//广播消息
	msg := handler.LeaveMsg(userid, t.seats[userid])
	//位置玩家离开
	if _, ok := t.seats[userid]; ok {
		t.broadcast(msg) //如果是位置玩家广播离开
	} else if v, ok := t.pids[userid]; ok {
		v.Tell(msg) //旁观玩家不广播
	}
	//清除数据
	if v, ok := t.pids[userid]; ok {
		delete(t.router, v.String())
		delete(t.pids, userid)
		//FIXME 为什么会出现key不同没有删除情况?
		for k, val := range t.router {
			if val == userid {
				delete(t.router, k)
				break
			}
		}
	}
	delete(t.seats, userid)
	delete(t.offline, userid)
	//有下注或者是庄家,且游戏没有结束,标注为离线状态
	if (t.roleBets[userid] > 0 || t.dealer == userid) && t.state != data.STATE_OVER {
		t.offline[userid] = true
	} else {
		delete(t.players, userid)
	}
	return pb.OK
}

//进入房间响应消息
func (t *Desk) roleList(ctx actor.Context) {
	rsp := new(pb.SHuiYinRoomRoles)
	//旁观玩家
	rsp.List = t.getRoomUsers(false)
	ctx.Respond(rsp)
}

//进入房间响应消息
func (t *Desk) enterMsg(userid string) interface{} {
	stoc := new(pb.SHuiYinEnterRoom)
	//房间数据
	roomInfo := handler.PackRoomInfo(t.DeskData)
	timer := t.nexttime - utils.Timestamp()
	if timer < 0 {
		timer = 0
	}
	glog.Debugf("enter msg %d, %d", t.nexttime, timer)
	roomInfo.Timer = uint32(timer)
	roomInfo.State = t.state
	roomInfo.Num = uint32(len(t.players))
	//roomInfo.Expect = t.HuiYinDeskData.expect
	//roomInfo.Opencode = t.HuiYinDeskData.opencode
	roomInfo.Expect = t.HuiYinDeskData.lastexpect
	roomInfo.Opencode = t.HuiYinDeskData.lastopencode
	roomInfo.Points = t.HuiYinDeskData.lastpower
	stoc.Roominfo = roomInfo
	//位置下注
	for k, v := range t.seatBets {
		betsinfo := &pb.HuiYinRoomBets{
			Seat: k,
			Bets: v,
		}
		stoc.Seatbets = append(stoc.Seatbets, betsinfo)
	}
	//玩家下注
	for k, m := range t.seatRoleBets {
		if v, ok := m[userid]; ok {
			betsinfo := &pb.HuiYinRoomBets{
				Seat: k,
				Bets: v,
			}
			stoc.Rolebets = append(stoc.Rolebets, betsinfo)
		}
	}
	//坐下玩家
	stoc.Userinfo = t.getRoomUsers(true)
	//庄信息
	stoc.Dealerinfo = t.pushDealerMsg()
	return stoc
}

//玩家列表,isSeat 是否坐下
func (t *Desk) getRoomUsers(isSeat bool) (l []*pb.RoomUser) {
	l = make([]*pb.RoomUser, 0)
	for k, v := range t.players {
		if isSeat && t.seats[k] != 0 {
			msg2 := new(pb.RoomUser)
			msg2.Seat = t.seats[k] //为坐下玩家
			msg2.Data = handler.PackUserData(v)
			l = append(l, msg2)
		} else if !isSeat && t.seats[k] == 0 {
			msg2 := new(pb.RoomUser)
			msg2.Seat = t.seats[k] //为坐下玩家
			msg2.Data = handler.PackUserData(v)
			l = append(l, msg2)
		}
	}
	return
}

//庄位置信息
func (t *Desk) pushBeDealerMsg() interface{} {
	msg1 := new(pb.SHuiYinPushBeDealer)
	msg1.Dealer = t.dealer
	msg1.Carry = t.carry
	msg1.Seat = t.dealerSeat
	msg1.List = t.dealerListMsg()
	if p, ok := t.players[t.dealer]; ok && p != nil {
		msg1.Nickname = p.GetNickname()
	}
	return msg1
}

//庄位置信息
func (t *Desk) pushDealerMsg() *pb.SHuiYinPushDealer {
	msg1 := new(pb.SHuiYinPushDealer)
	msg1.Dealer = t.dealer
	msg1.Carry = t.carry
	msg1.Seat = t.dealerSeat
	msg1.List = t.dealerListMsg()
	if p, ok := t.players[t.dealer]; ok && p != nil {
		msg1.Nickname = p.GetNickname()
	}
	msg1.Down = t.dealerDown
	return msg1
}

//上庄列表消息
func (t *Desk) dealerListMsg() []*pb.HuiYinDealerList {
	l := make([]*pb.HuiYinDealerList, 0)
	for _, m := range t.dealers {
		for k, v := range m {
			p := t.getPlayer(k)
			if p == nil {
				continue
			}
			list := &pb.HuiYinDealerList{
				Userid:   k,
				Nickname: p.GetNickname(),
				Photo:    p.GetPhoto(),
				Chip:     v,
			}
			l = append(l, list)
		}
	}
	return l
}

// 位置下注信息
func (t *Desk) seatBetInfo(seat uint32) (rsp *pb.SHuiYinDeskBetInfo) {
	rsp = new(pb.SHuiYinDeskBetInfo)
	rsp.Seat = seat
	rsp.Bets = t.seatBets[seat]
	if m, ok := t.seatRoleBets[seat]; ok {
		for k, v := range m {
			p := t.getPlayer(k)
			if p == nil {
				continue
			}
			val := &pb.BetInfo{
				Userid:   k,
				Nickname: p.GetNickname(),
				Photo:    p.GetPhoto(),
				Bets:     v,
			}
			rsp.List = append(rsp.List, val)
		}
	}
	return
}

// vim: set foldmethod=marker foldmarker=//',//.:
