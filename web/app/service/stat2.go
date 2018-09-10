package service

import (
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

//抽佣统计
func (this *loggerService) UserStat5() (n1, n2, n3 int64) {
	m1 := bson.M{
		"fee_all": bson.M{"$gt": 0},
	}
	m := bson.M{"$match": m1}
	n := bson.M{
		"$group": bson.M{
			"_id": 1,
			"fee_all": bson.M{
				"$sum": "$fee_all",
			},
		},
	}
	//统计
	operations := []bson.M{m, n}
	result := bson.M{}
	pipe := AgentFees.Pipe(operations)
	err := pipe.One(&result)
	if err != nil {
		beego.Error("UserStat5 fail err: ", err)
		//return
	}
	//代理所得抽佣总数
	if v, ok := result["fee_all"]; ok {
		n1 = v.(int64)
	}
	//今日抽佣
	//endTime := utils.Stamp2Time(utils.TimestampToday())
	endTime := TimeToday4()
	m1["ctime"] = bson.M{"$gte": endTime}
	m = bson.M{"$match": m1}
	//
	operations = []bson.M{m, n}
	result = bson.M{}
	pipe = AgentFees.Pipe(operations)
	err = pipe.One(&result)
	if err != nil {
		beego.Error("UserStat5 fail err: ", err)
		//return
	}
	if v, ok := result["fee_all"]; ok {
		n2 = v.(int64)
	}
	//昨天抽佣
	//startTime := utils.Stamp2Time(utils.TimestampYesterday())
	startTime := TimeYesterday4()
	m1["ctime"] = bson.M{"$gte": startTime, "$lt": endTime}
	m = bson.M{"$match": m1}
	//
	operations = []bson.M{m, n}
	result = bson.M{}
	pipe = AgentFees.Pipe(operations)
	err = pipe.One(&result)
	if err != nil {
		beego.Error("UserStat5 fail err: ", err)
		//return
	}
	if v, ok := result["fee_all"]; ok {
		n3 = v.(int64)
	}
	return
}
