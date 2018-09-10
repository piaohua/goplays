package entity

import "time"

//玩家数据
type PlayerUser struct {
	Userid   string `bson:"_id" json:"userid"`          // 用户id
	Nickname string `bson:"nickname" json:"nickname"`   // 用户昵称
	Photo    string `bson:"photo" json:"photo"`         // 头像
	Wxuid    string `bson:"wxuid" json:"wxuid"`         // 微信uid
	Sex      uint32 `bson:"sex" json:"sex"`             // 用户性别,男1 女2 非男非女3
	Phone    string `bson:"phone" json:"phone"`         // 绑定的手机号码
	Auth     string `bson:"auth" json:"auth"`           // 密码验证码
	Password string `bson:"password" json:"password"`   // MD5密码
	RegistIp string `bson:"regist_ip" json:"regist_ip"` // 注册账户时的IP地址
	LoginIp  string `bson:"login_ip" json:"login_ip"`   // 登录账户时的IP地址
	Diamond  int64  `bson:"diamond" json:"diamond"`     // 钻石
	Coin     int64  `bson:"coin" json:"coin"`           // 金币
	Chip     int64  `bson:"chip" json:"chip"`           // 筹码
	Card     int64  `bson:"card" json:"card"`           // 房卡
	Vip      uint32 `bson:"vip" json:"vip"`             // vip
	Status   uint32 `bson:"status" json:"status"`       // 正常1  锁定2  黑名单3
	Robot    bool   `bson:"robot" json:"robot"`         // 是否是机器人
	//战绩
	Win  uint32 `bson:"win" json:"win"`   // 赢
	Lost uint32 `bson:"lost" json:"lost"` // 输
	Ping uint32 `bson:"ping" json:"ping"` // 平
	//最高充值
	Money uint32 `bson:"money" json:"money"` // 充值总金额(分)
	//代理
	Agent string    `bson:"agent" json:"agent"` // 代理ID
	Atime time.Time `bson:"atime" json:"atime"` // 绑定代理时间
	//时间
	Ctime     time.Time `bson:"ctime" json:"ctime"`           // 注册时间
	LoginTime time.Time `bson:"login_time" json:"login_time"` // 最后登录时间
	//临时状态
	State   int     // 在线状态
	Agency  int     // 是否代理
	Rate    uint32  //提现率,可配置，百分值(比如:80表示80%_)
	FeeRate int64   //历史抽成给上级数量
	Chipf   float64 `bson:"-"` //数量
}

const (
	TradeSuccess = 0 //交易成功
	TradeFail    = 1 //交易失败
	Tradeing     = 2 //交易中(下单状态)
	TradeGoods   = 3 //发货失败
)

var TradeResult = map[int]string{
	TradeSuccess: "成功",
	TradeFail:    "交易失败",
	//Tradeing:     "交易中",
	TradeGoods: "发货失败",
}

// 交易记录
type TradeRecord struct {
	Id        string    `bson:"_id"`       //商户订单号(游戏内自定义订单号)
	Transid   string    `bson:"transid"`   //交易流水号(计费支付平台的交易流水号,微信订单号)
	Userid    string    `bson:"userid"`    //用户在商户应用的唯一标识(userid)
	Itemid    string    `bson:"itemid"`    //购买商品ID
	Amount    string    `bson:"amount"`    //购买商品数量
	Diamond   uint32    `bson:"diamond"`   //购买钻石数量
	Money     uint32    `bson:"money"`     //交易总金额(单位为分)
	Transtime string    `bson:"transtime"` //交易完成时间 yyyy-mm-dd hh24:mi:ss
	Result    int       `bson:"result"`    //交易结果(0–交易成功,1–交易失败,2-交易中,3-发货中)
	Waresid   uint32    `bson:"waresid"`   //商品编码(平台为应用内需计费商品分配的编码)
	Currency  string    `bson:"currency"`  //货币类型(RMB,CNY)
	Transtype int       `bson:"transtype"` //交易类型(0–支付交易)
	Feetype   int       `bson:"feetype"`   //计费方式(表示商品采用的计费方式)
	Paytype   uint32    `bson:"paytype"`   //支付方式(表示用户采用的支付方式,403-微信支付)
	Clientip  string    `bson:"clientip"`  //客户端ip
	Agent     string    `bson:"agent"`     //绑定的父级代理商游戏ID
	Atype     uint32    `bson:"atype"`     //代理包类型
	First     int       `bson:"first"`     //首次充值
	Utime     time.Time `bson:"utime"`     //本条记录更新unix时间戳
	DayStamp  time.Time `bson:"day_stamp"` //Time Today
	Ctime     time.Time `bson:"ctime"`     //本条记录生成unix时间戳
}

const (
	NOTICE_TYPE1 = 1 //活动公告
	NOTICE_TYPE2 = 2 //广播消息
)

const (
	NOTICE_ACT_TYPE0 = 0 //无操作消息
	NOTICE_ACT_TYPE1 = 1 //支付消息
	NOTICE_ACT_TYPE2 = 2 //活动消息
)

//公告
type Notice struct {
	Id      string    `bson:"_id"`
	Atype   uint32    `bson:"atype"`    //分包类型
	Rtype   int       `bson:"rtype"`    //类型,1=公告消息,2=广播消息
	Acttype int       `bson:"act_type"` //操作类型,0=无操作,1=支付,2=活动
	Top     int       `bson:"top"`      //置顶
	Num     int       `bson:"num"`      //广播次数
	Del     int       `bson:"del"`      //是否移除
	Content string    `bson:"content"`  //广播内容
	Etime   time.Time `bson:"etime"`    //过期时间
	Ctime   time.Time `bson:"ctime"`    //创建时间
}

//商城
type Shop struct {
	Id     string    `bson:"_id"`    //购买ID
	Atype  uint32    `bson:"atype"`  //分包类型
	Status int       `bson:"status"` //物品状态,1=热卖
	Propid int       `bson:"propid"` //兑换的物品,1=钻石
	Payway int       `bson:"payway"` //支付方式,1=RMB
	Number uint32    `bson:"number"` //兑换的数量
	Price  uint32    `bson:"price"`  //支付价格(单位元)
	Name   string    `bson:"name"`   //物品名字
	Info   string    `bson:"info"`   //物品信息
	Del    int       `bson:"del"`    //是否移除
	Etime  time.Time `bson:"etime"`  //过期时间
	Ctime  time.Time `bson:"ctime"`  //创建时间
}

const (
	EnvType1  = 1  //"注册赠送钻石",
	EnvType2  = 2  //"注册赠送金币",
	EnvType3  = 3  //"注册赠送筹码",
	EnvType4  = 4  //"注册赠送房卡",
	EnvType5  = 5  //"绑定赠送",
	EnvType6  = 6  //"首充送n倍",
	EnvType7  = 7  //"首充送金币",
	EnvType8  = 8  //"救济金次数",
	EnvType9  = 9  //"转盘抽奖次数",
	EnvType10 = 10 //"破产金额",
	EnvType11 = 11 //"救济金额",
	EnvType12 = 12 //"虚假人数",
	EnvType13 = 13 //"机器人分配1",
	EnvType14 = 14 //"机器人分配2",
	EnvType15 = 15 //"机器人下注AI",
)

var EnvTypeValue = map[int]string{
	EnvType1: "注册赠送钻石",
	EnvType2: "注册赠送金币",
	EnvType3: "注册赠送筹码",
	EnvType4: "注册赠送房卡",
	//EnvType5:  "绑定赠送",
	//EnvType6:  "首充送n倍",
	//EnvType7:  "首充送金币",
	//EnvType8:  "救济金次数",
	//EnvType9:  "转盘抽奖次数",
	//EnvType10: "破产金额",
	//EnvType11: "救济金额",
	EnvType12: "虚假人数",
	EnvType13: "机器人分配1",
	EnvType14: "机器人分配2",
	EnvType15: "机器人下注AI",
}

var EnvTypeKey = map[int]string{
	EnvType1:  "regist_diamond",
	EnvType2:  "regist_coin",
	EnvType3:  "regist_chip",
	EnvType4:  "regist_card",
	EnvType5:  "build",
	EnvType6:  "first_pay_multi",
	EnvType7:  "first_pay_coin",
	EnvType8:  "relieve",
	EnvType9:  "prizedraw",
	EnvType10: "bankrupt_coin",
	EnvType11: "relieve_coin",
	EnvType12: "robot_num",
	EnvType13: "robot_allot1",
	EnvType14: "robot_allot2",
	EnvType15: "robot_bet",
}

var EnvKeyType = map[string]int{
	"regist_diamond":  EnvType1,
	"regist_coin":     EnvType2,
	"regist_chip":     EnvType3,
	"regist_card":     EnvType4,
	"build":           EnvType5,
	"first_pay_multi": EnvType6,
	"first_pay_coin":  EnvType7,
	"relieve":         EnvType8,
	"prizedraw":       EnvType9,
	"bankrupt_coin":   EnvType10,
	"relieve_coin":    EnvType11,
	"robot_num":       EnvType12,
	"robot_allot1":    EnvType13,
	"robot_allot2":    EnvType14,
	"robot_bet":       EnvType15,
}

//vip
type Vip struct {
	Id     string    `bson:"_id"`    //ID
	Level  int       `bson:"level"`  //等级
	Number uint32    `bson:"number"` //等级充值金额数量限制(分)
	Pay    uint32    `bson:"pay"`    //充值赠送百分比5=赠送充值的5%
	Prize  uint32    `bson:"prize"`  //赠送抽奖次数
	Kick   int       `bson:"kick"`   //经典场可踢人次数
	Ctime  time.Time `bson:"ctime"`  //创建时间
}

type Env struct {
	Key   string `bson:"_id" json:"key"`     //key
	Value int32  `bson:"value" json:"value"` //value
	Name  string `bson:"-"`                  //key
}
