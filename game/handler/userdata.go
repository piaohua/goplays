package handler

import (
	"goplays/data"
	"goplays/game/config"
	"goplays/pb"
	"utils"
)

//获取货币信息
func GetCurrency(ctos *pb.CGetCurrency, p *data.User) (stoc *pb.SGetCurrency) {
	stoc = new(pb.SGetCurrency)
	if p == nil {
		return
	}
	stoc.Data = PackCurrency(p)
	return
}

//打包基础货币消息
func PackCurrency(p *data.User) (msg *pb.Currency) {
	return &pb.Currency{
		Coin:    p.GetCoin(),
		Diamond: p.GetDiamond(),
		Card:    p.GetCard(),
		Chip:    p.GetChip(),
	}
}

//变动货币消息
func CurrencyMsg(diamond, coin, card, chip int64) (msg *pb.Currency) {
	return &pb.Currency{
		Diamond: diamond,
		Coin:    coin,
		Card:    card,
		Chip:    chip,
	}
}

//货币变动推送消息
func PushCurrencyMsg(diamond, coin, card, chip int64,
	ltype int32) (msg *pb.SPushCurrency) {
	msg = new(pb.SPushCurrency)
	msg.Rtype = uint32(ltype)
	msg.Data = CurrencyMsg(diamond, coin, card, chip)
	return
}

//货币变动消息
func ChangeCurrencyMsg(diamond, coin, card, chip int64,
	ltype int32, userid string) (msg *pb.ChangeCurrency) {
	msg = new(pb.ChangeCurrency)
	msg.Type = ltype
	msg.Diamond = diamond
	msg.Coin = coin
	msg.Card = card
	msg.Chip = chip
	msg.Userid = userid
	return
}

//货币变动消息
func Pay2ChangeCurr(arg *pb.PayCurrency) (msg *pb.ChangeCurrency) {
	msg = &pb.ChangeCurrency{
		Userid:  arg.Userid,
		Type:    arg.Type,
		Coin:    arg.Coin,
		Diamond: arg.Diamond,
		Chip:    arg.Chip,
		Card:    arg.Card,
	}
	return
}

//心跳请求
func Ping(ctos *pb.CPing) (stoc *pb.SPing) {
	stoc = new(pb.SPing)
	stoc.Time = ctos.GetTime()
	return
}

// 公告列表
func GetNotice(atype uint32) (stoc *pb.SNotice) {
	stoc = new(pb.SNotice)
	list := config.GetNotices(atype)
	for _, v := range list {
		body := &pb.Notice{
			Rtype:   uint32(v.Rtype),
			Acttype: uint32(v.Acttype),
			Content: []byte(v.Content),
		}
		stoc.List = append(stoc.List, body)
	}
	return
}

//获取自己数据
func GetUserDataMsg(ctos *pb.CUserData, p *data.User) (stoc *pb.SUserData) {
	stoc = new(pb.SUserData)
	stoc.Data = new(pb.UserData)
	userid := ctos.GetUserid()
	if userid == "" {
		stoc.Error = pb.UsernameEmpty
		return
	}
	stoc.Data = PackUserData(p)
	stoc.Info = PackUserTop(p)
	return
}

//打包基础数据
func PackUserData(p *data.User) (stoc *pb.UserData) {
	return &pb.UserData{
		Userid:   p.GetUserid(),
		Nickname: p.GetNickname(),
		Phone:    p.GetPhone(),
		Sex:      p.GetSex(),
		Photo:    p.GetPhoto(),
		Agent:    p.GetAgent(),
		Coin:     p.GetCoin(),
		Diamond:  p.GetDiamond(),
		Card:     p.GetCard(),
		Chip:     p.GetChip(),
		Vip:      p.GetVip(),
	}
}

//玩家个人数据
func PackUserTop(p *data.User) (stoc *pb.TopInfo) {
	return &pb.TopInfo{
		//Topcoins:      p.GetTopCoins(),      //最高拥有金币总金额
		//Topdiamonds:   p.GetTopDiamonds(),   //最高拥有钻石总金额
		//Topcards:      p.GetTopCards(),      //最高拥有房卡总数
		Topchips: p.GetTopChips(), //最高拥有筹码总金额
		//Topwincoin:    p.GetTopWinCoin(),    //单局赢最高金币金额
		//Topwindiamond: p.GetTopWinDiamond(), //单局赢最高钻石金额
		Topwinchip: p.GetTopWinChip(),                              //单局赢最高筹码金额
		Registtime: utils.Format("Y-m-d H:i:s", p.GetRegistTime()), //加入游戏时间
		Logintime:  utils.Format("Y-m-d H:i:s", p.GetLoginTime()),  //最后登录时间
	}
}

//获取其它玩家数据
func GetUserData(p *data.User) (stoc *pb.GotUserData) {
	stoc = new(pb.GotUserData)
	if p == nil {
		stoc.Error = pb.UsernameEmpty
		return
	}
	//基本数据
	stoc.Data = PackUserData(p)
	stoc.Info = PackUserTop(p)
	return
}

//获取其它玩家数据消息
func UserDataMsg(p *pb.GotUserData) (stoc *pb.SUserData) {
	stoc = new(pb.SUserData)
	if p == nil {
		stoc.Error = pb.UsernameEmpty
		return
	}
	if p.Error != pb.OK {
		stoc.Error = p.Error
		return
	}
	//基本数据
	stoc.Data = p.GetData()
	stoc.Info = p.GetInfo()
	return
}

//打包日志消息
func LogDiamondMsg(num int64, ltype int32,
	p *data.User) (msg *pb.LogDiamond) {
	msg = &pb.LogDiamond{
		Userid: p.GetUserid(),
		Type:   int32(ltype),
		Num:    num,
		Rest:   p.GetDiamond(),
	}
	return
}

func LogCoinMsg(num int64, ltype int32,
	p *data.User) (msg *pb.LogCoin) {
	msg = &pb.LogCoin{
		Userid: p.GetUserid(),
		Type:   int32(ltype),
		Num:    num,
		Rest:   p.GetCoin(),
	}
	return
}

func LogCardMsg(num int64, ltype int32,
	p *data.User) (msg *pb.LogCard) {
	msg = &pb.LogCard{
		Userid: p.GetUserid(),
		Type:   int32(ltype),
		Num:    num,
		Rest:   p.GetCard(),
	}
	return
}

func LogChipMsg(num int64, ltype int32,
	p *data.User) (msg *pb.LogChip) {
	msg = &pb.LogChip{
		Userid: p.GetUserid(),
		Type:   int32(ltype),
		Num:    num,
		Rest:   p.GetChip(),
	}
	return
}
