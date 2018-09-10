package service

import (
	"errors"

	"goplays/web/app/entity"

	"gopkg.in/mgo.v2/bson"
)

// 添加商品
func (this *playerService) AddGame(shop *entity.Game) error {
	//shop.Id = utils.String(utils.Timestamp())
	shop.Id = bson.NewObjectId().Hex()
	shop.Ctime = bson.Now()
	shop.Chip *= 100 //转换为分
	if !Insert(Games, shop) {
		return errors.New("写入失败:" + shop.Id)
	}
	return nil
}

// 获取商品
func (this *playerService) GetGame(id string) (entity.Game, error) {
	shop := entity.Game{}
	Get(Games, id, &shop)
	if shop.Id == "" {
		return shop, errors.New("不存在")
	}
	return shop, nil
}

// 移除商品
func (this *playerService) DelGame(id string) error {
	if Update(Games, bson.M{"_id": id}, bson.M{"$set": bson.M{"del": 1}}) {
		return nil
	}
	return errors.New("移除失败")
}

// 更新商品
func (this *playerService) UpdateGame(id string, cost uint32) error {
	m := bson.M{"cost": cost}
	if Update(Games, bson.M{"_id": id}, bson.M{"$set": m}) {
		return nil
	}
	return errors.New("移除失败")
}

// 获取列表
func (this *playerService) GetGameList(page, pageSize int, m bson.M) ([]entity.Game, error) {
	var list []entity.Game
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	err := Games.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)

	//转换为分展示
	list = this.gameList2(list)
	return list, err
}

//转换为分展示
func (this *playerService) gameList2(list []entity.Game) []entity.Game {
	for k, v := range list {
		v.Chip = v.Chip / 100
		list[k] = v
	}
	return list
}

// 获取总数
func (this *playerService) GetGameListTotal(m bson.M) (int64, error) {
	return int64(Count(Games, m)), nil
}
