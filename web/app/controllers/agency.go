package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"goplays/pb"
	"goplays/web/app/entity"
	"goplays/web/app/libs"
	"goplays/web/app/service"
	"utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"gopkg.in/mgo.v2/bson"
)

type AgencyController struct {
	BaseController
}

// 代理商列表
func (this *AgencyController) AgencyList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "create_time", "create_time")
	m["status"] = status
	count, _ := service.AgencyService.GetAgencyTotal(m)
	list, _ := service.AgencyService.GetAgencyList(page, this.pageSize, m)

	//le := len(list)
	//for i := 0; i < le; i++ {
	//	list[i].Cash = list[i].Cash / 100       //转换为元
	//	list[i].Extract = list[i].Extract / 100 //转换为元
	//}

	this.Data["pageTitle"] = "总代列表"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgencyController.AgencyList", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 当前账号代理商列表
func (this *AgencyController) MineAgencyList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "create_time", "create_time")
	m["status"] = status

	//是否是当前账号添加的代理
	m["parent"] = this.auth.GetUserId()

	count, _ := service.AgencyService.GetAgencyTotal(m)
	list, _ := service.AgencyService.GetAgencyList(page, this.pageSize, m)

	//le := len(list)
	//for i := 0; i < le; i++ {
	//	list[i].Cash = list[i].Cash / 100       //转换为元
	//	list[i].Extract = list[i].Extract / 100 //转换为元
	//}

	this.Data["pageTitle"] = "总代列表"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgencyController.MineAgencyList", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 赠送/扣除钻石操作
func (this *AgencyController) AgencyGive() {
	id := this.GetString("id")
	if id == "" {
		this.checkError(fmt.Errorf("用户ID不能为空"))
	}
	user, err := service.UserService.GetUser(id, false)
	this.checkError(err)
	if this.isPost() {
		diamond, err := this.GetInt("diamond")
		if err != nil {
			// handle error
			this.checkError(err)
		}

		if user.Agent == "" {
			this.checkError(fmt.Errorf("代理商不存在"))
		}

		reqMsg := &entity.ReqMsg{
			Userid: user.Agent,
			Rtype:  int(entity.LogType9),
			Itemid: entity.ITYPE1,
			Amount: int32(diamond),
		}
		data, err1 := json.Marshal(reqMsg)
		this.checkError(err1)
		_, err2 := service.Gm("ReqMsg", string(data))
		this.checkError(err2)

		service.ActionService.UpdateDiamond(this.auth.GetUser().UserName,
			utils.String(entity.LogType9), user.Agent, utils.String(diamond))
		this.redirect(beego.URLFor("AgencyController.AgencyList"))
	}

	this.Data["user"] = user
	this.Data["pageTitle"] = "钻石操作"
	this.display()
}

// 代理编辑
func (this *AgencyController) AgencyEdit() {
	id := this.GetString("id")
	if id == "" {
		this.checkError(fmt.Errorf("用户ID不能为空"))
	}
	user, err := service.UserService.GetUser(id, false)
	this.checkError(err)
	if this.isPost() {
		rate, err := this.GetInt("rate")
		if err != nil {
			this.checkError(err)
		}
		sysrate, err := this.GetInt("sysrate")
		if err != nil {
			this.checkError(err)
		}

		if user.Agent == "" {
			this.checkError(fmt.Errorf("代理商不存在"))
		}

		if rate >= 0 && rate <= 100 {
			user.Rate = uint32(rate)
		} else {
			this.checkError(fmt.Errorf("提取率错误"))
		}
		if sysrate >= 0 && sysrate <= 100 {
			user.SysRate = uint32(sysrate)
		} else {
			this.checkError(fmt.Errorf("提取率错误"))
		}

		service.AgencyService.UpdateAgencyRate2(user)

		service.ActionService.UpdateAgency(this.auth.GetUser().UserName,
			user.Agent, utils.String(rate))
		this.redirect(beego.URLFor("AgencyController.AgencyList"))
	}

	this.Data["user"] = user
	this.Data["pageTitle"] = "编辑操作"
	this.display()
}

// 代理编辑
func (this *AgencyController) MineAgencyEdit() {
	id := this.GetString("id")
	if id == "" {
		this.checkError(fmt.Errorf("用户ID不能为空"))
	}
	user, err := service.UserService.GetUser(id, false)
	this.checkError(err)
	if this.isPost() {
		rate, err := this.GetInt("rate")
		if err != nil {
			this.checkError(err)
		}
		sysrate, err := this.GetInt("sysrate")
		if err != nil {
			this.checkError(err)
		}

		if user.Agent == "" {
			this.checkError(fmt.Errorf("代理商不存在"))
		}

		if rate >= 0 && rate <= 100 {
			user.Rate = uint32(rate)
		} else {
			this.checkError(fmt.Errorf("提取率错误"))
		}
		if sysrate >= 0 && sysrate <= 100 {
			user.SysRate = uint32(sysrate)
		} else {
			this.checkError(fmt.Errorf("提取率错误"))
		}

		service.AgencyService.UpdateAgencyRate2(user)

		service.ActionService.UpdateAgency(this.auth.GetUser().UserName,
			user.Agent, utils.String(rate))
		this.redirect(beego.URLFor("AgencyController.MineAgencyList"))
	}

	this.Data["user"] = user
	this.Data["pageTitle"] = "编辑操作"
	this.display()
}

// 代理编辑
func (this *AgencyController) SetDirectlyAgency() {
	id := this.GetString("id")
	user, err := service.UserService.GetUserByAgent(id)
	this.checkError(err)
	if this.isPost() {
		rate, err := this.GetInt("rate")
		if err != nil {
			this.checkError(err)
		}

		if rate >= 0 && rate <= 100 {
			user.Rate = uint32(rate)
		} else {
			this.checkError(fmt.Errorf("提取率错误"))
		}

		//代理关系
		agent := this.auth.GetUser().Agent
		//if service.AgencyService.Parental(agent, user.ParentAgent) {
		if user.ParentAgent == agent && agent != "" {
			service.AgencyService.UpdateAgencyRate(user)
		} else {
			beego.Error("SetDirectlyAgency fail err: ", agent, id, user.ParentAgent)
			this.checkError(fmt.Errorf("代理商不存在"))
		}

		service.ActionService.UpdateAgency(this.auth.GetUser().UserName,
			user.Agent, utils.String(rate))
		this.redirect(beego.URLFor("AgencyController.DirectlyAgency"))
	}

	this.Data["user"] = user
	this.Data["pageTitle"] = "设置操作"
	this.display()
}

/*
// 赠送/扣除钻石操作(只存在后台)
func (this *AgencyController) Give2Agency_bak() {
	id := this.GetString("id")
	if id == "" {
		this.checkError(fmt.Errorf("用户ID不能为空"))
	}
	user, err := service.UserService.GetUser(id, false)
	this.checkError(err)
	if this.isPost() {
		diamond, err := this.GetInt("diamond")
		if err != nil {
			// handle error
			this.checkError(err)
		} else if user.Agent == "" {
			this.checkError(fmt.Errorf("代理商不存在"))
		} else {
			if diamond < 0 {
				n := int32(user.Number) + int32(diamond)
				if n < 0 {
					n = 0
				}
				user.Number = uint32(n)
			} else {
				user.Number += uint32(diamond)
			}
			err3 := service.AgencyService.UpdateNumber(user)
			if err3 != nil {
				this.checkError(fmt.Errorf("发放失败"))
			} else {
				service.ActionService.UpdateNumber(this.auth.GetUser().UserName,
					utils.String(entity.LogType27), user.Agent, utils.String(diamond))
				this.redirect(beego.URLFor("AgencyController.AgencyList"))
			}
		}
	}

	this.Data["user"] = user
	this.Data["pageTitle"] = "发放钻石"
	this.display()
}
*/

// 赠送/扣除钻石操作(只存在后台)
func (this *AgencyController) Give2Agency() {
	id := this.GetString("id")
	if id == "" {
		this.checkError(fmt.Errorf("用户ID不能为空"))
	}
	user, err := service.UserService.GetUser(id, false)
	this.checkError(err)
	if this.isPost() {
		diamond, err := this.GetInt("diamond")
		if err != nil {
			// handle error
			this.checkError(err)
		} else if user.Agent == "" {
			this.checkError(fmt.Errorf("代理商不存在"))
		} else {
			err3 := service.AgencyService.SetChip(id, int64(diamond))
			if err3 != nil {
				this.checkError(fmt.Errorf("发放失败"))
			} else {
				service.ActionService.UpdateNumber(this.auth.GetUser().UserName,
					utils.String(entity.LogType27), user.Agent, utils.String(diamond))
				this.redirect(beego.URLFor("AgencyController.AgencyList"))
			}
		}
	}

	this.Data["user"] = user
	this.Data["pageTitle"] = "发放筹码"
	this.display()
}

// 代理商赠送/扣除钻石操作,同时扣除代理商金额(从后台扣除)
func (this *AgencyController) MyAgencyGive() {
	userid := this.GetString("id")
	if userid == "" {
		this.checkError(fmt.Errorf("用户ID不能为空"))
	} else if this.isPost() {
		diamond, err := this.GetInt("diamond")
		user := this.auth.GetUser()
		if err != nil {
			this.checkError(err)
		} else if user.Agent == "" || diamond <= 0 {
			this.checkError(fmt.Errorf("无法赠送"))
		} else {
			reqMsg := &entity.ReqMsg{
				Userid: userid,
				Rtype:  int(entity.LogType27),
				Itemid: entity.ITYPE1, //钻石
				Amount: int32(diamond),
			}
			data, err1 := json.Marshal(reqMsg)
			if err1 != nil {
				this.checkError(err1)
			} else if user.Number < uint32(diamond) {
				this.checkError(fmt.Errorf("余额不足,无法赠送"))
			} else {
				_, err3 := service.Gm("ReqMsg", string(data))
				if err3 != nil {
					this.checkError(err3)
				} else {
					user.Number -= uint32(diamond)
					user.Expend += uint32(diamond)
					err3 := service.AgencyService.UpdateExpend(user)
					this.checkError(err3)
					service.ActionService.UpdateExpend(user.UserName,
						utils.String(entity.LogType27), userid, utils.String(diamond))
				}
			}
		}
		this.redirect(beego.URLFor("AgencyController.MyAgencyList"))
	}

	p, err := service.PlayerService.GetPlayer(userid)
	this.checkError(err)

	this.Data["id"] = userid
	this.Data["nickname"] = p.Nickname
	this.Data["pageTitle"] = "发放钻石"
	this.display()
}

// 我的代理,绑定我的用户
func (this *AgencyController) MyAgencyList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "agent_time", "agent_time")
	username := this.auth.GetUser().UserName
	list, _ := service.AgencyService.GetMyAgencyList(username, page, this.pageSize, m)
	count, _ := service.AgencyService.GetMyAgencyTotal(username, m)

	ids := make([]string, 0)
	for _, v := range list {
		ids = append(ids, v.Userid)
	}
	ms := make(map[string]int)
	userList, _ := service.UserService.GetAgencyList(ids)
	for _, v := range userList {
		if v.Agent != "" {
			ms[v.Agent] = 1
		}
	}
	for k, v := range list {
		if _, ok := ms[v.Userid]; ok {
			v.Agency = 1 //是代理商
			list[k] = v
		}
	}

	this.Data["pageTitle"] = "我的代理"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgencyController.MyAgencyList", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 代理商赠送/扣除钻石操作,同时扣除代理商金额
func (this *AgencyController) MyAgencyEdit() {
	userid := this.GetString("id")
	if userid == "" {
		this.checkError(fmt.Errorf("用户ID不能为空"))
	} else if this.isPost() {
		diamond, err := this.GetInt("diamond")
		if err != nil {
			// handle error
			this.checkError(err)
		}

		agent := this.auth.GetUser().Agent
		if agent == "" || diamond <= 0 {
			this.checkError(fmt.Errorf("无法赠送"))
		} else {
			reqMsg := &entity.ReqGiveDiamondMsg{
				Userid: userid,
				Agent:  agent,
				Rtype:  int(entity.LogType9),
				Itemid: entity.ITYPE1,
				Amount: int32(diamond),
			}
			data, err1 := json.Marshal(reqMsg)
			this.checkError(err1)
			_, err2 := service.Gm("ReqGiveDiamondMsg", string(data))
			this.checkError(err2)

			service.ActionService.UpdateDiamond(this.auth.GetUser().UserName,
				utils.String(entity.LogType9), userid, utils.String(diamond))
		}
		this.redirect(beego.URLFor("AgencyController.MyAgencyList"))
	}

	p, err := service.PlayerService.GetPlayer(userid)
	this.checkError(err)

	this.Data["id"] = userid
	this.Data["nickname"] = p.Nickname
	this.Data["pageTitle"] = "赠送钻石"
	this.display()
}

// 我的充值日志
func (this *AgencyController) MyAgencyPay() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	id := this.GetString("id")
	if page < 1 {
		page = 1
	}

	p, err := service.PlayerService.GetPlayer(id)
	agent := this.auth.GetUser().Agent
	if id == "" || err != nil || p.Agent != agent || agent == "" || p.Agent == "" {
		this.checkError(fmt.Errorf("非法操作"))
	} else {
		m := service.FindByDate(startDate, endDate, "ctime", "ctime")
		m["userid"] = id
		m["result"] = entity.TradeSuccess
		count, _ := service.LoggerService.GetPayTotal(m)
		list, _ := service.LoggerService.GetPayList(page, this.pageSize, m)

		le := len(list)
		for i := 0; i < le; i++ {
			list[i].Money = list[i].Money / 100 //转换为元
		}

		this.Data["pageTitle"] = "充值记录"
		this.Data["count"] = count
		this.Data["list"] = list
		this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgencyController.MyAgencyPay", "status", status, "id", id, "start_date", startDate, "end_date", endDate), true).ToString()
		this.Data["startDate"] = startDate
		this.Data["endDate"] = endDate
	}
	this.display()
}

// 提现记录
func (this *AgencyController) CashList() {
	status, _ := this.GetInt("status")
	page, _ := this.GetInt("page")
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	m["status"] = status
	list, _ := service.AgencyService.GetCashList(page, this.pageSize, m)
	count, _ := service.AgencyService.GetCashListTotal(m)

	//le := len(list)
	//for i := 0; i < le; i++ {
	//	list[i].Cash = list[i].Cash / 100 //转换为元展示
	//}

	this.Data["pageTitle"] = "提现记录"
	this.Data["status"] = status
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgencyController.CashList", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 当前账号提现记录
func (this *AgencyController) MineCashList() {
	status, _ := this.GetInt("status")
	page, _ := this.GetInt("page")
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	m["status"] = status

	//是否是当前账号添加的代理
	user := this.auth.GetUser()
	if user == nil {
		this.checkError(fmt.Errorf("账号不存在"))
	}
	m["agent"] = bson.M{"$in": user.Child}

	list, _ := service.AgencyService.GetCashList(page, this.pageSize, m)
	count, _ := service.AgencyService.GetCashListTotal(m)

	//le := len(list)
	//for i := 0; i < le; i++ {
	//	list[i].Cash = list[i].Cash / 100 //转换为元展示
	//}

	this.Data["pageTitle"] = "提现记录"
	this.Data["status"] = status
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgencyController.MineCashList", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 我的提现
func (this *AgencyController) MyCashList() {
	status, _ := this.GetInt("status")
	page, _ := this.GetInt("page")
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	m["status"] = status
	username := this.auth.GetUser().UserName
	list, _ := service.AgencyService.GetMyCashList(username, page, this.pageSize, m)
	count, _ := service.AgencyService.GetMyCashListTotal(username, m)

	//le := len(list)
	//for i := 0; i < le; i++ {
	//	list[i].Cash = list[i].Cash / 100 //转换为元展示
	//}

	this.Data["pageTitle"] = "我的提现"
	this.Data["status"] = status
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgencyController.MyCashList", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 申请提现,让玩家自己输入银行卡号，开户行，姓名
func (this *AgencyController) MyCashAdd() {
	payway := map[int]string{
		entity.PAYWAY1: "微信",
		entity.PAYWAY2: "支付宝",
		entity.PAYWAY3: "银行账号",
	}

	username := this.auth.GetUser().UserName
	user2, err := service.UserService.GetUserByName(username)
	this.checkError(err)
	if this.isPost() {
		name := this.GetString("name")
		bankCard, _ := this.GetInt("bankCard")
		bankAddr := this.GetString("bankAddr")
		//money, _ := this.GetFloat("money") //单位为元
		money, _ := this.GetInt("fee") //单位为元
		valid := validation.Validation{}
		valid.Required(name, "name").Message("姓名不能为空")
		valid.Required(bankCard, "bankCard").Message("收款方式不能为空")
		valid.Required(bankAddr, "bankAddr").Message("收款账号不能为空")
		//valid.Required(money, "money").Message("提现金额错误")
		valid.Required(money, "fee").Message("提现金额错误")
		if payway[bankCard] == "" {
			this.checkError(fmt.Errorf("收款方式错误"))
		} else if !valid.HasErrors() {
			username := this.auth.GetUser().UserName
			money = money * 100 //元转换为分
			err := service.AgencyService.ApplyCashAdd(username, name,
				bankAddr, bankCard, money)
			if err != nil {
				this.showMsg(err.Error(), MSG_ERR)
			} else {
				service.ActionService.AddApplyCash(username, utils.String(money))
				this.redirect(beego.URLFor("AgencyController.MyCashList"))
			}
		} else {
			for _, err := range valid.Errors {
				this.showMsg(err.Message, MSG_ERR)
				break
			}
			this.Data["name"] = name
			this.Data["bankCard"] = bankCard
			this.Data["bankAddr"] = bankAddr
			this.Data["money"] = money
		}
	}

	var FeeAll float64
	if user2 != nil {
		FeeAll = service.Chip2Float(user2.FeeAll)
	}

	this.Data["payway"] = payway
	this.Data["pageTitle"] = "申请提现"
	//this.Data["user2"] = user2
	this.Data["FeeAll"] = FeeAll
	this.display()
}

// 提现处理
func (this *AgencyController) AgencyExtract() {
	orderid := this.GetString("id")
	if orderid == "" {
		this.checkError(fmt.Errorf("订单不存在"))
	}

	username := this.auth.GetUser().UserName
	err := service.AgencyService.ExtractCash(username, orderid)
	if err != nil {
		this.checkError(err)
	} else {
		service.ActionService.ExtractApplyCash(username, orderid)
	}

	this.redirect(beego.URLFor("AgencyController.CashList"))
}

// 我的充值日志
func (this *AgencyController) MyPayList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	if page < 1 {
		page = 1
	}

	agent := this.auth.GetUser().Agent
	if agent == "" {
		this.checkError(fmt.Errorf("非法操作"))
	} else {
		m := service.FindByDate(startDate, endDate, "ctime", "ctime")
		m["userid"] = agent
		m["result"] = entity.TradeSuccess
		count, _ := service.LoggerService.GetPayTotal(m)
		list, _ := service.LoggerService.GetPayList(page, this.pageSize, m)

		le := len(list)
		for i := 0; i < le; i++ {
			list[i].Money = list[i].Money / 100 //转换为元
		}

		this.Data["pageTitle"] = "我的充值"
		this.Data["count"] = count
		this.Data["list"] = list
		this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgencyController.MyPayList", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
		this.Data["startDate"] = startDate
		this.Data["endDate"] = endDate
	}
	this.display()
}

// 历史开奖
func (this *AgencyController) HistoryOpen() {
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

	this.Data["pageTitle"] = "历史开奖"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgencyController.HistoryOpen", "status", status, "gtype", gtype, "rtype", rtype, "ltype", ltype, "start_date", startDate, "end_date", endDate), true).ToString()
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

// 历史开奖详情
func (this *AgencyController) OpenDetails() {
	id := this.GetString("id")

	record, _ := service.LoggerService.GetGameRecord(id)

	this.Data["pageTitle"] = "开奖详情"
	this.Data["trend"] = record.Trend
	this.Data["result"] = record.Result
	this.Data["record"] = record.Record
	this.Data["details"] = record.Details
	this.Data["id"] = id
	//this.Data["gtypes"] = entity.GameTypes
	//this.Data["rtypes"] = entity.RoomTypes
	//this.Data["ltypes"] = entity.LotteryTypes
	this.display()
}

// 玩家记录
func (this *AgencyController) UserRecord() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	userid := this.GetString("userid")
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
	if len(userid) != 0 {
		m["userid"] = userid
	}
	if len(expect) != 0 {
		m["expect"] = expect
	}
	count, _ := service.LoggerService.GetUserRecordTotal(m)
	list, _ := service.LoggerService.GetUserRecordList(page, this.pageSize, m)

	this.Data["pageTitle"] = "玩家记录"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgencyController.UserRecord", "status", status, "gtype", gtype, "rtype", rtype, "ltype", ltype, "start_date", startDate, "end_date", endDate, "userid", userid), true).ToString()
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

// 玩家记录
func (this *AgencyController) MineUserRecord() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	userid := this.GetString("userid")
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
	if len(userid) != 0 {
		//m["userid"] = userid
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
	if len(expect) != 0 {
		m["expect"] = expect
	}
	count, _ := service.LoggerService.GetUserRecordTotal(m)
	list, _ := service.LoggerService.GetUserRecordList(page, this.pageSize, m)

	this.Data["pageTitle"] = "玩家记录"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgencyController.MineUserRecord", "status", status, "gtype", gtype, "rtype", rtype, "ltype", ltype, "start_date", startDate, "end_date", endDate, "userid", userid), true).ToString()
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

// 玩家记录详情
func (this *AgencyController) UserRecordDetails() {
	id := this.GetString("id")
	userid := this.GetString("userid")

	record, _ := service.LoggerService.GetUserRecord(id, userid)

	this.Data["pageTitle"] = "玩家记录详情"
	this.Data["details"] = record.Details
	//this.Data["gtypes"] = entity.GameTypes
	//this.Data["rtypes"] = entity.RoomTypes
	//this.Data["ltypes"] = entity.LotteryTypes
	this.display()
}

/*
// 直属代理 agent == userid
func (this *AgencyController) DirectlyAgency() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	userid := this.GetString("userid")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "agent_time", "agent_time")
	if userid != "" {
		m["_id"] = userid
	}
	id := this.GetString("id")
	typeId, _ := strconv.Atoi(this.GetString("type_id"))
	var list []entity.PlayerUser
	var count int64
	//
	if len(id) == 0 {
		id = this.auth.GetUser().Agent
	}
	list, _ = service.AgencyService.GetMyAgencyList2(id, page, this.pageSize, m)
	count, _ = service.AgencyService.GetMyAgencyTotal2(id, m)

	ids := make([]string, 0)
	for _, v := range list {
		ids = append(ids, v.Userid)
	}
	ms := make(map[string]int)
	fs := make(map[string]int64)
	rs := make(map[string]uint32)
	userList, _ := service.UserService.GetAgencyLists(ids)
	for _, v := range userList {
		if v2, ok := v["agent"]; ok {
			ms[v2.(string)] = 1
			if v3, ok := v["fee_rate"]; ok {
				fs[v2.(string)] = v3.(int64)
			}
			if v3, ok := v["rate"]; ok {
				rs[v2.(string)] = uint32(v3.(int))
			}
		}
	}
	for k, v := range list {
		if _, ok := ms[v.Userid]; ok {
			v.Agency = 1 //是代理商
			v.FeeRate = fs[v.Userid]
			v.Rate = rs[v.Userid]
			list[k] = v
		}
	}

	types2 := map[int]string{
		1: "代理",
		2: "玩家",
	}

	//this.Data["pageTitle"] = "直属代理"
	this.Data["pageTitle"] = "代理"
	this.Data["id"] = id
	this.Data["types2"] = types2
	this.Data["typeId"] = typeId
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgencyController.DirectlyAgency", "status", status, "start_date", startDate, "end_date", endDate, "id", id), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.Data["userid"] = userid
	this.display()
}
*/

// 下属代理
func (this *AgencyController) UnderlingAgency() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "agent_time", "agent_time")
	//username := this.auth.GetUser().UserName
	username := this.GetString("id")
	list, _ := service.AgencyService.GetMyAgencyList2(username, page, this.pageSize, m)
	count, _ := service.AgencyService.GetMyAgencyTotal2(username, m)

	ids := make([]string, 0)
	for _, v := range list {
		ids = append(ids, v.Userid)
	}
	ms := make(map[string]int)
	userList, _ := service.UserService.GetAgencyList(ids)
	for _, v := range userList {
		if v.Agent != "" {
			ms[v.Agent] = 1
		}
	}
	for k, v := range list {
		if _, ok := ms[v.Userid]; ok {
			v.Agency = 1 //是代理商
			list[k] = v
		}
	}

	this.Data["pageTitle"] = "下属代理"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgencyController.UnderlingAgency", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 赠送/扣除钻石操作
func (this *AgencyController) Give2DirectlyAgency() {
	id := this.GetString("id")
	if id == "" {
		this.checkError(fmt.Errorf("用户ID不能为空"))
	}
	user, err := service.PlayerService.GetPlayer(id)
	this.checkError(err)
	username := this.auth.GetUser().UserName
	user2, err := service.UserService.GetUserByName(username)
	this.checkError(err)
	if this.isPost() {
		diamond, err := this.GetInt("chip")
		if err != nil {
			// handle error
			this.checkError(err)
		} else if user.Agent == "" {
			this.checkError(fmt.Errorf("未绑定玩家"))
		} else if user2.Chip < int64(diamond) {
			this.checkError(fmt.Errorf("余额不足"))
		} else {

			msg := new(pb.PayCurrency)
			msg.Userid = user.Userid
			msg.Type = entity.LogType9
			msg.Chip = int64(diamond)
			if diamond > 0 {
				result, err := service.GmRequest(pb.WebGive, pb.CONFIG_UPSERT, msg)
				beego.Trace("result: ", result, err)
				if err != nil {
					this.checkError(err)
					beego.Trace("err: ", err)
					//flash := beego.NewFlash()
					//flash.Error(fmt.Sprintf("%v", err))
					//flash.Store(&this.Controller)
				} else {
					//更新
					err3 := service.AgencyService.SetChip(user2.Id, (-1 * int64(diamond)))
					if err3 != nil {
						this.checkError(fmt.Errorf("发放失败"))
					} else {
						service.ActionService.UpdateChip(username,
							utils.String(entity.LogType9), user.Agent,
							utils.String(diamond))
					}
				}
			}

		}
		this.redirect(beego.URLFor("AgencyController.DirectlyAgency"))
	}

	this.Data["user"] = user
	this.Data["user2"] = user2
	this.Data["pageTitle"] = "筹码操作"
	this.display()
}

// 推广海报
func (this *AgencyController) Poster() {
	agent := this.auth.GetUser().Agent
	if agent == "" {
		//this.Abort("404")
		//agent = "888888"
		agent = "请升级为代理"
	}
	this.Data["id"] = agent
	this.Data["pageTitle"] = "推广海报"
	this.TplName = "agency/poster.html"
}

// 添加代理
func (this *AgencyController) AddAgency() {
	if this.isPost() {
		valid := validation.Validation{}

		username := this.GetString("username")
		//email := this.GetString("email")
		agent := this.GetString("agent")
		sex, _ := this.GetInt("sex")
		rate, _ := this.GetInt("rate")
		password1 := this.GetString("password1")
		password2 := this.GetString("password2")

		if rate >= 0 && rate <= 100 {
		} else {
			this.checkError(fmt.Errorf("反佣抽成设置错误"))
		}

		list := this.GetStrings("role_ids")
		if len(list) == 0 {
			list = []string{"2"} //TODO 暂时默认代理商
		}

		valid.Required(username, "username").Message("请输入用户名")
		valid.Range(rate, 0, 100, "rate").Message("反佣抽成设置错误")
		//valid.Required(email, "email").Message("请输入Email")
		//valid.Email(email, "email").Message("Email无效")
		valid.Required(agent, "agent").Message("请输入6位数字邀请码")
		valid.Required(password1, "password1").Message("请输入密码")
		valid.Required(password2, "password2").Message("请输入确认密码")
		valid.MinSize(password1, 6, "password1").Message("密码长度不能小于6个字符")
		valid.Match(password1, regexp.MustCompile(`^`+regexp.QuoteMeta(password2)+`$`), "password2").Message("两次输入的密码不一致")
		if valid.HasErrors() {
			for _, err := range valid.Errors {
				this.showMsg(err.Message, MSG_ERR)
			}
		}

		level := service.AgencyService.AgentLevel(this.auth.GetUser().Agent)
		if level >= 3 {
			this.checkError(errors.New("没有操作权限"))
		}

		if this.auth.GetUser().Agent == "" && this.auth.GetUser().Id != "1" {
			//this.checkError(errors.New("没有操作权限"))
		}

		parent_agent := this.auth.GetUser().Agent

		user, err := service.UserService.AddUser2(username, agent, parent_agent, this.auth.GetUserId(), password1, sex, rate)
		if err == nil {
			service.ActionService.AddUser(this.auth.GetUser().UserName, username)
			service.UserService.UpdateChild(this.auth.GetUser(), agent)
		}
		this.checkError(err)

		// 更新角色
		roleIds := make([]string, 0)
		//for _, v := range this.GetStrings("role_ids") {
		for _, v := range list {
			//if roleId, _ := strconv.Atoi(v); roleId > 0 {
			//	roleIds = append(roleIds, roleId)
			//}
			roleIds = append(roleIds, v)
		}
		service.UserService.UpdateUserRoles(user.Id, roleIds)

		this.redirect(beego.URLFor("AgencyController.DirectlyAgency"))
	}

	roleList2, _ := service.RoleService.GetAllRoles()
	roleList := make([]entity.Role, 0)
	//TODO 暂时只能选择代理商
	for _, v := range roleList2 {
		if v.Id == "2" {
			roleList = append(roleList, v)
			break
		}
	}
	//过滤掉超级管理员
	for k, v := range roleList {
		if v.Id == "1" {
			roleList = append(roleList[:k], roleList[k+1:]...)
			break
		}
	}

	this.Data["pageTitle"] = "添加代理"
	this.Data["roleList"] = roleList
	this.display()
}

// 直属代理 userid != agent
func (this *AgencyController) DirectlyAgency() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	userid := this.GetString("userid") //邀请码
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "agent_time", "agent_time")
	if userid != "" {
		m["agent"] = userid
	}
	//
	id := this.GetString("id")
	if id == "" {
		id = this.auth.GetUser().Agent
	}
	//
	var list []entity.User
	var count int64
	list, _ = service.AgencyService.GetMyAgencyList3(id, page, this.pageSize, m)
	count, _ = service.AgencyService.GetMyAgencyTotal3(id, m)

	//this.Data["pageTitle"] = "直属代理"
	this.Data["pageTitle"] = "代理"
	this.Data["id"] = id
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgencyController.DirectlyAgency", "status", status, "start_date", startDate, "end_date", endDate, "id", id), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.Data["userid"] = userid
	this.display()
}

// 直属玩家 userid != agent
func (this *AgencyController) DirectlyPlayer() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	userid := this.GetString("userid") //邀请码
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "agent_time", "agent_time")
	if userid != "" {
		m["_id"] = userid
	}
	//
	id := this.GetString("id")
	if id == "" {
		id = this.auth.GetUser().Agent
	}
	//
	var list []entity.PlayerUser
	var count int64
	list, _ = service.AgencyService.GetMyAgencyList4(id, page, this.pageSize, m)
	count, _ = service.AgencyService.GetMyAgencyTotal4(id, m)

	//this.Data["pageTitle"] = "直属代理"
	this.Data["pageTitle"] = "玩家"
	this.Data["id"] = id
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgencyController.DirectlyPlayer", "status", status, "start_date", startDate, "end_date", endDate, "id", id), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.Data["userid"] = userid
	this.display()
}
