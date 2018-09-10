package service

import (
	"errors"

	"goplays/web/app/entity"

	"gopkg.in/mgo.v2/bson"
)

type playerService struct{}

func (this *playerService) MineChild(id string, child []string) error {
	agent, err := this.GetPlayerAgent(id)
	if err != nil {
		return err
	}
	for _, val := range child {
		if val == agent {
			return nil
		}
	}
	return errors.New("玩家不存在")
}

//获取上级代理
func (this *playerService) GetPlayerAgent(id string) (string, error) {
	q := bson.M{"_id": id}
	f := []string{"agent"}
	var feeNum bson.M
	GetByQWithFields(PlayerUsers, q, f, &feeNum)
	if v, ok := feeNum["agent"]; ok {
		return v.(string), nil
	}
	return "", errors.New("用户不存在")
}

// 获取
func (this *playerService) GetPlayer(userid string) (*entity.PlayerUser, error) {
	player := new(entity.PlayerUser)
	Get(PlayerUsers, userid, player)
	if player.Userid == "" {
		return player, errors.New("用户不存在")
	}
	return player, nil
}

// 获取所有
func (this *playerService) GetAllPlayer() ([]entity.PlayerUser, error) {
	return this.GetList(0, 1, -1, bson.M{})
}

// 获取所有下属玩家
func (this *playerService) GetAllBuilds(agent string) []bson.M {
	var list []bson.M
	if agent == "" {
		return list
	}
	q := bson.M{"agent": agent}
	f := []string{"_id"}
	ListByQWithFields(PlayerUsers, q, f, &list)
	return list
}

// 获取所有下属玩家
func (this *playerService) GetAllBuilds2(agent string) (list []string) {
	list = make([]string, 0)
	if agent == "" {
		return list
	}
	list2 := this.GetAllBuilds(agent)
	for _, v := range list2 {
		if v2, ok := v["_id"]; ok {
			list = append(list, v2.(string))
		}
	}
	return list
}

// 获取代理商总数
func (this *playerService) GetBuilds(agent string) int {
	m := bson.M{"agent": agent}
	return Count(PlayerUsers, m)
}

// 获取列表
func (this *playerService) GetList(typeId, page, pageSize int, m bson.M) ([]entity.PlayerUser, error) {
	var list []entity.PlayerUser
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "login_time", false)
	switch typeId {
	case 0: //注册用户
		m["phone"] = bson.M{"$ne": ""}
		m["robot"] = false
	case 1: //游客用户
		m["phone"] = bson.M{"$eq": ""}
		m["tourist"] = bson.M{"$ne": ""}
		m["robot"] = false
	case 2: //机器人
		m["phone"] = bson.M{"$ne": ""}
		m["robot"] = true
	case 3: //全部玩家
	}
	err := PlayerUsers.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)

	//转换为分展示
	list = this.chipList2(list)
	return list, err
}

//转换为分展示
func (this *playerService) chipList2(list []entity.PlayerUser) []entity.PlayerUser {
	for k, v := range list {
		v.Chipf = Chip2Float(int64(v.Chip))
		list[k] = v
	}
	return list
}

// 获取总数
func (this *playerService) GetTotal(typeId int, m bson.M) (int64, error) {
	switch typeId {
	case 0: //注册用户
		m["phone"] = bson.M{"$ne": ""}
		m["robot"] = false
	case 1: //游客用户
		m["phone"] = bson.M{"$eq": ""}
		m["tourist"] = bson.M{"$ne": ""}
		m["robot"] = false
	case 2: //机器人
		m["phone"] = bson.M{"$ne": ""}
		m["robot"] = true
	case 3: //全部玩家
	}
	return int64(Count(PlayerUsers, m)), nil
}

// 获取类型
func (this *playerService) GetPlayerType() ([]int, error) {
	var types []int
	types = []int{1}
	return types, nil
}

// 添加公告
func (this *playerService) AddNotice(notice *entity.Notice) error {
	notice.Id = bson.NewObjectId().Hex()
	notice.Ctime = bson.Now()
	if !Insert(Notices, notice) {
		return errors.New("写入失败:" + notice.Id)
	}
	return nil
}

// 获取
func (this *playerService) GetNotice(id string) (*entity.Notice, error) {
	notice := new(entity.Notice)
	Get(Notices, id, notice)
	if notice.Id == "" {
		return notice, errors.New("公告不存在")
	}
	return notice, nil
}

// 获取
func (this *playerService) DelNotice(id string) error {
	if Update(Notices, bson.M{"_id": id}, bson.M{"$set": bson.M{"del": 1}}) {
		return nil
	}
	return errors.New("移除失败")
}

// 获取列表
func (this *playerService) GetNoticeList(page, pageSize int, m bson.M) ([]entity.Notice, error) {
	var list []entity.Notice
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	err := Notices.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, err
}

// 获取总数
func (this *playerService) GetNoticeListTotal(m bson.M) (int64, error) {
	return int64(Count(Notices, m)), nil
}

// 添加商品
func (this *playerService) AddShop(shop *entity.Shop) error {
	//shop.Id = utils.String(utils.Timestamp())
	//shop.Id = bson.NewObjectId().Hex()
	shop.Ctime = bson.Now()
	if !Insert(Shops, shop) {
		return errors.New("写入失败:" + shop.Id)
	}
	return nil
}

// 获取商品
func (this *playerService) GetShop(id string) (*entity.Shop, error) {
	shop := new(entity.Shop)
	Get(Shops, id, shop)
	if shop.Id == "" {
		return shop, errors.New("商品不存在")
	}
	return shop, nil
}

// 移除商品
func (this *playerService) DelShop(id string) error {
	if Update(Shops, bson.M{"_id": id}, bson.M{"$set": bson.M{"del": 1}}) {
		return nil
	}
	return errors.New("移除失败")
}

// 获取列表
func (this *playerService) GetShopList(page, pageSize int, m bson.M) ([]entity.Shop, error) {
	var list []entity.Shop
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	err := Shops.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, err
}

// 获取总数
func (this *playerService) GetShopListTotal(m bson.M) (int64, error) {
	return int64(Count(Shops, m)), nil
}

// 添加VIP
func (this *playerService) AddVip(shop *entity.Vip) error {
	shop.Ctime = bson.Now()
	if !Insert(Vips, shop) {
		return errors.New("写入失败:" + shop.Id)
	}
	return nil
}

// 获取商品
func (this *playerService) GetVip(id string) (*entity.Vip, error) {
	shop := new(entity.Vip)
	Get(Vips, id, shop)
	if shop.Id == "" {
		return shop, errors.New("商品不存在")
	}
	return shop, nil
}

// 移除商品
func (this *playerService) DelVip(id string) error {
	if Delete(Vips, bson.M{"_id": id}) {
		return nil
	}
	return errors.New("移除失败")
}

// 获取列表
func (this *playerService) GetVipList(page, pageSize int, m bson.M) ([]entity.Vip, error) {
	var list []entity.Vip
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	err := Vips.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, err
}

// 获取总数
func (this *playerService) GetVipListTotal(m bson.M) (int64, error) {
	return int64(Count(Vips, m)), nil
}

// 获取列表
func (this *playerService) GetEnvList(page, pageSize int, m bson.M) ([]entity.Env, error) {
	var list []entity.Env
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	err := Envs.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, err
}

// 获取总数
func (this *playerService) GetEnvListTotal(m bson.M) (int64, error) {
	return int64(Count(Envs, m)), nil
}

// 添加VIP
func (this *playerService) AddEnv(shop *entity.Env) error {
	if !Insert(Envs, shop) {
		return errors.New("写入失败:" + shop.Key)
	}
	return nil
}

// 获取商品
func (this *playerService) GetEnv(id string) (*entity.Env, error) {
	shop := new(entity.Env)
	Get(Envs, id, shop)
	if shop.Key == "" {
		return shop, errors.New("商品不存在")
	}
	return shop, nil
}

// 移除商品
func (this *playerService) DelEnv(id string) error {
	if Delete(Envs, bson.M{"_id": id}) {
		return nil
	}
	return errors.New("移除失败")
}
