package service

import (
	"errors"

	"goplays/web/app/entity"
	"goplays/web/app/libs"
	"utils"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

type userService struct{}

// 获取代理商总数
func (this *userService) GetBuilds(agent string) int {
	m := bson.M{"parent_agent": agent}
	return Count(Users, m)
}

// 生成ID
func (this *userService) GetUserIDGen() (string, error) {
	gen := new(entity.UserIDGen)
	gen.Id = "last_user_id"
	Get(GenIDs, gen.Id, gen)
	if gen.LastUserId == "" {
		gen.LastUserId = "2"
	}
	id := gen.LastUserId
	gen.LastUserId = utils.StringAdd(id)
	if Upsert(GenIDs, bson.M{"_id": gen.Id}, gen) {
		return id, nil
	}
	return id, errors.New("生成错误")
}

// 根据用户ids获取代理商列表
func (this *userService) GetAgencyList(userIds []string) ([]entity.User, error) {
	var userList []entity.User
	if len(userIds) == 0 {
		return userList, nil
	}
	ListByQ(Users, bson.M{"agent": bson.M{"$in": userIds}}, &userList)
	return userList, nil
}

// 分页获取用户列表
func (this *userService) GetUserList2(page, pageSize int, getRoleInfo bool) ([]entity.User, error) {
	var users []entity.User
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "_id", false)
	err := Users.
		Find(nil).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&users)
	if getRoleInfo {
		for k, user := range users {
			users[k].RoleList, _ = this.GetUserRoleList(user.Id)
		}
	}
	return users, err
}

// 根据用户ids获取代理商列表
func (this *userService) GetAgencyLists(userIds []string) ([]bson.M, error) {
	var userList []bson.M
	if len(userIds) == 0 {
		return userList, nil
	}
	//userIds = []string{"1"}
	q := bson.M{"agent": bson.M{"$in": userIds}}
	f := []string{"agent", "rate", "fee_rate"}
	//var result2 []bson.M
	//ListByQWithFields(Users, q, f, &result2)
	//fmt.Println("result2 ", result2)
	//result2  [map[_id:1 user_name:admin]]
	ListByQWithFields(Users, q, f, &userList)
	//fmt.Println("userList ", userList)
	//
	//q = bson.M{"_id": "1"}
	//f = []string{"user_name"}
	//var userid bson.M
	//GetByQWithFields(Users, q, f, &userid)
	//fmt.Println("userid ", userid)
	return userList, nil
}

// 根据用户分包类型获取代理商列表,属于分包的代理商
func (this *userService) GetAgencyListByAtype(atype uint32) ([]entity.User, error) {
	var userList []entity.User
	if atype == 0 {
		return userList, errors.New("分包不存在")
	}
	m := bson.M{"belong": atype}
	m["agent"] = bson.M{"$ne": ""}
	ListByQ(Users, m, &userList)
	return userList, nil
}

// 根据用户id获取一个用户信息
func (this *userService) GetUser(userId string, getRoleInfo bool) (*entity.User, error) {
	user := new(entity.User)
	Get(Users, userId, user)
	if user.Id != "" && getRoleInfo {
		user.RoleList, _ = this.GetUserRoleList(user.Id)
		return user, nil
	}
	if user.Id != "" && !getRoleInfo {
		return user, nil
	}
	return user, errors.New("获取失败")
}

// 根据用户名获取用户信息
func (this *userService) GetUserByName(userName string) (*entity.User, error) {
	user := new(entity.User)
	GetByQ(Users, bson.M{"user_name": userName}, user)
	if user.Id != "" {
		return user, nil
	}
	return user, errors.New("获取失败")
}

// 获取用户总数
func (this *userService) GetTotal() (int64, error) {
	return int64(Count(Users, nil)), nil
}

// 分页获取用户列表
func (this *userService) GetUserList(page, pageSize int, getRoleInfo bool) ([]entity.User, error) {
	var users []entity.User
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "_id", false)
	err := Users.
		Find(nil).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&users)
	for k, user := range users {
		users[k].RoleList, _ = this.GetUserRoleList(user.Id)
	}
	return users, err
}

// 根据角色id获取用户列表
func (this *userService) GetUserListByRoleId(roleId string) ([]entity.User, error) {
	var users []entity.User
	var userRole []entity.UserRole
	ListByQ(UserRoles, bson.M{"role_id": roleId}, &userRole)
	if len(userRole) == 0 {
		return users, errors.New("角色不存在")
	}
	for _, v := range userRole {
		var user entity.User
		GetByQ(Users, bson.M{"_id": v.UserId}, &user)
		if user.Id != "" {
			users = append(users, user)
		}
	}
	return users, nil
}

// 获取某个用户的角色列表
// 为什么不直接连表查询role表？因为不想“越权”查询
func (this *userService) GetUserRoleList(userId string) ([]entity.Role, error) {
	var (
		roleRef  []entity.UserRole
		roleList []entity.Role
	)
	ListByQ(UserRoles, bson.M{"user_id": userId}, &roleRef)
	roleList = make([]entity.Role, 0, len(roleRef))
	for _, v := range roleRef {
		if role, err := RoleService.GetRole(v.RoleId); err == nil {
			roleList = append(roleList, *role)
		}
	}
	return roleList, nil
}

// 添加用户
func (this *userService) AddUser(userName, agent, password string, sex int) (*entity.User, error) {
	if exists, _ := this.GetUserByName(userName); exists.Id != "" {
		return nil, errors.New("用户名已存在")
	}
	if agent != "" {
		if exists, _ := this.GetUserByAgent(agent); exists.Id != "" {
			return nil, errors.New("邀请码已存在")
		}
	}

	user := new(entity.User)
	user.UserName = userName
	user.Sex = sex
	//user.Email = email
	user.Agent = agent
	user.Salt = string(utils.RandomCreateBytes(10))
	user.Password = libs.Md5([]byte(password + user.Salt))
	user.CreateTime = bson.Now()
	user.UpdateTime = bson.Now()
	user.LastLogin = bson.Now()
	// user.LastLogin = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	var err error
	user.Id, err = this.GetUserIDGen()
	if err != nil {
		return user, errors.New("添加失败")
	}
	if Insert(Users, user) {
		return user, nil
	}
	return user, errors.New("添加失败")
}

// 更新用户信息
func (this *userService) UpdateUser(user *entity.User, fileds bson.M) error {
	if len(fileds) < 1 {
		return errors.New("更新字段不能为空")
	}
	if Update(Users, bson.M{"_id": user.Id}, bson.M{"$set": fileds}) {
		return nil
	}
	return errors.New("更新失败")
}

// 修改密码
func (this *userService) ModifyPassword(userId string, password string) error {
	user, err := this.GetUser(userId, false)
	if err != nil {
		return err
	}
	user.Salt = string(utils.RandomCreateBytes(10))
	user.Password = libs.Md5([]byte(password + user.Salt))
	if Update(Users, bson.M{"_id": user.Id}, bson.M{"$set": bson.M{"salt": user.Salt, "password": user.Password, "update_time": bson.Now()}}) {
		return nil
	}
	return errors.New("更新失败")
}

// 删除用户
func (this *userService) DeleteUser(userId string) error {
	if userId == "1" {
		return errors.New("不允许删除用户ID为1的用户")
	}
	if Delete(Users, bson.M{"_id": userId}) {
		return nil
	}
	return errors.New("删除用户失败")
}

// 设置用户角色
func (this *userService) UpdateUserRoles(userId string, roleIds []string) error {
	if _, err := this.GetUser(userId, false); err != nil {
		return err
	}
	DeleteAll(UserRoles, bson.M{"user_id": userId})
	for _, v := range roleIds {
		Insert(UserRoles, &entity.UserRole{Id: userId + "." + v, UserId: userId, RoleId: v})
	}
	return nil
}

//代理

func (this *userService) GetUserByPhone(phone string) (*entity.User, error) {
	user := new(entity.User)
	GetByQ(Users, bson.M{"phone": phone}, user)
	if user.Phone == "" {
		return user, errors.New("代理商不存在")
	}
	return user, nil
}

func (this *userService) GetUserByAgent(agent string) (*entity.User, error) {
	user := new(entity.User)
	GetByQ(Users, bson.M{"agent": agent}, user)
	if user.Agent == "" {
		return user, errors.New("代理商不存在")
	}
	return user, nil
}

////获取上级代理 agent == userid
//func (this *userService) GetUseridByAgent(agent string) (string, error) {
//	q := bson.M{"_id": agent}
//	f := []string{"agent"}
//	var userid bson.M
//	GetByQWithFields(Users, q, f, &userid)
//	if v, ok := userid["agent"]; ok {
//		return v.(string), nil
//	}
//	return "", errors.New("不存在")
//}

////获取代理id
func (this *userService) GetUseridByAgent(agent string) (string, error) {
	q := bson.M{"agent": agent}
	f := []string{"_id"}
	var userid bson.M
	GetByQWithFields(Users, q, f, &userid)
	if v, ok := userid["_id"]; ok {
		return v.(string), nil
	}
	return "", errors.New("不存在")
}

func (this *userService) GetUserByAtype(atype uint32) (*entity.User, error) {
	user := new(entity.User)
	GetByQ(Users, bson.M{"atype": atype}, user)
	if user.Atype == 0 {
		return user, errors.New("代理商不存在")
	}
	return user, nil
}

// 添加代理商用户
func (this *userService) AddAgencyUser(userName, password, agent, phone, weixin, qq, address string, atype uint32) (*entity.User, error) {
	//if exists, _ := this.GetUserByName(userName); exists.Id != "" {
	//	return nil, errors.New("用户名已存在")
	//}

	user := new(entity.User)
	user.UserName = userName
	// user.Salt = salt
	// user.Password = password
	user.CreateTime = bson.Now()
	user.UpdateTime = bson.Now()
	user.LastLogin = bson.Now()
	// user.LastLogin = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	user.Salt = string(utils.RandomCreateBytes(10))
	user.Password = libs.Md5([]byte(password + user.Salt))
	user.Phone = phone
	user.Agent = agent
	user.Weixin = weixin
	user.QQ = qq
	user.Address = address
	user.Belong = atype
	user.Status = -1 //需要审核
	var err error
	user.Id, err = this.GetUserIDGen()
	if err != nil {
		return user, errors.New("添加失败")
	}
	if Insert(Users, user) {
		return user, nil
	}
	return user, errors.New("添加失败")
}

// 获取分包类型
func (this *userService) GetPlayerAtypes() (map[int]string, error) {
	var userList []entity.User
	ListByQ(Users, bson.M{"atype": bson.M{"$ne": 0}}, &userList)
	types := make(map[int]string)
	types[0] = "0"
	for _, v := range userList {
		types[int(v.Atype)] = utils.String(v.Atype)
	}
	return types, nil
}

// 添加用户
func (this *userService) AddUser2(userName, agent, parent_agent, parent_id, password string, sex, rate int) (*entity.User, error) {
	if exists, _ := this.GetUserByName(userName); exists.Id != "" {
		return nil, errors.New("用户名已存在")
	}
	if exists, _ := this.GetUserByAgent(agent); exists.Id != "" {
		return nil, errors.New("邀请码已存在")
	}

	user := new(entity.User)
	user.UserName = userName
	user.Sex = sex
	//user.Email = email
	user.Agent = agent
	user.Rate = uint32(rate)
	user.ParentAgent = parent_agent
	user.Salt = string(utils.RandomCreateBytes(10))
	user.Password = libs.Md5([]byte(password + user.Salt))
	user.CreateTime = bson.Now()
	user.UpdateTime = bson.Now()
	user.Parent = parent_id
	//user.LastLogin = bson.Now()
	// user.LastLogin = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	var err error
	user.Id, err = this.GetUserIDGen()
	if err != nil {
		return user, errors.New("添加失败")
	}
	if Insert(Users, user) {
		return user, nil
	}
	return user, errors.New("添加失败")
}

// 更新用户信息
func (this *userService) UpdateChild(user *entity.User, child_agent string) error {
	err := this.update2child(user, child_agent)
	if err != nil {
		return err
	}
	//更新所有上级
	return this.update2parent(user.Parent, child_agent)
}

func (this *userService) update2child(user *entity.User, child_agent string) error {
	//更新添加者
	for _, v := range user.Child {
		if child_agent == v {
			return nil
		}
	}
	user.Child = append(user.Child, child_agent)
	if !Update(Users, bson.M{"_id": user.Id}, bson.M{"$set": bson.M{"child": user.Child}}) {
		return errors.New("更新失败")
	}
	return nil
}

func (this *userService) updateParentId(user *entity.User) error {
	if !Update(Users, bson.M{"_id": user.Id}, bson.M{"$set": bson.M{"parent": user.Parent}}) {
		return errors.New("更新失败")
	}
	return nil
}

//更新所有上级
func (this *userService) update2parent(parent_id, child_agent string) error {
	if parent_id == "" {
		return nil
	}
	user, err := this.GetUser(parent_id, false)
	if err != nil {
		return err
	}
	err = this.update2child(user, child_agent)
	if err != nil {
		return err
	}
	return this.update2parent(user.Parent, child_agent)
}

// 定时统计
func initParentAgent() {
	list, err := UserService.GetUserList2(1, -1, false)
	if err != nil {
		beego.Error("stat err: ", err)
	}
	beego.Trace("initParentAgent list len : ", len(list))
	for _, v := range list {
		initParent(&v) //
	}
}

//启动时初始化之前账号关系
func initParent(agency *entity.User) {
	if agency.Agent == "" {
		return
	}
	if agency.Id == "1" { //超级账号
		return
	}
	if agency.Parent == "" {
		if agency.ParentAgent == "" {
			agency.Parent = "1" //默认超级账号添加
		} else {
			id, err := UserService.GetUseridByAgent(agency.ParentAgent)
			if err == nil {
				agency.Parent = id
			}
		}
		//更新
		UserService.updateParentId(agency)
	}
	if agency.Parent == "" {
		return
	}
	initChild(agency)
}

//启动时初始化之前账号关系
func initChild(agency *entity.User) {
	if agency.Agent == "" {
		return
	}
	//从最低级代理开始
	n, _ := AgencyService.GetMyAgencyTotal3(agency.Agent, bson.M{})
	if n != 0 {
		return
	}
	//更新所有上级
	UserService.update2parent(agency.Parent, agency.Agent)
}

// 定时统计
func initParentAgent2() {
	//TODO 优化
	list, err := UserService.GetUserList2(1, -1, false)
	if err != nil {
		beego.Error("stat err: ", err)
	}
	beego.Trace("initParentAgent2 list len : ", len(list))
	for _, v := range list {
		playersStat(&v) //统计玩家
	}
}

//统计玩家,TODO 优化
func playersStat(agency *entity.User) {
	if len(agency.Child) == 0 {
		return
	}
	endTime := utils.LocalTime()
	startTime := agency.PlayersTime
	q := bson.M{
		"robot": false,
		"atime": bson.M{"$gte": startTime, "$lt": endTime},
		"agent": bson.M{"$in": agency.Child},
	}
	f := []string{"_id"}
	var result []bson.M
	ListByQWithFields(PlayerUsers, q, f, &result)
	beego.Trace("PlayersStat result : ", result)
	if len(result) == 0 {
		return
	}
	list := make([]string, 0)
	for _, v := range result {
		if val, ok := v["_id"]; ok {
			if id, ok := val.(string); ok && id != "" {
				list = append(list, id)
			}
		}
	}
	if len(list) == 0 {
		return
	}
	beego.Trace("PlayersStat : ", agency.Id)
	beego.Trace("PlayersStat : ", agency.Agent, list)
	agency.Players = append(agency.Players, list...)
	agency.PlayersTime = endTime
	err := update2players(agency)
	if err != nil {
		beego.Error("PlayersStat err : ", agency.Agent, list, err)
	}
}

// 更新用户信息
func update2players(user *entity.User) error {
	if Update(Users, bson.M{"_id": user.Id}, bson.M{"$set": bson.M{"players_time": user.PlayersTime, "players": user.Players}}) {
		return nil
	}
	return errors.New("更新失败")
}
