package data

const (
	//游戏类型
	GAME_NIU uint32 = 1 //牛牛(二个位置)
	GAME_SAN uint32 = 2 //三公(五个位置)
	GAME_JIU uint32 = 3 //牌九(五个位置)
	//房间类型
	ROOM_TYPE0 uint32 = 0 //免佣
	ROOM_TYPE1 uint32 = 1 //抽佣
)

const (
	//房间状态
	STATE_READY uint32 = 0 //准备状态,获取期号和开奖时间
	STATE_BET   uint32 = 1 //下注状态
	STATE_SEAL  uint32 = 2 //封盘状态
	STATE_OVER  uint32 = 3 //结算状态
	STATE_STOP  uint32 = 4 //停止状态
	//结算状态
	//房间位置
	SEAT1 uint32 = 1 //庄
	SEAT2 uint32 = 2 //天
	SEAT3 uint32 = 3 //地
	SEAT4 uint32 = 4 //玄
	SEAT5 uint32 = 5 //黄
	//房间0下庄 1上庄 2补庄
	DEALER_DOWN uint32 = 0 //下
	DEALER_UP   uint32 = 1 //上
	DEALER_BU   uint32 = 2 //补
)

const (
	//开奖种类
	GAME_BJPK10 uint32 = 1 //bjpk10
	GAME_MLAFT  uint32 = 2 //mlaft

	BJPK10 string = "bjpk10"
	MLAFT  string = "mlaft"
)

//房间基础数据
type DeskData struct {
	Rid    string `json:"rid"`    //房间ID
	Unique string `json:"unique"` //配置表唯一ID
	Gtype  uint32 `json:"gtype"`  //游戏类型
	Rtype  uint32 `json:"rtype"`  //房间类型
	Ltype  uint32 `json:"ltype"`  //彩票类型
	Rname  string `json:"rname"`  //房间名字
	Count  uint32 `json:"count"`  //牌局人数限制
	Ante   uint32 `json:"ante"`   //底分
	Cost   uint32 `json:"cost"`   //抽佣百分比
	Vip    uint32 `json:"vip"`    //vip限制
	Chip   uint32 `json:"chip"`   //chip限制
	Deal   bool   `json:"deal"`   //房间是否可以上庄
	Carry  uint32 `json:"carry"`  //上庄携带限制
	Down   uint32 `json:"down"`   //下庄携带限制
	Top    uint32 `json:"top"`    //下庄最高携带限制
	Sit    uint32 `json:"sit"`    //房间内坐下限制
	Ctime  uint32 `json:"ctime"`  //创建时间
	//Code   string `json:"code"`   //房间邀请码
}
