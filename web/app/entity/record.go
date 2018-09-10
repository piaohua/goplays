package entity

import (
	"time"
)

//开奖结果记录
type Pk10Record struct {
	Expect        string    `bson:"_id"`
	Opencode      string    `bson:"opencode"`
	Opentime      string    `bson:"opentime"`
	Opentimestamp int64     `bson:"opentimestamp"`
	Code          string    `bson:"code"`
	Ctime         time.Time `bson:"ctime"`
}

//房间单局记录,Roomid = GenOrderid()
//var roomid string = data.GenCporderid(roomid)
type GameRecord struct {
	Roomid string `bson:"_id"` //唯一
	//Gametype    uint32         `bson:"gametype"`
	//Roomtype    uint32         `bson:"roomtype"`
	//Lotterytype uint32         `bson:"lotterytype"`
	Gametype    int            `bson:"gametype"`
	Roomtype    int            `bson:"roomtype"`
	Lotterytype int            `bson:"lotterytype"`
	Expect      string         `bson:"expect"`
	Opencode    string         `bson:"opencode"`
	Opentime    string         `bson:"opentime"`
	Num         uint32         `bson:"num"`        //参与人数
	RobotFee    int64          `bson:"robot_fee"`  //机器人抽佣数量
	PlayerFee   int64          `bson:"player_fee"` //玩家抽佣数量
	FeeNum      int64          `bson:"fee_num"`    //抽佣数量
	BetNum      int64          `bson:"bet_num"`    //下注总数量
	WinNum      int64          `bson:"win_num"`    //赢总数量
	LoseNum     int64          `bson:"lose_num"`   //输总数量
	RefundNum   int64          `bson:"refund_num"` //退款数量
	Trend       []TrendResult  `bson:"seats"`      //位置结果
	Result      []ResultRecord `bson:"result"`     //全部玩家输赢结果
	Record      []FeeResult    `bson:"record"`     //玩家抽佣明细
	Details     []FeeDetails   `bson:"details"`    //位置上玩家抽佣明细
	Ctime       time.Time      `bson:"ctime"`

	RobotFeef  float64 `bson:"-"` //机器人抽佣数量
	PlayerFeef float64 `bson:"-"` //玩家抽佣数量
	FeeNumf    float64 `bson:"-"` //抽佣数量
	BetNumf    float64 `bson:"-"` //下注总数量
	WinNumf    float64 `bson:"-"` //赢总数量
	LoseNumf   float64 `bson:"-"` //输总数量
	RefundNumf float64 `bson:"-"` //退款数量
}

//开牌结果
type TrendResult struct {
	Rank  uint32   `bson:"rank"`  //排名(大小排行1->5)
	Seat  uint32   `bson:"seat"`  //位置(门内,第n门)
	Point uint32   `bson:"point"` //点数
	Cards []uint32 `bson:"cards"` //牌
}

//全部玩家信息
type ResultRecord struct {
	Userid string `bson:"userid"`
	Bets   int64  `bson:"bets"`   //下注总额
	Wins   int64  `bson:"wins"`   //输赢总额(不含本金)
	Refund int64  `bson:"refund"` //退款

	Betsf   float64 `bson:"-"` //下注总额
	Winsf   float64 `bson:"-"` //输赢总额(不含本金)
	Refundf float64 `bson:"-"` //退款
}

//玩家抽佣明细
type FeeResult struct {
	Userid string `bson:"userid"`
	Fee    int64  `bson:"fee"` //抽佣数量

	Feef float64 `bson:"-"` //抽佣数量
}

//位置抽佣明细
type FeeDetails struct {
	Seat   uint32      `bson:"seat"`   //位置
	Fee    int64       `bson:"fee"`    //位置抽佣数量
	Record []FeeResult `bson:"record"` //玩家抽佣明细

	Feef float64 `bson:"-"` //位置抽佣数量
}

//个人单局记录
type UserRecord struct {
	//Id     string    `bson:"_id"`
	Roomid string `bson:"roomid"` //唯一
	//Gametype    uint32          `bson:"gametype"`
	//Roomtype    uint32          `bson:"roomtype"`
	//Lotterytype uint32          `bson:"lotterytype"`
	Gametype    int             `bson:"gametype"`
	Roomtype    int             `bson:"roomtype"`
	Lotterytype int             `bson:"lotterytype"`
	Expect      string          `bson:"expect"`
	Userid      string          `bson:"userid"`
	Robot       bool            `bson:"robot"` // 是否是机器人
	Rest        int64           `bson:"rest"`
	Bets        int64           `bson:"bets"`
	Profits     int64           `bson:"profits"`
	Fee         int64           `bson:"fee"` //抽佣
	Details     []UseridDetails `bson:"details"`
	Ctime       time.Time       `bson:"ctime"`

	Restf    float64 `bson:"-"`
	Betsf    float64 `bson:"-"`
	Profitsf float64 `bson:"-"`
	Feef     float64 `bson:"-"` //抽佣
}

//个人详细结果
type UseridDetails struct {
	Seat   uint32 `bson:"seat"`   //位置
	Bets   int64  `bson:"bets"`   //位置个人下注总额
	Wins   int64  `bson:"wins"`   //位置个人输赢总额(不含本金)
	Refund int64  `bson:"refund"` //退款

	Betsf   float64 `bson:"-"` //位置个人下注总额
	Winsf   float64 `bson:"-"` //位置个人输赢总额(不含本金)
	Refundf float64 `bson:"-"` //退款
}
