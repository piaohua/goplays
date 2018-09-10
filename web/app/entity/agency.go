package entity

import "time"

/*
//0正常 1等待审核 2未通过审核
const (
	AGENCY_STATUS0 = 0
	AGENCY_STATUS1 = 1
	AGENCY_STATUS2 = 2
)

//代理管理(代理ID为游戏内ID)
type Agency struct {
	UserName string    `bson:"_id"`       //后台账户(代理商手机号码注册)
	Password string    `bson:"password"`  //账号密码
	Salt     string    `bson:"salt"`      //密码盐
	Phone    string    `bson:"phone"`     //绑定的手机号码(备用:非手机号注册时或多个手机时)
	Agent    string    `bson:"agent"`     //代理ID==Userid
	Level    int       `bson:"level"`     //代理等级ID:1级,2级...
	Weixin   string    `bson:"weixin"`    //微信ID
	Alipay   string    `bson:"alipay"`    //支付宝ID
	QQ       string    `bson:"qq"`        //qq号码
	Address  string    `bson:"address"`   //详细地址
	Status   int       `bson:"status"`    //状态,0正常 1等待审核 2未通过审核
	Number   uint32    `bson:"number"`    //当前余额
	Expend   uint32    `bson:"expend"`    //总消耗
	Cash     float32   `bson:"cash"`      //当前可提取额
	Extract  float32   `bson:"extract"`   //已经提取额
	CashTime time.Time `bson:"cash_time"` //提取指定时间前所有
	Created  time.Time `bson:"created"`   //加入时间
	Updated  time.Time `bson:"updated"`   //更新时间
}
*/

//1: "微信", 2: "支付宝", 3: "银行账号",
const (
	PAYWAY1 = 1
	PAYWAY2 = 2
	PAYWAY3 = 3
)

//提现/申请记录
type ApplyCash struct {
	Id       string    `bson:"_id"`       //单号
	Agent    string    `bson:"agent"`     //申请人
	Cash     float32   `bson:"cash"`      //申请提现金额
	Fee      int64     `bson:"fee"`       //申请提现金额
	Rest     int64     `bson:"rest"`      //剩余可提现金额
	Status   int       `bson:"status"`    //0表示已经处理,1表示等待处理
	UserName string    `bson:"user_name"` //处理人,后台账号
	RealName string    `bson:"real_name"` //姓名
	BankCard int       `bson:"bank_card"` //收款方式,1: "微信", 2: "支付宝", 3: "银行账号",
	BankAddr string    `bson:"bank_addr"` //收款账号,银行卡号 开户行
	Message  string    `bson:"message"`   //处理描述
	Ctime    time.Time `bson:"ctime"`     //申请时间
	Utime    time.Time `bson:"utime"`     //处理时间
	Feef     float64   `bson:"-"`         //申请提现金额
	Restf    float64   `bson:"-"`         //剩余可提现金额
}
