package entity

import "time"

const (
	USER_STATUS0 = 0  //正常
	USER_STATUS1 = -1 //禁用
)

// 账号自增id
type UserIDGen struct {
	Id         string `bson:"_id"`
	LastUserId string `bson:"last_user_id"`
}

// 账号
type User struct {
	Id         string    `bson:"_id"`         // AUTO_INCREMENT, PRIMARY KEY (`id`),
	UserName   string    `bson:"user_name"`   // 用户名, UNIQUE KEY `user_name` (`user_name`)
	Password   string    `bson:"password"`    // 密码
	Salt       string    `bson:"salt"`        // 密码盐
	Sex        int       `bson:"sex"`         // 性别
	Email      string    `bson:"email"`       // 邮箱
	LastLogin  time.Time `bson:"last_login"`  // 最后登录时间
	LastIp     string    `bson:"last_ip"`     // 最后登录IP
	Status     int       `bson:"status"`      // 状态，0正常 -1禁用
	CreateTime time.Time `bson:"create_time"` // 创建时间
	UpdateTime time.Time `bson:"update_time"` // 更新时间
	RoleList   []Role    `bson:"role_list"`   // 角色列表
	//代理
	Phone      string    `bson:"phone"`       //绑定的手机号码(备用:非手机号注册时或多个手机时)
	Agent      string    `bson:"agent"`       //代理ID==Userid
	Level      int       `bson:"level"`       //代理等级ID:1级,2级...
	Weixin     string    `bson:"weixin"`      //微信ID
	Alipay     string    `bson:"alipay"`      //支付宝ID
	QQ         string    `bson:"qq"`          //qq号码
	Address    string    `bson:"address"`     //详细地址
	Number     uint32    `bson:"number"`      //当前余额
	Expend     uint32    `bson:"expend"`      //总消耗
	Cash       float32   `bson:"cash"`        //当前可提取额(分)
	Extract    float32   `bson:"extract"`     //已经提取额(分)
	Rate       uint32    `bson:"rate"`        //提现率,可配置，百分值(比如:80表示80%_),给上级的提成比例
	Atype      uint32    `bson:"atype"`       //代理登录包类型,分包的人才有
	Belong     uint32    `bson:"belong"`      //属于某个代理登录包类型
	Builds     uint32    `bson:"builds"`      //绑定我的人数
	BuildsTime time.Time `bson:"builds_time"` //统计指定时间前所有
	CashTime   time.Time `bson:"cash_time"`   //提取指定时间前所有
	//反佣统计
	FeeAll      int64     `bson:"fee_all"`      //所有反佣(全部收益)
	FeeTime     time.Time `bson:"fee_time"`     //最后统计时间
	FeeTop      int64     `bson:"fee_top"`      //历史反佣
	FeeExtract  int64     `bson:"fee_extract"`  //已经提取反佣
	FeeRate     int64     `bson:"fee_rate"`     //历史抽成给上级数量
	SysFeeRate  int64     `bson:"sys_fee_rate"` //历史抽成给系统数量
	Chip        int64     `bson:"chip"`         //当前筹码余额
	ParentAgent string    `bson:"parent_agent"` //上级代理ID==Userid
	SysRate     uint32    `bson:"sys_rate"`     //系统抽成比例,百分比
	//
	FeeAllf     float64 `bson:"-"` //所有反佣(全部收益)
	FeeTopf     float64 `bson:"-"` //历史反佣
	FeeExtractf float64 `bson:"-"` //已经提取反佣
	FeeRatef    float64 `bson:"-"` //历史抽成给上级数量
	SysFeeRatef float64 `bson:"-"` //历史抽成给系统数量
	//TODO 优化下面字段操作, child, players 数据会过大
	Parent      string    `bson:"parent"`       //创建者id,用总代列表只展示自己添加的账号关系
	Child       []string  `bson:"child"`        //已创建的代理列表(邀请码id,包换所有下级代理的邀请码,主要用来区分玩家列表)
	Players     []string  `bson:"players"`      //已发展出来的玩家(玩家id,包换所有下级代理的玩家,主要用来区分数据查询)
	PlayersTime time.Time `bson:"players_time"` //统计指定时间前所有
}

// 账号属于分组(可属于多个分组)
type UserRole struct {
	Id     string `bson:"_id"`     // UNIQUE KEY `user_id` (`user_id`,`role_id`)
	UserId string `bson:"user_id"` // 用户id
	RoleId string `bson:"role_id"` // 角色id
}
