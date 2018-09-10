package service

var (
	UserService   *userService   // 用户服务
	RoleService   *roleService   // 角色服务
	MailService   *mailService   // 邮件服务
	ActionService *actionService // 系统动态
	SystemService *systemService
	//---
	PlayerService *playerService // 玩家管理
	LoggerService *loggerService // 日志管理
	AgencyService *agencyService // 代理管理
	ChartsService *chartsService // 统计图表
)

func initService() {
	UserService = &userService{}
	RoleService = &roleService{}
	MailService = &mailService{}
	ActionService = &actionService{}
	SystemService = &systemService{}
	PlayerService = &playerService{}
	LoggerService = &loggerService{}
	AgencyService = &agencyService{}
	ChartsService = &chartsService{}
}
