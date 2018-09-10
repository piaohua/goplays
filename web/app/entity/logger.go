package entity

import "time"

//1注册赠送,2开房消耗,3房间解散返还,
//4充值购买,5下注,7上庄，8下庄
//8下庄, 9后台操作,11破产补助
//18商城购买,19绑定赠送,20首充赠送
//23进入房间消耗
//24通比牛牛,25看牌抢庄,26牛牛抢庄
//27代理发放,28vip赠送
//38退款,39本金返还,40输赢,41坐庄输赢
//42异常退款,43抽佣
//44机器人破产补助
//45庄家抽佣
const (
	LogType1  int32 = 1
	LogType2  int32 = 2
	LogType3  int32 = 3
	LogType4  int32 = 4
	LogType5  int32 = 5
	LogType7  int32 = 7
	LogType8  int32 = 8
	LogType9  int32 = 9
	LogType11 int32 = 11
	LogType18 int32 = 18
	LogType19 int32 = 19
	LogType20 int32 = 20
	LogType23 int32 = 23
	LogType24 int32 = 24
	LogType25 int32 = 25
	LogType26 int32 = 26
	LogType27 int32 = 27
	LogType28 int32 = 28
	LogType38 int32 = 38
	LogType39 int32 = 39
	LogType40 int32 = 40
	LogType41 int32 = 41
	LogType42 int32 = 42
	LogType43 int32 = 43
	LogType44 int32 = 44
	LogType45 int32 = 45
)

var LogType = map[int32]string{
	0:         "全部",
	LogType1:  "注册赠送",
	LogType2:  "开房消耗",
	LogType3:  "房间解散返还",
	LogType4:  "充值购买",
	LogType5:  "下注",
	LogType7:  "上庄",
	LogType8:  "下庄",
	LogType9:  "后台操作",
	LogType11: "破产补助",
	LogType18: "商城购买",
	LogType19: "绑定赠送",
	LogType20: "首充赠送",
	LogType23: "进入房间消耗",
	LogType24: "通比牛牛",
	LogType25: "看牌抢庄",
	LogType26: "牛牛抢庄",
	LogType27: "代理发放",
	LogType28: "vip赠送",
	LogType38: "退款",
	LogType39: "本金返还",
	LogType40: "输赢",
	LogType41: "坐庄输赢",
	LogType42: "异常退款",
	LogType43: "抽佣",
	LogType44: "机器人破产补助",
	LogType45: "庄家抽佣",
}

//注册日志
type LogRegist struct {
	Id       string    `bson:"_id"`
	Userid   string    `bson:"userid"`    //账户ID
	Nickname string    `bson:"nickname"`  //账户名称
	Ip       string    `bson:"ip"`        //注册IP
	DayStamp time.Time `bson:"day_stamp"` //regist Time Today
	DayDate  int       `bson:"day_date"`  //regist day date
	Ctime    time.Time `bson:"ctime"`     //create Time
	Atype    uint32    `bson:"atype"`     //regist type
}

//登录日志
type LogLogin struct {
	Id         string    `bson:"_id"`
	Userid     string    `bson:"userid"`      //账户ID
	Event      int       `bson:"event"`       //事件：0=登录,1=正常退出,2＝系统关闭时被迫退出,3＝被动退出,4＝其它情况导致的退出
	Ip         string    `bson:"ip"`          //登录IP
	DayStamp   time.Time `bson:"day_stamp"`   //login Time Today
	LoginTime  time.Time `bson:"login_time"`  //login Time
	LogoutTime time.Time `bson:"logout_time"` //logout Time
	Atype      uint32    `bson:"atype"`       //regist type
}

//钻石日志
type LogDiamond struct {
	Id     string    `bson:"_id"`
	Userid string    `bson:"userid"` //账户ID
	Type   int       `bson:"type"`   //类型
	Num    int32     `bson:"num"`    //数量
	Rest   uint32    `bson:"rest"`   //剩余数量
	Ctime  time.Time `bson:"ctime"`  //create Time
}

//金币日志
type LogCoin struct {
	Id     string    `bson:"_id"`
	Userid string    `bson:"userid"` //账户ID
	Type   int       `bson:"type"`   //类型
	Num    int32     `bson:"num"`    //数量
	Rest   uint32    `bson:"rest"`   //剩余数量
	Ctime  time.Time `bson:"ctime"`  //create Time
}

//筹码日志
type LogChip struct {
	Id     string    `bson:"_id"`
	Userid string    `bson:"userid"` //账户ID
	Type   int       `bson:"type"`   //类型
	Num    int32     `bson:"num"`    //数量
	Rest   uint32    `bson:"rest"`   //剩余数量
	Ctime  time.Time `bson:"ctime"`  //create Time
	Numf   float64   `bson:"-"`      //数量
	Restf  float64   `bson:"-"`      //剩余数量
}

//绑定日志
type LogBuildAgency struct {
	Id       string    `bson:"_id"`
	Userid   string    `bson:"userid"`    //账户ID
	Agent    string    `bson:"agent"`     //绑定ID
	DayStamp time.Time `bson:"day_stamp"` //regist Time Today
	Day      int       `bson:"day"`       //regist day
	Month    int       `bson:"month"`     //regist month
	Ctime    time.Time `bson:"ctime"`     //create Time
}

//在线日志
type LogOnline struct {
	Id       string    `bson:"_id"`
	Num      int       `bson:"num"`       //online count
	DayStamp time.Time `bson:"day_stamp"` //Time Today
	Ctime    time.Time `bson:"ctime"`     //create Time
}

//做牌日志
type LogSetHand struct {
	Id       string    `bson:"_id"`
	Rid      string    `bson:"rid"`      //房间
	Round    int       `bson:"round"`    //
	Userid   string    `bson:"userid"`   //
	Nickname string    `bson:"nickname"` //昵称
	SetHands []uint32  `bson:"sethands"` //设置手牌
	Hands    []uint32  `bson:"hands"`    //手牌
	Niu      int       `bson:"niu"`      //牌力
	Score    int32     `bson:"score"`    //得分
	Ctime    time.Time `bson:"ctime"`    //create Time
}

//今日充值
type LogPayToday struct {
	Id       string    `bson:"_id"`
	Num      int       `bson:"num"`       //
	Count    int       `bson:"count"`     //
	Money    uint32    `bson:"money"`     //
	Diamond  uint32    `bson:"diamond"`   //
	DayStamp time.Time `bson:"day_stamp"` //Time Today
	Ctime    time.Time `bson:"ctime"`     //create Time
}

//今日注册
type LogRegistToday struct {
	Id       string    `bson:"_id"`
	Num      int       `bson:"num"`       //
	DayStamp time.Time `bson:"day_stamp"` //Time Today
	Ctime    time.Time `bson:"ctime"`     //create Time
}

//今日盈亏
type LogChipToday struct {
	Id          string    `bson:"_id"`
	Gametype    int       `bson:"gametype"`
	Roomtype    int       `bson:"roomtype"`
	Lotterytype int       `bson:"lotterytype"`
	RobotNum    int64     `bson:"robot_num"` //机器人盈亏
	RolesNum    int64     `bson:"roles_num"` //玩家盈亏
	DayStamp    time.Time `bson:"day_stamp"` //Time Today
	Ctime       time.Time `bson:"ctime"`     //create Time
	RobotNumf   float64   `bson:"-"`         //机器人盈亏
	RolesNumf   float64   `bson:"-"`         //玩家盈亏
}
