package main

import (
	"math"
	"sort"

	"goplays/data"
	"goplays/game/algo"
	"goplays/game/handler"
	"goplays/glog"
	"goplays/pb"
	"utils"
)

//开始游戏
func (t *Desk) gameStart() {
	t.gameInit()     //重置
	t.beComeDealer() //成为庄家
}

//封盘
func (t *Desk) gameSeal() {
}

/*
//jiu
第一门牌由1.2名开出的号码组成，
第二门由3.4名，
第三门由5.6名，
第四门由7.8名，
第五门由9.10名组成，

//niu
前区由1,2,3,4,5名
前区由6,7,8,9,10名

//san
第一门牌由1.2.3名开出的号码组成，
第二门由3.4.5名，
第三门由5.6.7名，
第四门由7.8.9名，
第五门由9.10.1名组成，

//2018-03-21 17:05:04 做如下修改:
牛牛号码组合：举例1、2、3、4、5、6、7、8、9、10
组合一：1-2-3-4-5，组合二：2-3-4-5-6，组合三：3-4-5-6-7，组合四：4-5-6-7-8，组合五：5-6-7-8-9，最后一位不要算，去掉。

*/

//' 发牌
func (t *Desk) dealCard() bool {
	if t.state != data.STATE_OVER {
		return false
	}
	codes := utils.Split(t.opencode, ",")
	if len(codes) != 10 {
		glog.Errorf("deal card err %s, %s", t.expect, t.opencode)
		return false
	}
	cards := make([]uint32, len(codes))
	for k, v := range codes {
		i, err1 := utils.Str2Int(v)
		if err1 != nil {
			glog.Errorf("codes err %v, codes %v", err1, codes)
			return false
		}
		cards[k] = uint32(i)
	}
	switch t.DeskData.Gtype {
	case data.GAME_NIU:
		//t.handCards[data.SEAT1] = cards[:5]
		//t.handCards[data.SEAT2] = cards[5:]
		t.handCards[data.SEAT1] = cards[:5]
		t.handCards[data.SEAT2] = cards[1:6]
		t.handCards[data.SEAT3] = cards[2:7]
		t.handCards[data.SEAT4] = cards[3:8]
		t.handCards[data.SEAT5] = cards[4:9]
		return true
	case data.GAME_SAN:
		t.handCards[data.SEAT1] = cards[:3]
		t.handCards[data.SEAT2] = cards[2:5]
		t.handCards[data.SEAT3] = cards[4:7]
		t.handCards[data.SEAT4] = cards[6:9]
		t.handCards[data.SEAT5] = append(cards[8:], cards[0])
		return true
	case data.GAME_JIU:
		t.handCards[data.SEAT1] = cards[:2]
		t.handCards[data.SEAT2] = cards[2:4]
		t.handCards[data.SEAT3] = cards[4:6]
		t.handCards[data.SEAT4] = cards[6:8]
		t.handCards[data.SEAT5] = cards[8:]
		return true
	default:
		glog.Error("Gtype err")
	}
	return false
}

//.

/*
赔付规则:
有庄：每个闲家分别和庄家比(大吃小)
无庄：通比规则(大吃小),依次赔付,赔付到平为止
      (最大吃平,第二大继续吃,依次吃到最后)
      退款只退给当前位置,不所有输家位置平分
*/

//结束游戏
func (t *Desk) gameOver() {
	//发牌
	if !t.dealCard() {
		//TODO 异常时退款
		t.unusualRefund()
	} else {
		//结算
		//有庄(庄赔),无庄(小赔大)
		if t.isDeal() {
			//有庄
			t.paili1() //牌力
			//有庄抽佣
			t.roomNewFeeDeal()
			//结算,闲家赔付
			t.payment1()
		} else {
			//无庄
			t.paili2() //牌力
			//抽佣
			t.roomNewFee()
			//结算,闲家赔付
			t.payment2()
		}
	}
	//统计玩家结果
	t.roleOver()
	//发放奖励
	t.settlement()
	//日志记录
	t.overRecord()
	//个人记录
	t.printOver() //打印信息
	//结束消息
	t.overMsg()
	//踢除离线玩家,庄家离线处理
	t.checkOffline()
	//下庄
	t.checkDealer()
	//TODO 推送庄家信息
	t.gameInit() //重置
}

//是否有庄
func (t *Desk) isDeal() bool {
	return t.dealer != "" // && t.dealerNum != 0
}

//离线处理
func (t *Desk) checkOffline() {
	for k, _ := range t.players {
		if !t.offline[k] {
			continue
		}
		errcode := t.Leave(k)
		if errcode != pb.OK {
			continue
		}
		//离开消息
		msg2 := &pb.LeaveDesk{
			Roomid: t.id,
			Userid: k,
			Type:   1,
		}
		nodePid.Tell(msg2)
	}
}

//' 统计玩家结果
func (t *Desk) roleOver() {
	//本金
	for _, m := range t.seatRoleCost {
		for userid, val := range m {
			t.roleCost[userid] += val
		}
	}
	//退款
	for seat, m := range t.seatRoleRefund {
		for userid, val := range m {
			t.roleRefund[userid] += val
			//位置总退款
			t.seatRefund[seat] += val
		}
	}
	//收款
	for _, m := range t.seatRoleWins {
		for userid, val := range m {
			if val <= 0 { //输家
				t.roleLose[userid] += val
				continue
			}
			//实际赢利所得
			t.roleWins[userid] += val
		}
	}
	//位置输赢
	for _, v := range t.result {
		//位置存在退款只能是输家
		if val, ok := t.seatRefund[v.Seat]; ok {
			//位置输掉金额=退款 - 下注额
			t.seatWins[v.Seat] = val - t.seatBets[v.Seat]
		} else {
			//位置赢利,实际赢利
			if m, ok := t.seatRoleWins[v.Seat]; ok {
				for _, val2 := range m {
					t.seatWins[v.Seat] += val2
				}
			}
		}
	}
	//抽佣,TODO 无庄场暂时不再抽赢家,改为结算前先抽输家再分配钱
	//if t.isDeal() {
	//	//TODO 有庄暂时按结果抽佣
	//	t.roomFee()
	//}
	//玩家总输赢,结果展示用
	for k, v := range t.roleBets {
		//赢利+退款+本金-下注额
		t.roleProfits[k] = t.roleWins[k] + t.roleRefund[k] + t.roleCost[k] - v
	}
}

//抽佣
func (t *Desk) countFee(v int64) (n int64) {
	//抽佣房间
	if t.DeskData.Rtype != data.ROOM_TYPE1 {
		return
	}
	if v <= 0 {
		return
	}
	n = int64(math.Trunc((float64(t.DeskData.Cost) / float64(100)) * float64(v)))
	return
}

//抽佣，抽取赢家的纯得利百分比
func (t *Desk) roomFee() {
	//抽佣房间
	if t.DeskData.Rtype != data.ROOM_TYPE1 {
		return
	}
	//TODO 存在不抽佣情况,如果玩家输赢持平时候就会不抽佣
	//比较自己的下注赔付给自己了
	for k, v := range t.roleWins {
		//输多于赢的钱
		val := t.roleLose[k] + v
		if val <= 0 {
			continue
		}
		//赢钱抽佣(抽赢到的钱)
		n := t.countFee(val)
		if n <= 0 {
			//t.roleProfits[k] = val
			continue
		}
		//t.roleProfits[k] = val - n
		//只记录日志
		p := t.getPlayer(k)
		if p != nil { //抽佣日志
			msg2 := handler.LogChipMsg(n, data.LogType43, p)
			nodePid.Tell(msg2)
			//机器人抽佣
			if p.GetRobot() {
				t.HuiYinDeskData.robotNum += n
			} else {
				t.HuiYinDeskData.playerNum += n
			}
		}
		//赢利
		t.roleWins[k] = v - n
		//抽佣数量
		t.HuiYinDeskData.feeNum += n
		//玩家抽佣明细
		t.fees[k] += n
	}
}

//.

//' 有庄场结算前抽佣
//有庄按下注比例抽庄家,从庄家身上扣除
//分配方式和无庄一样先分配位置再分配玩家
func (t *Desk) roomNewFeeDeal() {
	//抽庄家不区分
	t.roomNewFee5()
	//抽佣(抽下注总额)
	n := t.countFee(t.HuiYinDeskData.betNum)
	if n <= 0 {
		return
	}
	p := t.getPlayer(t.dealer)
	if p == nil {
		glog.Errorf("roomNewFeeDeal err dealer %s", t.dealer)
		return
	}
	//总抽佣
	//t.HuiYinDeskData.feeNum = n
	//TODO 暂时算庄家位置抽佣,记录中没有位置上玩家抽佣明细
	t.seatFees[t.dealerSeat] = n
	//扣除庄家
	t.carry -= n
	if t.carry < 0 {
		glog.Errorf("carry not enough %d, n %d", t.carry, n)
		t.carry = 0
	}
	msg3 := handler.LogChipMsg(n, data.LogType45, p)
	nodePid.Tell(msg3)
	//庄家输赢数量
	t.roleProfits[t.dealer] = -1 * n
	//反佣分配,按位置比例分配反佣,位置上玩家平均反佣
	t.roomNewFee6()
	//反佣给庄家自己
	//t.setRoleFees(t.dealer, n)
	//日志记录
	t.roomNewFee4()
}

//.

//' 结算前抽佣

//新的抽佣算法,结算前按照下注总数抽佣,
//抽佣数量按输家位置顺序分配,
//抽佣总数平均到每个位置,位置数量再平均到每个玩家
func (t *Desk) roomNewFee() {
	//无人下注或没有输赢不记录
	if len(t.roleBets) == 0 || t.betNum <= 0 {
		return
	}
	//抽佣房间
	if t.DeskData.Rtype != data.ROOM_TYPE1 {
		t.roomNewFee5()
		return
	}
	//抽佣(抽下注总额)
	n := t.countFee(t.HuiYinDeskData.betNum)
	if n <= 0 {
		return
	}
	//总抽佣
	//t.HuiYinDeskData.feeNum = n
	//位置分配抽佣
	seatLen := len(t.seatPower)
	for j := seatLen - 1; j >= 0; j-- {
		seat := t.seatPower[j]
		//位置没有下注
		if _, ok := t.seatBets[seat]; !ok {
			continue
		}
		if t.seatBets[seat] >= n {
			t.seatFees[seat] = n
			t.seatBets[seat] -= n //位置下注额相应减少
			break
		} else {
			t.seatFees[seat] = t.seatBets[seat]
			n -= t.seatBets[seat]
			//t.seatBets[seat] = 0 //位置下注额相应减少
			delete(t.seatBets, seat)
		}
	}
	//位置上玩家平均抽佣
	for k, v := range t.seatFees {
		t.roomNewFee2(k, v)
	}
	//抽佣后下注位置剩余金额
	t.roomNewFee3()
	//反佣分配
	t.roomNewFee6()
	//日志记录
	t.roomNewFee4()
}

//免佣房间,直接复制下注结果
func (t *Desk) roomNewFee5() {
	for s, m := range t.seatRoleBets {
		if _, ok := t.seatRoleFeeBets[s]; !ok {
			m2 := make(map[string]int64)
			t.seatRoleFeeBets[s] = m2
		}
		for k, v := range m {
			t.seatRoleFeeBets[s][k] = v
		}
	}
}

//日志记录
func (t *Desk) roomNewFee4() {
	for k, n := range t.fees {
		//只记录日志
		p := t.getPlayer(k)
		if p != nil { //抽佣日志
			msg2 := handler.LogChipMsg(n, data.LogType43, p)
			nodePid.Tell(msg2)
			//机器人抽佣
			if p.GetRobot() {
				t.HuiYinDeskData.robotNum += n
			} else {
				t.HuiYinDeskData.playerNum += n
			}
		}
		//总抽佣
		t.HuiYinDeskData.feeNum += n
	}
}

//抽佣后下注位置剩余金额
func (t *Desk) roomNewFee3() {
	for s, m := range t.seatRoleBets {
		if _, ok := t.seatRoleFeeBets[s]; !ok {
			m2 := make(map[string]int64)
			t.seatRoleFeeBets[s] = m2
		}
		if m2, ok := t.seatRoleFees[s]; ok {
			//有抽佣
			for k, v := range m {
				//剩余 = 下注额 - 抽佣数量
				if v > m2[k] {
					t.seatRoleFeeBets[s][k] = v - m2[k]
				}
				//已经抽完不写入
			}
		} else {
			for k, v := range m {
				t.seatRoleFeeBets[s][k] = v
			}
		}
	}
}

//按比例分配每个玩家抽佣金额
func (t *Desk) roomNewFee2(seat uint32, feeNum int64) {
	bets := t.seatBets[seat] //位置下注总额
	//位置上所有下注玩家
	ids := make([]string, 0)
	for k, _ := range t.seatRoleBets[seat] {
		ids = append(ids, k)
	}
	//数量
	val := feeNum
	for {
		if len(ids) == 0 || val <= 0 {
			break
		}
		//TODO 如果最后一个玩家不足
		if len(ids) == 1 {
			t.setSeatRoleFees(seat, ids[0], val)
			ids = ids[1:]
			break
		}
		//(玩家下注额 / 总额) * 总抽佣额
		num := t.seatRoleBets[seat][ids[0]]
		n := int64(math.Trunc((float64(num) / float64(bets)) * float64(feeNum))) //feeNum不变
		t.setSeatRoleFees(seat, ids[0], n)
		val -= n
		ids = ids[1:]
	}
}

//设置每个玩家分配到的抽佣
func (t *Desk) setSeatRoleFees(seat uint32, userid string, num int64) {
	if num <= 0 {
		return
	}
	//seatRoleFees 作用是用于计算抽佣后玩家剩余下注额
	//实际反佣给玩家是按整个房间下注比例重新分配
	//所以fees不能在这里计算
	//玩家抽佣明细
	//t.fees[userid] += num
	//位置明细
	if m, ok := t.seatRoleFees[seat]; ok {
		m[userid] = num
		t.seatRoleFees[seat] = m
	} else {
		m := make(map[string]int64)
		m[userid] = num
		t.seatRoleFees[seat] = m
	}
}

//实际反佣给玩家是按整个房间下注比例重新分配
func (t *Desk) roomNewFee6() {
	//t.HuiYinDeskData.feeNum = n
	var feeNum int64
	for _, v := range t.seatFees {
		feeNum += v
	}
	bets := t.betNum //位置下注总额
	if bets <= 0 {
		return
	}
	//位置上所有下注玩家
	ids := make([]string, 0)
	for k, _ := range t.roleBets {
		ids = append(ids, k)
	}
	//数量
	val := feeNum
	for {
		if len(ids) == 0 || val <= 0 {
			break
		}
		if len(ids) == 1 {
			t.setRoleFees(ids[0], val)
			ids = ids[1:]
			break
		}
		//(玩家下注额 / 总额) * 总反佣额
		num := t.roleBets[ids[0]]
		n := int64(math.Trunc((float64(num) / float64(bets)) * float64(feeNum))) //feeNum不变
		t.setRoleFees(ids[0], n)
		val -= n
		ids = ids[1:]
	}
}

//设置每个玩家分配到的反佣
func (t *Desk) setRoleFees(userid string, num int64) {
	if num <= 0 {
		return
	}
	//实际反佣给玩家是按整个房间下注比例重新分配
	t.fees[userid] += num
}

//.

//' 结算奖励发放

//异常时退款,TODO 是否直接走退款消息
func (t *Desk) unusualRefund() {
	//for k, v := range t.roleBets {
	//	t.sendChip(k, v, data.LogType42)
	//}
	//直接走退款消息
	for seat, m := range t.seatRoleBets {
		for userid, val := range m {
			t.setSeatRoleRefund(seat, userid, val)
		}
	}
}

//关闭时退款,TODO 是否直接走退款消息
func (t *Desk) closeRefund() {
	glog.Errorf("closeRefund roleBets %#v", t.roleBets)
	for k, v := range t.roleBets {
		t.sendChip(k, v, data.LogType42)
	}
	glog.Errorf("closeRefund dealer %s, carry %d", t.dealer, t.carry)
	//庄家退款
	if t.dealer != "" {
		p := t.players[t.dealer]
		t.delBeDealer(t.dealer, p)
	}
	glog.Errorf("closeRefund dealers %#v", t.dealers)
	//上庄列表退款
	for _, m := range t.dealers {
		for k, v := range m {
			t.sendChip(k, v, data.LogType8)
		}
	}
	t.dealers = make([]map[string]int64, 0)
}

//结算发放奖励
func (t *Desk) settlement() {
	//退款
	t.sendRefund()
	//本金
	t.sendCost()
	//实际赢金额
	t.sendWin()
}

//正常退款
func (t *Desk) sendRefund() {
	for k, v := range t.roleRefund {
		t.sendChip(k, v, data.LogType38)
	}
}

//本金返还
func (t *Desk) sendCost() {
	for k, v := range t.roleCost {
		t.sendChip(k, v, data.LogType39)
	}
}

//庄家输赢结果和日志, dealerCarry 庄家结算后金额
func (t *Desk) sendDealerWin(dealerCarry int64) {
	glog.Debugf("sendDealerWin dealer %s, %d, %d", t.dealer, dealerCarry, t.carry)
	//庄家实际输赢金额=结算后金额-开始时携带
	num := dealerCarry - t.carry
	//庄家输赢数量
	t.roleProfits[t.dealer] += num
	//剩余金额
	if dealerCarry < 0 {
		t.carry = 0
	} else {
		t.carry = dealerCarry
	}
	p := t.getPlayer(t.dealer)
	if p == nil {
		glog.Errorf("sendDealerWin err dealerCarry %d, dealer %s", dealerCarry, t.dealer)
		return
	}
	//输赢日志
	if num != 0 {
		msg2 := handler.LogChipMsg(num, data.LogType41, p)
		nodePid.Tell(msg2)
	}
}

//庄家输赢结果和日志
/*
func (t *Desk) sendDealerWin_Deprecated(dealerCarry int64) {
	num := dealerCarry - t.carry
	var n int64
	if num <= 0 {
		t.carry = 0
	} else {
		//赢钱抽拥
		n = t.countFee(num)
		if n > 0 {
			t.carry = dealerCarry - n
			num -= n
			//抽佣数量,TODO 庄家暂时按输赢抽
			t.HuiYinDeskData.feeNum += n
		} else {
			t.carry = dealerCarry
		}
	}
	p := t.getPlayer(t.dealer)
	if p == nil {
		glog.Errorf("sendDealerWin err dealerCarry %d, dealer %s", dealerCarry, t.dealer)
		return
	}
	//机器人抽佣 TODO 庄家暂时按输赢抽
	if p.GetRobot() {
		t.HuiYinDeskData.robotNum += n
	} else {
		t.HuiYinDeskData.playerNum += n
	}
	//玩家抽佣明细 TODO 庄家暂时按输赢抽
	t.fees[t.dealer] += n
	// TODO 庄家暂时按输赢抽
	if n > 0 { //抽佣日志
		msg3 := handler.LogChipMsg(n, data.LogType43, p)
		nodePid.Tell(msg3)
	}
	//输赢日志
	if num != 0 {
		msg2 := handler.LogChipMsg(num, data.LogType41, p)
		nodePid.Tell(msg2)
	}
}
*/

//赢家
func (t *Desk) sendWin() {
	for k, v := range t.roleWins {
		t.sendChip(k, v, data.LogType40)
	}
}

//.

//' 结果日志记录

//结算日志记录
func (t *Desk) overRecord() {
	//上局赢家记录
	t.winerRecord()
	//记录房间趋势
	t.trendRecord()
	//记录日志
	t.saveRecord()
	//个人记录
	t.setRecord()
}

//赢家记录
func (t *Desk) winerRecord() {
	t.Winers = make([]*data.Winer, 0)
	for k, v := range t.roleProfits {
		//for k, v := range t.roleWins {
		//if v <= 0 {
		//	continue
		//}
		p := t.getPlayer(k)
		if p == nil {
			continue
		}
		w := new(data.Winer)
		w.Userid = k
		w.Nickname = p.GetNickname()
		w.Photo = p.GetPhoto()
		w.Chip = v
		if k == t.dealer {
			w.Dealer = true
		}
		t.Winers = append(t.Winers, w)
	}
}

//趋势记录
func (t *Desk) trendRecord() {
	trend := new(data.Trend)
	trend.Expect = t.HuiYinDeskData.expect
	trend.Opencode = t.HuiYinDeskData.opencode
	trend.Opentime = t.HuiYinDeskData.opentime
	trend.Result = t.HuiYinDeskData.result
	t.Trends = append(t.Trends, trend)
	//限制长度
	if len(t.Trends) >= 15 {
		t.Trends = t.Trends[1:]
	}
	//存储记录
	msg := new(pb.Pk10TrendLog)
	msg.Expect = trend.Expect
	msg.Opencode = trend.Opencode
	msg.Opentime = trend.Opentime
	msg.Result = t.trendResult()
	nodePid.Tell(msg)
}

//趋势结果
func (t *Desk) trendResult() []*pb.TrendResult {
	l := make([]*pb.TrendResult, 0)
	for _, v := range t.result {
		d := &pb.TrendResult{
			Rank:  v.Rank,
			Seat:  v.Seat,
			Point: v.Point,
			Cards: v.Cards,
		}
		l = append(l, d)
	}
	return l
}

//结果记录
func (t *Desk) saveRecord() {
	//无人下注或没有输赢不记录
	if len(t.roleBets) == 0 && !t.isDeal() {
		return
	}
	var roomid string = data.GenCporderid(t.id)
	msg1 := new(pb.Pk10GameLog)
	msg1.Roomid = roomid
	msg1.Gametype = t.DeskData.Gtype
	msg1.Roomtype = t.DeskData.Rtype
	msg1.Lotterytype = t.DeskData.Ltype
	msg1.Expect = t.HuiYinDeskData.expect
	msg1.Opencode = t.HuiYinDeskData.opencode
	msg1.Opentime = t.HuiYinDeskData.opentime
	msg1.Num = uint32(len(t.roleBets))          //参与人数
	msg1.FeeNum = t.HuiYinDeskData.feeNum       //抽佣数量
	msg1.RobotFee = t.HuiYinDeskData.robotNum   //抽佣数量
	msg1.PlayerFee = t.HuiYinDeskData.playerNum //抽佣数量
	msg1.Trend = t.trendResult()
	t.resultRecord(msg1)
	t.feeResult(msg1)
	nodePid.Tell(msg1)
	//个人
	for k, v := range t.roleBets {
		msg2 := new(pb.Pk10UseridLog)
		msg2.Roomid = roomid
		msg2.Userid = k
		msg2.Bets = v
		msg2.Profits = t.roleProfits[k]
		msg2.Gametype = t.DeskData.Gtype
		msg2.Roomtype = t.DeskData.Rtype
		msg2.Lotterytype = t.DeskData.Ltype
		msg2.Expect = t.HuiYinDeskData.expect
		if p, ok := t.players[k]; ok {
			msg2.Rest = p.GetChip()
			msg2.Robot = p.GetRobot()
		}
		msg2.Fee = t.HuiYinDeskData.fees[k]
		msg2.Details = t.roleDetails(k)
		msg2.Dealer = t.dealer
		msg2.Dealerseat = t.dealerSeat
		nodePid.Tell(msg2)
	}
	//庄家记录
	if t.isDeal() {
		k := t.dealer
		msg2 := new(pb.Pk10UseridLog)
		msg2.Roomid = roomid
		msg2.Userid = k
		msg2.Profits = t.roleProfits[k]
		msg2.Gametype = t.DeskData.Gtype
		msg2.Roomtype = t.DeskData.Rtype
		msg2.Lotterytype = t.DeskData.Ltype
		msg2.Expect = t.HuiYinDeskData.expect
		if p, ok := t.players[k]; ok {
			msg2.Rest = p.GetChip()
			msg2.Robot = p.GetRobot()
		}
		msg2.Fee = t.HuiYinDeskData.fees[k]
		msg2.Details = t.roleDetails(k)
		msg2.Dealer = t.dealer
		msg2.Dealerseat = t.dealerSeat
		nodePid.Tell(msg2)
	}
}

//个人位置上结果
func (t *Desk) roleDetails(k string) []*pb.UseridDetails {
	l := make([]*pb.UseridDetails, 0)
	for s, m := range t.seatRoleBets {
		if v, ok := m[k]; ok {
			d := &pb.UseridDetails{
				Seat: s,
				Bets: v,
			}
			if m2, ok2 := t.seatRoleWins[s]; ok2 {
				d.Wins = m2[k]
			}
			//无庄场输赢结果需要加上抽佣数量
			if !t.isDeal() {
				if m2, ok2 := t.seatRoleFees[s]; ok2 {
					d.Wins -= m2[k]
				}
			}
			if m2, ok2 := t.seatRoleRefund[s]; ok2 {
				d.Refund = m2[k]
			}
			l = append(l, d)
		}
	}
	//庄家位置输赢
	if t.isDeal() && k == t.dealer {
		for s, v := range t.dealerSeatWins {
			d := &pb.UseridDetails{
				Seat: s,
				Wins: v,
			}
			l = append(l, d)
		}
	}
	return l
}

//玩家总明细结果
func (t *Desk) resultRecord(msg1 *pb.Pk10GameLog) {
	l := make([]*pb.ResultRecord, 0)
	for k, v := range t.roleBets {
		d := &pb.ResultRecord{
			Userid: k,
			Bets:   v,
			Wins:   t.roleProfits[k], //输赢(纯利)
			Refund: t.roleRefund[k],  //退款
		}
		l = append(l, d)
		//统计
		msg1.BetNum += v
		if t.roleProfits[k] > 0 {
			msg1.WinNum += t.roleProfits[k] //赢(纯利)
		} else {
			msg1.LoseNum += t.roleProfits[k] //输(纯利)
		}
		msg1.RefundNum += t.roleRefund[k] //退款
	}
	//庄家记录
	if t.isDeal() {
		k := t.dealer
		d := &pb.ResultRecord{
			Userid: k,
			Wins:   t.roleProfits[k], //输赢(纯利)
		}
		l = append(l, d)
		//统计
		if t.roleProfits[k] > 0 {
			msg1.WinNum += t.roleProfits[k] //赢(纯利)
		} else {
			msg1.LoseNum += t.roleProfits[k] //输(纯利)
		}
	}
	msg1.Result = l
}

//玩家抽佣明细
func (t *Desk) feeResult(msg1 *pb.Pk10GameLog) {
	l := make([]*pb.FeeResult, 0)
	for k, v := range t.fees {
		d := &pb.FeeResult{
			Userid: k,
			Fee:    v,
		}
		l = append(l, d)
	}
	msg1.Record = l
	//位置上玩家抽佣明细
	l2 := make([]*pb.FeeDetails, 0)
	for k, v := range t.seatFees {
		d := &pb.FeeDetails{
			Seat: k,
			Fee:  v,
		}
		if m, ok := t.seatRoleFees[k]; ok {
			for k2, v2 := range m {
				d2 := &pb.FeeResult{
					Userid: k2,
					Fee:    v2,
				}
				d.Record = append(d.Record, d2)
			}
		}
		l2 = append(l2, d)
	}
	msg1.Details = l2
}

//个人记录
func (t *Desk) setRecord() {
	for k, v := range t.roleProfits {
		p := t.getPid(k)
		if p == nil {
			continue
		}
		user := t.getPlayer(k)
		if user == nil {
			continue
		}
		msg2 := new(pb.SetRecord)
		if v > 0 {
			msg2.Rtype = 1
		} else if v < 0 {
			msg2.Rtype = -1
		} else {
			msg2.Rtype = 0
		}
		//更新游戏内数据
		user.SetRecord(msg2.Rtype)
		//更新节点数据
		p.Tell(msg2)
	}
}

//.

//' 结算消息

//结算消息,TODO 是否只广播坐下的人的数据
func (t *Desk) overMsg() {
	msg := new(pb.SHuiYinGameover)
	msg.Dealer = t.dealer
	msg.DealerSeat = t.dealerSeat
	msg.Carry = t.carry
	msg.Expect = t.expect
	msg.Opencode = t.opencode
	//下注位置结算明细
	msg.Seats = t.seatOverMsg()
	//玩家结算明细
	msg.Data = t.roleOverMsg()
	//glog.Debugf("over msg %#v", msg)
	//glog.Debugf("over pids len %d", len(t.pids))
	//自己结算明细,单独通知
	for k, p := range t.pids {
		msg2 := new(pb.SHuiYinGameover)
		*msg2 = *msg
		msg2.Data = make([]*pb.HuiYinRoomOver, 0)
		msg2.Seats = make([]*pb.HuiYinSeatOver, 0)
		msg2.Data = msg.Data
		//单个玩家
		if msg4, ok := t.roleOverMsg1(k); ok {
			msg2.Data = append(msg2.Data, msg4)
		}
		//下注位置列表
		for _, val2 := range msg.Seats {
			val3 := new(pb.HuiYinSeatOver)
			*val3 = *val2
			val3.List = make([]*pb.HuiYinRoomWins, 0)
			val3.List = val2.List
			//单个玩家
			if msg4, ok := t.seatRoleOverMsg1(val2.Seat, k); ok {
				val3.List = append(val3.List, msg4)
			}
			msg2.Seats = append(msg2.Seats, val3)
		}
		glog.Debugf("over msg %#v", msg2)
		if p == nil {
			continue
		}
		p.Tell(msg2)
	}
}

//玩家(k)自己结算明细且没有坐下
func (t *Desk) roleOverMsg1(k string) (*pb.HuiYinRoomOver, bool) {
	//坐下玩家跳过
	if t.seats[k] != 0 {
		return nil, false
	}
	if v, ok := t.roleBets[k]; ok {
		msg4 := &pb.HuiYinRoomOver{
			Userid: k,
			Seat:   t.seats[k],
			Bets:   v,
			Cost:   t.roleCost[k],
			Wins:   t.roleProfits[k], //输赢(纯利)
			Refund: t.roleRefund[k],
		}
		return msg4, true
	}
	//庄家数据
	if t.isDeal() && k == t.dealer {
		msg4 := &pb.HuiYinRoomOver{
			Userid: k,
			Seat:   t.seats[k],
			Wins:   t.roleProfits[k], //输赢(纯利)
		}
		return msg4, true
	}
	return nil, false
}

//玩家结算明细
func (t *Desk) roleOverMsg() []*pb.HuiYinRoomOver {
	//坐下玩家结算明细
	l := make([]*pb.HuiYinRoomOver, 0)
	for k, v := range t.roleBets {
		if t.seats[k] == 0 {
			continue
		}
		msg4 := &pb.HuiYinRoomOver{
			Userid: k,
			Seat:   t.seats[k],
			Bets:   v,
			Cost:   t.roleCost[k],
			Wins:   t.roleProfits[k], //输赢(纯利)
			Refund: t.roleRefund[k],
		}
		l = append(l, msg4)
	}
	//庄家位置输赢
	if t.isDeal() {
		k := t.dealer
		msg4 := &pb.HuiYinRoomOver{
			Userid: k,
			Seat:   t.seats[k],
			Wins:   t.roleProfits[k], //输赢(纯利)
		}
		l = append(l, msg4)
	}
	return l
}

//位置下注明细
func (t *Desk) seatOverMsg() []*pb.HuiYinSeatOver {
	//下注位置结算明细
	l := make([]*pb.HuiYinSeatOver, 0)
	for _, v := range t.result {
		msg3 := &pb.HuiYinSeatOver{
			Rank:  v.Rank,
			Seat:  v.Seat,
			Cards: v.Cards,
			Point: v.Point,
			Bets:  t.seatBets[v.Seat],
		}
		//位置上总输赢结果
		//位置本金输赢
		if m, ok := t.seatRoleCost[v.Seat]; ok {
			for _, val2 := range m {
				if val2 > 0 {
					msg3.Cost += val2 //赢家本金
				}
			}
		}
		//庄家位置输赢
		if t.isDeal() && v.Seat == t.dealerSeat {
			msg3.WinNum = t.roleProfits[t.dealer]
			//位置存在退款只能是输家
		} else if val, ok := t.seatRefund[v.Seat]; ok {
			//位置输掉金额=退款 - 下注额
			msg3.WinNum = val - t.seatBets[v.Seat]
			//位置退款
			msg3.Refund = t.seatRefund[v.Seat]
		} else {
			//位置赢利,实际赢利
			if m, ok := t.seatRoleWins[v.Seat]; ok {
				for _, val2 := range m {
					msg3.WinNum += val2
				}
			}
		}
		msg3.List = t.seatRoleOverMsg(v.Seat)
		l = append(l, msg3)
	}
	return l
}

//下注位置上坐下的玩家明细
func (t *Desk) seatRoleOverMsg(seat uint32) []*pb.HuiYinRoomWins {
	l := make([]*pb.HuiYinRoomWins, 0)
	for k, v := range t.seatRoleBets[seat] {
		if t.seats[k] == 0 {
			continue
		}
		msg4 := &pb.HuiYinRoomWins{
			Userid: k,
			Seat:   t.seats[k],
			Bets:   v,
		}
		if m, ok := t.seatRoleWins[seat]; ok {
			msg4.Wins = m[k]
		}
		if m, ok := t.seatRoleRefund[seat]; ok {
			msg4.Refund = m[k]
		}
		l = append(l, msg4)
	}
	return l
}

//下注位置上单个玩家明细
func (t *Desk) seatRoleOverMsg1(seat uint32, k string) (*pb.HuiYinRoomWins, bool) {
	if t.seats[k] != 0 {
		return nil, false
	}
	if m2, ok := t.seatRoleBets[seat]; ok {
		if v, ok2 := m2[k]; ok2 {
			msg4 := &pb.HuiYinRoomWins{
				Userid: k,
				Seat:   t.seats[k],
				Bets:   v, //玩家在位置上下注数量
			}
			if m, ok := t.seatRoleWins[seat]; ok {
				msg4.Wins = m[k]
			}
			if m, ok := t.seatRoleRefund[seat]; ok {
				msg4.Refund = m[k]
			}
			return msg4, true
		}
	}
	return nil, false
}

//.

//' 有庄,庄家赔付

//比较大小
func (t *Desk) paili1() {
	for k, v := range t.handCards {
		t.power[k] = algo.Point(v)
		//结果记录
		result := data.TrendResult{
			Seat:  k,
			Point: t.power[k],
			Cards: v,
		}
		t.result = append(t.result, result)
	}
	cs1 := t.getHandCards(t.dealerSeat) //庄家牌
	a1 := t.power[t.dealerSeat]
	for k, v := range t.handCards {
		if k == t.dealerSeat {
			//庄家位置
			continue
		}
		//闲家输赢倍率
		t.multiple[k] = muliti(a1, t.power[k], cs1, v)
	}
	t.setRank1()
}

//返回庄家赢倍数,a1庄家牌力,an闲家牌力,庄家赢返回正数,输返回负数
func muliti(a1, an uint32, cs1, csn []uint32) int64 {
	switch {
	case a1 > an:
		return int64(algo.Multiple(a1))
	case a1 < an:
		return -1 * int64(algo.Multiple(an))
	case a1 == an:
		if a1 == 0 {
			//有庄时0点算庄家赢
			return int64(algo.Multiple(a1))
		}
		if algo.Compare(cs1, csn) {
			return int64(algo.Multiple(a1))
		}
	}
	//庄家输
	return -1 * int64(algo.Multiple(an))
}

//闲家赔付,有庄结算
func (t *Desk) payment1() {
	//没人下注
	if len(t.roleBets) == 0 {
		return
	}
	var winNum int64 //闲家赢总额
	//庄家携带,需要赔付总额为庄家携带总额
	dealerCarry := t.carry
	//输掉的位置
	loseSeat := make([]uint32, 0)
	//赢家的位置
	winSeat := make([]uint32, 0)
	glog.Debugf("multiple %#v", t.multiple)
	for k, v := range t.multiple {
		if v < 0 { //表示庄家输
			winNum += t.seatBets[k]
			winSeat = append(winSeat, k)
		} else { //表示庄家赢
			loseSeat = append(loseSeat, k)
		}
	}
	//先收钱再赔付
	//闲家赔付
	glog.Debugf("loseSeat %#v, winSeat %#v", loseSeat, winSeat)
	//庄家收钱
	for _, s := range loseSeat {
		dealerCarry += t.payment1lose(s)
	}
	//没有赢家,庄家通吃
	if len(winSeat) == 0 {
		t.sendDealerWin(dealerCarry)
		return
	}
	//赢家本金退还
	for _, s := range winSeat {
		t.payment1cost(s)
	}
	glog.Debugf("winNum %#v, dealerCarry %#v", winNum, dealerCarry)
	//不够赔付
	if dealerCarry < winNum {
		//按位置分配dealerCarry,返回每个位置上分到的金额
		seatsWin := t.payment1SetWin(dealerCarry, winNum, winSeat)
		for s, v := range seatsWin {
			//位置上按下注比例分配
			t.payment1loseNum(s, v, t.seatBets[s])
			//庄家位置输赢
			t.dealerSeatWins[s] = (-1 * v)
		}
		//庄家日志,庄家赔光
		t.sendDealerWin(0)
		return
	}
	//足够赔付
	for _, s := range winSeat {
		//按1比1赔付
		t.payment1win(s)
		//庄家位置输赢,有庄场抽庄家的,所以可以直接用seatBets
		//t.dealerSeatWins[s] = (-1 * t.seatBets[s])
		//或者
		for _, v := range t.seatRoleFeeBets[s] {
			t.dealerSeatWins[s] += (-1 * v)
		}
	}
	//庄家日志
	t.sendDealerWin(dealerCarry - winNum)
}

//按位置分配庄家全部金额dealerCarry
func (t *Desk) payment1SetWin(dealerCarry, winNum int64,
	winSeat []uint32) (m map[uint32]int64) {
	m = make(map[uint32]int64)
	val := dealerCarry
	for {
		if len(winSeat) == 0 || val <= 0 {
			break
		}
		if len(winSeat) == 1 {
			m[winSeat[0]] = val
			winSeat = winSeat[1:]
			break
		}
		//(位置下注总额 / 赢家总额) * 庄家总额
		sbet := t.seatBets[winSeat[0]]
		n := int64(math.Trunc((float64(sbet) / float64(winNum)) * float64(dealerCarry)))
		m[winSeat[0]] = n //位置分配到的金额
		val -= n
		winSeat = winSeat[1:]
	}
	return
}

//位置上玩家按下注比例分配loseNum, winNum位置seat下注总额
func (t *Desk) payment1loseNum(seat uint32, loseNum, winNum int64) {
	ids := make([]string, 0)
	for k, _ := range t.seatRoleFeeBets[seat] {
		ids = append(ids, k)
	}
	val := loseNum
	for {
		if len(ids) == 0 || val <= 0 {
			break
		}
		if len(ids) == 1 {
			//t.seatRoleWins[seat][ids[0]] += val
			t.setSeatRoleWins(seat, ids[0], val)
			ids = ids[1:]
			break
		}
		//(玩家下注额 / 总额) * 总退款额
		num := t.seatRoleFeeBets[seat][ids[0]]
		n := int64(math.Trunc((float64(num) / float64(winNum)) * float64(loseNum)))
		//t.seatRoleWins[seat][ids[0]] += n
		t.setSeatRoleWins(seat, ids[0], n)
		val -= n
		ids = ids[1:]
	}
}

//有庄闲家赢本金返还
func (t *Desk) payment1cost(seat uint32) {
	//位置上没有下注
	if _, ok := t.seatRoleFeeBets[seat]; !ok {
		return
	}
	if _, ok := t.seatBets[seat]; !ok {
		return
	}
	//返还本金
	for k, v := range t.seatRoleFeeBets[seat] {
		t.setSeatRoleCost(seat, k, v)
	}
}

//有庄庄家赔付,足够赔付
func (t *Desk) payment1win(seat uint32) {
	//位置上没有下注
	if _, ok := t.seatRoleFeeBets[seat]; !ok {
		return
	}
	if _, ok := t.seatBets[seat]; !ok {
		return
	}
	//按1比1赔付
	for k, v := range t.seatRoleFeeBets[seat] {
		t.setSeatRoleWins(seat, k, v)
	}
}

//有庄输家赔付
func (t *Desk) payment1lose(seat uint32) (num int64) {
	//位置上没有下注
	if _, ok := t.seatRoleFeeBets[seat]; !ok {
		return
	}
	if _, ok := t.seatBets[seat]; !ok {
		return
	}
	//不足赔付
	if t.seatBets[seat] <= t.carry {
		//全部赔付给庄家
		t.payment1loseSeat(seat)
		//庄家位置输赢
		t.dealerSeatWins[seat] = t.seatBets[seat]
		num = t.seatBets[seat]
		return
	}
	//足够赔付
	//庄家位置输赢
	t.dealerSeatWins[seat] = t.carry
	num = t.carry
	//位置退款金额
	refundNum := t.seatBets[seat] - t.carry
	//每个玩家退款金额
	t.payment1seatRefund(seat, refundNum)
	//每个玩家输掉金额
	t.payment2loseSeat(seat)
	return
}

//部分赔付给庄家或赢家
func (t *Desk) payment2loseSeat(seat uint32) {
	for userid, val := range t.seatRoleFeeBets[seat] {
		//退款金额-下注金额
		rNum := t.getSeatRoleRefund(seat, userid)
		t.setSeatRoleWins(seat, userid, (rNum - val))
	}
}

//全部赔付给庄家或赢家
func (t *Desk) payment1loseSeat(seat uint32) {
	for k, v := range t.seatRoleFeeBets[seat] {
		t.setSeatRoleWins(seat, k, (-1 * v))
	}
}

//按下注比例分配每个玩家退款金额
func (t *Desk) payment1seatRefund(seat uint32, refundNum int64) {
	bets := t.seatBets[seat] //位置下注总额
	//位置上所有下注玩家
	ids := make([]string, 0)
	for k, _ := range t.seatRoleFeeBets[seat] {
		ids = append(ids, k)
	}
	//数量
	val := refundNum
	for {
		if len(ids) == 0 || val <= 0 {
			break
		}
		if len(ids) == 1 {
			t.setSeatRoleRefund(seat, ids[0], val)
			ids = ids[1:]
			break
		}
		//(玩家下注额 / 总额) * 总退款额
		num := t.seatRoleFeeBets[seat][ids[0]]
		n := int64(math.Trunc((float64(num) / float64(bets)) * float64(refundNum))) //refundNum不变
		t.setSeatRoleRefund(seat, ids[0], n)
		val -= n
		ids = ids[1:]
	}
}

//获取玩家退款
func (t *Desk) getSeatRoleRefund(seat uint32, userid string) (num int64) {
	if m, ok := t.seatRoleRefund[seat]; ok {
		num = m[userid]
	}
	return
}

//设置每个玩家退款
func (t *Desk) setSeatRoleRefund(seat uint32, userid string, num int64) {
	if m, ok := t.seatRoleRefund[seat]; ok {
		m[userid] = num
		t.seatRoleRefund[seat] = m
	} else {
		m := make(map[string]int64)
		m[userid] = num
		t.seatRoleRefund[seat] = m
	}
}

//设置每个玩家输赢
func (t *Desk) setSeatRoleWins(seat uint32, userid string, num int64) {
	if m, ok := t.seatRoleWins[seat]; ok {
		m[userid] = num //不能多次设置,会出现覆盖
		t.seatRoleWins[seat] = m
	} else {
		m := make(map[string]int64)
		m[userid] = num
		t.seatRoleWins[seat] = m
	}
}

//设置每个玩家本金返还
func (t *Desk) setSeatRoleCost(seat uint32, userid string, num int64) {
	if m, ok := t.seatRoleCost[seat]; ok {
		m[userid] = num
		t.seatRoleCost[seat] = m
	} else {
		m := make(map[string]int64)
		m[userid] = num
		t.seatRoleCost[seat] = m
	}
}

//.

//' 无庄,通比赔付

//比较大小,无庄(小赔大)
func (t *Desk) paili2() {
	for k, v := range t.handCards {
		//点数大小
		t.power[k] = algo.Point(v)
		//结果记录
		result := data.TrendResult{
			Seat:  k,
			Point: t.power[k],
			Cards: v,
		}
		t.result = append(t.result, result)
		//位置上没有下注无法赔付
		if t.seatBets[k] <= 0 {
			continue
		}
	}
	t.setRank()
}

//位置大小排序
func (t *Desk) setRank1() {
	sort.Slice(t.result, func(i, j int) bool {
		if t.result[i].Point == t.result[j].Point {
			if t.result[i].Point == 0 {
				//有庄时0点算庄家赢
				if t.result[i].Seat == t.dealerSeat {
					return true
				}
				if t.result[j].Seat == t.dealerSeat {
					return false
				}
			}
			//不是庄家按牌大小排序
			return algo.Compare(t.result[i].Cards, t.result[j].Cards)
		}
		return t.result[i].Point > t.result[j].Point
	})
	//设置排名
	for k, v := range t.result {
		//排名名次(1 - 5)
		v.Rank = uint32(k + 1)
		t.result[k] = v
		//设置位置大小排名,从大到小排序
		t.seatPower = append(t.seatPower, v.Seat)
	}
}

//位置大小排序
func (t *Desk) setRank() {
	sort.Slice(t.result, func(i, j int) bool {
		if t.result[i].Point == t.result[j].Point {
			return algo.Compare(t.result[i].Cards, t.result[j].Cards)
		}
		return t.result[i].Point > t.result[j].Point
	})
	//设置排名
	for k, v := range t.result {
		//排名名次(1 - 5)
		v.Rank = uint32(k + 1)
		t.result[k] = v
		//设置位置大小排名,从大到小排序
		t.seatPower = append(t.seatPower, v.Seat)
	}
}

//闲家赔付,无庄结算
func (t *Desk) payment2() {
	//TODO 一个人下注时也抽佣
	//只下一个位置或只一个人下注直接退款
	//if len(t.seatBets) == 1 || len(t.roleBets) == 1 {
	//	t.payment2refund()
	//	return
	//}
	//大吃小
	refunds, wins := t.payment2bets()
	glog.Debugf("refunds %v, wins %v", refunds, wins)
	//退款
	for k, v := range refunds {
		if v <= 0 {
			continue
		}
		if v >= t.seatBets[k] {
			//没有输赢,全额退款
			t.payment2setRefund(k)
		} else {
			//按比例分配
			t.payment1seatRefund(k, v)
		}
	}
	//分款
	for k, v := range wins {
		if v <= 0 {
			//每个玩家输掉金额
			t.payment2loseSeat(k)
			continue
		}
		//赢家本金退还
		t.payment1cost(k)
		if v >= t.seatBets[k] {
			//足够赔付,按1比1赔付
			t.payment1win(k)
		} else {
			//不足赔付,位置上按下注比例分配
			t.payment1loseNum(k, v, t.seatBets[k])
		}
	}
}

//大吃小依次赔付
func (t *Desk) payment2bets() (seatBets1, seatBets2 map[uint32]int64) {
	//赔付,存在剩余时退款
	seatBets1 = make(map[uint32]int64)
	//收款,大于0时存在收款
	seatBets2 = make(map[uint32]int64)
	for k, v := range t.seatBets {
		seatBets1[k] = v
		seatBets2[k] = 0
	}
	seatLen := len(t.seatPower)
	for i := 0; i < seatLen; i++ {
		for j := seatLen - 1; j > i; j-- {
			k := t.seatPower[i]
			v := t.seatPower[j]
			if seatBets1[k] <= seatBets1[v] {
				seatBets2[k] += seatBets1[k] //赔付
				seatBets1[v] -= seatBets1[k] //剩余
				seatBets1[k] = 0             //赔付完成
				break
			} else {
				//不足
				seatBets2[k] += seatBets1[v] //全部交付
				seatBets1[k] -= seatBets1[v] //剩余赔付
				seatBets1[v] = 0             //交付全部
			}
		}
	}
	//赢家不存在退款
	for k, v := range seatBets1 {
		if v <= 0 {
			continue
		}
		//赢家不存在退款
		if v2, ok := seatBets2[k]; ok && v2 > 0 {
			seatBets1[k] = 0
		}
	}
	return
}

//直接退款
func (t *Desk) payment2refund() {
	for seat, m := range t.seatRoleFeeBets {
		for userid, val := range m {
			t.setSeatRoleRefund(seat, userid, val)
		}
	}
}

//直接退款
func (t *Desk) payment2setRefund(seat uint32) {
	if m, ok := t.seatRoleFeeBets[seat]; ok {
		for userid, val := range m {
			t.setSeatRoleRefund(seat, userid, val)
		}
	}
}

//.

// vim: set foldmethod=marker foldmarker=//',//.:
