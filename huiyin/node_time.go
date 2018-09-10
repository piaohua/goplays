package main

import (
	"goplays/data"
	"goplays/game/handler"
	"goplays/glog"
	"goplays/pb"
	"utils"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//日志记录
func (a *DeskActor) logPk10(expect, opencode, opentime string,
	opentimestamp int64) {
	if expect == "" {
		glog.Errorf("logPk10 err %s, %s, %d")
		return
	}
	glog.Debugf("record %s, %s, %s, %d", expect, opencode,
		opentime, opentimestamp)
	msg2 := &pb.Pk10RecordLog{
		Expect:        expect,
		Opencode:      opencode,
		Opentime:      opentime,
		Opentimestamp: opentimestamp,
		Code:          a.code,
	}
	a.dbmsPid.Tell(msg2)
}

//数据获取地址
func getPk10ApiUrl(n int, code string) (bjpk10Api string) {
	bjpk10Api = cfg.Section(code).Key("hoapiplus").Value()
	if n != 0 {
		//备用地址
		bjpk10Api = cfg.Section(code).Key("zapiplus").Value()
	}
	return
}

//获取数据
func getPk10Codes(n int, code string) (d []data.Bjpk10) {
	bjpk10Api := getPk10ApiUrl(n, code)
	var err1 error
	d, err1 = handler.GetPk10Api2(bjpk10Api)
	glog.Debugf("getPk10 %v, %v", d, err1)
	if err1 != nil {
		glog.Errorf("getPk10 err %v", err1)
		return
	}
	return
}

//设置开奖号码
func (a *DeskActor) setPk10Code(n int) bool {
	if (a.timer % 3) != 0 {
		return false
	}
	d := getPk10Codes(n, a.code)
	if len(d) == 0 {
		return false
	}
	glog.Debugf("getPk10 n %d, d %v", n, d)
	//TODO 开奖时间验证,确认取到的是当前开奖结果
	if d[0].Expect <= a.lastexpect {
		glog.Errorf("getPk10 lastexpect %s, d %#v", a.lastexpect, d)
		return false
	}
	glog.Errorf("getPk10 lastexpect %s, expect %s", a.lastexpect, d[0].Expect)
	a.expect = d[0].Expect
	a.opencode = d[0].Opencode
	a.opentime = d[0].Opentime
	a.opentimestamp = d[0].Opentimestamp
	return true
}

//启动设置开奖号码
func (a *DeskActor) initSetPk10Code(n int) bool {
	d := getPk10Codes(n, a.code)
	if len(d) == 0 || d == nil {
		return false
	}
	glog.Debugf("getPk10 %v", d)
	//启动设置上轮结果,封盘后会开奖所以不设置
	//上期号码必须严格验证且获取到,不然下期就会出错
	if a.lastexpect == "" {
		a.lastexpect = d[0].Expect
		a.lastopencode = d[0].Opencode
	}
	glog.Debugf("getPk10 lastexpect %s", a.lastexpect)
	//默认取5条,全部写数据库
	for _, v := range d {
		a.logPk10(v.Expect, v.Opencode, v.Opentime, v.Opentimestamp)
	}
	return true
}

//设置上期期号
func (a *DeskActor) initSetLast(now, st int64) {
	//如果是封盘时期启动
	if now < st {
		return
	}
	switch a.code {
	case data.BJPK10:
		//TODO 存在放假,暂时不这样处理
		//a.lastexpect = setLastPk10(now, st)
	case data.MLAFT:
		//设置mlaft上期期号
		today := a.startTime.Year()*10000 +
			int(a.startTime.Month())*100 + a.startTime.Day()
		n := today * 1000
		//st := utils.Time2Stamp(a.startTime)
		n += int((now - st) / 300)
		a.lastexpect = utils.String(n)
	}
}

//设置pk10上期期号
func setLastPk10(now, st int64) string {
	//从2号开始算
	var startCode int64 = 664951
	//startCodeTime, _ := utils.Str2Local("2018-02-02 09:02:00")
	startCodeToday, _ := utils.Str2Local("2018-02-02 00:00:00")
	today := utils.TimestampToday()
	//每天开奖179期
	n := (today - startCodeToday) / 86400 * 179
	opened := ((now - st) / 300) + n + startCode
	return utils.String(opened)
}

//开奖时间计算
func leftCount(now, st, et int64) (nextopentime, left int64) {
	//一天总期数
	all := (et - st) / 300
	//当天已经开奖几次
	alreadyOpen := (now - st) / 300
	//上期开奖时间
	lastopentime := st + (alreadyOpen * 300)
	//当天停盘
	if now < (st + 42 + 30) {
		nextopentime = (st + 300 - now)
		left = all
		return
	}
	//下期开奖时间
	nextopentime = (lastopentime + 300 - now)
	if nextopentime < 0 {
		nextopentime = 0
	}
	//剩余期数
	left = all - alreadyOpen
	return
}

//设置上期期号
func (a *DeskActor) setLastExpect() {
	if len(a.expect) != 0 {
		glog.Debugf("initPk10 lastexpect %s, expect %s", a.lastexpect, a.expect)
		a.lastexpect = a.expect
		a.lastopencode = a.opencode
	}
}

//开始时重置记录
func (a *DeskActor) initPk10() {
	//防止启动时清除掉上局记录
	//a.setLastExpect()
	a.expect = ""
	a.opencode = ""
	a.opentime = ""
	a.opentimestamp = 0
}

//启动计时,TODO 时间校验
func (a *DeskActor) initTicker(ctx actor.Context) {
	now := utils.Timestamp()
	//初始化时间状态
	a.initTime(now)
	st := utils.Time2Stamp(a.startTime)
	glog.Debugf("now %s", utils.Unix2Str(now))
	//初始状态
	a.initState(now, st)
	glog.Debugf("state %d, timer %d", a.state, a.timer)
	//设置上期期号
	a.initSetLast(now, (st + 42 + 30))
	//获取数据,使用备用地址
	a.initSetPk10Code(1)
	//更新初始状态
	a.pushDeskState()
	//计时器
	go a.ticker(ctx)
}

//第一局开始时间 2018-01-12 09:02 秒开始结算
//第一局开奖时间 2018-01-12 09:07
//最后一局结束时间 2018-01-11 23:57
//初始化时间状态
func (a *DeskActor) initTime(now int64) {
	if a.state != data.STATE_READY {
		return
	}
	year, month, day := utils.DateTime()
	switch a.code {
	case data.BJPK10:
		a.startTime = utils.DateLocal(year, month, day, 9, 2, 0, 0)
		a.endTime = utils.DateLocal(year, month, day, 23, 57, 0, 0)
	case data.MLAFT:
		//第一局开始时间 2018-01-12 13:04 秒开始结算
		//第一局开奖时间 2018-01-12 13:09
		//最后一局结束时间 2018-01-13 04:04
		a.startTime = utils.DateLocal(year, month, day, 13, 4, 0, 0)
		a.endTime = utils.DateLocal(year, month, day, 04, 4, 0, 0)
		am := utils.DateLocal(year, month, day, 0, 0, 0, 0)
		et := utils.Time2Stamp(a.endTime)
		bt := utils.Time2Stamp(am)
		if now < et && now > bt {
			a.startTime = a.startTime.AddDate(0, 0, -1)
		} else {
			a.endTime = a.endTime.AddDate(0, 0, 1)
		}
	}
}

//初始化状态
func (a *DeskActor) initState(now, st int64) {
	//now := utils.Timestamp()
	//st := utils.Time2Stamp(a.startTime)
	if now < (st + 42 + 30) {
		a.state = data.STATE_SEAL
		a.nexttime = st
		return
	}
	//初始时间,已经走的时间[0, 300)
	a.timer = useTime(now, st)
	if a.timer < 42 {
		a.state = data.STATE_SEAL
		a.nexttime = now + int64((42 - a.timer))
	} else if a.timer < (42 + 30) {
		a.state = data.STATE_OVER
		a.nexttime = now + int64((42 + 30 - a.timer))
	} else {
		a.state = data.STATE_BET
		a.nexttime = now + int64((300 - a.timer))
	}
}

//下轮时间
func (a *DeskActor) nextTimes(now int64) bool {
	if now < utils.Time2Stamp(a.endTime) {
		return false
	}
	//最后一场再判断并切换,在23:57-23:59切
	a.startTime = a.startTime.AddDate(0, 0, 1)
	a.endTime = a.endTime.AddDate(0, 0, 1)
	return true
}

//已经开始时间
func useTime(now, st int64) int64 {
	//开始时间之前启动(凌晨 +86400)
	return (now - st + 86400) % 300
}

//TODO 优化处理放假不开奖问题
//状态循环
func (a *DeskActor) handing(ctx actor.Context) {
	now := utils.Timestamp()
	st := utils.Time2Stamp(a.startTime)
	//停盘时间
	if (now < (st + 42 + 30)) && a.state == data.STATE_SEAL {
		return
	}
	state := a.state
	if (now == (st + 42 + 30)) && a.state == data.STATE_SEAL {
		//第一局时切换下注
		a.state = data.STATE_BET
	}
	a.timer = useTime(now, st)
	switch a.state {
	case data.STATE_SEAL:
		if a.timer <= 42 {
			//封盘中
		} else if (a.timer < (42 + 30 + 156)) && a.expect == "" {
			//抓取开奖结果,每三秒抓取一次
			var n int
			if a.timer >= 126 { //时间过半
				n = 1 //启动备用
			}
			if a.setPk10Code(n) {
				//官方开奖记录
				a.logPk10(a.expect, a.opencode, a.opentime, a.opentimestamp)
				//设置上期期号
				a.setLastExpect()
			}
		} else {
			a.state = data.STATE_OVER
			a.nexttime = now + 30
		}
	case data.STATE_OVER:
		if now >= a.nexttime {
			//结算时判断是否继续或封盘
			//检测切换下轮
			if a.nextTimes(now) {
				a.state = data.STATE_SEAL
				a.nexttime = utils.Time2Stamp(a.startTime)
			} else {
				a.state = data.STATE_BET
				a.nexttime = now + 300 - a.timer
			}
		}
	case data.STATE_BET:
		if a.timer == 0 {
			a.state = data.STATE_SEAL
			a.nexttime = now + 42
			//重置
			a.initPk10()
		}
	}
	///切换状态消息
	if a.state != state {
		glog.Debugf("now %s, nexttime %s\n", utils.Unix2Str(now), utils.Unix2Str(a.nexttime))
		glog.Debugf("state %d, timer %d\n", a.state, a.timer)
		//切换状态时更新
		a.pushDeskState()
		glog.Debugf("getPk10 lastexpect %s", a.lastexpect)
	}
	//召唤机器人
	a.robotAllot(now)
}
