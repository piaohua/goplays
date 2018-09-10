package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/astaxie/beego"

	"goplays/web/app/entity"
	"utils"

	"gopkg.in/mgo.v2/bson"
)

type loggerService struct{}

// 获取注册日志列表
func (this *loggerService) SaveAgentFee(f *entity.AgentFee) bool {
	f.Ctime = bson.Now()
	return Insert(AgentFees, f)
}

// 获取注册日志列表
func (this *loggerService) SaveAgentFeeLog(f *entity.AgentFeeLog) bool {
	f.Ctime = bson.Now()
	return Insert(AgentFeeLogs, f)
}

// 账务统计日志列表
func (this *loggerService) GetAccountingList(page, pageSize int, m bson.M) ([]entity.AccountingLog, error) {
	var list []entity.AccountingLog
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	AccountingLogs.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)

	//转换为分展示
	list = this.accountingList2(list)
	return list, nil
}

//转换为分展示
func (this *loggerService) accountingList2(list []entity.AccountingLog) []entity.AccountingLog {
	for k, v := range list {
		v.Chipsf = Chip2Float(int64(v.Chips))                                 //玩家当天剩余总筹码
		v.Betsf = Chip2Float(int64(v.Bets))                                   //玩家当天总投注额
		v.Paysf = Chip2Float(int64(v.Pays))                                   //玩家当天总充值
		v.AllFeef = Chip2Float(int64(v.AllFee))                               //抽佣总数
		v.AgentAllFeef = Chip2Float(int64(v.AgentAllFee))                     //代理抽佣总数
		v.YesterdayFeef = Chip2Float(int64(v.YesterdayFee))                   //昨天总抽佣
		v.YesterdayAgentFeef = Chip2Float(int64(v.YesterdayAgentFee))         //昨天代理抽佣
		v.RobotProfitsYesterdayf = Chip2Float(int64(v.RobotProfitsYesterday)) //机器人昨日盈亏
		v.SysProfitsYesterdayf = Chip2Float(int64(v.SysProfitsYesterday))     //系统盈亏 = 当天总抽佣 - 当天代理抽佣 + 当天机器人盈亏
		list[k] = v
	}
	return list
}

// 账务统计日志总数
func (this *loggerService) GetAccountingTotal(m bson.M) (int64, error) {
	return int64(Count(AccountingLogs, m)), nil
}

// 获取注册日志列表
func (this *loggerService) GetRegistList(page, pageSize int, m bson.M) ([]entity.LogRegist, error) {
	var list []entity.LogRegist
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	RegistLogs.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, nil
}

// 获取注册日志总数
func (this *loggerService) GetRegistTotal(m bson.M) (int64, error) {
	return int64(Count(RegistLogs, m)), nil
}

// 获取登录日志列表
func (this *loggerService) GetLoginList(page, pageSize int, m bson.M) ([]entity.LogLogin, error) {
	var list []entity.LogLogin
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "login_time", false)
	LoginLogs.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, nil
}

// 获取登录日志总数
func (this *loggerService) GetLoginTotal(m bson.M) (int64, error) {
	return int64(Count(LoginLogs, m)), nil
}

// 获取充值列表
func (this *loggerService) GetPayList(page, pageSize int, m bson.M) ([]entity.TradeRecord, error) {
	var list []entity.TradeRecord
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	if _, ok := m["result"]; !ok {
		m["result"] = bson.M{"$ne": entity.Tradeing}
	}
	TradeRecords.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, nil
}

// 获取绑定代理充值数量
func (this *loggerService) GetAgencyPayTotal(username string) (int64, error) {
	var count int64
	if username == "" {
		return count, nil
	}
	agency, err := UserService.GetUserByName(username)
	if err != nil {
		return count, err
	}
	if agency.Agent == "" {
		return count, errors.New("代理不存在")
	}
	return int64(Count(TradeRecords, bson.M{"agent": agency.Agent, "result": entity.TradeSuccess})), nil
}

// 获取充值总数
func (this *loggerService) GetPayTotal(m bson.M) (int64, error) {
	var count int64
	if _, ok := m["result"]; !ok {
		m["result"] = bson.M{"$ne": entity.Tradeing}
	}
	count = int64(Count(TradeRecords, m))
	return count, nil
}

// 获取钻石日志列表
func (this *loggerService) GetDiamondList(page, pageSize int, m bson.M) ([]entity.LogDiamond, error) {
	var list []entity.LogDiamond
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	DiamondLogs.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, nil
}

// 获取注册日志总数
func (this *loggerService) GetDiamondTotal(m bson.M) (int64, error) {
	return int64(Count(DiamondLogs, m)), nil
}

// 获取金币日志列表
func (this *loggerService) GetCoinList(page, pageSize int, m bson.M) ([]entity.LogCoin, error) {
	var list []entity.LogCoin
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	CoinLogs.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, nil
}

// 获取金币日志总数
func (this *loggerService) GetCoinTotal(m bson.M) (int64, error) {
	return int64(Count(CoinLogs, m)), nil
}

// 获取筹码日志列表
func (this *loggerService) GetChipList(page, pageSize int, m bson.M) ([]entity.LogChip, error) {
	var list []entity.LogChip
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	ChipLogs.
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
func (this *loggerService) chipList2(list []entity.LogChip) []entity.LogChip {
	for k, v := range list {
		v.Numf = Chip2Float(int64(v.Num))
		v.Restf = Chip2Float(int64(v.Rest))
		list[k] = v
	}
	return list
}

// 获取筹码日志总数
func (this *loggerService) GetChipTotal(m bson.M) (int64, error) {
	return int64(Count(ChipLogs, m)), nil
}

// 获取绑定日志列表
func (this *loggerService) GetBuildList(page, pageSize int, m bson.M) ([]entity.LogBuildAgency, error) {
	var list []entity.LogBuildAgency
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	LogBuildAgencys.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, nil
}

// 获取注册日志总数
func (this *loggerService) GetBuildTotal(m bson.M) (int64, error) {
	return int64(Count(LogBuildAgencys, m)), nil
}

// 绑定统计
func (this *loggerService) GetPubStat(username, rangeType string) map[int]int {
	result := make(map[int]int)

	if username == "" {
		return result
	}
	agency, err := UserService.GetUserByName(username)
	if err != nil {
		return result
	}
	if agency.Agent == "" {
		return result
	}

	m := bson.M{}
	var n bson.M

	switch rangeType {
	case "this_month":
		year, month, _ := time.Now().Date()
		startTimeStr := fmt.Sprintf("%d-%02d-01 00:00:00", year, month)
		endTimeStr := fmt.Sprintf("%d-%02d-31 23:59:59", year, month)
		startTime, _ := utils.Str2Unix(startTimeStr)
		endTime, _ := utils.Str2Unix(endTimeStr)
		m = bson.M{
			"$match": bson.M{
				"agent":     agency.Agent,
				"day_stamp": bson.M{"$gte": startTime, "$lte": endTime},
			},
		}
		n = bson.M{
			"$group": bson.M{
				"date": "$day",
				"count": bson.M{
					"$sum": 1,
				},
			},
		}
	case "last_month":
		year, month, _ := time.Now().AddDate(0, -1, 0).Date()
		startTimeStr := fmt.Sprintf("%d-%02d-01 00:00:00", year, month)
		endTimeStr := fmt.Sprintf("%d-%02d-31 23:59:59", year, month)
		startTime, _ := utils.Str2Unix(startTimeStr)
		endTime, _ := utils.Str2Unix(endTimeStr)
		m = bson.M{
			"$match": bson.M{
				"agent":     agency.Agent,
				"day_stamp": bson.M{"$gte": startTime, "$lte": endTime},
			},
		}
		n = bson.M{
			"$group": bson.M{
				"date": "$day",
				"count": bson.M{
					"$sum": 1,
				},
			},
		}
	case "this_year":
		year := time.Now().Year()
		startTimeStr := fmt.Sprintf("%d-01-01 00:00:00", year)
		endTimeStr := fmt.Sprintf("%d-12-31 23:59:59", year)
		startTime, _ := utils.Str2Unix(startTimeStr)
		endTime, _ := utils.Str2Unix(endTimeStr)
		m = bson.M{
			"$match": bson.M{
				"agent":     agency.Agent,
				"day_stamp": bson.M{"$gte": startTime, "$lte": endTime},
			},
		}
		n = bson.M{
			"$group": bson.M{
				"date": "$month",
				"count": bson.M{
					"$sum": 1,
				},
			},
		}
	case "last_year":
		year := time.Now().Year() - 1
		startTimeStr := fmt.Sprintf("%d-01-01 00:00:00", year)
		endTimeStr := fmt.Sprintf("%d-12-31 23:59:59", year)
		startTime, _ := utils.Str2Unix(startTimeStr)
		endTime, _ := utils.Str2Unix(endTimeStr)
		m = bson.M{
			"$match": bson.M{
				"agent":     agency.Agent,
				"day_stamp": bson.M{"$gte": startTime, "$lte": endTime},
			},
		}
		n = bson.M{
			"$group": bson.M{
				"date": "$month",
				"count": bson.M{
					"$sum": 1,
				},
			},
		}
	}

	operations := []bson.M{m, n}
	maps := []bson.M{}
	pipe := LogBuildAgencys.Pipe(operations)
	err1 := pipe.All(&maps)

	if err1 == nil && len(maps) > 0 {
		for _, v := range maps {
			date, _ := utils.Str2Int(v["date"].(string))
			count, _ := utils.Str2Int(v["count"].(string))
			result[date] = count
		}
	}
	return result
}

// 获取日志列表
func (this *loggerService) GetRegTodayList(page, pageSize int, m bson.M) ([]entity.LogRegistToday, error) {
	var list []entity.LogRegistToday
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	LogRegistTodays.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, nil
}

// 获取日志总数
func (this *loggerService) GetRegTodayTotal(m bson.M) (int64, error) {
	return int64(Count(LogRegistTodays, m)), nil
}

// 获取日志列表
func (this *loggerService) GetPayTodayList(page, pageSize int, m bson.M) ([]entity.LogPayToday, error) {
	var list []entity.LogPayToday
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	LogPayTodays.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, nil
}

// 获取日志总数
func (this *loggerService) GetPayTodayTotal(m bson.M) (int64, error) {
	return int64(Count(LogPayTodays, m)), nil
}

// 获取盈亏统计日志列表
func (this *loggerService) GetChipTodayList(page, pageSize int, m bson.M) ([]entity.LogChipToday, error) {
	var list []entity.LogChipToday
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	LogChipTodays.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)

	//转换为分展示
	list = this.chipToday2(list)
	return list, nil
}

//转换为分展示
func (this *loggerService) chipToday2(list []entity.LogChipToday) []entity.LogChipToday {
	for k, v := range list {
		v.RobotNumf = Chip2Float(int64(v.RobotNum))
		v.RolesNumf = Chip2Float(int64(v.RolesNum))
		list[k] = v
	}
	return list
}

// 获取盈亏统计日志总数
func (this *loggerService) GetChipTodayTotal(m bson.M) (int64, error) {
	return int64(Count(LogChipTodays, m)), nil
}

// 获取开奖结果记录日志列表
func (this *loggerService) GetExpectList(page, pageSize int, m bson.M) ([]entity.Pk10Record, error) {
	//var list []entity.Pk10Record
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	var list = make([]entity.Pk10Record, 0)
	err := Pk10Records.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, err
}

// 获取开奖结果记录日志总数
func (this *loggerService) GetExpectTotal(m bson.M) (int64, error) {
	return int64(Count(Pk10Records, m)), nil
}

// 房间单局记录
func (this *loggerService) GetGameRecordList(page, pageSize int, m bson.M) ([]entity.GameRecord, error) {
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	var list = make([]entity.GameRecord, 0)
	err := GameRecords.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)

	//转换为分展示
	list = this.gameRecordList2(list)
	return list, err
}

//转换为分展示
func (this *loggerService) gameRecordList2(list []entity.GameRecord) []entity.GameRecord {
	for k, v := range list {
		v.RobotFeef = Chip2Float(int64(v.RobotFee))
		v.PlayerFeef = Chip2Float(int64(v.PlayerFee))
		v.FeeNumf = Chip2Float(int64(v.FeeNum))
		v.BetNumf = Chip2Float(int64(v.BetNum))
		v.WinNumf = Chip2Float(int64(v.WinNum))
		v.LoseNumf = Chip2Float(int64(v.LoseNum))
		v.RefundNumf = Chip2Float(int64(v.RefundNum))
		list[k] = v
	}
	return list
}

func (this *loggerService) GetGameRecordTotal(m bson.M) (int64, error) {
	return int64(Count(GameRecords, m)), nil
}

func (this *loggerService) GetGameRecord(id string) (*entity.GameRecord, error) {
	user := new(entity.GameRecord)
	GetByQ(GameRecords, bson.M{"_id": id}, user)
	if user.Roomid == "" {
		return user, errors.New("不存在")
	}

	//转换为分展示
	user = this.gameRecord2(user)
	return user, nil
}

//转换为分展示
func (this *loggerService) gameRecord2(v *entity.GameRecord) *entity.GameRecord {
	v.RobotFeef = Chip2Float(int64(v.RobotFee))
	v.PlayerFeef = Chip2Float(int64(v.PlayerFee))
	v.FeeNumf = Chip2Float(int64(v.FeeNum))
	v.BetNumf = Chip2Float(int64(v.BetNum))
	v.WinNumf = Chip2Float(int64(v.WinNum))
	v.LoseNumf = Chip2Float(int64(v.LoseNum))
	v.RefundNumf = Chip2Float(int64(v.RefundNum))
	//
	for k, val := range v.Result {
		val.Betsf = Chip2Float(int64(val.Bets))
		val.Winsf = Chip2Float(int64(val.Wins))
		val.Refundf = Chip2Float(int64(val.Refund))
		v.Result[k] = val
	}
	//
	for k, val := range v.Record {
		val.Feef = Chip2Float(int64(val.Fee))
		v.Record[k] = val
	}
	//
	for k, val := range v.Details {
		val.Feef = Chip2Float(int64(val.Fee))
		//
		for k2, val2 := range val.Record {
			val2.Feef = Chip2Float(int64(val2.Fee))
			val.Record[k2] = val2
		}
		v.Details[k] = val
	}
	return v
}

// 个人房间单局记录
func (this *loggerService) GetUserRecordList(page, pageSize int, m bson.M) ([]entity.UserRecord, error) {
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	var list = make([]entity.UserRecord, 0)
	err := UserRecords.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)

	//转换为分展示
	list = this.userRecordList2(list)
	return list, err
}

//转换为分展示
func (this *loggerService) userRecordList2(list []entity.UserRecord) []entity.UserRecord {
	for k, v := range list {
		v.Restf = Chip2Float(int64(v.Rest))
		v.Betsf = Chip2Float(int64(v.Bets))
		v.Profitsf = Chip2Float(int64(v.Profits))
		v.Feef = Chip2Float(int64(v.Fee))
		list[k] = v
	}
	return list
}

func (this *loggerService) GetUserRecordTotal(m bson.M) (int64, error) {
	return int64(Count(UserRecords, m)), nil
}

func (this *loggerService) GetUserRecord(id, userid string) (*entity.UserRecord, error) {
	user := new(entity.UserRecord)
	GetByQ(UserRecords, bson.M{"roomid": id, "userid": userid}, user)
	if user.Roomid == "" {
		return user, errors.New("不存在")
	}

	//转换为分展示
	user = this.userRecord2(user)
	return user, nil
}

//转换为分展示
func (this *loggerService) userRecord2(v *entity.UserRecord) *entity.UserRecord {
	v.Restf = Chip2Float(int64(v.Rest))
	v.Betsf = Chip2Float(int64(v.Bets))
	v.Profitsf = Chip2Float(int64(v.Profits))
	v.Feef = Chip2Float(int64(v.Fee))
	for k, val := range v.Details {
		val.Betsf = Chip2Float(int64(val.Bets))
		val.Winsf = Chip2Float(int64(val.Wins))
		val.Refundf = Chip2Float(int64(val.Refundf))
		v.Details[k] = val
	}
	return v
}

func userStat2f(stat *entity.UserStat) {
	stat.Chipsf = Chip2Float(int64(stat.Chips))
	stat.TodayFeef = Chip2Float(int64(stat.TodayFee))
	stat.YesterdayFeef = Chip2Float(int64(stat.YesterdayFee))
	stat.AllFeef = Chip2Float(int64(stat.AllFee))
	stat.TodayAgentFeef = Chip2Float(int64(stat.TodayAgentFee))
	stat.YesterdayAgentFeef = Chip2Float(int64(stat.YesterdayAgentFee))
	stat.AgentAllFeef = Chip2Float(int64(stat.AgentAllFee))
}

func agencyStat2f(stat *entity.AgencyStat) {
	stat.Chipf = Chip2Float(int64(stat.Chip))
	stat.TodayBetsf = Chip2Float(int64(stat.TodayBets))
	stat.TodayProfitsf = Chip2Float(int64(stat.TodayProfits))
	stat.AllProfitsf = Chip2Float(int64(stat.AllProfits))
	stat.TopProfitsf = Chip2Float(int64(stat.TopProfits))
	stat.ExtractProfitsf = Chip2Float(int64(stat.ExtractProfits))
}

//首页统计
func (this *loggerService) GetIndexStat(username string) (stat1 *entity.UserStat,
	stat2 *entity.AgencyStat, isAgent int) {
	stat1 = new(entity.UserStat)
	stat2 = new(entity.AgencyStat)
	if username == "" {
		return
	}
	agency, err := UserService.GetUserByName(username)
	if err != nil {
		return
	}
	if agency.Agent != "" {
		isAgent = 1
		this.getAgencyStat(agency.Agent, stat2)
		agencyStat2f(stat2)
	} else {
		this.getUserStat(stat1)
		userStat2f(stat1)
	}
	return
}

//超级管理员后台
func (this *loggerService) getUserStat(stat *entity.UserStat) {
	//stat := new(entity.UserStat)
	//会员统计
	n1, n2, n3 := this.UserStat1()
	stat.TodayNewPlayers = n2
	stat.YesterdayNewPlayers = n3
	stat.AllPlayers = n1
	//用户持有财富统计
	n1, n2, n3 = this.UserStat2()
	stat.Chips = n1
	stat.Cards = n2
	stat.Coins = n3
	//抽佣统计
	n1, n2, n3 = this.UserStat3()
	stat.AllFee = n1
	stat.TodayFee = n2
	stat.YesterdayFee = n3
	//代理统计
	n1, n2, n3 = this.UserStat4()
	stat.TodayNewAgent = n2
	stat.YesterdayNewAgent = n3
	stat.AllAgent = n1
	//抽佣统计
	n1, n2, n3 = this.UserStat5()
	stat.AgentAllFee = n1
	stat.TodayAgentFee = n2
	stat.YesterdayAgentFee = n3
	return
}

//机器人今日输赢，昨日输赢，总输赢
func (this *loggerService) GetStatChipToday() (int64, int64, int64) {
	num1 := this.getRobotProfitsToday()
	num2 := this.getRobotProfitsYesterday()
	num3 := this.getRobotProfitsAll()
	beego.Trace("GetStatChipToday : ", num1, num2, num3)
	beego.Error("GetStatChipToday : ", num1, num2, num3)
	//TODO 总输赢=num3 + 今日输赢
	return num1, num2, num3
}

//统计机器人昨日盈亏
func (this *loggerService) getRobotProfitsYesterday() (num int64) {
	//startTime := utils.Stamp2Time(utils.TimestampYesterday())
	startTime := TimeYesterday4()
	//now := time.Now()
	//startTime := time.Date(now.Year(), now.Month()-1, 22, 04, 30, 0, 0, time.Local)
	beego.Trace("getRobotProfitsYesterday : ", startTime.String())
	m := bson.M{
		"$match": bson.M{
			"day_stamp": startTime,
			"robot_num": bson.M{"$ne": 0},
		},
	}
	beego.Trace("getRobotProfitsYesterday : ", m)
	return this.getRobotProfits(m)
}

//统计机器人总盈亏
func (this *loggerService) getRobotProfitsAll() (num int64) {
	m := bson.M{
		"$match": bson.M{
			"robot_num": bson.M{"$ne": 0},
		},
	}
	return this.getRobotProfits(m)
}

//统计机器人盈亏
func (this *loggerService) getRobotProfits(m bson.M) (num int64) {
	n := bson.M{
		"$group": bson.M{
			"_id": 1,
			"robot_num": bson.M{
				"$sum": "$robot_num",
			},
		},
	}
	operations := []bson.M{m, n}
	result := bson.M{}
	pipe := LogChipTodays.Pipe(operations)
	err := pipe.One(&result)
	if err != nil {
		beego.Error("getRobotProfitsAll fail err: ", err)
		return
	}
	beego.Trace("getRobotProfitsAll : ", result)
	if v, ok := result["robot_num"]; ok {
		num = v.(int64)
		return
	}
	return
}

//统计机器人今日盈亏
func (this *loggerService) getRobotProfitsToday() (num int64) {
	//endTime := utils.Stamp2Time(utils.TimestampToday())
	endTime := TimeToday4()
	m := bson.M{
		"$match": bson.M{
			"robot": true,
			"ctime": bson.M{"$gte": endTime},
		},
	}
	n := bson.M{
		"$group": bson.M{
			"_id": 1,
			"profits": bson.M{
				"$sum": "$profits",
			},
		},
	}
	operations := []bson.M{m, n}
	result := bson.M{}
	pipe := UserRecords.Pipe(operations)
	err := pipe.One(&result)
	if err != nil {
		beego.Error("getStatProfitsToday fail err: ", err)
		return
	}
	beego.Trace("getStatProfitsToday : ", result)
	if v, ok := result["profits"]; ok {
		num = v.(int64)
		return
	}
	return
}

//会员统计
func (this *loggerService) UserStat1() (n1, n2, n3 int64) {
	//会员总数
	m := bson.M{}
	n1, _ = PlayerService.GetTotal(0, m)
	//今日新增
	//endTime := utils.Stamp2Time(utils.TimestampToday())
	endTime := TimeToday4()
	m["ctime"] = bson.M{"$gte": endTime}
	n2, _ = PlayerService.GetTotal(0, m)
	//昨天新增
	//startTime := utils.Stamp2Time(utils.TimestampYesterday())
	startTime := TimeYesterday4()
	m["ctime"] = bson.M{"$gte": startTime, "$lt": endTime}
	n3, _ = PlayerService.GetTotal(0, m)
	return
}

//用户持有财富统计
func (this *loggerService) UserStat2() (n1, n2, n3 int64) {
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
			"card": bson.M{
				"$sum": "$card",
			},
			"coin": bson.M{
				"$sum": "$coin",
			},
		},
	}
	//统计
	operations := []bson.M{m, n}
	result := bson.M{}
	pipe := PlayerUsers.Pipe(operations)
	err := pipe.One(&result)
	if err != nil {
		beego.Error("UserStat2 fail err: ", err)
		return
	}
	//用户持有筹码总数
	if v, ok := result["chip"]; ok {
		n1 = v.(int64)
	}
	//用户持有房卡总数
	if v, ok := result["card"]; ok {
		n2 = v.(int64)
	}
	//用户持有金币总数
	if v, ok := result["coin"]; ok {
		n3 = v.(int64)
	}
	return
}

//抽佣统计
func (this *loggerService) UserStat3() (n1, n2, n3 int64) {
	m1 := bson.M{
		"fee_num":  bson.M{"$gt": 0},
		"roomtype": 1, //抽佣房间
	}
	m := bson.M{"$match": m1}
	n := bson.M{
		"$group": bson.M{
			"_id": 1,
			"fee_num": bson.M{
				"$sum": "$fee_num",
			},
		},
	}
	//统计
	operations := []bson.M{m, n}
	result := bson.M{}
	pipe := GameRecords.Pipe(operations)
	err := pipe.One(&result)
	if err != nil {
		beego.Error("UserStat3 fail err: ", err)
		//return
	}
	//抽佣总数
	if v, ok := result["fee_num"]; ok {
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
	pipe = GameRecords.Pipe(operations)
	err = pipe.One(&result)
	if err != nil {
		beego.Error("UserStat3 fail err: ", err)
		//return
	}
	if v, ok := result["fee_num"]; ok {
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
	pipe = GameRecords.Pipe(operations)
	err = pipe.One(&result)
	if err != nil {
		beego.Error("UserStat3 fail err: ", err)
		//return
	}
	if v, ok := result["fee_num"]; ok {
		n3 = v.(int64)
	}
	return
}

//代理统计
func (this *loggerService) UserStat4() (n1, n2, n3 int64) {
	//代理总数
	m := bson.M{}
	m["agent"] = bson.M{"$ne": ""}
	n1, _ = AgencyService.GetAgencyTotal(m)
	//今日新增
	//endTime := utils.Stamp2Time(utils.TimestampToday())
	endTime := TimeToday4()
	m["create_time"] = bson.M{"$gte": endTime}
	n2, _ = AgencyService.GetAgencyTotal(m)
	//昨天新增
	//startTime := utils.Stamp2Time(utils.TimestampYesterday())
	startTime := TimeYesterday4()
	m["create_time"] = bson.M{"$gte": startTime, "$lt": endTime}
	n3, _ = AgencyService.GetAgencyTotal(m)
	return
}

//代理登录后台
func (this *loggerService) getAgencyStat(userid string,
	stat *entity.AgencyStat) {
	//stat := new(entity.AgencyStat)
	//我的财富
	//n1, n2, n3 := this.AgencyStat1(userid)
	//stat.Chip = n1
	//stat.Card = n2
	//stat.Coin = n3
	//下属代理统计
	this.AgencyStat2(userid, stat)
	//收益统计
	a1, a2, a3, a4, a5 := this.AgencyStat5(userid)
	stat.AllProfits = a3
	stat.TodayProfits = a1
	stat.TodayBets = a2
	stat.TopProfits = a4
	stat.ExtractProfits = a5
	return
}

//我的财富
func (this *loggerService) AgencyStat1(userid string) (n1, n2, n3 int64) {
	//TODO 现在邀请码不是绑定玩家id userid != agent
	return
	p, _ := PlayerService.GetPlayer(userid)
	n1 = p.Chip
	n2 = p.Card
	n3 = p.Coin
	return
}

//下属代理统计
func (this *loggerService) AgencyStat2(userid string,
	stat *entity.AgencyStat) {
	stat.UnderlingPlayer = PlayerService.GetBuilds(userid)
	stat.UnderlingAgency = UserService.GetBuilds(userid)
}

/*
//下属代理统计 agent == userid
func (this *loggerService) AgencyStat2(userid string,
	stat *entity.AgencyStat) {
	list2 := PlayerService.GetAllBuilds2(userid)
	if len(list2) == 0 {
		return
	}
	for _, v := range list2 {
		if AgencyService.IsAgent(v) {
			stat.UnderlingAgency += 1
			this.AgencyStat2(v, stat)
		} else {
			stat.UnderlingPlayer += 1
		}
	}
}
*/

/*
//总收益统计(对象为管理员)
func (this *loggerService) AgencyStat3() (n1, n2, n3 int64) {
	m1 := bson.M{
		"fee_num":  bson.M{"$gt": 0},
		"roomtype": 1, //抽佣房间
	}
	m := bson.M{"$match": m1}
	m2 := bson.M{
		"_id": 1,
		"fee_num": bson.M{
			"$sum": "$fee_num",
		},
	}
	n := bson.M{"$group": m2}
	//统计
	operations := []bson.M{m, n}
	result := bson.M{}
	pipe := GameRecords.Pipe(operations)
	err := pipe.One(&result)
	if err != nil {
		beego.Trace("AgencyStat3 fail err: ", err)
		//return
	}
	//全部收益
	if v, ok := result["fee_num"]; ok {
		n1 = v.(int64)
	}
	//今日统计
	endTime := utils.Stamp2Time(utils.TimestampToday())
	m1["ctime"] = bson.M{"$gte": endTime}
	m2["bet_num"] = bson.M{"$sum": "$bet_num"}
	m = bson.M{"$match": m1}
	n = bson.M{"$group": m2}
	//
	operations = []bson.M{m, n}
	result = bson.M{}
	pipe = GameRecords.Pipe(operations)
	err = pipe.One(&result)
	if err != nil {
		beego.Trace("AgencyStat3 fail err: ", err)
		//return
	}
	//今日收益
	if v, ok := result["fee_num"]; ok {
		n2 = v.(int64)
	}
	//今日注额
	if v, ok := result["bet_num"]; ok {
		n3 = v.(int64)
	}
	return
}
*/

/*
//总收益统计(对象为代理)
func (this *loggerService) AgencyStat4(agent string) (n1, n2, n3 int64) {
	m1 := bson.M{
		"fee_num":  bson.M{"$gt": 0},
		"roomtype": 1, //抽佣房间
	}
	m := bson.M{"$match": m1}
	m2 := bson.M{
		"_id": 1,
		"fee_num": bson.M{
			"$sum": "$fee_num",
		},
	}
	n := bson.M{"$group": m2}
	//统计
	operations := []bson.M{m, n}
	result := bson.M{}
	pipe := GameRecords.Pipe(operations)
	err := pipe.One(&result)
	if err != nil {
		beego.Trace("AgencyStat4 fail err: ", err)
		//return
	}
	//全部收益
	if v, ok := result["fee_num"]; ok {
		n1 = v.(int64)
	}
	//今日统计
	endTime := utils.Stamp2Time(utils.TimestampToday())
	m1["ctime"] = bson.M{"$gte": endTime}
	m2["bet_num"] = bson.M{"$sum": "$bet_num"}
	m = bson.M{"$match": m1}
	n = bson.M{"$group": m2}
	//
	operations = []bson.M{m, n}
	result = bson.M{}
	pipe = GameRecords.Pipe(operations)
	err = pipe.One(&result)
	if err != nil {
		beego.Trace("AgencyStat4 fail err: ", err)
		//return
	}
	//今日收益
	if v, ok := result["fee_num"]; ok {
		n2 = v.(int64)
	}
	//今日注额
	if v, ok := result["bet_num"]; ok {
		n3 = v.(int64)
	}
	return
}
*/

//绑定我的玩家
func buildList(agent string) (list []string) {
	return PlayerService.GetAllBuilds2(agent)
}

//今日收益 TODO 优化
func (this *loggerService) AgencyStat5(agent string) (n1, n2, n3, n4, n5 int64) {
	if agent == "" {
		return
	}
	//绑定我的玩家
	list := buildList(agent)
	beego.Trace("AgencyStat5 list: ", list)
	//统计条件
	nowTime := utils.LocalTime()
	//endTime := utils.Stamp2Time(utils.TimestampToday())
	endTime := TimeToday4()
	//startTime := utils.Stamp2Time(utils.TimestampYesterday())
	//收益
	//今日收益,只统计自己所得
	n1 = feeStatByLog(agent, endTime, nowTime)
	//这里是统计所有产出,不能这样展示
	//n1 = feeStat(list, endTime, nowTime)
	//昨日收益
	//n2 = feeStat(list, startTime, endTime)
	//今日注额
	n2 = betStat(list, endTime, nowTime)
	//全部收益
	//n3, _ = AgencyService.GetFeeAll(agent)
	//全部收益
	n3, n4, n5 = AgencyService.GetFeeAll2(agent)
	return
}

//总收益统计
func feeStat1(agency *entity.User) {
	if agency.Agent == "" {
		return
	}
	//绑定我的玩家
	list := buildList(agency.Agent)
	beego.Trace("feeStat1 list: ", list)
	//统计条件
	endTime := utils.LocalTime()
	startTime := agency.FeeTime
	//收益统计
	feeNum := feeStat(list, startTime, endTime)
	if feeNum <= 0 {
		return
	}

	//上级代理

	var feeNum1 int64 //上级分成
	var feeNum2 int64 //系统分成
	var feeNum3 int64 //自己所得

	feeNum1, feeNum2, feeNum3 = feeStat5(feeNum, agency)

	if agency.Rate > 0 && len(agency.ParentAgent) > 0 {
		//上级处理
		feeStat3(feeNum1, agency.ParentAgent)
	}

	//抽成日志记录,TODO 上级分成最后可能会被系统提成
	feeStat6(feeNum, feeNum1, feeNum2, feeNum3, agency.Agent,
		agency.ParentAgent, agency.Rate, agency.SysRate)

	//更新自己代理数据
	q := bson.M{"agent": agency.Agent}
	n := bson.M{"fee_time": endTime}
	c := bson.M{"fee_all": feeNum3, "fee_top": feeNum3,
		"fee_rate": feeNum1, "sys_fee_rate": feeNum2}
	i := bson.M{"$set": n, "$inc": c}
	//更新当前统计代理统计时间
	feeStat2(q, i)
	//日志
	feeStat7(feeNum1, feeNum2, feeNum3, agency.Agent)
}

//feeNum1 上级分成 feeNum2 系统分成 feeNum3 自己所得
func feeStat4(feeNum1, feeNum2, feeNum3 int64, agent string) {
	//更新上级代理数据
	q := bson.M{"agent": agent}
	c := bson.M{"fee_all": feeNum3, "fee_top": feeNum3,
		"fee_rate": feeNum1, "sys_fee_rate": feeNum2}
	i := bson.M{"$inc": c}
	feeStat2(q, i)
	//日志
	feeStat7(feeNum1, feeNum2, feeNum3, agent)
}

//抽成日志
//feeNum1 上级分成 feeNum2 系统分成 feeNum3 自己所得
func feeStat7(feeNum1, feeNum2, feeNum3 int64, agent string) {
	agentFee := &entity.AgentFeeLog{
		Agent:      agent,
		FeeAll:     feeNum3,
		FeeRate:    feeNum1,
		SysFeeRate: feeNum2,
	}
	LoggerService.SaveAgentFeeLog(agentFee)
}

//feeNum1 上级分成 feeNum2 系统分成 feeNum3 自己所得
func feeStat5(feeNum int64, agency *entity.User) (feeNum1, feeNum2, feeNum3 int64) {
	//系统
	feeNum2 = getFeeRate1(feeNum, agency.SysRate)
	feeNum3 = feeNum - feeNum2

	//if agency.Rate > 0 && len(agency.ParentAgent) > 0 {
	if agency.Rate > 0 {
		//上级
		feeNum1 = getFeeRate(feeNum3, agency.Rate)
		feeNum3 = feeNum3 - feeNum1
	}
	return
}

//上级代理抽成
func feeStat3(feeNum int64, agent string) {
	if feeNum <= 0 {
		return
	}
	agency, err := UserService.GetUserByAgent(agent)
	if err != nil {
		//获取出错时全部给自己
		feeStat4(0, 0, feeNum, agent)
		return
	}

	var feeNum1 int64 //上级分成
	//var feeNum2 int64 //系统分成
	var feeNum3 int64 //自己所得

	//最高一级抽成给系统或者上级,最高级没有上级
	//if agency.Rate > 0 && len(agency.ParentAgent) > 0 {
	if agency.Rate > 0 {
		//上级
		feeNum1 = getFeeRate(feeNum, agency.Rate)
	}
	//自己所得
	feeNum3 = feeNum - feeNum1

	//TODO 没有上级情况下出现抽成时feeNum1相当于被系统抽成
	if len(agency.ParentAgent) == 0 && agency.Rate > 0 && feeNum1 > 0 {
		//记录系统抽成日志记录, 剩余全部算系统的
		feeStat6(feeNum, 0, feeNum1, feeNum3, agency.Agent,
			agency.ParentAgent, agency.Rate, agency.SysRate)
		feeStat4(feeNum1, 0, feeNum3, agent)
		return
	}
	if len(agency.ParentAgent) == 0 || agency.Rate == 0 {
		//记录系统抽成日志记录
		feeStat6(feeNum, feeNum1, 0, feeNum3, agency.Agent,
			agency.ParentAgent, agency.Rate, agency.SysRate)
		//没有上级或没有上级提成
		feeStat4(feeNum1, 0, feeNum3, agent)
		return
	}

	//记录系统抽成日志记录
	feeStat6(feeNum, feeNum1, 0, feeNum3, agency.Agent,
		agency.ParentAgent, agency.Rate, agency.SysRate)

	//更新当前代理
	feeStat4(feeNum1, 0, feeNum3, agent)
	//循环上级
	feeStat3(feeNum1, agency.ParentAgent)
}

//抽成日志
//feeNum1 上级分成 feeNum2 系统分成 feeNum3 自己所得
func feeStat6(feeNum, feeNum1, feeNum2, feeNum3 int64,
	agent, parent_agent string, rate, sysRate uint32) {
	agentFee := &entity.AgentFee{
		Agent:       agent,
		ParentAgent: parent_agent,
		FeeNum:      feeNum,
		Fee:         feeNum3,
		ParentFee:   feeNum1,
		//FeeAll:      feeNum1 + feeNum3,
		FeeAll:  feeNum3,
		SysFee:  feeNum2,
		Rate:    rate,
		SysRate: sysRate,
	}
	LoggerService.SaveAgentFee(agentFee)
}

//更新数据
func feeStat2(q, i bson.M) {
	if !Update(Users, q, i) {
		beego.Error("feeStat2 err: ", q, i)
		return
	}
	beego.Trace("feeStat2 ok: ", q, i)
}

//系统抽成
func getFeeRate1(feeNum int64, rate uint32) int64 {
	return int64(math.Trunc(float64(feeNum) * (float64(rate) / 100)))
}

//抽成
func getFeeRate(feeNum int64, rate uint32) int64 {
	return int64(math.Trunc(float64(feeNum) * (1 - (float64(rate) / 100))))
}

//管理后台首页的数据统计
//本周 上周 上上周
//
//机器人注额
//机器人盈亏
//机器人抽佣
//
//玩家注额
//玩家盈亏
//玩家抽佣

func (this *loggerService) GetWeekStat(robot bool) (m1, m2, m3 string) {
	list := make([]string, 0)
	n1, n2, n3 := this.GetWeekStat2(robot, list)
	d1, _ := json.Marshal(n1)
	d2, _ := json.Marshal(n2)
	d3, _ := json.Marshal(n3)
	m1 = string(d1)
	m2 = string(d2)
	m3 = string(d3)
	return
}

//代理
func (this *loggerService) GetAgentWeekStat(agent string) (m1, m2, m3 string) {
	list := buildList(agent)
	if len(list) == 0 {
		return
	}
	n1, n2, n3 := this.GetWeekStat2(false, list)
	d1, _ := json.Marshal(n1)
	d2, _ := json.Marshal(n2)
	d3, _ := json.Marshal(n3)
	m1 = string(d1)
	m2 = string(d2)
	m3 = string(d3)
	return
}

func (this *loggerService) GetWeekStat2(robot bool, list []string) (m1, m2, m3 []float64) {
	//本周
	start, end := utils.ThisWeek()
	t1, t2, t3 := this.weekStat(robot, start, end, list)
	//m1 = []int64{t1, t2, t3}
	t1f := Chip2Float(int64(t1))
	t2f := Chip2Float(int64(t2))
	t3f := Chip2Float(int64(t3))
	m1 = []float64{t1f, t2f, t3f}
	//上周
	start = start.AddDate(0, 0, -7)
	end = end.AddDate(0, 0, -7)
	l1, l2, l3 := this.weekStat(robot, start, end, list)
	//m2 = []int64{l1, l2, l3}
	l1f := Chip2Float(int64(l1))
	l2f := Chip2Float(int64(l2))
	l3f := Chip2Float(int64(l3))
	m2 = []float64{l1f, l2f, l3f}
	//上上周
	start = start.AddDate(0, 0, -7)
	end = end.AddDate(0, 0, -7)
	n1, n2, n3 := this.weekStat(robot, start, end, list)
	//m3 = []int64{n1, n2, n3}
	n1f := Chip2Float(int64(n1))
	n2f := Chip2Float(int64(n2))
	n3f := Chip2Float(int64(n3))
	m3 = []float64{n1f, n2f, n3f}
	return
}

func (this *loggerService) weekStat(robot bool, startTime,
	endTime time.Time, list []string) (n1, n2, n3 int64) {
	m1 := bson.M{
		//"roomtype": 1,    //抽佣房间
		"robot": robot, //是否机器人
		"ctime": bson.M{"$gte": startTime, "$lt": endTime},
	}
	if len(list) != 0 {
		m1["userid"] = bson.M{"$in": list}
	}
	m := bson.M{"$match": m1}
	m2 := bson.M{
		"_id": 1,
		"bets": bson.M{
			"$sum": "$bets",
		},
		"profits": bson.M{
			"$sum": "$profits",
		},
		"fee": bson.M{
			"$sum": "$fee",
		},
	}
	n := bson.M{"$group": m2}
	//统计
	operations := []bson.M{m, n}
	result := bson.M{}
	pipe := UserRecords.Pipe(operations)
	err := pipe.One(&result)
	if err != nil {
		beego.Error("weekStat fail err: ", err)
		return
	}
	//注额
	if v, ok := result["bets"]; ok {
		//n1 = v.(int64)
		if val, ok := v.(int64); ok {
			n1 = val
		}
	}
	//盈亏
	if v, ok := result["profits"]; ok {
		//n2 = v.(int64)
		if val, ok := v.(int64); ok {
			n2 = val
		}
	}
	//抽佣
	if v, ok := result["fee"]; ok {
		//n3 = v.(int64)
		if val, ok := v.(int64); ok {
			n3 = val
		}
	}
	return
}

//后台筹码充值赠送统计
func (this *loggerService) GetGiveStat() (n1 int64) {
	m1 := bson.M{
		"type": entity.LogType9, //后台操作筹码
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
		beego.Error("GetGiveStat fail err: ", err)
		return
	}
	//beego.Trace("GetGiveStat result: ", result)
	//注额
	if v, ok := result["num"]; ok {
		if val, ok := v.(int64); ok {
			n1 = val
		}
	}
	return
}

//首页统计收益
func (this *loggerService) GetFeeStat() (n1, n2 int64) {
	m1 := bson.M{
		"agent": bson.M{"$ne": ""},
	}
	m := bson.M{"$match": m1}
	m2 := bson.M{
		"_id": 1,
		"profits": bson.M{
			"$sum": "$fee_extract",
		},
		"fee": bson.M{
			"$sum": "$fee_all",
		},
	}
	n := bson.M{"$group": m2}
	//统计
	operations := []bson.M{m, n}
	result := bson.M{}
	pipe := Users.Pipe(operations)
	err := pipe.One(&result)
	if err != nil {
		beego.Error("GetFeeStat fail err: ", err)
		return
	}
	//可提现
	if v, ok := result["fee"]; ok {
		if val, ok := v.(int64); ok {
			n1 = val
		}
	}
	//已提现
	if v, ok := result["profits"]; ok {
		if val, ok := v.(int64); ok {
			n2 = val
		}
	}
	return
}
