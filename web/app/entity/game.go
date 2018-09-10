package entity

import (
	"time"
)

//游戏房间
type Game struct {
	Id     string    `bson:"_id" json:"id"`        //unique ID
	Gtype  int       `bson:"gtype" json:"gtype"`   //游戏类型1 niu,2 san,3 jiu
	Rtype  int       `bson:"rtype" json:"rtype"`   //房间类型0免佣,1抽佣
	Ltype  int       `bson:"ltype" json:"ltype"`   //彩票类型1bjpk10,1mlaft
	Name   string    `bson:"name" json:"name"`     //房间名称
	Status int       `bson:"status" json:"status"` //房间状态1打开,2关闭,3隐藏
	Count  uint32    `bson:"count" json:"count"`   //房间限制人数
	Ante   uint32    `bson:"ante" json:"ante"`     //房间底分
	Cost   uint32    `bson:"cost" json:"cost"`     //房间抽佣百分比
	Vip    uint32    `bson:"vip" json:"vip"`       //房间vip限制
	Chip   uint32    `bson:"chip" json:"chip"`     //房间进入筹码限制
	Deal   bool      `bson:"deal" json:"deal"`     //房间是否可以上庄
	Carry  uint32    `bson:"carry" json:"carry"`   //房间上庄最小携带筹码限制
	Down   uint32    `bson:"down" json:"down"`     //房间下庄最小携带筹码限制
	Top    uint32    `bson:"top" json:"top"`       //下庄最高携带限制
	Sit    uint32    `bson:"sit" json:"sit"`       //房间内坐下限制
	Del    int       `bson:"del" json:"del"`       //是否移除
	Node   string    `bson:"node" json:"node"`     //所在节点(game.huiyin1|game.huiyin2)
	Ctime  time.Time `bson:"ctime" json:"ctime"`   //创建时间
	//Num   uint32    `bson:"num" json:"num"`      //启动房间数量
}

//所属节点
var GameNodesName = map[string]string{
	"game.huiyin1": "节点1区赛车彩种",
	"game.huiyin2": "节点2区飞艇彩种",
}

//所属节点
var GameNodes = map[int]string{
	1: "节点1区赛车彩种",
	2: "节点2区飞艇彩种",
}
var Game2Nodes = map[int]string{
	1: "game.huiyin1",
	2: "game.huiyin2",
}

//彩种类型
var LotteryTypes = map[int]string{
	1: "赛车彩种",
	2: "飞艇彩种",
}

//玩法类型
var GameTypes = map[int]string{
	1: "牛牛",
	2: "三公",
	3: "牌九",
}

//房间类型
var RoomTypes = map[int]string{
	0: "免佣",
	1: "抽佣",
}

//房间状态
var RoomStatus = map[int]string{
	1: "打开",
	2: "关闭",
	3: "隐藏",
}

//是否上庄
var IsDeal = map[int]string{
	0: "否",
	1: "是",
}
var Is2Deal = map[int]bool{
	0: false,
	1: true,
}

//彩种类型
var LotteryCodes = map[int]string{
	1: "bjpk10",
	2: "mlaft",
}

//是否机器人
var IsRobot = map[bool]string{
	false: "玩家",
	true:  "机器人",
}
