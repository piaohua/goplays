package controllers

import (
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"goplays/web/app/entity"
	"goplays/web/app/libs"
	"goplays/web/app/service"

	"github.com/astaxie/beego"
)

type LoggerController struct {
	BaseController
}

// 注册日志
func (this *LoggerController) RegistList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	userid := this.GetString("userid")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	if userid != "" {
		m["userid"] = userid
	}
	count, _ := service.LoggerService.GetRegistTotal(m)
	list, _ := service.LoggerService.GetRegistList(page, this.pageSize, m)

	this.Data["pageTitle"] = "注册日志"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("LoggerController.RegistList", "status", status, "userid", userid, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 登录日志
func (this *LoggerController) LoginList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	userid := this.GetString("userid")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "login_time", "login_time")
	if userid != "" {
		m["userid"] = userid
	}
	count, _ := service.LoggerService.GetLoginTotal(m)
	list, _ := service.LoggerService.GetLoginList(page, this.pageSize, m)

	this.Data["pageTitle"] = "登录日志"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("LoggerController.LoginList", "status", status, "userid", userid, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 充值日志
func (this *LoggerController) PayList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	userid := this.GetString("userid")
	agent := this.GetString("agent")
	typeId, _ := this.GetInt("type_id")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	if userid != "" {
		m["userid"] = userid
	}
	if agent != "" {
		m["agent"] = agent
	}
	m["result"] = typeId
	count, _ := service.LoggerService.GetPayTotal(m)
	list, _ := service.LoggerService.GetPayList(page, this.pageSize, m)

	le := len(list)
	for i := 0; i < le; i++ {
		list[i].Money = list[i].Money / 100 //转换为元
	}

	typeList := entity.TradeResult

	this.Data["pageTitle"] = "充值日志"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["typeList"] = typeList
	this.Data["typeId"] = typeId
	this.Data["userid"] = userid
	this.Data["agent"] = agent
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("LoggerController.PayList", "status", status, "userid", userid, "agent", agent, "type_id", typeId, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 钻石日志
func (this *LoggerController) DiamondList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	userid := this.GetString("userid")
	typeId, _ := this.GetInt("type_id")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	if userid != "" {
		m["userid"] = userid
	}
	if typeId != 0 {
		m["type"] = typeId
	}
	count, _ := service.LoggerService.GetDiamondTotal(m)
	list, _ := service.LoggerService.GetDiamondList(page, this.pageSize, m)

	typeList := entity.LogType

	this.Data["pageTitle"] = "钻石日志"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["typeList"] = typeList
	this.Data["typeId"] = typeId
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("LoggerController.DiamondList", "status", status, "userid", userid, "type_id", typeId, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 金币日志
func (this *LoggerController) CoinList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	userid := this.GetString("userid")
	typeId, _ := this.GetInt("type_id")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	if userid != "" {
		m["userid"] = userid
	}
	if typeId != 0 {
		m["type"] = typeId
	}
	count, _ := service.LoggerService.GetCoinTotal(m)
	list, _ := service.LoggerService.GetCoinList(page, this.pageSize, m)

	typeList := entity.LogType

	this.Data["pageTitle"] = "金币日志"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["typeList"] = typeList
	this.Data["typeId"] = typeId
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("LoggerController.CoinList", "status", status, "userid", userid, "type_id", typeId, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 筹码日志
func (this *LoggerController) ChipList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	userid := this.GetString("userid")
	typeId, _ := this.GetInt("type_id")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	if userid != "" {
		m["userid"] = userid
	}
	if typeId != 0 {
		m["type"] = typeId
	}
	count, _ := service.LoggerService.GetChipTotal(m)
	list, _ := service.LoggerService.GetChipList(page, this.pageSize, m)

	typeList := entity.LogType

	this.Data["pageTitle"] = "筹码日志"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["typeList"] = typeList
	this.Data["typeId"] = typeId
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("LoggerController.ChipList", "status", status, "userid", userid, "type_id", typeId, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 抽佣筹码日志
func (this *AgencyController) FeeChipList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	userid := this.GetString("userid")
	agent := this.GetString("agent")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	if userid != "" {
		m["userid"] = userid
	}
	//代理
	if agent != "" {
		list3 := service.PlayerService.GetAllBuilds2(agent)
		m["userid"] = bson.M{"$in": list3}
	}
	m["type"] = entity.LogType43 //抽佣类型
	count, _ := service.LoggerService.GetChipTotal(m)
	list, _ := service.LoggerService.GetChipList(page, this.pageSize, m)

	//typeList := entity.LogType

	this.Data["pageTitle"] = "抽佣明细"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["userid"] = userid
	this.Data["agent"] = agent
	//this.Data["typeList"] = typeList
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("LoggerController.FeeChipList", "status", status, "userid", userid, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 抽佣筹码日志
func (this *AgencyController) MineFeeChipList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	userid := this.GetString("userid")
	agent := this.GetString("agent")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	//if userid != "" {
	//	m["userid"] = userid
	//}
	//代理
	//if agent != "" {
	//	list3 := service.PlayerService.GetAllBuilds2(agent)
	//	m["userid"] = bson.M{"$in": list3}
	//}
	m["type"] = entity.LogType43 //抽佣类型
	if userid != "" {
		list := this.auth.GetUser().Players
		var fined bool
		for _, v := range list {
			if v == userid {
				fined = true
				break
			}
		}
		if fined {
			m["userid"] = userid
		} else {
			m["userid"] = bson.M{"$in": []string{}}
		}
	} else {
		m["userid"] = bson.M{"$in": this.auth.GetUser().Players}
	}
	count, _ := service.LoggerService.GetChipTotal(m)
	list, _ := service.LoggerService.GetChipList(page, this.pageSize, m)

	//typeList := entity.LogType

	this.Data["pageTitle"] = "抽佣明细"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["userid"] = userid
	this.Data["agent"] = agent
	//this.Data["typeList"] = typeList
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("LoggerController.MineFeeChipList", "status", status, "userid", userid, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 绑定日志
func (this *LoggerController) BuildList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	userid := this.GetString("userid")
	agent := this.GetString("agent")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	if userid != "" {
		m["userid"] = userid
	}
	if agent != "" {
		m["agent"] = agent
	}
	count, _ := service.LoggerService.GetBuildTotal(m)
	list, _ := service.LoggerService.GetBuildList(page, this.pageSize, m)

	this.Data["pageTitle"] = "绑定日志"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("LoggerController.BuildList", "status", status, "userid", userid, "agent", agent, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 注册统计日志
func (this *LoggerController) RegistTodayList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	count, _ := service.LoggerService.GetRegTodayTotal(m)
	list, _ := service.LoggerService.GetRegTodayList(page, this.pageSize, m)

	this.Data["pageTitle"] = "注册统计"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("LoggerController.RegistTodayList", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 充值统计日志
func (this *LoggerController) PayTodayList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	count, _ := service.LoggerService.GetPayTodayTotal(m)
	list, _ := service.LoggerService.GetPayTodayList(page, this.pageSize, m)

	this.Data["pageTitle"] = "充值统计"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("LoggerController.PayTodayList", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 官方开奖结果
func (this *LoggerController) ExpectList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	code, _ := this.GetInt("code")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	if v, ok := entity.LotteryCodes[code]; ok {
		m["code"] = v
	}
	count, _ := service.LoggerService.GetExpectTotal(m)
	list, _ := service.LoggerService.GetExpectList(page, this.pageSize, m)

	this.Data["pageTitle"] = "开奖日志"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("LoggerController.ExpectList", "status", status, "code", code, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.Data["codes"] = entity.LotteryTypes
	this.Data["code"] = code
	this.display()
}

// 房间单局记录
func (this *LoggerController) GameRecordList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	expect := this.GetString("expect")
	gtype, _ := this.GetInt("gtype")
	rtype, _ := this.GetInt("rtype")
	ltype, _ := this.GetInt("ltype")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	m["gametype"] = gtype
	m["roomtype"] = rtype
	m["lotterytype"] = ltype
	if len(expect) != 0 {
		m["expect"] = expect
	}
	count, _ := service.LoggerService.GetGameRecordTotal(m)
	list, _ := service.LoggerService.GetGameRecordList(page, this.pageSize, m)

	this.Data["pageTitle"] = "房间日志"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("LoggerController.GameRecordList", "status", status, "gtype", gtype, "rtype", rtype, "ltype", ltype, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.Data["gtypes"] = entity.GameTypes
	this.Data["rtypes"] = entity.RoomTypes
	this.Data["ltypes"] = entity.LotteryTypes
	this.Data["gtype"] = gtype
	this.Data["rtype"] = rtype
	this.Data["ltype"] = ltype
	this.display()
}

// 个人单局记录
func (this *LoggerController) UserRecordList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	userid := this.GetString("userid")
	expect := this.GetString("expect")
	gtype, _ := this.GetInt("gtype")
	rtype, _ := this.GetInt("rtype")
	ltype, _ := this.GetInt("ltype")
	robot, _ := this.GetBool("robot")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	m["gametype"] = gtype
	m["roomtype"] = rtype
	m["lotterytype"] = ltype
	m["robot"] = robot
	if len(userid) != 0 {
		m["userid"] = userid
	}
	if len(expect) != 0 {
		m["expect"] = expect
	}
	count, _ := service.LoggerService.GetUserRecordTotal(m)
	list, _ := service.LoggerService.GetUserRecordList(page, this.pageSize, m)

	this.Data["pageTitle"] = "单局日志"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("LoggerController.UserRecordList", "status", status, "gtype", gtype, "rtype", rtype, "ltype", ltype, "robot", robot, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.Data["gtypes"] = entity.GameTypes
	this.Data["rtypes"] = entity.RoomTypes
	this.Data["ltypes"] = entity.LotteryTypes
	this.Data["robots"] = entity.IsRobot
	this.Data["gtype"] = gtype
	this.Data["rtype"] = rtype
	this.Data["ltype"] = ltype
	this.Data["robot"] = robot
	this.display()
}

// 盈亏统计日志
func (this *LoggerController) ChipTodayList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	gtype, _ := this.GetInt("gtype")
	rtype, _ := this.GetInt("rtype")
	ltype, _ := this.GetInt("ltype")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	if gtype > 0 {
		m["gametype"] = gtype
	}
	if rtype >= 0 {
		m["roomtype"] = rtype
	}
	if ltype > 0 {
		m["lotterytype"] = ltype
	}
	gtypes := make(map[int]string)
	rtypes := make(map[int]string)
	ltypes := make(map[int]string)
	gtypes[-1] = "全部"
	for k, v := range entity.GameTypes {
		gtypes[k] = v
	}
	rtypes[-1] = "全部"
	for k, v := range entity.RoomTypes {
		rtypes[k] = v
	}
	ltypes[-1] = "全部"
	for k, v := range entity.LotteryTypes {
		ltypes[k] = v
	}
	count, _ := service.LoggerService.GetChipTodayTotal(m)
	list, _ := service.LoggerService.GetChipTodayList(page, this.pageSize, m)

	this.Data["pageTitle"] = "盈亏统计"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("LoggerController.ChipTodayList", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.Data["gtypes"] = gtypes
	this.Data["rtypes"] = rtypes
	this.Data["ltypes"] = ltypes
	this.Data["gtype"] = gtype
	this.Data["rtype"] = rtype
	this.Data["ltype"] = ltype
	this.display()
}

// 后台充值赠送日志
func (this *LoggerController) GiveRecordList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	userid := this.GetString("userid")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	if userid != "" {
		m["userid"] = userid
	}
	m["type"] = entity.LogType9
	count, _ := service.LoggerService.GetChipTotal(m)
	list, _ := service.LoggerService.GetChipList(page, this.pageSize, m)

	this.Data["pageTitle"] = "后台充值记录"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("LoggerController.GiveRecordList", "userid", userid, "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 当前账号后台充值赠送日志
//TODO 写入记录的时候直接记录代理邀请码,方便查找
func (this *LoggerController) MineGiveRecordList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	userid := this.GetString("userid")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	m["type"] = entity.LogType9
	m["userid"] = bson.M{"$in": this.auth.GetUser().Players}
	if userid != "" {
		for _, v := range this.auth.GetUser().Players {
			if v == userid {
				m["userid"] = userid
				break
			}
		}
	}
	count, _ := service.LoggerService.GetChipTotal(m)
	list, _ := service.LoggerService.GetChipList(page, this.pageSize, m)

	this.Data["pageTitle"] = "后台充值记录"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("LoggerController.MineGiveRecordList", "userid", userid, "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 账务统计日志
func (this *LoggerController) AccountingList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	count, _ := service.LoggerService.GetAccountingTotal(m)
	list, _ := service.LoggerService.GetAccountingList(page, this.pageSize, m)

	this.Data["pageTitle"] = "账务统计"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("LoggerController.AccountingList", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}
