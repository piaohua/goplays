package data

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//TODO 添加房间id,人数同步后台显示

//游戏
type Game struct {
	Id     string    `bson:"_id" json:"id"`        //unique ID
	Gtype  uint32    `bson:"gtype" json:"gtype"`   //游戏类型1 niu,2 san,3 jiu
	Rtype  uint32    `bson:"rtype" json:"rtype"`   //房间类型0免佣,1抽佣
	Ltype  uint32    `bson:"ltype" json:"ltype"`   //彩票类型1bjpk10,1mlaft
	Name   string    `bson:"name" json:"name"`     //房间名称
	Status uint32    `bson:"status" json:"status"` //房间状态1打开,2关闭,3隐藏
	Count  uint32    `bson:"count" json:"count"`   //房间限制人数
	Ante   uint32    `bson:"ante" json:"ante"`     //房间底分
	Cost   uint32    `bson:"cost" json:"cost"`     //房间抽佣百分比
	Vip    uint32    `bson:"vip" json:"vip"`       //房间vip限制
	Chip   uint32    `bson:"chip" json:"chip"`     //房间进入筹码限制
	Deal   bool      `bson:"deal" json:"deal"`     //房间是否可以上庄
	Carry  uint32    `bson:"carry" json:"carry"`   //房间上庄最小携带筹码限制
	Down   uint32    `bson:"down" json:"down"`     //房间下庄最小携带筹码限制
	Top    uint32    `bson:"top" json:"Top"`       //房间下庄最大携带筹码限制
	Sit    uint32    `bson:"sit" json:"sit"`       //房间内坐下限制
	Del    int       `bson:"del" json:"del"`       //是否移除
	Node   string    `bson:"node" json:"node"`     //所在节点(game.huiyin1|game.huiyin2)
	Ctime  time.Time `bson:"ctime" json:"ctime"`   //创建时间
	//Num   uint32    `bson:"num" json:"num"`      //启动房间数量
}

func GetGameList() []Game {
	var list []Game
	ListByQ(Games, bson.M{"del": 0}, &list)
	return list
}
