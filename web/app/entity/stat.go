package entity

import (
	"time"
)

//statistics

//每日赢亏统计
type ProfitStat struct {
	//Id     string    `bson:"_id"`
	Userid    string    `bson:"userid"`
	Robot     bool      `bson:"robot"` // 是否是机器人
	Day       int       `bson:"day"`   //20180205
	Month     int       `bson:"month"` //201802
	Yesterday int64     `bson:"yesterday"`
	Seven     int64     `bson:"seven"`
	Thirty    int64     `bson:"thirty"`
	DayStamp  time.Time `bson:"day_stamp"` //Time Today
	Ctime     time.Time `bson:"ctime"`
}

//赢亏统计,按玩家单个统计
type UserProfitStat struct {
	Userid    string    `bson:"_id"`
	Robot     bool      `bson:"robot"` // 是否是机器人
	Yesterday int64     `bson:"yesterday"`
	Seven     int64     `bson:"seven"`
	Thirty    int64     `bson:"thirty"`
	All       int64     `bson:"all"`
	DayStamp  time.Time `bson:"day_stamp"` //Time Today
	Utime     time.Time `bson:"utime"`     //update time
	Ctime     time.Time `bson:"ctime"`
}

//TODO 昨日数据统计
type YesterdayStat struct {
	//Id     string    `bson:"_id"`
	Day      int       `bson:"day"`       //20180205
	Month    int       `bson:"month"`     //201802
	Agency   int       `bson:"agency"`    //昨天新增代理
	Player   int       `bson:"player"`    //昨天新增玩家
	Fee      int64     `bson:"fee"`       //昨天抽佣
	DayStamp time.Time `bson:"day_stamp"` //Time Today
	Ctime    time.Time `bson:"ctime"`
}

//实时统计 TODO 优化
type UserStat struct {
	//Id     string    `bson:"_id"`
	TodayNewPlayers     int64 `bson:"today_new_players"`     //今日新增会员(玩家)
	YesterdayNewPlayers int64 `bson:"yesterday_new_players"` //昨天新增
	AllPlayers          int64 `bson:"all_players"`           //会员总数
	Chips               int64 `bson:"chips"`                 //用户持有筹码总数
	Cards               int64 `bson:"cards"`                 //用户持有房卡总数
	Coins               int64 `bson:"coins"`                 //用户持有金币总数
	TodayFee            int64 `bson:"today_fee"`             //今日抽佣
	YesterdayFee        int64 `bson:"yesterday_fee"`         //昨天抽佣
	AllFee              int64 `bson:"all_fee"`               //抽佣总数
	TodayNewAgent       int64 `bson:"today_new_agent"`       //今日新增代理
	YesterdayNewAgent   int64 `bson:"yesterday_new_agent"`   //昨天新增代理
	AllAgent            int64 `bson:"all_agent"`             //代理总数
	//
	TodayAgentFee     int64 `bson:"today_agent_fee"`     //今日代理抽佣
	YesterdayAgentFee int64 `bson:"yesterday_agent_fee"` //昨天代理抽佣
	AgentAllFee       int64 `bson:"agent_all_fee"`       //代理抽佣总数

	Chipsf        float64 `bson:"-"` //用户持有筹码总数
	TodayFeef     float64 `bson:"-"` //今日抽佣
	YesterdayFeef float64 `bson:"-"` //昨天抽佣
	AllFeef       float64 `bson:"-"` //抽佣总数
	//
	TodayAgentFeef     float64 `bson:"-"` //今日代理抽佣
	YesterdayAgentFeef float64 `bson:"-"` //昨天代理抽佣
	AgentAllFeef       float64 `bson:"-"` //代理抽佣总数
}

//实时统计 TODO 优化
type AgencyStat struct {
	//Id     string    `bson:"_id"`
	Chip            int64 `bson:"chip"`             //用户持有筹码数
	Card            int64 `bson:"card"`             //用户持有房卡数
	Coin            int64 `bson:"coin"`             //用户持有金币数
	UnderlingAgency int   `bson:"underling_agency"` //下属代理数
	UnderlingPlayer int   `bson:"underling_player"` //下属玩家数
	TodayBets       int64 `bson:"today_bets"`       //今日注额
	TodayProfits    int64 `bson:"today_profits"`    //今日收益
	AllProfits      int64 `bson:"all_profits"`      //全部收益(剩余可提取)
	TopProfits      int64 `bson:"top_profits"`      //历史最高
	ExtractProfits  int64 `bson:"extract_profits"`  //已提取

	Chipf           float64 `bson:"-"` //用户持有筹码数
	TodayBetsf      float64 `bson:"-"` //今日注额
	TodayProfitsf   float64 `bson:"-"` //今日收益
	AllProfitsf     float64 `bson:"-"` //全部收益(剩余可提取)
	TopProfitsf     float64 `bson:"-"` //历史最高
	ExtractProfitsf float64 `bson:"-"` //已提取
}

//抽佣明细
type AgentFee struct {
	//Id     string    `bson:"_id"`
	Agent       string `bson:"agent"`        //代理邀请码
	ParentAgent string `bson:"parent_agent"` //上级代理邀请码
	FeeNum      int64  `bson:"fee_num"`      //抽佣总数量
	Fee         int64  `bson:"fee"`          //代理所得
	ParentFee   int64  `bson:"parent_fee"`   //上级代理所得
	//FeeAll      int64     `bson:"fee_all"`      //代理总所得(上级ParentFee+当前代理所得Fee)
	FeeAll  int64     `bson:"fee_all"`  //代理所得
	SysFee  int64     `bson:"sys_fee"`  //系统所得
	Rate    uint32    `bson:"rate"`     //代理分成比例
	SysRate uint32    `bson:"sys_rate"` //系统分成比例
	Ctime   time.Time `bson:"ctime"`
}

//抽佣明细
type AgentFeeLog struct {
	//Id     string    `bson:"_id"`
	Agent      string    `bson:"agent"`        //代理邀请码
	FeeAll     int64     `bson:"fee_all"`      //代理所得变动
	FeeRate    int64     `bson:"fee_rate"`     //抽成给上级数量
	SysFeeRate int64     `bson:"sys_fee_rate"` //抽成给系统数量
	Ctime      time.Time `bson:"ctime"`
}

//账务统计
type AccountingLog struct {
	//Id     string    `bson:"_id"`
	Chips                 int64     `bson:"chips"`                   //玩家当天剩余总筹码
	Bets                  int64     `bson:"bets"`                    //玩家当天总投注额
	Pays                  int64     `bson:"pays"`                    //玩家当天总充值
	AllFee                int64     `bson:"all_fee"`                 //抽佣总数
	AgentAllFee           int64     `bson:"agent_all_fee"`           //代理抽佣总数
	YesterdayFee          int64     `bson:"yesterday_fee"`           //昨天总抽佣
	YesterdayAgentFee     int64     `bson:"yesterday_agent_fee"`     //昨天代理抽佣
	RobotProfitsYesterday int64     `bson:"robot_profits_yesterday"` //机器人昨日盈亏
	SysProfitsYesterday   int64     `bson:"sys_profits_yesterday"`   //系统盈亏 = 当天总抽佣 - 当天代理抽佣 + 当天机器人盈亏
	DayStamp              time.Time `bson:"day_stamp"`               //Time Today
	Ctime                 time.Time `bson:"ctime"`

	//
	Chipsf                 float64 `bson:"-"` //玩家当天剩余总筹码
	Betsf                  float64 `bson:"-"` //玩家当天总投注额
	Paysf                  float64 `bson:"-"` //玩家当天总充值
	AllFeef                float64 `bson:"-"` //抽佣总数
	AgentAllFeef           float64 `bson:"-"` //代理抽佣总数
	YesterdayFeef          float64 `bson:"-"` //昨天总抽佣
	YesterdayAgentFeef     float64 `bson:"-"` //昨天代理抽佣
	RobotProfitsYesterdayf float64 `bson:"-"` //机器人昨日盈亏
	SysProfitsYesterdayf   float64 `bson:"-"` //系统盈亏 = 当天总抽佣 - 当天代理抽佣 + 当天机器人盈亏
}
