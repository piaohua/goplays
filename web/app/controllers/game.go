package controllers

import (
	"errors"
	"fmt"

	"goplays/pb"
	"goplays/web/app/entity"
	"goplays/web/app/libs"
	"goplays/web/app/service"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

// 游戏列表
func (this *PlayerController) GameList() {
	status, _ := this.GetInt("status")
	page, _ := this.GetInt("page")
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	m["del"] = status
	list, _ := service.PlayerService.GetGameList(page, this.pageSize, m)
	count, _ := service.PlayerService.GetGameListTotal(m)

	le := len(list)
	for i := 0; i < le; i++ {
		list[i].Chip = list[i].Chip / 100   //转换为元
		list[i].Down = list[i].Down / 100   //转换为元
		list[i].Carry = list[i].Carry / 100 //转换为元
		list[i].Top = list[i].Top / 100     //转换为元
	}

	this.Data["pageTitle"] = "游戏列表"
	this.Data["status"] = status
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("PlayerController.GameList", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 添加公告
func (this *PlayerController) GameAdd() {
	if this.isPost() {
		node, _ := this.GetInt("node")
		gtype, _ := this.GetInt("gtype")
		rtype, _ := this.GetInt("rtype")
		roomstatus, _ := this.GetInt("roomstatus")
		ltype, _ := this.GetInt("ltype")
		name := this.GetString("name")
		count, _ := this.GetInt("count")
		ante, _ := this.GetInt("ante")
		cost, _ := this.GetInt("cost")
		vip, _ := this.GetInt("vip")
		chip, _ := this.GetInt("chip")
		deal, _ := this.GetInt("deal")
		carry, _ := this.GetInt("carry")
		down, _ := this.GetInt("down")
		top, _ := this.GetInt("top")
		sit, _ := this.GetInt("sit")
		game := new(entity.Game)
		game.Node = entity.Game2Nodes[node]
		game.Deal = entity.Is2Deal[deal]
		game.Name = name
		game.Gtype = int(gtype)
		game.Rtype = int(rtype)
		game.Ltype = int(ltype)
		game.Status = int(roomstatus)
		game.Count = uint32(count)
		game.Ante = uint32(ante)
		game.Cost = uint32(cost)
		game.Vip = uint32(vip)
		game.Chip = uint32(chip * 100)
		game.Carry = uint32(carry * 100)
		game.Down = uint32(down * 100)
		game.Top = uint32(top * 100)
		game.Sit = uint32(sit)
		fmt.Printf("game %#v\n", game)
		err := this.validGame(game)
		this.checkError(err)
		if err == nil {
			err = service.PlayerService.AddGame(game)
			this.checkError(err)
			service.ActionService.AddGame(this.auth.GetUser().UserName, game.Id)
		}
		this.redirect(beego.URLFor("PlayerController.GameList"))
	}

	this.Data["pageTitle"] = "添加房间"
	this.Data["types1"] = entity.GameNodes
	this.Data["types2"] = entity.GameTypes
	this.Data["types3"] = entity.RoomTypes
	this.Data["types4"] = entity.RoomStatus
	this.Data["types5"] = entity.LotteryTypes
	this.Data["types6"] = entity.IsDeal
	this.display()
}

func (this *PlayerController) validGame(shop *entity.Game) error {
	valid := validation.Validation{}
	valid.Required(shop.Name, "name").Message("名字不能为空")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return errors.New(err.Message)
		}
	}

	return nil
}

// 发布
func (this *PlayerController) Game() {
	id := this.GetString("id")

	shop, err := service.PlayerService.GetGame(id)
	if err != nil {
		this.checkError(err)
	} else {

		b := make(map[string]entity.Game)
		b[shop.Id] = shop
		atype := pb.CONFIG_UPSERT
		if shop.Del == 1 {
			atype = pb.CONFIG_DELETE
		}
		_, err = service.GmRequest(pb.WebGame, atype, b)
		if err != nil {
			flash := beego.NewFlash()
			flash.Error(fmt.Sprintf("%v", err))
			flash.Store(&this.Controller)
		}
	}

	service.ActionService.Game(this.auth.GetUserName(), id)

	this.redirect(beego.URLFor("PlayerController.GameList"))
}

// 移除公告广播
func (this *PlayerController) GameDel() {
	id := this.GetString("id")

	shop, err := service.PlayerService.GetGame(id)
	if err != nil {
		this.checkError(err)
	} else {
		shop.Del = 1
		b := make(map[string]entity.Game)
		b[shop.Id] = shop
		_, err = service.GmRequest(pb.WebGame, pb.CONFIG_DELETE, b)
		if err != nil {
			flash := beego.NewFlash()
			flash.Error(fmt.Sprintf("%v", err))
			flash.Store(&this.Controller)
		} else {
			err3 := service.PlayerService.DelGame(id)
			this.checkError(err3)
		}
	}

	service.ActionService.DelGame(this.auth.GetUserName(), id)

	this.redirect(beego.URLFor("PlayerController.GameList"))
}

// 编辑游戏
func (this *PlayerController) GameEdit() {
	id := this.GetString("id")

	if this.isPost() && id != "" {
		cost, _ := this.GetInt("cost")

		err3 := service.PlayerService.UpdateGame(id, uint32(cost))
		this.checkError(err3)

		service.ActionService.EditGame(this.auth.GetUser().UserName, id)
		this.redirect(beego.URLFor("PlayerController.GameList"))
	} else {
		p, err := service.PlayerService.GetGame(id)
		this.checkError(err)

		this.Data["game"] = p
		this.Data["pageTitle"] = "编辑游戏"
		this.display()
	}
}
