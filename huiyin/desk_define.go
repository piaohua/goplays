package main

import (
	"goplays/data"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//房间牌桌数据结构
type Desk struct {
	//房间id
	id string
	//房间状态,0准备中,1游戏中,2休息时间,3结算
	state uint32
	//下一个状态时间点
	nexttime int64
	//计时
	timer int
	//房间操作数据
	*HuiYinDeskData
	//房间类型基础数据
	*data.DeskData
	//房间无座玩家,数据实时同步
	players map[string]*data.User
	//userid:seat (seat:1~8)
	seats map[string]uint32
	//userid-playerPid
	pids map[string]*actor.PID
	//playerPid-userid
	router map[string]string
	//离线状态userid-bool
	offline map[string]bool
	//关闭通道
	stopCh chan struct{}
	//name
	Name string
	//trends 趋势,开奖结果
	Trends []*data.Trend
	//上局赢家
	Winers []*data.Winer
}

//房间操作数据
type HuiYinDeskData struct {
	lastexpect        string   //上期期号
	lastopencode      string   //上期号码
	lastopentime      string   //上期开奖时间
	lastopentimestamp int64    //上期开奖时间
	lastpower         []uint32 //上期牌力,按位置排序
	//
	expect        string //期号
	opencode      string //号码
	opentime      string //开奖时间
	opentimestamp int64  //开奖时间
	//庄
	dealerNum  uint32 //做庄次数
	dealer     string //庄家
	dealerSeat uint32 //庄家位置
	carry      int64  //庄家的携带
	dealerDown bool   //结束后庄家是否下庄
	betNum     int64  //当前局下注总数
	//上庄列表,userid: carry
	dealers []map[string]int64
	//userid:num, 玩家下注金额
	roleBets map[string]int64
	//seat:num, 位置下注金额
	seatBets map[uint32]int64
	//手牌 seat:cards,seat=(1,2,3,4,5)
	handCards map[uint32][]uint32
	//位置(1-5)对应牌力
	power map[uint32]uint32
	//位置大小排序
	seatPower []uint32
	//结果 seat:num,seat=(1,2,3,4,5),倍数
	multiple map[uint32]int64
	//位置(1-5)上每个玩家下注seat - userid - chip
	seatRoleBets map[uint32]map[string]int64
	//本金返还(赢家位置成本,无返还就无记录)seat - userid - chip
	seatRoleCost map[uint32]map[string]int64
	//位置(1-5)上每个玩家输赢(不含本金,负表示输,输家不含退款)seat - userid - chip
	seatRoleWins map[uint32]map[string]int64
	//位置退款明细seat - userid - chip
	seatRoleRefund map[uint32]map[string]int64
	//庄家在每个位置上的输赢,记录庄家赔付明细(暂时只记录)seat-chip
	dealerSeatWins map[uint32]int64
	//最后开始统计位置和玩家
	//每个玩家总赢利所得(不含退款和本金)userid-chip
	roleWins map[string]int64
	//每个玩家输掉金额userid-chip
	roleLose map[string]int64
	//每个玩家总输赢(纯利,含退款和本金)userid-chip
	roleProfits map[string]int64
	//成本返还(赢家位置)userid-chip
	roleCost map[string]int64
	//玩家退款userid-chip
	roleRefund map[string]int64
	//位置退款总额seat - chip
	seatRefund map[uint32]int64
	//位置(1-5)输赢总量seat-chip
	seatWins map[uint32]int64
	//开牌结果
	result []data.TrendResult
	//抽佣数量
	feeNum int64
	//抽佣数量
	robotNum int64
	//抽佣数量
	playerNum int64
	//抽佣明细userid - fee
	fees map[string]int64
	//位置抽佣总额seat - chip
	seatFees map[uint32]int64
	//位置(1-5)上每个玩家抽佣明细seat - userid - chip
	seatRoleFees map[uint32]map[string]int64
	//位置(1-5)上每个玩家下注(抽佣后)seat - userid - chip
	seatRoleFeeBets map[uint32]map[string]int64
}
