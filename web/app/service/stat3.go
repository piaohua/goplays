package service

import (
	"time"

	"goplays/web/app/entity"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

//账务统计
func statAccountingLog(dayStamp time.Time) {
	//已经统计过
	m2 := bson.M{"day_stamp": dayStamp}
	if Count(AccountingLogs, m2) > 0 {
		return
	}
	//
	a := new(entity.AccountingLog)
	a.DayStamp = dayStamp
	aStatChips(a)
	aStatBets(a)
	aStatPays(a)
	aStatFeeAll(a)
	aStatAgentFee(a)
	aRobotProfits(a)
	aSysProfits(a)
	if !aSave(a) {
		beego.Error("statAccountingLog save failed : ", a)
	} else {
		beego.Trace("statAccountingLog save success : ", a)
	}
}

// 写入数据
func aSave(this *entity.AccountingLog) bool {
	this.Ctime = bson.Now()
	return Insert(AccountingLogs, this)
}

//玩家当天剩余总筹码
func aStatChips(this *entity.AccountingLog) {
	m := bson.M{
		"$match": bson.M{
			"phone": bson.M{"$ne": ""},
			"robot": false,
		},
	}
	n := bson.M{
		"$group": bson.M{
			"_id": 1,
			"chip": bson.M{
				"$sum": "$chip",
			},
		},
	}
	//统计
	operations := []bson.M{m, n}
	result := bson.M{}
	pipe := PlayerUsers.Pipe(operations)
	err := pipe.One(&result)
	if err != nil {
		beego.Error("aStatChips fail err: ", err)
		return
	}
	//用户持有筹码总数
	if v, ok := result["chip"]; ok {
		this.Chips = v.(int64)
	}
}

//玩家当天总投注额
func aStatBets(this *entity.AccountingLog) {
	startTime := TimeYesterday4()
	endTime := TimeToday4()
	m1 := bson.M{
		"bets":  bson.M{"$gt": 0},
		"robot": false,
		"ctime": bson.M{"$gte": startTime, "$lt": endTime},
	}
	m := bson.M{"$match": m1}
	m2 := bson.M{
		"_id": 1,
		"num": bson.M{
			"$sum": "$bets",
		},
	}
	n := bson.M{"$group": m2}
	//统计
	operations := []bson.M{m, n}
	result := bson.M{}
	pipe := UserRecords.Pipe(operations)
	err := pipe.One(&result)
	if err != nil {
		beego.Error("aStatBets fail err: ", err)
		return
	}
	//收益
	if v, ok := result["num"]; ok {
		this.Bets = v.(int64)
	}
}

//玩家当天总充值
func aStatPays(this *entity.AccountingLog) {
	startTime := TimeYesterday4()
	endTime := TimeToday4()
	m1 := bson.M{
		"type":  entity.LogType9, //后台操作筹码
		"ctime": bson.M{"$gte": startTime, "$lt": endTime},
	}
	m := bson.M{"$match": m1}
	m2 := bson.M{
		"_id": 1,
		"num": bson.M{
			"$sum": "$num",
		},
	}
	n := bson.M{"$group": m2}
	//统计
	operations := []bson.M{m, n}
	result := bson.M{}
	pipe := ChipLogs.Pipe(operations)
	err := pipe.One(&result)
	if err != nil {
		beego.Error("aStatPays fail err: ", err)
		return
	}
	//注额
	if v, ok := result["num"]; ok {
		if val, ok := v.(int64); ok {
			this.Pays = val
		}
	}
}

//抽佣总数
func aStatFeeAll(this *entity.AccountingLog) {
	//抽佣统计
	n1, n2, n3 := LoggerService.UserStat3()
	beego.Error("aStatFeeAll n1, n2, n3: ", n1, n2, n3)
	this.AllFee = n1
	this.YesterdayFee = n3
}

//代理抽佣总数
func aStatAgentFee(this *entity.AccountingLog) {
	//抽佣统计
	n1, n2, n3 := LoggerService.UserStat5()
	beego.Error("aStatAgentFee n1, n2, n3: ", n1, n2, n3)
	this.AgentAllFee = n1
	this.YesterdayAgentFee = n3
}

//机器人昨日盈亏
func aRobotProfits(this *entity.AccountingLog) {
	this.RobotProfitsYesterday = LoggerService.getRobotProfitsYesterday()
}

//系统盈亏
func aSysProfits(this *entity.AccountingLog) {
	//系统盈亏 = 当天总抽佣 - 当天代理抽佣 + 当天机器人盈亏
	this.SysProfitsYesterday = this.YesterdayFee - this.YesterdayAgentFee + this.RobotProfitsYesterday
}
