package service

import (
	"errors"
	"fmt"
	"strings"

	"goplays/web/app/entity"
	"goplays/web/app/libs"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

// 登录验证服务
type AuthService struct {
	loginUser *entity.User    // 当前登录用户
	permMap   map[string]bool // 当前用户权限表
	openPerm  map[string]bool // 公开的权限
}

func NewAuth() *AuthService {
	return new(AuthService)
}

// 初始化开放权限
func (this *AuthService) initOpenPerm() {
	this.openPerm = map[string]bool{
		"main.index":      true,
		"main.profile":    true,
		"main.login":      true,
		"main.logout":     true,
		"main.getpubstat": true,
		"main.regist":     true,
		"main.servers":    true,
		"main.files":      true,
	}
}

// 获取当前登录的用户对象
func (this *AuthService) GetUser() *entity.User {
	return this.loginUser
}

// 获取当前登录的用户id
func (this *AuthService) GetUserId() string {
	if this.IsLogined() {
		return this.loginUser.Id
	}
	return ""
}

// 获取当前登录的用户名
func (this *AuthService) GetUserName() string {
	if this.IsLogined() {
		return this.loginUser.UserName
	}
	return ""
}

// 初始化
func (this *AuthService) Init(token string) {
	this.initOpenPerm()
	arr := strings.Split(token, "|")
	beego.Trace("登录验证, token: ", token)
	if len(arr) == 2 {
		idstr, password := arr[0], arr[1]
		//userId, _ := strconv.Atoi(idstr)
		//if userId > 0 {
		if idstr != "" {
			user, err := UserService.GetUser(idstr, true)
			//beego.Trace("验证，用户信息: ", user, err)
			if err == nil && password == libs.Md5([]byte(gmKey+user.Password+user.Salt)) {
				this.loginUser = user
				this.initPermMap()
				//beego.Trace("验证成功，用户信息: ", user)
			}
		}
	}
}

// 初始化权限表
func (this *AuthService) initPermMap() {
	this.permMap = make(map[string]bool)
	for _, role := range this.loginUser.RoleList {
		for _, perm := range role.PermList {
			this.permMap[perm.Key] = true
		}
	}
}

// 检查是否有某个权限
func (this *AuthService) HasAccessPerm(module, action string) bool {
	key := module + "." + action
	if !this.IsLogined() {
		return false
	}
	if this.loginUser.Id == "1" || this.isOpenPerm(key) {
		return true
	}
	if _, ok := this.permMap[key]; ok {
		return true
	}
	return false
}

// 检查是否登录
func (this *AuthService) IsLogined() bool {
	return this.loginUser != nil && this.loginUser.Id != ""
}

// 是否公开访问的操作
func (this *AuthService) isOpenPerm(key string) bool {
	if _, ok := this.openPerm[key]; ok {
		return true
	}
	return false
}

// 用户登录
func (this *AuthService) Login(userName, password string) (string, error) {
	user, err := UserService.GetUserByName(userName)
	if err != nil {
		//if err == orm.ErrNoRows {
		//	return "", errors.New("帐号或密码错误")
		//} else {
		//	return "", errors.New("系统错误")
		//}
		return "", errors.New("帐号或密码错误")
	}

	if user.Password != libs.Md5([]byte(password+user.Salt)) {
		return "", errors.New("帐号或密码错误")
	}
	if user.Status == -1 {
		//return "", errors.New("该帐号已禁用")
		return "", errors.New("帐号等待审核中")
	}

	user.LastLogin = bson.Now()
	UserService.UpdateUser(user, bson.M{"last_login": user.LastLogin})
	this.loginUser = user

	token := fmt.Sprintf("%s|%s", user.Id, libs.Md5([]byte(gmKey+user.Password+user.Salt)))
	return token, nil
}

// 退出登录
func (this *AuthService) Logout() error {
	return nil
}

// 用户注册成为代理
func (this *AuthService) Regist(username, password, agent, phone, weixin, qq, address string, atype uint32) error {
	user, err := UserService.GetUserByName(username)
	if err == nil || user.Id != "" {
		return errors.New("账号名字已经被注册")
	}
	user, err = UserService.GetUserByPhone(phone)
	if err == nil || user.Phone != "" {
		return errors.New("电话号码已经被注册")
	}
	user, err = UserService.GetUserByAgent(agent)
	if err == nil || user.Agent != "" {
		return errors.New("代理ID已经被注册")
	}
	user, err = UserService.AddAgencyUser(username, password, agent, phone, weixin, qq, address, atype)
	if err != nil {
		return err
	}
	return nil
}
