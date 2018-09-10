package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
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

type PlayerController struct {
	BaseController
}

// 玩家列表
func (this *PlayerController) List() {
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
		m["_id"] = userid
	}
	count, _ := service.PlayerService.GetTotal(typeId, m)
	list, _ := service.PlayerService.GetList(typeId, page, this.pageSize, m)

	ids := make([]string, 0)
	for _, v := range list {
		ids = append(ids, v.Userid)
	}

	//请求服务器
	if len(ids) != 0 {
		result, err := service.GmRequest(pb.WebOnline, pb.CONFIG_UPSERT, ids)
		if err != nil {
			flash := beego.NewFlash()
			flash.Error(fmt.Sprintf("%v", err))
			flash.Store(&this.Controller)
		} else {
			if resp, ok := result.(map[string]int); ok {
				for k, v := range list {
					list[k].State = resp[v.Userid]
				}
			}
		}
	}

	typeList := map[int]string{
		0: "注册用户",
		1: "游客用户",
		2: "机器人",
		3: "全部玩家",
		//4: "其它用户",
		//5: "微信用户",
	}

	this.Data["pageTitle"] = "玩家列表"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["typeList"] = typeList
	this.Data["typeId"] = typeId
	this.Data["userid"] = userid
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("PlayerController.List", "status", status, "type_id", typeId, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 钻石操作
func (this *PlayerController) Edit() {
	userid := this.GetString("userid")

	types2 := map[int]string{
		//1: "钻石",
		//2: "金币",
		//3: "房卡",
		4: "筹码",
		//5: "VIP(元)",
	}

	if this.isPost() {
		diamond, err := this.GetInt("diamond")
		propid, _ := this.GetInt("propid")
		if userid == "" {
			this.checkError(fmt.Errorf("用户ID不能为空"))
		}
		if err != nil {
			// handle error
			this.checkError(err)
		}
		msg := new(pb.PayCurrency)
		msg.Userid = userid
		msg.Type = entity.LogType9
		switch propid {
		case 1:
			msg.Diamond = int64(diamond)
		case 2:
			msg.Coin = int64(diamond)
		case 3:
			msg.Card = int64(diamond)
		case 4:
			diamond = diamond * 100 //转换为分
			msg.Chip = int64(diamond)
		case 5:
			diamond = diamond * 100 //转换为分
			msg.Money = int64(diamond)
		}
		if diamond != 0 {
			result, err := service.GmRequest(pb.WebGive, pb.CONFIG_UPSERT, msg)
			beego.Trace("result: ", result)
			if err != nil {
				this.checkError(err)
				//flash := beego.NewFlash()
				//flash.Error(fmt.Sprintf("%v", err))
				//flash.Store(&this.Controller)
			} else {
				//arg := &entity.LogGiveRecord{
				//	UserName: this.auth.GetUser().UserName,
				//	Userid:   userid,
				//	Propid:   propid,
				//	Num:      diamond,
				//}
				//service.LoggerService.SaveGiveRecord(arg)
				service.ActionService.UpdateDiamond(this.auth.GetUser().UserName,
					utils.String(entity.LogType9), userid, utils.String(diamond))
			}
		}
		this.redirect(beego.URLFor("PlayerController.List"))
	} else {
		p, err := service.PlayerService.GetPlayer(userid)
		this.checkError(err)

		this.Data["player"] = p
		this.Data["types2"] = types2
		this.Data["pageTitle"] = "货币操作"
		this.display()
	}
}

// 当前账号添加的代理玩家列表
func (this *PlayerController) MineList() {
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
		m["_id"] = userid
	}

	//是否是当前账号添加的代理
	user := this.auth.GetUser()
	if user == nil {
		this.checkError(fmt.Errorf("用户不存在"))
	}
	m["agent"] = bson.M{"$in": user.Child}

	count, _ := service.PlayerService.GetTotal(typeId, m)
	list, _ := service.PlayerService.GetList(typeId, page, this.pageSize, m)

	ids := make([]string, 0)
	for _, v := range list {
		ids = append(ids, v.Userid)
	}

	//请求服务器
	if len(ids) != 0 {
		result, err := service.GmRequest(pb.WebOnline, pb.CONFIG_UPSERT, ids)
		if err != nil {
			flash := beego.NewFlash()
			flash.Error(fmt.Sprintf("%v", err))
			flash.Store(&this.Controller)
		} else {
			if resp, ok := result.(map[string]int); ok {
				for k, v := range list {
					list[k].State = resp[v.Userid]
				}
			}
		}
	}

	typeList := map[int]string{
		0: "注册用户",
		1: "游客用户",
		2: "机器人",
		3: "全部玩家",
		//4: "其它用户",
		//5: "微信用户",
	}

	this.Data["pageTitle"] = "玩家列表"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["typeList"] = typeList
	this.Data["typeId"] = typeId
	this.Data["userid"] = userid
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("PlayerController.MineList", "status", status, "type_id", typeId, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 钻石操作
func (this *PlayerController) MineEdit() {
	userid := this.GetString("userid")

	types2 := map[int]string{
		//1: "钻石",
		//2: "金币",
		//3: "房卡",
		4: "筹码",
		//5: "VIP(元)",
	}

	if this.isPost() {
		diamond, err := this.GetInt("diamond")
		propid, _ := this.GetInt("propid")
		if userid == "" {
			this.checkError(fmt.Errorf("用户ID不能为空"))
		}
		if err != nil {
			// handle error
			this.checkError(err)
		}
		err2 := service.PlayerService.MineChild(userid, this.auth.GetUser().Child)
		if err2 != nil {
			this.checkError(err2)
		}
		msg := new(pb.PayCurrency)
		msg.Userid = userid
		msg.Type = entity.LogType9
		switch propid {
		case 1:
			msg.Diamond = int64(diamond)
		case 2:
			msg.Coin = int64(diamond)
		case 3:
			msg.Card = int64(diamond)
		case 4:
			diamond = diamond * 100 //转换为分
			msg.Chip = int64(diamond)
		case 5:
			diamond = diamond * 100 //转换为分
			msg.Money = int64(diamond)
		}
		if diamond != 0 {
			result, err := service.GmRequest(pb.WebGive, pb.CONFIG_UPSERT, msg)
			beego.Trace("result: ", result)
			if err != nil {
				this.checkError(err)
				//flash := beego.NewFlash()
				//flash.Error(fmt.Sprintf("%v", err))
				//flash.Store(&this.Controller)
			} else {
				//arg := &entity.LogGiveRecord{
				//	UserName: this.auth.GetUser().UserName,
				//	Userid:   userid,
				//	Propid:   propid,
				//	Num:      diamond,
				//}
				//service.LoggerService.SaveGiveRecord(arg)
				service.ActionService.UpdateDiamond(this.auth.GetUser().UserName,
					utils.String(entity.LogType9), userid, utils.String(diamond))
			}
		}
		this.redirect(beego.URLFor("PlayerController.MineList"))
	} else {
		p, err := service.PlayerService.GetPlayer(userid)
		this.checkError(err)

		this.Data["player"] = p
		this.Data["types2"] = types2
		this.Data["pageTitle"] = "货币操作"
		this.display()
	}
}

// 打印房间数据
func (this *PlayerController) Desk() {
	userid := this.GetString("userid")
	rtype, _ := this.GetInt("rtype")

	reqMsg := &entity.ReqRoomMsg{
		Userid: userid,
		Rtype:  rtype,
	}
	data, err1 := json.Marshal(reqMsg)
	this.checkError(err1)
	b, err2 := service.Gm("ReqRoomMsg", string(data))
	if err2 == nil {
		resp := new(entity.RespRoomMsg)
		err3 := json.Unmarshal([]byte(b), resp)
		fmt.Printf("resp %#v, err3 %v\n", resp, err3)
		beego.Trace("resp %#v, err3 %v", resp, err3)
	}
	this.checkError(err2)

	this.redirect(beego.URLFor("PlayerController.List"))
}

// 绑定代理操作
func (this *PlayerController) Build() {
	userid := this.GetString("userid")
	if this.isPost() {
		agent := this.GetString("agent")
		if userid == "" {
			this.checkError(fmt.Errorf("用户ID不能为空"))
		} else {
			reqMsg := &entity.ReqBuildMsg{
				Userid: userid,
				Agent:  agent,
			}
			data, err1 := json.Marshal(reqMsg)
			this.checkError(err1)
			_, err2 := service.Gm("ReqBuildMsg", string(data))
			this.checkError(err2)

			service.ActionService.UpdateBuild(this.auth.GetUser().UserName, userid, agent)
		}
		this.redirect(beego.URLFor("PlayerController.List"))
	} else {
		p, err := service.PlayerService.GetPlayer(userid)
		this.checkError(err)

		this.Data["player"] = p
		this.Data["pageTitle"] = "绑定操作"
		this.display()
	}
}

// 公告广播
func (this *PlayerController) NoticeList() {
	status, _ := this.GetInt("status")
	page, _ := this.GetInt("page")
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	// 过期处理
	if status == 0 { //未过期
		m["etime"] = bson.M{"$gte": bson.Now()}
	} else { //已过期
		m["etime"] = bson.M{"$lt": bson.Now()}
	}
	m["del"] = status
	list, _ := service.PlayerService.GetNoticeList(page, this.pageSize, m)
	count, _ := service.PlayerService.GetNoticeListTotal(m)

	this.Data["pageTitle"] = "公告列表"
	this.Data["status"] = status
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("PlayerController.NoticeList", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 添加公告
func (this *PlayerController) NoticeAdd() {
	if this.isPost() {
		notice := new(entity.Notice)
		rtype, _ := this.GetInt("rtype")
		atype, _ := this.GetInt("atype")
		acttype, _ := this.GetInt("acttype")
		top, _ := this.GetInt("top")
		num, _ := this.GetInt("num")
		content := this.GetString("content")
		etime := this.GetString("end_date")
		e := fmt.Sprintf("%s 23:59:59", etime)
		endTime := utils.Str2Time(e) //过期时间
		//fmt.Println("e : ", e, endTime)
		if endTime.IsZero() {
			this.checkError(errors.New("参数错误"))
		}
		notice.Rtype = rtype
		notice.Atype = uint32(atype)
		notice.Acttype = acttype
		notice.Top = top
		notice.Num = num
		notice.Content = content
		notice.Etime = endTime
		err := this.validNotice(notice)
		this.checkError(err)

		err = service.PlayerService.AddNotice(notice)
		this.checkError(err)
		service.ActionService.AddNotice(this.auth.GetUser().UserName, notice.Id)
		this.redirect(beego.URLFor("PlayerController.NoticeList"))
	}

	types1 := map[int]string{
		0: "显示消息",
		1: "支付消息",
		2: "活动消息",
	}

	types2 := map[int]string{
		1: "活动公告",
		2: "广播消息",
	}

	tops := map[int]string{
		0: "否",
		1: "是",
	}

	types4, _ := service.UserService.GetPlayerAtypes()

	this.Data["pageTitle"] = "添加公告"
	this.Data["types1"] = types1
	this.Data["types2"] = types2
	this.Data["tops"] = tops
	this.Data["types4"] = types4
	this.display()
}

func (this *PlayerController) validNotice(notice *entity.Notice) error {
	valid := validation.Validation{}
	valid.Required(notice.Rtype, "rtype").Message("消息类型不能为空")
	valid.Required(notice.Acttype, "acttype").Message("操作类型不能为空")
	valid.Required(notice.Num, "num").Message("公告次数不能为空")
	valid.Required(notice.Content, "content").Message("公告内容不能为空")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return errors.New(err.Message)
		}
	}

	return nil
}

// 公告广播
func (this *PlayerController) Notice() {
	id := this.GetString("id")

	notice, err := service.PlayerService.GetNotice(id)
	this.checkError(err)

	reqMsg := &entity.ReqNoticeMsg{
		Id:      notice.Id,
		Rtype:   notice.Rtype,
		Atype:   notice.Atype,
		Acttype: notice.Acttype,
		Top:     notice.Top,
		Num:     notice.Num,
		Del:     notice.Del, //是否移除
		Content: notice.Content,
		Etime:   notice.Etime,
		Ctime:   notice.Ctime,
	}
	data, err1 := json.Marshal(reqMsg)
	this.checkError(err1)
	_, err2 := service.Gm("ReqNoticeMsg", string(data))
	this.checkError(err2)

	service.ActionService.Notice(this.auth.GetUserName(), id)

	this.redirect(beego.URLFor("PlayerController.NoticeList"))
}

// 移除公告广播
func (this *PlayerController) NoticeDel() {
	id := this.GetString("id")

	notice, err := service.PlayerService.GetNotice(id)
	this.checkError(err)

	reqMsg := &entity.ReqNoticeMsg{
		Id:      notice.Id,
		Rtype:   notice.Rtype,
		Atype:   notice.Atype,
		Acttype: notice.Acttype,
		Top:     notice.Top,
		Num:     notice.Num,
		Del:     1, //是否移除
		Content: notice.Content,
		Etime:   notice.Etime,
		Ctime:   notice.Ctime,
	}
	data, err1 := json.Marshal(reqMsg)
	this.checkError(err1)
	_, err2 := service.Gm("ReqNoticeMsg", string(data))
	this.checkError(err2)

	if err2 == nil {
		err3 := service.PlayerService.DelNotice(id)
		this.checkError(err3)
	}

	service.ActionService.DelNotice(this.auth.GetUserName(), id)

	this.redirect(beego.URLFor("PlayerController.NoticeList"))
}

// 商城列表
func (this *PlayerController) ShopList() {
	status, _ := this.GetInt("status")
	page, _ := this.GetInt("page")
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	// 过期处理
	if status == 0 { //未过期
		m["etime"] = bson.M{"$gte": bson.Now()}
	} else { //已过期
		m["etime"] = bson.M{"$lt": bson.Now()}
	}
	m["del"] = status
	list, _ := service.PlayerService.GetShopList(page, this.pageSize, m)
	count, _ := service.PlayerService.GetShopListTotal(m)

	this.Data["pageTitle"] = "商城列表"
	this.Data["status"] = status
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("PlayerController.ShopList", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 添加公告
func (this *PlayerController) ShopAdd() {
	if this.isPost() {
		id, _ := this.GetInt("id")
		status, _ := this.GetInt("status")
		propid, _ := this.GetInt("propid")
		payway, _ := this.GetInt("payway")
		atype, _ := this.GetInt("atype")
		number, _ := this.GetInt("number")
		price, _ := this.GetInt("price")
		name := this.GetString("name")
		info := this.GetString("info")
		etime := this.GetString("end_date")
		e := fmt.Sprintf("%s 23:59:59", etime)
		endTime := utils.Str2Time(e) //过期时间
		//fmt.Println("e : ", e, endTime)
		if endTime.IsZero() {
			this.checkError(errors.New("参数错误"))
		}
		if id > 0 {
			shop := new(entity.Shop)
			shop.Id = utils.String(id)
			shop.Atype = uint32(atype) //分包类型
			shop.Status = status
			shop.Propid = propid
			shop.Payway = payway
			shop.Number = uint32(number)
			shop.Price = uint32(price)
			shop.Name = name
			shop.Info = info
			shop.Etime = endTime
			err := this.validShop(shop)
			this.checkError(err)

			err = service.PlayerService.AddShop(shop)
			this.checkError(err)
			service.ActionService.AddShop(this.auth.GetUser().UserName, shop.Id)
			this.redirect(beego.URLFor("PlayerController.ShopList"))
		} else {
			this.checkError(errors.New("ID错误"))
		}
	}

	types1 := map[int]string{
		1: "热卖",
		2: "普通",
	}

	types2 := map[int]string{
		1: "钻石",
		2: "金币",
	}

	types3 := map[int]string{
		1: "RMB",
		2: "钻石",
	}

	types4, _ := service.UserService.GetPlayerAtypes()

	this.Data["pageTitle"] = "添加商品"
	this.Data["types1"] = types1
	this.Data["types2"] = types2
	this.Data["types3"] = types3
	this.Data["types4"] = types4
	this.display()
}

func (this *PlayerController) validShop(shop *entity.Shop) error {
	valid := validation.Validation{}
	valid.Required(shop.Name, "name").Message("物品名称不能为空")
	valid.Required(shop.Info, "info").Message("物品描述不能为空")
	valid.Required(shop.Number, "number").Message("购买数量不能为空")
	//valid.Range(shop.Number, 1, 5000, "number").Message("购买数量不对")
	valid.Required(shop.Price, "price").Message("购买价格不能为空")
	//valid.Range(shop.Price, 1, 5000, "price").Message("购买价格不对")
	valid.Required(shop.Propid, "propid").Message("购买的物品不能为空")
	//valid.Range(shop.Propid, 1, 10, "propid").Message("购买的物品不对")
	valid.Required(shop.Payway, "payway").Message("支付方式不能为空")
	//valid.Range(shop.Payway, 1, 10, "payway").Message("支付方式不对")
	valid.Required(shop.Status, "status").Message("物品状态不能为空")
	//valid.Range(shop.Status, 1, 100, "status").Message("物品状态不对")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return errors.New(err.Message)
		}
	}

	return nil
}

// 公告广播
func (this *PlayerController) Shop() {
	id := this.GetString("id")

	shop, err := service.PlayerService.GetShop(id)
	this.checkError(err)

	reqMsg := &entity.ReqShopMsg{
		Id:     shop.Id,     //购买ID
		Atype:  shop.Atype,  //分包类型
		Status: shop.Status, //物品状态,1=热卖
		Propid: shop.Propid, //兑换的物品,1=钻石
		Payway: shop.Payway, //支付方式,1=RMB
		Number: shop.Number, //兑换的数量
		Price:  shop.Price,  //支付价格
		Name:   shop.Name,   //物品名字
		Info:   shop.Info,   //物品信息
		Del:    shop.Del,    //是否移除
		Etime:  shop.Etime,  //过期时间
		Ctime:  shop.Ctime,  //创建时间
	}
	data, err1 := json.Marshal(reqMsg)
	this.checkError(err1)
	_, err2 := service.Gm("ReqShopMsg", string(data))
	this.checkError(err2)

	service.ActionService.Shop(this.auth.GetUserName(), id)

	this.redirect(beego.URLFor("PlayerController.ShopList"))
}

// 移除公告广播
func (this *PlayerController) ShopDel() {
	id := this.GetString("id")

	shop, err := service.PlayerService.GetShop(id)
	this.checkError(err)

	reqMsg := &entity.ReqShopMsg{
		Id:     shop.Id,     //购买ID
		Atype:  shop.Atype,  //分包类型
		Status: shop.Status, //物品状态,1=热卖
		Propid: shop.Propid, //兑换的物品,1=钻石
		Payway: shop.Payway, //支付方式,1=RMB
		Number: shop.Number, //兑换的数量
		Price:  shop.Price,  //支付价格
		Name:   shop.Name,   //物品名字
		Info:   shop.Info,   //物品信息
		Del:    1,           //是否移除
		Etime:  shop.Etime,  //过期时间
		Ctime:  shop.Ctime,  //创建时间
	}
	data, err1 := json.Marshal(reqMsg)
	this.checkError(err1)
	_, err2 := service.Gm("ReqShopMsg", string(data))
	this.checkError(err2)

	if err2 == nil {
		err3 := service.PlayerService.DelShop(id)
		this.checkError(err3)
	}

	service.ActionService.DelShop(this.auth.GetUserName(), id)

	this.redirect(beego.URLFor("PlayerController.ShopList"))
}

// 设置变量
func (this *PlayerController) EnvAdd() {

	types1 := entity.EnvTypeValue

	if this.isPost() {
		key, _ := this.GetInt("key")
		value, _ := this.GetInt("value")

		types2 := entity.EnvTypeKey
		if v, ok := types2[key]; ok {

			b := make(map[string]int32)
			b[v] = int32(value)
			result, err := service.GmRequest(pb.WebEnv, pb.CONFIG_UPSERT, b)
			beego.Trace("result: ", result, err)
			if err != nil {
				this.checkError(err)
			} else {
				service.ActionService.EnvAdd(this.auth.GetUserName(), v)
				s := new(entity.Env)
				s.Key = v
				s.Value = int32(value)
				err3 := service.PlayerService.AddEnv(s)
				this.checkError(err3)
			}
		} else {
			this.checkError(errors.New("参数错误"))
		}
		this.redirect(beego.URLFor("PlayerController.EnvList"))
	}

	this.Data["pageTitle"] = "添加变量"
	this.Data["types1"] = types1
	this.display()
}

// 删除变量
func (this *PlayerController) EnvDel() {
	key := this.GetString("key")

	var value int32
	b := make(map[string]int32)
	b[key] = int32(value)
	result, err := service.GmRequest(pb.WebEnv, pb.CONFIG_DELETE, b)
	beego.Trace("result: ", result, err)
	if err != nil {
		this.checkError(err)
	} else {
		service.ActionService.EnvDel(this.auth.GetUserName(), key)
		err3 := service.PlayerService.DelEnv(key)
		this.checkError(err3)
	}

	this.redirect(beego.URLFor("PlayerController.EnvList"))
}

// 变量列表
func (this *PlayerController) EnvList() {
	//status, _ := this.GetInt("status")
	page, _ := this.GetInt("page")
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")

	list, _ := service.PlayerService.GetEnvList(page, this.pageSize, m)

	for k, v := range list {
		name := entity.EnvKeyType[v.Key]
		v.Name = entity.EnvTypeValue[name]
		list[k] = v
	}

	this.Data["pageTitle"] = "环境变量"
	this.Data["list"] = list
	this.display()
}

// vip列表
func (this *PlayerController) VipList() {
	status, _ := this.GetInt("status")
	page, _ := this.GetInt("page")
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	list, _ := service.PlayerService.GetVipList(page, this.pageSize, m)
	count, _ := service.PlayerService.GetVipListTotal(m)

	le := len(list)
	for i := 0; i < le; i++ {
		list[i].Number = list[i].Number / 100 //转换为元
	}

	this.Data["pageTitle"] = "VIP列表"
	this.Data["status"] = status
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("PlayerController.VipList", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 添加
func (this *PlayerController) VipAdd() {
	if this.isPost() {
		id := this.GetString("id")
		level, _ := this.GetInt("level")
		number, _ := this.GetInt("number")
		pay, _ := this.GetInt("pay")
		prize, _ := this.GetInt("prize")
		kick, _ := this.GetInt("kick")
		s := new(entity.Vip)
		s.Id = id
		s.Level = level
		s.Number = uint32(number * 100) //转换为分
		s.Pay = uint32(pay)
		s.Prize = uint32(prize)
		s.Kick = kick
		err := this.validVip(s)
		this.checkError(err)

		if err == nil {
			err = service.PlayerService.AddVip(s)
			this.checkError(err)
			service.ActionService.AddVip(this.auth.GetUser().UserName, s.Id)
		}
		this.redirect(beego.URLFor("PlayerController.VipList"))
	}

	this.Data["pageTitle"] = "添加VIP"
	this.display()
}

func (this *PlayerController) validVip(s *entity.Vip) error {
	valid := validation.Validation{}
	valid.Required(s.Id, "id").Message("Id不能为空")
	valid.Required(s.Level, "level").Message("等级不能为空")
	valid.Required(s.Number, "number").Message("金额不能为空")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return errors.New(err.Message)
		}
	}

	return nil
}

// 发布
func (this *PlayerController) Vip() {
	id := this.GetString("id")

	shop, err := service.PlayerService.GetVip(id)
	this.checkError(err)

	reqMsg := &entity.ReqVipMsg{
		Id:     shop.Id,
		Level:  shop.Level,
		Number: shop.Number,
		Pay:    shop.Pay,
		Prize:  shop.Prize,
		Kick:   shop.Kick,
		Ctime:  shop.Ctime,
	}
	data, err1 := json.Marshal(reqMsg)
	this.checkError(err1)
	_, err2 := service.Gm("ReqVipMsg", string(data))
	this.checkError(err2)

	service.ActionService.Vip(this.auth.GetUserName(), id)

	this.redirect(beego.URLFor("PlayerController.VipList"))
}

// 移除公告广播
func (this *PlayerController) VipDel() {
	id := this.GetString("id")

	shop, err := service.PlayerService.GetVip(id)
	this.checkError(err)

	reqMsg := &entity.ReqVipMsg{
		Id:     shop.Id,
		Level:  shop.Level,
		Number: shop.Number,
		Pay:    shop.Pay,
		Prize:  shop.Prize,
		Kick:   shop.Kick,
		Del:    1,
		Ctime:  shop.Ctime,
	}
	data, err1 := json.Marshal(reqMsg)
	this.checkError(err1)
	_, err2 := service.Gm("ReqVipMsg", string(data))
	this.checkError(err2)

	if err2 == nil {
		err3 := service.PlayerService.DelVip(id)
		this.checkError(err3)
	}

	service.ActionService.DelVip(this.auth.GetUserName(), id)

	this.redirect(beego.URLFor("PlayerController.VipList"))
}
