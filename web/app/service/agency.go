package service

import (
	"errors"
	"fmt"
	"time"

	"goplays/web/app/entity"
	"utils"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

type agencyService struct{}

//是否是代理关系 userid != agent
func (this *agencyService) Parental(parent_agent, agent string) bool {
	if agent == "" || parent_agent == "" {
		return false
	}
	q := bson.M{"agent": agent}
	//不存在
	if !Has(Users, q) {
		return false
	}
	//上级
	f := []string{"parent_agent"}
	var result bson.M
	GetByQWithFields(Users, q, f, &result)
	if v, ok := result["parent_agent"]; ok {
		if val, ok := v.(string); ok {
			if val == "" {
				return false
			}
			if val == parent_agent {
				return true
			}
			return this.Parental(parent_agent, val)
		}
	}
	return false
}

//代理等级
func (this *agencyService) AgentLevel(agent string) int {
	if agent == "" {
		return 0
	}
	q := bson.M{"agent": agent}
	//不存在
	if !Has(Users, q) {
		return 0
	}
	return this.AgentLevel2(agent, 1)
}

//上级
func (this *agencyService) AgentLevel2(agent string, level int) int {
	q := bson.M{"agent": agent}
	f := []string{"parent_agent"}
	var result bson.M
	GetByQWithFields(Users, q, f, &result)
	if v, ok := result["parent_agent"]; ok {
		if val, ok := v.(string); ok {
			if val == "" {
				return level
			}
			level += 1
			return this.AgentLevel2(val, level)
		}
	}
	return level
}

//玩家是否是代理
func (this *agencyService) IsAgent(agent string) bool {
	q := bson.M{"agent": agent}
	return Has(Users, q)
}

//获取上级代理
func (this *agencyService) GetFeeAll(agent string) (int64, error) {
	q := bson.M{"agent": agent}
	f := []string{"fee_all"}
	var feeNum bson.M
	GetByQWithFields(Users, q, f, &feeNum)
	if v, ok := feeNum["fee_all"]; ok {
		return v.(int64), nil
	}
	return 0, nil
}

//获取上级代理
func (this *agencyService) GetFeeAll2(agent string) (n1, n2, n3 int64) {
	q := bson.M{"agent": agent}
	f := []string{"fee_all", "fee_top", "fee_extract"}
	var feeNum bson.M
	GetByQWithFields(Users, q, f, &feeNum)
	if v, ok := feeNum["fee_all"]; ok {
		n1 = v.(int64)
	}
	if v, ok := feeNum["fee_top"]; ok {
		n2 = v.(int64)
	}
	if v, ok := feeNum["fee_extract"]; ok {
		n3 = v.(int64)
	}
	return
}

// 获取代理商列表
func (this *agencyService) GetAgencyList(page, pageSize int, m bson.M) ([]entity.User, error) {
	var list []entity.User
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "update_time", false)
	m["agent"] = bson.M{"$ne": ""}
	err := Users.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)

	//转换为分展示
	list = this.agencyList2(list)
	return list, err
}

//转换为分展示
func (this *agencyService) agencyList2(list []entity.User) []entity.User {
	for k, v := range list {
		v.FeeAllf = Chip2Float(int64(v.FeeAll))
		v.FeeTopf = Chip2Float(int64(v.FeeTop))
		v.FeeExtractf = Chip2Float(int64(v.FeeExtract))
		v.FeeRatef = Chip2Float(int64(v.FeeRate))
		v.SysFeeRatef = Chip2Float(int64(v.SysFeeRate))
		list[k] = v
	}
	return list
}

// 获取代理商总数
func (this *agencyService) GetAgencyTotal(m bson.M) (int64, error) {
	return int64(Count(Users, m)), nil
}

// 获取代理商类型
func (this *agencyService) GetAgencyType() ([]int, error) {
	var types []int
	ListByQ(Users, bson.M{"$group": bson.M{"level": "$level"}}, &types)
	return types, nil
}

// 获取绑定我的玩家列表
func (this *agencyService) GetMyAgencyList(username string, page, pageSize int, m bson.M) ([]entity.PlayerUser, error) {
	var list []entity.PlayerUser
	if pageSize == -1 {
		pageSize = 100000
	}
	if username == "" {
		return list, nil
	}
	agency, err := UserService.GetUserByName(username)
	if err != nil {
		return list, err
	}
	if agency.Agent == "" {
		return list, nil
	}
	m["agent"] = agency.Agent
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "create_time", false)
	PlayerUsers.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)

	//转换为分展示
	list = this.chipList2(list)
	return list, nil
}

//转换为分展示
func (this *agencyService) chipList2(list []entity.PlayerUser) []entity.PlayerUser {
	for k, v := range list {
		v.Chipf = Chip2Float(int64(v.Chip))
		list[k] = v
	}
	return list
}

// 获取绑定我的玩家总数
func (this *agencyService) GetMyAgencyTotal(username string, m bson.M) (int64, error) {
	var count int64
	if username == "" {
		return count, nil
	}
	agency, err := UserService.GetUserByName(username)
	if err != nil || agency == nil {
		return count, err
	}
	return int64(agency.Builds), nil
}

// 获取绑定我的玩家列表
func (this *agencyService) GetMyAgencyList2(username string, page, pageSize int, m bson.M) ([]entity.PlayerUser, error) {
	var list []entity.PlayerUser
	if pageSize == -1 {
		pageSize = 100000
	}
	if username == "" {
		return list, nil
	}
	agency, err := UserService.GetUserByAgent(username)
	if err != nil {
		return list, err
	}
	if agency.Agent == "" {
		return list, nil
	}
	m["agent"] = agency.Agent
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "create_time", false)
	PlayerUsers.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)

	//转换为分展示
	list = this.chipList2(list)
	return list, nil
}

// 获取绑定我的玩家总数
func (this *agencyService) GetMyAgencyTotal2(username string, m bson.M) (int64, error) {
	var count int64
	if username == "" {
		return count, nil
	}
	agency, err := UserService.GetUserByAgent(username)
	if err != nil || agency == nil {
		return count, err
	}
	return int64(agency.Builds), nil
}

// 获取绑定我的代理
func (this *agencyService) GetMyAgencyList3(id string, page, pageSize int, m bson.M) ([]entity.User, error) {
	var list []entity.User
	if pageSize == -1 {
		pageSize = 100000
	}
	if id == "" {
		return list, nil
	}
	m["parent_agent"] = id
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "create_time", false)
	Users.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)

	//转换为分展示
	list = this.agencyList2(list)
	return list, nil
}

// 获取绑定我的代理
func (this *agencyService) GetMyAgencyTotal3(id string, m bson.M) (int64, error) {
	if id == "" {
		return 0, nil
	}
	m["parent_agent"] = id
	return int64(Count(Users, m)), nil
}

// 获取绑定我的玩家列表
func (this *agencyService) GetMyAgencyList4(id string, page, pageSize int, m bson.M) ([]entity.PlayerUser, error) {
	var list []entity.PlayerUser
	if pageSize == -1 {
		pageSize = 100000
	}
	if id == "" {
		return list, nil
	}
	m["agent"] = id
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	PlayerUsers.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)

	//转换为分展示
	list = this.chipList2(list)
	return list, nil
}

// 获取绑定我的玩家总数
func (this *agencyService) GetMyAgencyTotal4(id string, m bson.M) (int64, error) {
	if id == "" {
		return 0, nil
	}
	m["agent"] = id
	return int64(Count(PlayerUsers, m)), nil
}

// 获取全部提现记录
func (this *agencyService) GetCashList(page, pageSize int, m bson.M) ([]entity.ApplyCash, error) {
	var list []entity.ApplyCash
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)

	ApplyCashs.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)

	//转换为分展示
	list = this.cashList2(list)
	return list, nil
}

//转换为分展示
func (this *agencyService) cashList2(list []entity.ApplyCash) []entity.ApplyCash {
	for k, v := range list {
		v.Feef = Chip2Float(int64(v.Fee))
		v.Restf = Chip2Float(int64(v.Rest))
		list[k] = v
	}
	return list
}

// 获取全部提现记录总数
func (this *agencyService) GetCashListTotal(m bson.M) (int64, error) {
	var count int64
	count = int64(Count(ApplyCashs, m))
	return count, nil
}

// 获取我的提现记录
func (this *agencyService) GetMyCashList(username string, page, pageSize int, m bson.M) ([]entity.ApplyCash, error) {
	var list []entity.ApplyCash
	if pageSize == -1 {
		pageSize = 100000
	}
	if username == "" {
		return list, nil
	}
	agency, err := UserService.GetUserByName(username)
	if err != nil {
		return list, err
	}
	if agency.Agent == "" {
		return list, nil
	}

	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	m["agent"] = agency.Agent

	ApplyCashs.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)

	//转换为分展示
	list = this.cashList2(list)
	return list, nil
}

// 获取我的提现记录总数
func (this *agencyService) GetMyCashListTotal(username string, m bson.M) (int64, error) {
	var count int64
	if username == "" {
		return count, nil
	}
	agency, err := UserService.GetUserByName(username)
	if err != nil {
		return count, err
	}
	if agency.Agent == "" {
		return count, nil
	}
	m["agent"] = agency.Agent
	count = int64(Count(ApplyCashs, m))
	return count, nil
}

// 获取我的总的已提现金额
func (this *agencyService) GetMyExtractTotal(username string) (float32, error) {
	var count float32
	if username == "" {
		return count, nil
	}
	agency, err := UserService.GetUserByName(username)
	if err != nil {
		return count, err
	}
	if agency.Agent == "" {
		return count, nil
	}
	return agency.Extract, nil
}

//代理获取可以提现金额
func (this *agencyService) GetMyCashTotal(username string) (float32, error) {
	if username == "" {
		return 0, nil
	}
	agency, err := UserService.GetUserByName(username)
	if err != nil {
		return 0, err
	}
	if agency.Agent == "" {
		return 0, errors.New("代理不存在")
	}
	return agency.Cash, nil
}

// 申请提现,注册代理后7天才能提现
func (this *agencyService) ApplyCashAdd(username, name, bankAddr string,
	bankCard int, fee int) error {
	agency, err := UserService.GetUserByName(username)
	if agency.Agent == "" || err != nil {
		return errors.New("代理商不存在")
	}
	if int64(fee) > agency.FeeAll {
		return errors.New("金额不足")
	}
	if fee <= 0 {
		return errors.New("金额错误")
	}
	//if cash > float64(agency.Cash) {
	//	return errors.New("金额不足")
	//}
	//now := utils.LocalTime().Unix()
	//if (now-agency.CreateTime.Unix()) < 7*86400 ||
	//	agency.CreateTime.IsZero() {
	//	return errors.New("成为代理7天后方可提现")
	//}
	applyCash := new(entity.ApplyCash)
	applyCash.Id = bson.NewObjectId().Hex()
	applyCash.Agent = agency.Agent
	//applyCash.Cash = float32(cash)
	applyCash.Fee = int64(fee)
	applyCash.Rest = (agency.FeeAll - int64(fee))
	applyCash.Status = 1 //表示等待处理
	applyCash.RealName = name
	applyCash.BankCard = bankCard
	applyCash.BankAddr = bankAddr
	applyCash.Ctime = bson.Now()
	if Insert(ApplyCashs, applyCash) {
		err1 := this.updateFeeAll(agency.Id, (-1 * applyCash.Fee))
		if err1 != nil {
			Delete(ApplyCashs, bson.M{"_id": applyCash.Id})
			return errors.New("提取失败")
		}
		return nil
	}
	return errors.New("提取失败")
}

// 提现处理
func (this *agencyService) ExtractCash(username, orderid string) error {
	agency, err := UserService.GetUserByName(username)
	if err != nil || agency == nil {
		return err
	}
	//if agency.Agent == "" {
	//	return errors.New("代理ID不存在")
	//}
	m := bson.M{"_id": orderid}
	n := bson.M{"status": 0, "user_name": username, "utime": bson.Now()}
	//更新状态为已经处理
	if Update(ApplyCashs, m, bson.M{"$set": n}) {
		return nil
	}
	return errors.New("提现失败")
}

// 更新账户金额
func (this *agencyService) UpdateNumber(agency *entity.User) error {
	m := bson.M{"_id": agency.Id}
	n := bson.M{"number": agency.Number, "update_time": bson.Now()}
	if Update(Users, m, bson.M{"$set": n}) {
		return nil
	}
	return errors.New("更新失败")
}

// 更新账户金额
func (this *agencyService) SetChip(id string, chip int64) error {
	m := bson.M{"_id": id}
	n := bson.M{"$inc": bson.M{"chip": chip},
		"$set": bson.M{"update_time": bson.Now()}}
	if Update(Users, m, n) {
		return nil
	}
	return errors.New("更新失败")
}

// 更新账户消耗金额(赠送时更新)
func (this *agencyService) UpdateExpend(agency *entity.User) error {
	m := bson.M{"_id": agency.Id}
	n := bson.M{"number": agency.Number, "expend": agency.Expend, "update_time": bson.Now()}
	if Update(Users, m, bson.M{"$set": n}) {
		return nil
	}
	return errors.New("更新失败")
}

// 更新提现率
func (this *agencyService) UpdateAgencyRate2(agency *entity.User) error {
	m := bson.M{"user_name": agency.UserName, "agent": agency.Agent}
	n := bson.M{"sys_rate": agency.SysRate, "rate": agency.Rate, "update_time": bson.Now()}
	if Update(Users, m, bson.M{"$set": n}) {
		return nil
	}
	return errors.New("更新失败")
}

// 更新提现率
func (this *agencyService) UpdateAgencyRate(agency *entity.User) error {
	m := bson.M{"user_name": agency.UserName, "agent": agency.Agent}
	n := bson.M{"rate": agency.Rate, "update_time": bson.Now()}
	if Update(Users, m, bson.M{"$set": n}) {
		return nil
	}
	return errors.New("更新失败")
}

// 更新提现金额
func (this *agencyService) updateFeeAll(id string, fee int64) error {
	m := bson.M{"_id": id}
	n := bson.M{"update_time": bson.Now()}
	c := bson.M{"fee_all": fee, "fee_extract": (-1 * fee)}
	if Update(Users, m, bson.M{"$set": n, "$inc": c}) {
		return nil
	}
	return errors.New("更新失败")
}

// 更新提现金额
func (this *agencyService) updateCash(id string, cash float32) error {
	m := bson.M{"_id": id}
	n := bson.M{"update_time": bson.Now()}
	c := bson.M{"cash": cash}
	if Update(Users, m, bson.M{"$set": n, "$inc": c}) {
		return nil
	}
	return errors.New("更新失败")
}

// 更新提现金额
func (this *agencyService) updateAgencyCash(agency *entity.User, cash float32) error {
	m := bson.M{"_id": agency.Id}
	n := bson.M{"cash_time": agency.CashTime, "update_time": bson.Now()}
	c := bson.M{"cash": cash}
	if Update(Users, m, bson.M{"$set": n, "$inc": c}) {
		return nil
	}
	return errors.New("更新失败")
}

// 更新提现金额
func (this *agencyService) updateAtypeCash(atype uint32, cash float32) error {
	m := bson.M{"atype": atype}
	n := bson.M{"update_time": bson.Now()}
	c := bson.M{"cash": cash}
	if Update(Users, m, bson.M{"$set": n, "$inc": c}) {
		return nil
	}
	return errors.New("更新失败")
}

// 定时统计
func (this *agencyService) stat() {
	//TODO 优化
	list, err := this.GetAgencyList(1, -1, bson.M{})
	if err != nil {
		beego.Error("stat err: ", err)
	}
	for _, v := range list {
		this.statBuilds(&v) //绑定更新
		this.statCash(&v)   //提现更新
		feeStat1(&v)        //收益统计
		playersStat(&v)     //统计玩家
	}
}

// 获取我的总的可提现金额,代理抽成50%,代理如果属于分包,包主抽成15%
func (this *agencyService) statCash(agency *entity.User) {
	if agency.Agent == "" {
		return
	}
	endTime := utils.LocalTime()
	startTime := agency.CashTime
	//统计属于代理的
	money, err2 := this.statAgentCash(startTime, endTime, agency.Agent)
	if err2 != nil {
		beego.Error("statCash err2: ", agency.Id, err2)
		return
	}
	if agency.Rate > 0 {
		money = money * float32(agency.Rate) / 100 //返现60%
	} else {
		money = money * 0.6 //返现60%
	}
	//包主抽成
	if agency.Belong != 0 && money > 0 { //属于分包
		//分成给包主15%
		money1 := money * 0.05 //抽成15%
		money1 = float32(utils.Float64(fmt.Sprintf("%.2f", money1)))
		if money1 != 0 {
			err := this.updateAtypeCash(agency.Belong, money1)
			if err != nil {
				beego.Error("statCash Belong err: ", agency.Id, agency.Belong, err)
				return
			}
		}
		beego.Trace("statCash Belong ok: ", agency.Id, agency.Belong, money1)
		//自己分到85%
		money = money - money1 //提取金额(单位:分)
	} else if agency.Atype != 0 { //是包主
		//统计分包没绑定订单
		money2, err2 := this.statAtypeCash(startTime, endTime, agency.Atype)
		if err2 != nil {
			beego.Error("statCash Atype err2: ", agency.Id, agency.Atype, money, err2)
			return
		}
		money2 = money2 * 0.5 //返现50%
		money += money2       //提取金额(单位:分)
	}
	money = float32(utils.Float64(fmt.Sprintf("%.2f", money)))
	if money == 0 {
		return
	}
	//存在数据中
	agency.CashTime = endTime //截至统计时间
	err1 := this.updateAgencyCash(agency, money)
	if err1 != nil {
		beego.Error("statCash err1: ", agency.Id, err1)
		return
	}
	beego.Trace("statCash ok: ", agency.Id, money)
}

// 获取绑定我的玩家总数
func (this *agencyService) statBuilds(agency *entity.User) {
	if agency.Agent == "" {
		return
	}
	endTime := utils.LocalTime()
	startTime := agency.BuildsTime
	m := bson.M{"agent": agency.Agent}
	m["atime"] = bson.M{"$gte": startTime, "$lt": endTime}
	count, _ := PlayerService.GetTotal(0, m)
	//TODO 优化统计
	if count == 0 {
		return
	}
	// 更新绑定人数
	q := bson.M{"_id": agency.Id}
	n := bson.M{"builds_time": endTime}
	c := bson.M{"builds": uint32(count)}
	if !Update(Users, q, bson.M{"$set": n, "$inc": c}) {
		beego.Error("statBuilds err: ", agency.Id, count)
		return
	}
	beego.Trace("statBuilds ok: ", agency.Id, count)
}

//统计属于代理的
func (this *agencyService) statAgentCash(startTime, endTime time.Time, agent string) (float32, error) {
	m := bson.M{
		"$match": bson.M{
			"agent":  agent,
			"result": entity.TradeSuccess,
			"ctime":  bson.M{"$gte": startTime, "$lt": endTime},
		},
	}
	n := bson.M{
		"$group": bson.M{
			"_id": "$agent",
			"money": bson.M{
				"$sum": "$money",
			},
		},
	}
	operations := []bson.M{m, n}
	result := bson.M{}
	pipe := TradeRecords.Pipe(operations)
	err := pipe.One(&result)
	if err != nil {
		if err.Error() == "not found" {
			return 0, nil
		}
		return 0, err
	}
	if v, ok := result["money"]; ok {
		return float32(v.(int)), nil
	}
	return 0, nil
}

//不属于代理，但是属于分包
func (this *agencyService) statAtypeCash(startTime, endTime time.Time, atype uint32) (float32, error) {
	m := bson.M{
		"$match": bson.M{
			"agent":  "",
			"atype":  atype,
			"result": entity.TradeSuccess,
			"ctime":  bson.M{"$gte": startTime, "$lt": endTime},
		},
	}
	n := bson.M{
		"$group": bson.M{
			"_id": "$agent",
			"money": bson.M{
				"$sum": "$money",
			},
		},
	}
	operations := []bson.M{m, n}
	result := bson.M{}
	pipe := TradeRecords.Pipe(operations)
	err := pipe.One(&result)
	if err != nil {
		if err.Error() == "not found" {
			return 0, nil
		}
		return 0, err
	}
	if v, ok := result["money"]; ok {
		return float32(v.(int)), nil //返现50%
	}
	return 0, nil
}

//db.col_trade_record.aggregate([{$group : {_id : "$agent", num_tutorial : {$sum : "$money"}}}])
//db.t_user.update({'agent':'28405'},{$set:{'cash':26640}})

//统计

//统计昨日注册人数
func (this *agencyService) statRegToday(dayStamp time.Time) {
	//dayStamp = utils.Stamp2Time(utils.TimestampYesterday())
	m := bson.M{"day_stamp": dayStamp}
	if Count(LogRegistTodays, m) > 0 {
		return
	}
	num := Count(RegistLogs, m)
	reg := new(entity.LogRegistToday)
	reg.Id = bson.NewObjectIdWithTime(dayStamp).Hex()
	reg.Num = num
	reg.DayStamp = dayStamp
	reg.Ctime = bson.Now()
	if !Insert(LogRegistTodays, reg) {
		beego.Error("statRegToday fail reg: ", dayStamp, num)
	}
	beego.Trace("statRegToday ok reg: ", dayStamp, num)
}

//统计昨日充值人数
func (this *agencyService) statPayToday(dayStamp time.Time) {
	m2 := bson.M{"day_stamp": dayStamp}
	if Count(LogPayTodays, m2) > 0 {
		return
	}
	//startTime := utils.Stamp2Time(utils.TimestampYesterday())
	startTime := TimeYesterday4()
	//endTime := utils.Stamp2Time(utils.TimestampToday())
	endTime := TimeToday4()
	//
	m := bson.M{
		"$match": bson.M{
			//"day_stamp": dayStamp,
			"result": bson.M{"$ne": entity.Tradeing},
			"ctime":  bson.M{"$gte": startTime, "$lt": endTime},
		},
	}
	n := bson.M{
		"$group": bson.M{
			"_id": "$day_stamp",
			"num": bson.M{
				"$sum": 1,
			},
			"money": bson.M{
				"$sum": "$money",
			},
			"diamond": bson.M{
				"$sum": "$diamond",
			},
		},
	}
	//
	operations := []bson.M{m, n}
	result := bson.M{}
	pipe := TradeRecords.Pipe(operations)
	err := pipe.One(&result)
	if err != nil {
		beego.Error("statPayToday fail err: ", dayStamp, err)
	}
	//
	h := bson.M{
		"$group": bson.M{
			"_id": "$userid",
			"num": bson.M{
				"$sum": 1,
			},
		},
	}
	operations2 := []bson.M{m, h}
	result2 := []bson.M{}
	pipe2 := TradeRecords.Pipe(operations2)
	err2 := pipe2.All(&result2)
	if err2 != nil {
		beego.Error("statPayToday fail err2: ", dayStamp, err2)
	}
	beego.Trace("statPayToday fail err3: ", result2)
	//
	pay := new(entity.LogPayToday)
	if v, ok := result["num"]; ok {
		pay.Num = int(v.(int))
	}
	if v, ok := result["money"]; ok {
		pay.Money = uint32(v.(int))
	}
	if v, ok := result["diamond"]; ok {
		pay.Diamond = uint32(v.(int))
	}
	pay.Count = len(result2)
	pay.Id = bson.NewObjectIdWithTime(dayStamp).Hex()
	pay.DayStamp = dayStamp
	pay.Ctime = bson.Now()
	if !Insert(LogPayTodays, pay) {
		beego.Error("statPayToday fail pay: ", dayStamp, pay)
	}
	beego.Trace("statPayToday ok pay: ", dayStamp, pay)
}

//统计昨日盈亏
func (this *agencyService) statChipToday(dayStamp time.Time) {
	//dayStamp = utils.Stamp2Time(utils.TimestampYesterday())
	//now := time.Now()
	//dayStamp = time.Date(now.Year(), now.Month()-1, 22, 04, 30, 0, 0, time.Local)
	m := bson.M{"day_stamp": dayStamp}
	if Count(LogChipTodays, m) > 0 {
		return
	}
	//startTime := utils.Stamp2Time(utils.TimestampYesterday())
	startTime := TimeYesterday4()
	//startTime := time.Date(now.Year(), now.Month()-1, 22, 04, 30, 0, 0, time.Local)
	//endTime := utils.Stamp2Time(utils.TimestampToday())
	endTime := TimeToday4()
	//endTime := time.Date(now.Year(), now.Month()-1, 23, 04, 30, 0, 0, time.Local)
	beego.Trace("statChipToday dayStamp : ", dayStamp.String())
	beego.Trace("statChipToday startTime : ", startTime.String())
	beego.Trace("statChipToday endTime : ", endTime.String())
	result := this.statProfitsToday(startTime, endTime)
	beego.Trace("statChipToday result : ", result)
	//
	list := make(map[string]*entity.LogChipToday)
	for _, v := range result {
		if val2, ok := v["_id"]; ok {
			val := val2.(bson.M)
			//gametype := val["gametype"].(int)
			//roomtype := val["roomtype"].(int)
			//lotterytype := val["lotterytype"].(int)
			//robot := val["robot"].(bool)
			//num := v["profits"].(int64)
			//
			var gametype, roomtype, lotterytype int
			var robot bool
			var num int64
			if val3, ok := val["gametype"]; ok {
				gametype = val3.(int)
			}
			if val3, ok := val["roomtype"]; ok {
				roomtype = val3.(int)
			}
			if val3, ok := val["lotterytype"]; ok {
				lotterytype = val3.(int)
			}
			if val3, ok := val["robot"]; ok {
				robot = val3.(bool)
			}
			if val3, ok := v["profits"]; ok {
				num = val3.(int64)
			}
			//
			str := utils.String(gametype) + utils.String(roomtype) + utils.String(lotterytype)
			if val1, ok := list[str]; ok {
				if robot {
					val1.RobotNum = num
				} else {
					val1.RolesNum = num
				}
				list[str] = val1
			} else {
				chipLog := new(entity.LogChipToday)
				chipLog.Id = bson.NewObjectId().Hex()
				chipLog.Gametype = gametype
				chipLog.Roomtype = roomtype
				chipLog.Lotterytype = lotterytype
				chipLog.DayStamp = dayStamp
				chipLog.Ctime = bson.Now()
				if robot {
					chipLog.RobotNum = num
				} else {
					chipLog.RolesNum = num
				}
				list[str] = chipLog
			}
		}
	}
	beego.Trace("statChipToday list : ", list)
	for k, v := range list {
		if !Insert(LogChipTodays, v) {
			beego.Error("statChipToday fail chipLog: ", dayStamp, k)
		}
		beego.Trace("statChipToday ok chipLog: ", dayStamp, k)
	}
}

//统计昨日盈亏
func (this *agencyService) statProfitsToday(startTime, endTime time.Time) (result []bson.M) {
	m := bson.M{
		"$match": bson.M{
			"ctime": bson.M{"$gte": startTime, "$lt": endTime},
		},
	}
	n := bson.M{
		"$group": bson.M{
			"_id": bson.M{"gametype": "$gametype",
				"roomtype": "$roomtype", "lotterytype": "$lotterytype",
				"robot": "$robot"},
			"profits": bson.M{
				"$sum": "$profits",
			},
		},
	}
	operations := []bson.M{m, n}
	result = []bson.M{}
	pipe := UserRecords.Pipe(operations)
	err := pipe.All(&result)
	if err != nil {
		beego.Error("statProfitsToday fail err: ", err)
		return
	}
	beego.Trace("statProfitsToday : ", result)
	return
}
