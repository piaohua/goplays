package service

import (
	"time"

	"goplays/web/app/entity"
	"utils"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

//统计游戏内玩家日、周、月赢亏数据展示
func statProfit(dayStamp time.Time) {
	//已经统计过
	m2 := bson.M{"day_stamp": dayStamp}
	if Count(StatRecords, m2) > 0 {
		return
	}
	//结果统计
	yesterdayResult := statYesterdayProfit()
	sevenResult := statSevenProfit()
	thirtyResult := statThirtyProfit()
	beego.Trace("yesterdayResult : ", len(yesterdayResult))
	beego.Trace("sevenResult : ", len(sevenResult))
	beego.Trace("thirtyResult : ", len(thirtyResult))
	//保存数据
	ps := make(map[string]*entity.ProfitStat)
	for _, v1 := range yesterdayResult {
		var userid string
		if val1, ok := v1["_id"]; ok {
			val3 := val1.(bson.M)
			if val2, ok := val3["userid"]; ok {
				userid = string(val2.(string))
				if _, ok := ps[userid]; !ok {
					val := new(entity.ProfitStat)
					val.Userid = userid
					ps[userid] = val
				}
			}
			if val2, ok := val3["robot"]; ok {
				if val, ok := ps[userid]; ok {
					val.Robot = bool(val2.(bool))
				}
			}
		}
		if val1, ok := v1["profits"]; ok {
			if val, ok := ps[userid]; ok {
				val.Yesterday = int64(val1.(int64))
				val.Seven = int64(val1.(int64))
				val.Thirty = int64(val1.(int64))
			}
		}
	}
	for _, v1 := range sevenResult {
		var userid string
		if val1, ok := v1["_id"]; ok {
			val3 := val1.(bson.M)
			if val2, ok := val3["userid"]; ok {
				userid = string(val2.(string))
				if _, ok := ps[userid]; !ok {
					val := new(entity.ProfitStat)
					val.Userid = userid
					ps[userid] = val
				}
			}
			if val2, ok := val3["robot"]; ok {
				if val, ok := ps[userid]; ok {
					val.Robot = bool(val2.(bool))
				}
			}
		}
		if val1, ok := v1["yesterday"]; ok {
			if val, ok := ps[userid]; ok {
				val.Seven += int64(val1.(int64))
			}
		}
	}
	for _, v1 := range thirtyResult {
		var userid string
		if val1, ok := v1["_id"]; ok {
			val3 := val1.(bson.M)
			if val2, ok := val3["userid"]; ok {
				userid = string(val2.(string))
				if _, ok := ps[userid]; !ok {
					val := new(entity.ProfitStat)
					val.Userid = userid
					ps[userid] = val
				}
			}
			if val2, ok := val3["robot"]; ok {
				if val, ok := ps[userid]; ok {
					val.Robot = bool(val2.(bool))
				}
			}
		}
		if val1, ok := v1["yesterday"]; ok {
			if val, ok := ps[userid]; ok {
				val.Thirty += int64(val1.(int64))
			}
		}
	}
	day := utils.Time2DayDate(dayStamp)
	month := utils.Time2MonthDate(dayStamp)
	beego.Trace("statProfit : ", dayStamp, day, month)
	//保存
	for k, v := range ps {
		beego.Trace("statProfit userid : ", k)
		v.Ctime = bson.Now()
		v.Day = day
		v.Month = month
		v.DayStamp = dayStamp
		//每日记录
		if !Insert(StatRecords, v) {
			beego.Error("statProfit fail : ", dayStamp, v)
		} else {
			beego.Trace("statProfit ok : ", dayStamp, v)
		}
		//玩家单个记录,TODO 重复统计 bson.M{"day_stamp": bson.M{"$ne": dayStamp}}
		q := bson.M{"_id": k}
		if Has(UserStatRecords, q) {
			n := bson.M{"yesterday": v.Yesterday, "seven": v.Seven,
				"thirty": v.Thirty, "day_stamp": v.DayStamp,
				"utime": bson.Now()}
			c := bson.M{"all": v.Yesterday}
			i := bson.M{"$set": n, "$inc": c}
			if !Update(UserStatRecords, q, i) {
				us := &entity.UserProfitStat{
					Userid:    v.Userid,
					Robot:     v.Robot,
					Yesterday: v.Yesterday,
					Seven:     v.Seven,
					Thirty:    v.Thirty,
					All:       v.Yesterday,
					DayStamp:  dayStamp,
					Utime:     bson.Now(),
					Ctime:     bson.Now(),
				}
				if !Insert(UserStatRecords, us) {
					beego.Error("statProfit failed : ", k, us)
				} else {
					beego.Trace("statProfit ok : ", k, us)
				}
			} else {
				beego.Trace("statProfit ok : ", k, n)
			}
		} else {
			us := &entity.UserProfitStat{
				Userid:    v.Userid,
				Robot:     v.Robot,
				Yesterday: v.Yesterday,
				Seven:     v.Seven,
				Thirty:    v.Thirty,
				All:       v.Yesterday,
				DayStamp:  dayStamp,
				Utime:     bson.Now(),
				Ctime:     bson.Now(),
			}
			if !Insert(UserStatRecords, us) {
				beego.Error("statProfit failed : ", k, us)
			} else {
				beego.Trace("statProfit ok : ", k, us)
			}
		}
	}
}

//统计昨日赢亏
func statYesterdayProfit() (result []bson.M) {
	//统计时间段
	//startTime := utils.Stamp2Time(utils.TimestampYesterday())
	startTime := TimeYesterday4()
	//endTime := utils.Stamp2Time(utils.TimestampToday())
	endTime := TimeToday4()
	//统计aggregation
	m := bson.M{
		"$match": bson.M{
			"ctime": bson.M{"$gte": startTime, "$lt": endTime},
		},
	}
	n := bson.M{
		"$group": bson.M{
			"_id": bson.M{"userid": "$userid", "robot": "$robot"},
			"profits": bson.M{
				"$sum": "$profits",
			},
		},
	}
	//统计
	operations := []bson.M{m, n}
	result = []bson.M{}
	pipe := UserRecords.Pipe(operations)
	err := pipe.All(&result)
	if err != nil {
		beego.Error("statTodayProfit fail err: ", err)
		return
	}
	return
}

//统计七日赢亏
func statSevenProfit() (result []bson.M) {
	//统计时间段
	//endTime := utils.Stamp2Time(utils.TimestampToday())
	//endTime := utils.Stamp2Time(utils.TimestampYesterday())
	endTime := TimeYesterday4()
	startTime := endTime.AddDate(0, 0, -6)
	//统计aggregation
	m := bson.M{
		"$match": bson.M{
			"ctime": bson.M{"$gte": startTime, "$lt": endTime},
		},
	}
	n := bson.M{
		"$group": bson.M{
			"_id": bson.M{"userid": "$userid", "robot": "$robot"},
			"yesterday": bson.M{
				"$sum": "$yesterday",
			},
		},
	}
	//统计
	operations := []bson.M{m, n}
	result = []bson.M{}
	pipe := StatRecords.Pipe(operations)
	err := pipe.All(&result)
	if err != nil {
		beego.Error("statSevenProfit fail err: ", err)
		return
	}
	return
}

//统计三十日赢亏
func statThirtyProfit() (result []bson.M) {
	//统计时间段
	//endTime := utils.Stamp2Time(utils.TimestampYesterday())
	endTime := TimeYesterday4()
	startTime := endTime.AddDate(0, 0, -29)
	//统计aggregation
	m := bson.M{
		"$match": bson.M{
			"ctime": bson.M{"$gte": startTime, "$lt": endTime},
		},
	}
	n := bson.M{
		"$group": bson.M{
			"_id": bson.M{"userid": "$userid", "robot": "$robot"},
			"yesterday": bson.M{
				"$sum": "$yesterday",
			},
		},
	}
	//统计
	operations := []bson.M{m, n}
	result = []bson.M{}
	pipe := StatRecords.Pipe(operations)
	err := pipe.All(&result)
	if err != nil {
		beego.Error("statThirtyProfit fail err: ", err)
		return
	}
	return
}

//注额统计
func betStat(list []string, startTime, endTime time.Time) (feeNum int64) {
	m1 := bson.M{
		"bets":   bson.M{"$gt": 0},
		"userid": bson.M{"$in": list},
		"ctime":  bson.M{"$gte": startTime, "$lt": endTime},
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
		beego.Error("betStat fail err: ", err)
		return
	}
	//收益
	if v, ok := result["num"]; ok {
		feeNum = v.(int64)
	}
	return
}

//收益统计(所有玩家抽佣产出,包含上级和系统的)
func feeStat(list []string, startTime, endTime time.Time) (feeNum int64) {
	m1 := bson.M{
		"num":    bson.M{"$gt": 0},
		"type":   entity.LogType43, //抽佣类型
		"userid": bson.M{"$in": list},
		"ctime":  bson.M{"$gte": startTime, "$lt": endTime},
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
		beego.Error("feeStat fail err: ", err)
		return
	}
	//收益
	if v, ok := result["num"]; ok {
		feeNum = v.(int64)
	}
	return
}

//统计当前代理抽佣所得,从日志中统计
func feeStatByLog(agent string, startTime, endTime time.Time) (feeNum int64) {
	m1 := bson.M{
		"fee_all": bson.M{"$gt": 0},
		"agent":   agent, //代理
		"ctime":   bson.M{"$gte": startTime, "$lt": endTime},
	}
	m := bson.M{"$match": m1}
	m2 := bson.M{
		"_id": 1,
		"num": bson.M{
			"$sum": "$fee_all",
		},
	}
	n := bson.M{"$group": m2}
	//统计
	operations := []bson.M{m, n}
	result := bson.M{}
	pipe := AgentFeeLogs.Pipe(operations)
	err := pipe.One(&result)
	if err != nil {
		beego.Error("feeStatByLog fail err: ", err)
		return
	}
	//收益
	if v, ok := result["num"]; ok {
		feeNum = v.(int64)
	}
	return
}
