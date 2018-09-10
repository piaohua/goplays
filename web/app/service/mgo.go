package service

import (
	"fmt"
	"os"
	"time"

	"utils"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Init mgo and the common DAO

// 数据连接
var Session *mgo.Session
var mgoCloseCh chan bool

// 各个表的Collection对象
var Actions *mgo.Collection
var Perms *mgo.Collection
var Roles *mgo.Collection
var RolePerms *mgo.Collection
var Users *mgo.Collection
var UserRoles *mgo.Collection
var MailTpls *mgo.Collection
var GenIDs *mgo.Collection

var TradeRecords *mgo.Collection
var PlayerUsers *mgo.Collection
var GameRecords *mgo.Collection
var UserRecords *mgo.Collection
var StatRecords *mgo.Collection
var UserStatRecords *mgo.Collection
var Pk10Records *mgo.Collection

var Agencys *mgo.Collection
var RegistLogs *mgo.Collection
var LoginLogs *mgo.Collection
var DiamondLogs *mgo.Collection
var CoinLogs *mgo.Collection
var ChipLogs *mgo.Collection
var LogBuildAgencys *mgo.Collection
var LogOnlines *mgo.Collection
var LogPayTodays *mgo.Collection
var LogRegistTodays *mgo.Collection
var ApplyCashs *mgo.Collection
var LogChipTodays *mgo.Collection
var AgentFees *mgo.Collection
var AgentFeeLogs *mgo.Collection
var AccountingLogs *mgo.Collection

var Notices *mgo.Collection
var Shops *mgo.Collection
var Envs *mgo.Collection
var Vips *mgo.Collection
var Games *mgo.Collection

// 初始化时连接数据库
func InitMgo() {
	// get db config from host, port, username, password
	dbHost := beego.AppConfig.String("mdb.host")
	dbPort := beego.AppConfig.String("mdb.port")
	dbUser := beego.AppConfig.String("mdb.user")
	dbPassword := beego.AppConfig.String("mdb.password")
	dbName := beego.AppConfig.String("mdb.name")
	usernameAndPassword := dbUser + ":" + dbPassword + "@"
	if dbUser == "" || dbPassword == "" {
		usernameAndPassword = ""
	}
	if dbPort == "" {
		dbPort = "27017"
	}
	url := "mongodb://" + usernameAndPassword + dbHost + ":" + dbPort + "/" + dbName

	// [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
	// mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb
	var err error
	Session, err = mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	go ticker()

	// Optional. Switch the session to a monotonic behavior.
	Session.SetMode(mgo.Monotonic, true)

	// niuadmin
	Actions = Session.DB(dbName).C("t_action")
	Perms = Session.DB(dbName).C("t_perm")
	Roles = Session.DB(dbName).C("t_role")
	RolePerms = Session.DB(dbName).C("t_role_perm")
	Users = Session.DB(dbName).C("t_user")
	UserRoles = Session.DB(dbName).C("t_user_role")
	MailTpls = Session.DB(dbName).C("t_mail_tpl")
	GenIDs = Session.DB(dbName).C("t_last_id")

	// trade_record
	TradeRecords = Session.DB(dbName).C("col_trade_record")
	// user
	PlayerUsers = Session.DB(dbName).C("col_user")
	// record
	GameRecords = Session.DB(dbName).C("col_game_record")
	UserRecords = Session.DB(dbName).C("col_user_record")
	StatRecords = Session.DB(dbName).C("col_stat_record")
	UserStatRecords = Session.DB(dbName).C("col_user_stat_record")

	Pk10Records = Session.DB(dbName).C("col_pkten")

	//
	Agencys = Session.DB(dbName).C("col_agency")
	//
	RegistLogs = Session.DB(dbName).C("col_log_regist")
	LoginLogs = Session.DB(dbName).C("col_log_login")
	DiamondLogs = Session.DB(dbName).C("col_log_diamond")
	CoinLogs = Session.DB(dbName).C("col_log_coin")
	ChipLogs = Session.DB(dbName).C("col_log_chip")
	LogBuildAgencys = Session.DB(dbName).C("col_log_build_agency")
	LogOnlines = Session.DB(dbName).C("col_log_online")
	LogPayTodays = Session.DB(dbName).C("col_log_pay_today")
	LogRegistTodays = Session.DB(dbName).C("col_log_regist_today")
	LogChipTodays = Session.DB(dbName).C("col_log_chip_today")
	AgentFees = Session.DB(dbName).C("col_log_agent_fee")
	AgentFeeLogs = Session.DB(dbName).C("col_log_agent_fee_log")
	AccountingLogs = Session.DB(dbName).C("col_log_accounting_log")
	//
	ApplyCashs = Session.DB(dbName).C("col_apply_cash")
	//
	Notices = Session.DB(dbName).C("col_notice")
	Shops = Session.DB(dbName).C("col_shop")
	Envs = Session.DB(dbName).C("col_env")
	Vips = Session.DB(dbName).C("col_vip")
	Games = Session.DB(dbName).C("col_game")

	//init
	initService()

	//创建文件目录
	os.MkdirAll(GetFilePath(), 0755)
	//TODO test
	//dayStamp := utils.Stamp2Time(utils.TimestampYesterday())
	//dayStamp := TimeYesterday4()
	//statProfit(dayStamp)
	//AgencyService.statChipToday(dayStamp) //定时统计更新
	//initParentAgent()
	//initParentAgent2()
	//num := LoggerService.getRobotProfitsYesterday()
	//beego.Trace("getRobotProfitsYesterday : ", num)
}

// 创建文件目录
func GetFilePath() string {
	return fmt.Sprintf(beego.AppConfig.String("files_dir"))
}

func Close() {
	Session.Close()
	if mgoCloseCh != nil {
		close(mgoCloseCh)
	}
}

// common DAO
// 公用方法

//----------------------

func Insert(collection *mgo.Collection, i interface{}) bool {
	err := collection.Insert(i)
	return Err(err)
}

//----------------------

// 适合一条记录全部更新
func Update(collection *mgo.Collection, query interface{}, i interface{}) bool {
	err := collection.Update(query, i)
	return Err(err)
}
func Upsert(collection *mgo.Collection, query interface{}, i interface{}) bool {
	_, err := collection.Upsert(query, i)
	return Err(err)
}
func UpdateAll(collection *mgo.Collection, query interface{}, i interface{}) bool {
	_, err := collection.UpdateAll(query, i)
	return Err(err)
}
func UpdateByIdAndUserId(collection *mgo.Collection, id, userId string, i interface{}) bool {
	err := collection.Update(GetIdAndUserIdQ(id, userId), i)
	return Err(err)
}

func UpdateByIdAndUserId2(collection *mgo.Collection, id, userId bson.ObjectId, i interface{}) bool {
	err := collection.Update(GetIdAndUserIdBsonQ(id, userId), i)
	return Err(err)
}
func UpdateByIdAndUserIdField(collection *mgo.Collection, id, userId, field string, value interface{}) bool {
	return UpdateByIdAndUserId(collection, id, userId, bson.M{"$set": bson.M{field: value}})
}
func UpdateByIdAndUserIdMap(collection *mgo.Collection, id, userId string, v bson.M) bool {
	return UpdateByIdAndUserId(collection, id, userId, bson.M{"$set": v})
}

func UpdateByIdAndUserIdField2(collection *mgo.Collection, id, userId bson.ObjectId, field string, value interface{}) bool {
	return UpdateByIdAndUserId2(collection, id, userId, bson.M{"$set": bson.M{field: value}})
}
func UpdateByIdAndUserIdMap2(collection *mgo.Collection, id, userId bson.ObjectId, v bson.M) bool {
	return UpdateByIdAndUserId2(collection, id, userId, bson.M{"$set": v})
}

//
func UpdateByQField(collection *mgo.Collection, q interface{}, field string, value interface{}) bool {
	_, err := collection.UpdateAll(q, bson.M{"$set": bson.M{field: value}})
	return Err(err)
}
func UpdateByQI(collection *mgo.Collection, q interface{}, v interface{}) bool {
	_, err := collection.UpdateAll(q, bson.M{"$set": v})
	return Err(err)
}

// 查询条件和值
func UpdateByQMap(collection *mgo.Collection, q interface{}, v interface{}) bool {
	_, err := collection.UpdateAll(q, bson.M{"$set": v})
	return Err(err)
}

//------------------------

// 删除一条
func Delete(collection *mgo.Collection, q interface{}) bool {
	err := collection.Remove(q)
	return Err(err)
}
func DeleteByIdAndUserId(collection *mgo.Collection, id, userId string) bool {
	err := collection.Remove(GetIdAndUserIdQ(id, userId))
	return Err(err)
}
func DeleteByIdAndUserId2(collection *mgo.Collection, id, userId bson.ObjectId) bool {
	err := collection.Remove(GetIdAndUserIdBsonQ(id, userId))
	return Err(err)
}

// 删除所有
func DeleteAllByIdAndUserId(collection *mgo.Collection, id, userId string) bool {
	_, err := collection.RemoveAll(GetIdAndUserIdQ(id, userId))
	return Err(err)
}
func DeleteAllByIdAndUserId2(collection *mgo.Collection, id, userId bson.ObjectId) bool {
	_, err := collection.RemoveAll(GetIdAndUserIdBsonQ(id, userId))
	return Err(err)
}

func DeleteAll(collection *mgo.Collection, q interface{}) bool {
	_, err := collection.RemoveAll(q)
	return Err(err)
}

//-------------------------

func Get(collection *mgo.Collection, id string, i interface{}) {
	collection.FindId(id).One(i)
}
func Get2(collection *mgo.Collection, id bson.ObjectId, i interface{}) {
	collection.FindId(id).One(i)
}
func Get3(collection *mgo.Collection, id string, i interface{}) {
	collection.FindId(bson.ObjectIdHex(id)).One(i)
}

func GetByQ(collection *mgo.Collection, q interface{}, i interface{}) {
	collection.Find(q).One(i)
}
func ListByQ(collection *mgo.Collection, q interface{}, i interface{}) {
	collection.Find(q).All(i)
}

func ListByQLimit(collection *mgo.Collection, q interface{}, i interface{}, limit int) {
	collection.Find(q).Limit(limit).All(i)
}

// 查询某些字段, q是查询条件, fields是字段名列表
func GetByQWithFields(collection *mgo.Collection, q bson.M, fields []string, i interface{}) {
	selector := make(bson.M, len(fields))
	for _, field := range fields {
		selector[field] = true
	}
	collection.Find(q).Select(selector).One(i)
}

// 查询某些字段, q是查询条件, fields是字段名列表
func ListByQWithFields(collection *mgo.Collection, q bson.M, fields []string, i interface{}) {
	selector := make(bson.M, len(fields))
	for _, field := range fields {
		selector[field] = true
	}
	collection.Find(q).Select(selector).All(i)
}
func GetByIdAndUserId(collection *mgo.Collection, id, userId string, i interface{}) {
	collection.Find(GetIdAndUserIdQ(id, userId)).One(i)
}
func GetByIdAndUserId2(collection *mgo.Collection, id, userId bson.ObjectId, i interface{}) {
	collection.Find(GetIdAndUserIdBsonQ(id, userId)).One(i)
}

// 按field去重
func Distinct(collection *mgo.Collection, q bson.M, field string, i interface{}) {
	collection.Find(q).Distinct(field, i)
}

//----------------------

func Count(collection *mgo.Collection, q interface{}) int {
	cnt, err := collection.Find(q).Count()
	if err != nil {
		Err(err)
	}
	return cnt
}

func Has(collection *mgo.Collection, q interface{}) bool {
	if Count(collection, q) > 0 {
		return true
	}
	return false
}

//-----------------

// 得到主键和userId的复合查询条件
func GetIdAndUserIdQ(id, userId string) bson.M {
	return bson.M{"_id": bson.ObjectIdHex(id), "UserId": bson.ObjectIdHex(userId)}
}
func GetIdAndUserIdBsonQ(id, userId bson.ObjectId) bson.M {
	return bson.M{"_id": id, "UserId": userId}
}

// DB处理错误
func Err(err error) bool {
	if err != nil {
		//fmt.Println(err)
		// 删除时, 查找
		if err.Error() == "not found" {
			return true
		}
		return false
	}
	return true
}

// 检查mognodb是否lost connection
// 每个请求之前都要检查!!
func CheckMongoSessionLost() {
	fmt.Println("检查CheckMongoSessionLostErr")
	err := Session.Ping()
	if err != nil {
		fmt.Println("Lost connection to db!")
		Session.Refresh()
		err = Session.Ping()
		if err == nil {
			fmt.Println("Reconnect to db successful.")
		} else {
			fmt.Println("重连失败!!!! 警告")
		}
	}
}

//计时器
func ticker() {
	defer Close()
	mgoCloseCh = make(chan bool, 1)
	tick := time.Tick(3 * time.Minute)
	for {
		select {
		case <-mgoCloseCh:
			return
		default:
		}
		select {
		case <-tick:
			CheckMongoSessionLost()
			AgencyService.stat() //定时统计更新
			//凌晨4点统计昨天数据
			//TODO 以4点30为一天算,所以必须4点30后统计
			switch utils.Hour() {
			case 4:
				min := utils.Minute()
				if min > 30 && min < 38 {
					//dayStamp := utils.Stamp2Time(utils.TimestampYesterday())
					dayStamp := TimeYesterday4()
					AgencyService.statRegToday(dayStamp)  //统计昨日注册人数
					AgencyService.statPayToday(dayStamp)  //统计昨日充值人数
					AgencyService.statChipToday(dayStamp) //统计昨日盈亏
				} else if utils.Minute() > 40 && utils.Minute() < 48 {
					dayStamp := TimeYesterday4()
					statAccountingLog(dayStamp) //账务统计
				} else if utils.Minute() > 50 && utils.Minute() < 58 {
					//dayStamp := utils.Stamp2Time(utils.TimestampYesterday())
					dayStamp := TimeYesterday4()
					statProfit(dayStamp) //统计游戏内玩家日、周、月赢亏数据展示
				}
			}
		case <-mgoCloseCh:
			return
		}
	}
}

// 分页, 排序处理
func parsePageAndSort(pageNumber, pageSize int, sortField string, isAsc bool) (skipNum int, sortFieldR string) {
	skipNum = (pageNumber - 1) * pageSize
	if skipNum < 0 {
		skipNum = 0
	}
	if sortField == "" {
		sortField = "UpdatedTime"
	}
	if !isAsc {
		sortFieldR = "-" + sortField
	} else {
		sortFieldR = sortField
	}
	return
}

// 时间查询
func FindByDate(startDate, endDate, startField, endField string) bson.M {
	m := bson.M{}
	if startDate == "" && endDate == "" {
		return m
	}
	s := fmt.Sprintf("%s 00:00:00", startDate)
	startTime := utils.Str2Time(s)
	e := fmt.Sprintf("%s 23:59:59", endDate)
	endTime := utils.Str2Time(e)
	if !startTime.IsZero() && !endTime.IsZero() &&
		startDate != "" && endDate != "" &&
		startField == endField && startField != "" {
		m[startField] = bson.M{"$gte": startTime, "$lte": endTime}
	} else if !startTime.IsZero() && startField != "" &&
		startDate != "" {
		m[startField] = bson.M{"$gte": startTime}
	} else if !endTime.IsZero() && endField != "" &&
		endDate != "" {
		m[endField] = bson.M{"$lte": endTime}
	}
	return m
}

//筹码转换为分展示
func Chip2Float(chip int64) float64 {
	return (float64(chip) / 100)
}

//昨日凌晨4点30分
func TimestampYesterday4() int64 {
	return TimestampToday4() - 86400
}

//今日凌晨4点30分
func TimestampToday4() int64 {
	return utils.TimestampTodayTime().Unix() + 14400 + 1800
}

//以凌晨4点算昨日开始
func TimeYesterday4() time.Time {
	return TimeToday4().AddDate(0, 0, -1)
}

//以凌晨4点算今日开始
func TimeToday4() time.Time {
	now := utils.Timestamp()
	t4 := TimestampToday4()
	//0点到4点30分之间
	if now < t4 {
		//昨日凌晨4点30分
		return utils.Stamp2Time(TimestampYesterday4())
	}
	//今日凌晨4点30分
	return utils.Stamp2Time(t4)
}
