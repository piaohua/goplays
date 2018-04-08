/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2018-01-22 17:06:19
 * Filename      : sender.go
 * Description   : 机器人
 * *******************************************************/
package main

import (
	"crypto/md5"
	"encoding/hex"
	"sync"

	"goplays/data"
	"goplays/glog"
	"goplays/pb"
	"utils"
)

//' 登录
// 发送注册请求
func (c *Robot) SendRegist() {
	ctos := new(pb.CRegist)
	ctos.Phone = c.data.Phone
	ctos.Nickname = c.data.Nickname
	h := md5.New()
	passwd := cfg.Section("robot").Key("passwd").Value()
	h.Write([]byte(passwd)) // 需要加密的字符串为
	pwd := hex.EncodeToString(h.Sum(nil))
	ctos.Password = pwd
	c.Sender(ctos)
}

// 发送登录请求
func (c *Robot) SendLogin() {
	ctos := new(pb.CLogin)
	ctos.Phone = c.data.Phone
	h := md5.New()
	passwd := cfg.Section("robot").Key("passwd").Value()
	h.Write([]byte(passwd)) // 需要加密的字符串为
	pwd := hex.EncodeToString(h.Sum(nil))
	ctos.Password = pwd
	//glog.Infof("ctos -> %#v", ctos)
	utils.Sleep(2)
	c.Sender(ctos)
}

// 获取玩家数据
func (c *Robot) SendUserData() {
	ctos := new(pb.CUserData)
	ctos.Userid = c.data.Userid
	c.Sender(ctos)
}

// 解散
func (c *Robot) SendPing() {
	ctos := new(pb.CPing)
	//ctos.Time := uint32(utils.Timestamp())
	ctos.Time = 1
	//glog.Debugf("ping : %#v", ctos)
	c.Sender(ctos)
}

func (c *Robot) AddCurrency() {
	msg4 := &pb.PayCurrency{
		Userid: c.data.Userid,
		Type:   data.LogType44,
		Chip:   200000,
	}
	hallPid.Tell(msg4)
}

//.

//' huiyin

// 游戏列表
func (c *Robot) SendGames() {
	ctos := new(pb.CHuiYinGames)
	c.Sender(ctos)
}

// 房间列表
func (c *Robot) SendRoomList() {
	ctos := new(pb.CHuiYinRoomList)
	//lottery type 1 赛车, 2 飞艇
	ctos.Ltype = c.ltype
	c.Sender(ctos)
}

// 离开
func (c *Robot) SendLeave() {
	ctos := new(pb.CHuiYinLeave)
	c.Sender(ctos)
}

// 进入房间
func (c *Robot) SendEntryRoom(roomid string) {
	ctos := new(pb.CHuiYinEnterRoom)
	ctos.Roomid = roomid
	glog.Debugf("enter roomid %s", roomid)
	utils.Sleep(2)
	c.Sender(ctos)
}

// 玩家入坐
func (c *Robot) SendSitDown() {
	seat := uint32(utils.RandInt32N(4) + 1) //随机
	ctos := &pb.CHuiYinSit{
		State: true,
		Seat:  seat,
	}
	c.sits++ //尝试次数
	utils.Sleep(2)
	//TODO 暂时取消坐下
	if false {
		c.Sender(ctos)
	}
}

// 玩家离坐
func (c *Robot) SendStandup() {
	ctos := &pb.CHuiYinSit{
		State: false,
		Seat:  c.seat,
	}
	utils.Sleep(2)
	//TODO 暂时取消坐下
	if false {
		c.Sender(ctos)
	}
	utils.Sleep(2)
	c.SendLeave()
	utils.Sleep(2)
	c.Close() //下线
}

// 玩家下注
func (c *Robot) SendRoomBet() {
	if c.envBet > 0 {
		//c.SendRoomBet2()
		//c.SendRoomBet3()
		//return
	}
	//不同游戏位置不同
	var a1 []uint32 = []uint32{1, 2, 3, 4, 5}
	var c1 []uint32 = []uint32{10, 50, 100, 500, 1000, 3000}
	var coin uint32 = uint32(c.data.Chip) / 4
	var i2 int
	for i := 5; i >= 0; i-- {
		if coin >= c1[i] {
			i2 = i
			break
		}
	}
	var val int
	switch i2 {
	case 0:
		val = i2
	default:
		val = int(utils.RandInt32N(int32(i2))) //随机
	}
	var i1 int32 = utils.RandInt32N(5) //随机
	ctos := &pb.CHuiYinRoomBet{
		Value:   c1[val],
		Seatbet: a1[i1],
	}
	c.bits -= 1
	c.bitNum -= c1[val]
	var t1 int = utils.RandIntN(3) + 1 //随机
	utils.Sleep(t1)
	c.Sender(ctos)
}

// 玩家下注
func (c *Robot) SendRoomBet2() {
	//不同游戏位置不同
	var a1 []uint32 = []uint32{1, 2, 3, 4, 5}
	var val uint32
	var r1 int = utils.RandIntN(1018) //随机
	if r1 < 200 {
	} else if r1 < 650 {
		if c.data.Chip < 10000 {
			val = 10
		} else {
			val = 100
		}
	} else if r1 < 998 {
		if c.data.Chip < 50000 {
			val = 50
		} else {
			val = 500
		}
	} else {
		//chip 单位为分
		if c.data.Chip < 100000 {
			val = 100
		} else {
			val = 1000
		}
	}
	var i1 int32 = utils.RandInt32N(5) //随机
	ctos := &pb.CHuiYinRoomBet{
		Value:   val,
		Seatbet: a1[i1],
	}
	c.bits -= 1
	c.bitNum -= val
	var t1 int = utils.RandIntN(10) + 1 //随机
	utils.Sleep(t1)
	if val > 0 {
		c.Sender(ctos)
	}
}

// 玩家下注
func (c *Robot) SendRoomBet3() {
	//不同游戏位置不同
	var val uint32 = c.betValue()
	var i1 int32 = utils.RandInt32N(5) //随机
	ctos := &pb.CHuiYinRoomBet{
		Value:   val,
		Seatbet: uint32(i1 + 1),
	}
	c.bits -= 1
	c.bitNum -= val
	var t1 int = utils.RandIntN(10) + 1 //随机
	utils.Sleep(t1)
	if val > 0 {
		c.Sender(ctos)
	}
}

// 玩家下注
func (c *Robot) betValue() (val uint32) {
	//chip 单位为分
	if c.data.Chip > 300000 {
		val = 3000
	} else if c.data.Chip > 100000 {
		val = 1000
	} else if c.data.Chip > 50000 {
		val = 500
	} else if c.data.Chip > 10000 {
		val = 100
	} else if c.data.Chip > 5000 {
		val = 50
	} else if c.data.Chip > 1000 {
		val = 10
	}
	return
}

//.

//' 新的下注规则

//下注限制管理, TODO 优化在房间内主动限制
type RobotBet struct {
	sync.Mutex //互斥锁
	//
	state uint32 //状态
	//
	RoomBets map[string]int64  //房间已经下注数量
	RoomType map[string]uint32 //房间类型
	BetLimit map[string]int64  //房间下注限制数量
}

var rbet *RobotBet

func init() {
	rbet = new(RobotBet)
	rbet.RoomBets = make(map[string]int64)
	rbet.RoomType = make(map[string]uint32)
	rbet.BetLimit = make(map[string]int64)
}

func (r *RobotBet) SetRoom(list []*pb.HuiYinRoom) {
	r.Lock()
	defer r.Unlock()
	for _, v := range list {
		roomid := v.GetInfo().Roomid
		rtype := v.GetInfo().Rtype
		if _, ok := r.RoomType[roomid]; !ok {
			r.RoomType[roomid] = rtype
			r.RoomBets[roomid] = 0
			r.BetLimit[roomid] = 0
		}
	}
}

func (r *RobotBet) SetBet(roomid string, value int64) bool {
	r.Lock()
	defer r.Unlock()
	r.RoomBets[roomid] += value
	return r.RoomBets[roomid] >= r.BetLimit[roomid]
}

func (r *RobotBet) SetState() {
	r.Lock()
	defer r.Unlock()
	r.state = data.STATE_BET
}

func (r *RobotBet) Reset() {
	r.Lock()
	defer r.Unlock()
	if r.state == data.STATE_OVER {
		return
	}
	r.RoomBets = make(map[string]int64)
	//重置限制
	for k, v := range r.RoomType {
		switch v {
		case data.ROOM_TYPE0:
			r.BetLimit[k] = int64(betNum0[utils.RandIntN(len(betNum0))])
		case data.ROOM_TYPE1:
			r.BetLimit[k] = int64(betNum1[utils.RandIntN(len(betNum1))])
		}
	}
	r.state = data.STATE_OVER
}

// 玩家下注限制
var betNum0 = []uint32{10000, 10080, 10100, 10500, 10800, 11000, 11200, 11400, 11800, 12000, 12400, 12600, 13000, 13200, 13500, 13800, 14600, 15000}
var betNum1 = []uint32{1000, 1200, 1500, 1800, 2500, 2600, 3000, 3200, 3600, 3800, 4200, 4600, 4800, 5000, 5500, 6800, 5800, 6300, 6500, 7000}
var bet0 = []uint32{3000, 100, 1000, 500}
var bet1 = []uint32{50, 100, 10, 1000, 500}

//开始下注
func (r *Robot) SendRoomBet4() {
	//达到下注限制,TODO 优化，只有在房间内才能真正限制到
	if rbet.SetBet(r.roomid, 0) {
		return
	}
	switch r.rtype {
	case data.ROOM_TYPE0: //免佣
		//随机下注次数
		r.bits = 1
	case data.ROOM_TYPE1: //抽佣
		//随机下注次数
		r.bits = uint32(utils.RandIntN(5) + 1)
	default:
		r.SendStandup()
	}
	if r.betSeat == 0 {
		r.setBetSeat() // 位置固定
	}
	r.SendRoomBet5()
}

func (r *Robot) SendRoomBet5() {
	val := r.getBetVal()
	if val == 0 {
		return
	}
	ctos := &pb.CHuiYinRoomBet{
		Value:   val,
		Seatbet: r.betSeat,
	}
	var t1 int = utils.RandIntN(10) + 1 //随机
	utils.Sleep(t1)
	r.Sender(ctos)
}

// 玩家选择下注金额
func (r *Robot) getBetVal() uint32 {
	var c1 []uint32
	switch r.rtype {
	case data.ROOM_TYPE0: //免佣
		for _, v := range bet0 {
			if r.data.Chip > int64(v) {
				c1 = append(c1, v)
			}
		}
	case data.ROOM_TYPE1: //抽佣
		for _, v := range bet1 {
			if r.data.Chip > int64(v) {
				c1 = append(c1, v)
			}
		}
	}
	if len(c1) == 0 {
		return 0
	}
	return c1[utils.RandIntN(len(c1))] //随机
}

//.

// vim: set foldmethod=marker foldmarker=//',//.:
