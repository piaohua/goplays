package service

import (
	"fmt"

	"goplays/web/app/entity"

	"gopkg.in/mgo.v2/bson"
)

// 系统动态
type actionService struct{}

// 添加记录
func (this *actionService) Add(action, actor, objectType string, objectId string, extra string) bool {
	act := new(entity.Action)
	act.Action = action
	act.Actor = actor
	act.ObjectType = objectType
	act.ObjectId = objectId
	act.Extra = extra
	act.CreateTime = bson.Now()
	act.Id = bson.NewObjectId().Hex()
	Insert(Actions, act)
	return true
}

// 登录动态
func (this *actionService) Login(userName string, userId string, ip string) {
	this.Add("login", userName, "user", userId, ip)
}

// 退出登录
func (this *actionService) Logout(userName string, userId string, ip string) {
	this.Add("logout", userName, "user", userId, ip)
}

// 更新个人信息
func (this *actionService) UpdateProfile(userName string, userId string) {
	this.Add("update_profile", userName, "user", userId, "")
}

// 更新个人钻石,otype操作类型,oid操作目标id,extra操作的数量
func (this *actionService) UpdateDiamond(userName, otype, oid, extra string) {
	this.Add("update_diamond", userName, otype, oid, extra)
}
func (this *actionService) UpdateNumber(userName, otype, oid, extra string) {
	this.Add("update_number", userName, otype, oid, extra)
}
func (this *actionService) UpdateExpend(userName, otype, oid, extra string) {
	this.Add("update_expend", userName, otype, oid, extra)
}
func (this *actionService) UpdateChip(userName, otype, oid, extra string) {
	this.Add("update_chip", userName, otype, oid, extra)
}

// 注册动态
func (this *actionService) Regist(userName string, agent string, ip string) {
	this.Add("regist", userName, "agent", agent, ip)
}

// 角色动态
func (this *actionService) AddRole(userName string, roleName string) {
	this.Add("add_role", userName, "role_name", roleName, "")
}
func (this *actionService) DelRole(userName string, roleId string) {
	this.Add("del_role", userName, "role_id", roleId, "")
}
func (this *actionService) UpdateRole(userName string, roleName string) {
	this.Add("update_role", userName, "role_name", roleName, "")
}
func (this *actionService) PermRole(userName string, roleName string) {
	this.Add("perm_role", userName, "role_name", roleName, "")
}

// 用户动态
func (this *actionService) AddUser(userName string, user_name string) {
	this.Add("add_user", userName, "user_name", user_name, "")
}
func (this *actionService) UpdateUser(userName string, id string) {
	this.Add("update_user", userName, "user_id", id, "")
}
func (this *actionService) DelUser(userName string, id string) {
	this.Add("del_user", userName, "user_id", id, "")
}

// 代理动态
func (this *actionService) AddAgency(userName string, phone, agency string) {
	this.Add("add_agency", userName, "phone", phone, agency)
}
func (this *actionService) UpdateBuild(userName string, userid, agent string) {
	this.Add("build_agency", userName, "userid", userid, agent)
}
func (this *actionService) AddApplyCash(userName string, money string) {
	this.Add("apply_cash", userName, "money", money, "")
}
func (this *actionService) ExtractApplyCash(userName string, orderid string) {
	this.Add("extract_cash", userName, "orderid", orderid, "")
}
func (this *actionService) UpdateAgency(userName string, agent, rate string) {
	this.Add("update_agency", userName, "rate", rate, agent)
}

// 公告动态
func (this *actionService) AddNotice(userName string, notice_id string) {
	this.Add("add_notice", userName, "notice_id", notice_id, "")
}
func (this *actionService) Notice(userName string, notice_id string) {
	this.Add("notice", userName, "notice_id", notice_id, "")
}
func (this *actionService) DelNotice(userName string, notice_id string) {
	this.Add("del_notice", userName, "notice_id", notice_id, "")
}

// 商城动态
func (this *actionService) AddShop(userName string, shop_id string) {
	this.Add("add_shop", userName, "shop_id", shop_id, "")
}
func (this *actionService) Shop(userName string, shop_id string) {
	this.Add("shop", userName, "shop_id", shop_id, "")
}
func (this *actionService) DelShop(userName string, shop_id string) {
	this.Add("del_shop", userName, "shop_id", shop_id, "")
}

// 变量动态
func (this *actionService) EnvAdd(userName string, key string) {
	this.Add("add_env", userName, "key", key, "")
}
func (this *actionService) EnvDel(userName string, key string) {
	this.Add("del_env", userName, "key", key, "")
}

// VIP动态
func (this *actionService) AddVip(userName string, classic_id string) {
	this.Add("add_vip", userName, "vip_id", classic_id, "")
}
func (this *actionService) Vip(userName string, classic_id string) {
	this.Add("vip", userName, "vip_id", classic_id, "")
}
func (this *actionService) DelVip(userName string, classic_id string) {
	this.Add("del_vip", userName, "vip_id", classic_id, "")
}

// 商城动态
func (this *actionService) AddGame(userName string, game_id string) {
	this.Add("add_game", userName, "game_id", game_id, "")
}
func (this *actionService) Game(userName string, game_id string) {
	this.Add("game", userName, "game_id", game_id, "")
}
func (this *actionService) DelGame(userName string, game_id string) {
	this.Add("del_game", userName, "game_id", game_id, "")
}
func (this *actionService) EditGame(userName string, game_id string) {
	this.Add("edit_game", userName, "game_id", game_id, "")
}

// 获取动态列表
func (this *actionService) GetList(userName string, page, pageSize int) ([]entity.Action, error) {
	var list []entity.Action
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "create_time", false)
	//m := bson.M{"create_time": bson.M{"$gte": bson.Now()}}
	err := Actions.
		Find(bson.M{"actor": userName}).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	if err == nil {
		num := len(list)
		for i := 0; i < num; i++ {
			this.format(&list[i])
		}
	}
	return list, err
}

// 格式化
func (this *actionService) format(action *entity.Action) {
	switch action.Action {
	case "login":
		action.Message = fmt.Sprintf("<b>%s</b> 登录系统，IP为 <b>%s</b>。", action.Actor, action.Extra)
	case "logout":
		action.Message = fmt.Sprintf("<b>%s</b> 退出系统。", action.Actor)
	case "update_profile":
		action.Message = fmt.Sprintf("<b>%s</b> 更新了个人资料。", action.Actor)
	case "create_task":
		action.Message = fmt.Sprintf("<b>%s</b> 创建了编号为 <b class='blue'>%d</b> 的发布单。", action.Actor, action.ObjectId)
	case "regist":
		action.Message = fmt.Sprintf("<b>%s</b> 注册成功，IP为 <b>%s</b>。", action.Actor, action.Extra)
	case "add_role":
		action.Message = fmt.Sprintf("<b>%s</b> 添加角色，角色名称为 <b>%s</b>。", action.Actor, action.ObjectId)
	case "del_role":
		action.Message = fmt.Sprintf("<b>%s</b> 删除角色，角色ID为 <b>%s</b>。", action.Actor, action.ObjectId)
	case "update_role":
		action.Message = fmt.Sprintf("<b>%s</b> 更新角色，角色名称为 <b>%s</b>。", action.Actor, action.ObjectId)
	case "perm_role":
		action.Message = fmt.Sprintf("<b>%s</b> 更新角色权限，角色名称为 <b>%s</b>。", action.Actor, action.ObjectId)
	case "add_user":
		action.Message = fmt.Sprintf("<b>%s</b> 添加账号，账号名称为 <b>%s</b>。", action.Actor, action.ObjectId)
	case "update_user":
		action.Message = fmt.Sprintf("<b>%s</b> 更新账号，账号ID为 <b>%s</b>。", action.Actor, action.ObjectId)
	case "del_user":
		action.Message = fmt.Sprintf("<b>%s</b> 删除账号，账号ID为 <b>%s</b>。", action.Actor, action.ObjectId)
	}
}
