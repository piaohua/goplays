/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-12-17 18:22:30
 * Filename      : free.go
 * Description   : 自由场协议消息请求
 * *******************************************************/
package handler

import (
	"time"

	"goplays/data"
	"goplays/game/config"
	"goplays/glog"

	"gopkg.in/mgo.v2/bson"

	jsoniter "github.com/json-iterator/go"
)

//打包
func Desk2Data(deskData *data.DeskData) []byte {
	result, err := jsoniter.Marshal(deskData)
	if err != nil {
		glog.Errorf("Desk2Data Marshal err %v", err)
		return []byte{}
	}
	return result
}

//解析
func Data2Desk(deskDataStr []byte) *data.DeskData {
	deskData := new(data.DeskData)
	err := jsoniter.Unmarshal(deskDataStr, deskData)
	if err != nil {
		glog.Errorf("Data2Desk Unmarshal err %v", err)
		return nil
	}
	return deskData
}

//打包
func NewDeskData(d *data.Game) *data.DeskData {
	return &data.DeskData{
		Unique: d.Id,
		Gtype:  d.Gtype,
		Rtype:  d.Rtype,
		Ltype:  d.Ltype,
		Rname:  d.Name,
		Count:  d.Count,
		Ante:   d.Ante,
		Cost:   d.Cost,
		Vip:    d.Vip,
		Chip:   d.Chip,
		Deal:   d.Deal,
		Carry:  d.Carry,
		Down:   d.Down,
		Top:    d.Top,
		Sit:    d.Sit,
	}
}

//测试数据
func SetGameList() {
	g1 := data.Game{
		Id:     bson.NewObjectId().Hex(),
		Gtype:  data.GAME_NIU,
		Rtype:  data.ROOM_TYPE0,
		Ltype:  data.GAME_BJPK10,
		Name:   "牛牛1区",
		Status: 1,
		Count:  100,
		Ante:   1,
		Cost:   5,
		Vip:    0,
		Chip:   0,
		Deal:   true,
		Carry:  20000,
		Down:   10000,
		Top:    60000,
		Sit:    20000,
		Del:    1,
		Ctime:  time.Now(),
	}
	g2 := data.Game{
		Id:     bson.NewObjectId().Hex(),
		Gtype:  data.GAME_SAN,
		Rtype:  data.ROOM_TYPE0,
		Ltype:  data.GAME_BJPK10,
		Name:   "三公1区",
		Status: 1,
		Count:  100,
		Ante:   1,
		Cost:   5,
		Vip:    0,
		Chip:   0,
		Deal:   true,
		Carry:  20000,
		Down:   10000,
		Top:    60000,
		Sit:    20000,
		Del:    1,
		Ctime:  time.Now(),
	}
	g3 := data.Game{
		Id:     bson.NewObjectId().Hex(),
		Gtype:  data.GAME_JIU,
		Rtype:  data.ROOM_TYPE0, //免佣房间
		Ltype:  data.GAME_BJPK10,
		Name:   "北京赛车1区",
		Status: 1,
		Count:  100,
		Ante:   1,
		Cost:   5,
		Vip:    0,
		Chip:   0,
		Deal:   false, //无庄
		Carry:  20000,
		Down:   10000,
		Top:    60000,
		Sit:    20000,
		Del:    0,
		Node:   "game.huiyin1",
		Ctime:  time.Now(),
	}
	g4 := data.Game{
		Id:     bson.NewObjectId().Hex(),
		Gtype:  data.GAME_JIU,
		Rtype:  data.ROOM_TYPE1, //抽佣房间
		Ltype:  data.GAME_BJPK10,
		Name:   "北京赛车2区",
		Status: 1,
		Count:  100,
		Ante:   1,
		Cost:   5,
		Vip:    0,
		Chip:   0,
		Deal:   false, //无庄
		Carry:  20000,
		Down:   10000,
		Top:    60000,
		Sit:    20000,
		Del:    0,
		Node:   "game.huiyin1",
		Ctime:  time.Now(),
	}
	g5 := data.Game{
		Id:     bson.NewObjectId().Hex(),
		Gtype:  data.GAME_JIU,
		Rtype:  data.ROOM_TYPE0, //免佣房间
		Ltype:  data.GAME_MLAFT,
		Name:   "幸运飞艇1区",
		Status: 1,
		Count:  100,
		Ante:   1,
		Cost:   5,
		Vip:    0,
		Chip:   0,
		Deal:   false, //有庄
		Carry:  20000,
		Down:   10000,
		Top:    60000,
		Sit:    20000,
		Del:    0,
		Node:   "game.huiyin2",
		Ctime:  time.Now(),
	}
	g6 := data.Game{
		Id:     bson.NewObjectId().Hex(),
		Gtype:  data.GAME_JIU,
		Rtype:  data.ROOM_TYPE1, //抽佣房间
		Ltype:  data.GAME_MLAFT,
		Name:   "幸运飞艇2区",
		Status: 1,
		Count:  100,
		Ante:   1,
		Cost:   5,
		Vip:    0,
		Chip:   0,
		Deal:   false, //有庄
		Carry:  20000,
		Down:   10000,
		Top:    60000,
		Sit:    20000,
		Del:    0,
		Node:   "game.huiyin2",
		Ctime:  time.Now(),
	}
	config.SetGame(g1)
	config.SetGame(g2)
	config.SetGame(g3)
	config.SetGame(g4)
	config.SetGame(g5)
	config.SetGame(g6)
}

func SetShopList() {
	s1 := data.Shop{
		Id:     "111",
		Status: 1,
		Propid: 1,
		Payway: 1,
		Number: 10,
		Price:  1000,
		Name:   "筹码",
		Info:   "筹码",
		Del:    0,
		Etime:  time.Now().AddDate(0, 0, 5),
		Ctime:  time.Now(),
	}
	s2 := data.Shop{
		Id:     "112",
		Status: 1,
		Propid: 1,
		Payway: 1,
		Number: 10,
		Price:  1000,
		Name:   "筹码",
		Info:   "筹码",
		Del:    0,
		Etime:  time.Now().AddDate(0, 0, 5),
		Ctime:  time.Now(),
	}
	s3 := data.Shop{
		Id:     "113",
		Status: 1,
		Propid: 1,
		Payway: 1,
		Number: 10,
		Price:  1000,
		Name:   "筹码",
		Info:   "筹码",
		Del:    0,
		Etime:  time.Now().AddDate(0, 0, 5),
		Ctime:  time.Now(),
	}
	config.SetShop(s1)
	config.SetShop(s2)
	config.SetShop(s3)
}
