/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-11-19 11:32:23
 * Filename      : robots.go
 * Description   : 机器人
 * *******************************************************/
package main

import (
	"time"

	"goplays/glog"
	"goplays/pb"
	"utils"
)

//消息通知
func Msg2Robots(msg interface{}, num uint32) {
	for num > 0 {
		rbs.Send2rbs(msg)
		num--
	}
}

//登录成功
func RegistRoom(roomid string, ltype uint32) {
	msg := &pb.RobotRoomList{
		Roomid: roomid,
		Ltype:  ltype,
	}
	glog.Debugf("regist room %s ltype %d", roomid, ltype)
	Msg2Robots(msg, 1)
}

//登录成功
func EnterRoom(phone, roomid string, ltype uint32) {
	msg := &pb.RobotEnterRoom{
		Phone:  phone,
		Roomid: roomid,
		Ltype:  ltype,
	}
	Msg2Robots(msg, 1)
}

//登录成功
func Logined(phone string, ltype uint32) {
	msg := &pb.RobotLogin{
		Phone: phone,
		Ltype: ltype,
	}
	Msg2Robots(msg, 1)
}

//登出成功
func Logout(roomid, phone, code string, chip int64) {
	msg := &pb.RobotLogout{
		Roomid: roomid,
		Phone:  phone,
		Code:   code,
		Chip:   chip,
	}
	Msg2Robots(msg, 1)
}

// 已经注册,重新登录
func ReLogined(roomid, phone, code string, rtype uint32, envBet int32) {
	msg := &pb.RobotReLogin{
		Roomid: roomid,
		Phone:  phone,
		Code:   code,
		Rtype:  rtype,
		EnvBet: envBet,
	}
	Msg2Robots(msg, 1)
}

//发送消息
func (r *RobotServer) Send2rbs(msg interface{}) {
	if r.msgCh == nil {
		glog.Errorf("server msg channel closed %#v", msg)
		return
	}
	if len(r.msgCh) == cap(r.msgCh) {
		//FIXME send msg channel full -> 100
		glog.Errorf("send msg channel full -> %d", len(r.msgCh))
		return
	}
	select {
	case <-r.stopCh:
		return
	default:
	}
	select {
	case <-r.stopCh:
		return
	default:
		r.msgCh <- msg
	}
}

//开始前5分钟
func pkBetTime5() bool {
	year, month, day := utils.DateTime()
	startTime := utils.DateLocal(year, month, day, 9, 2, 0, 0)
	endTime := utils.DateLocal(year, month, day, 9, 9, 0, 0)
	now := utils.Timestamp()
	st := utils.Time2Stamp(startTime)
	et := utils.Time2Stamp(endTime)
	//TODO
	//startTime := utils.Str2Time("2018-02-22 09:02:00")
	//endTime := utils.Str2Time("2018-02-22 23:57:00")
	if now > st && now < et {
		return true
	}
	return false
}

//开始前5分钟
func ftBetTime5() bool {
	year, month, day := utils.DateTime()
	startTime := utils.DateLocal(year, month, day, 13, 4, 0, 0)
	endTime := utils.DateLocal(year, month, day, 13, 9, 0, 0)
	now := utils.Timestamp()
	st := utils.Time2Stamp(startTime)
	et := utils.Time2Stamp(endTime)
	if now > st && now < et {
		return true
	}
	return false
}

//一、虚假人数：界面上显示的房间人数并非真实数字，根据《北京赛车pk10》（或幸运飞艇）的时间有所不同
//*每天刚开始的前五分钟内，显示人数=虚假人数+真实人数
//*虚假人数初始值=2，之后每40s随机增加2-10
//机器人测试
func (r *RobotServer) runPkFake() {
	glog.Infof("runPkFake started phone -> %s", r.phone)
	tick := time.Tick(40 * time.Second)
	//lottery type 1 赛车, 2 飞艇
	for {
		select {
		case <-tick:
			if hallPid == nil {
				continue
			}
			if pkBetTime5() {
				msg4 := &pb.RobotFake{
					Ltype: 1,
					Type:  1,
				}
				hallPid.Tell(msg4)
			}
			if ftBetTime5() {
				msg4 := &pb.RobotFake{
					Ltype: 2,
					Type:  1,
				}
				hallPid.Tell(msg4)
			}
		case <-r.stopCh:
			return
		}
	}
}

//*从第5分钟开始，虚假人数的值每10分钟判断一次是否  真实人数≥40
//*是，则将虚假人数的取值重置为0
//*否。则将虚假人数的重新随机取值，取值范围为30-50
//机器人测试
func (r *RobotServer) runFtFake() {
	glog.Infof("runFtFake started phone -> %s", r.phone)
	tick := time.Tick(10 * time.Minute)
	for {
		select {
		case <-tick:
			if hallPid == nil {
				continue
			}
			if pkBetTime() {
				msg4 := &pb.RobotFake{
					Ltype: 1,
					Type:  2,
				}
				hallPid.Tell(msg4)
			}
			if ftBetTime() {
				msg4 := &pb.RobotFake{
					Ltype: 2,
					Type:  2,
				}
				hallPid.Tell(msg4)
			}
		case <-r.stopCh:
			return
		}
	}
}

//机器人测试
func (r *RobotServer) runPkTest() {
	glog.Infof("runPkTest started phone -> %s", r.phone)
	tick := time.Tick(3 * time.Minute)
	//lottery type 1 赛车, 2 飞艇
	msg4 := &pb.RobotMsg{
		Ltype: 1,
	}
	for {
		select {
		case <-tick:
			if !pkBetTime() {
				//Msg2Robots(new(pb.RobotStop), 1)
				continue
			}
			glog.Infof("r.online -> %d\n", len(r.online))
			glog.Infof("r.offline -> %d\n", len(r.offline))
			glog.Infof("r.phone -> %s\n", r.phone)
			//TODO:优化,按时间段运行
			//运行指定数量机器人(每个创建一个牌局)
			//code = "create" 表示机器人创建房间
			//go Msg2Robots(msg1, 5)
			if len(r.online) < 30 {
				go Msg2Robots(msg4, 5)
			}
		case <-r.stopCh:
			return
		}
	}
}

//下注时间
func pkBetTime() bool {
	year, month, day := utils.DateTime()
	startTime := utils.DateLocal(year, month, day, 9, 2, 72, 0)
	endTime := utils.DateLocal(year, month, day, 23, 57, 0, 0)
	now := utils.Timestamp()
	st := utils.Time2Stamp(startTime)
	et := utils.Time2Stamp(endTime)
	//TODO
	//startTime := utils.Str2Time("2018-02-22 09:02:00")
	//endTime := utils.Str2Time("2018-02-22 23:57:00")
	if now > st && now < et {
		return true
	}
	return false
}

//机器人测试
func (r *RobotServer) runFtTest() {
	glog.Infof("runFtTest started phone -> %s", r.phone)
	tick := time.Tick(3 * time.Minute)
	msg4 := &pb.RobotMsg{
		Ltype: 2,
	}
	for {
		select {
		case <-tick:
			if !ftBetTime() {
				//Msg2Robots(new(pb.RobotStop), 1)
				continue
			}
			glog.Infof("r.online -> %d\n", len(r.online))
			glog.Infof("r.offline -> %d\n", len(r.offline))
			glog.Infof("r.phone -> %s\n", r.phone)
			//TODO:优化,按时间段运行
			//运行指定数量机器人(每个创建一个牌局)
			//code = "create" 表示机器人创建房间
			//go Msg2Robots(msg1, 5)
			if len(r.online) < 30 {
				go Msg2Robots(msg4, 5)
			}
		case <-r.stopCh:
			return
		}
	}
}

//下注时间
func ftBetTime() bool {
	year, month, day := utils.DateTime()
	startTime := utils.DateLocal(year, month, day, 4, 4, 0, 0)
	endTime := utils.DateLocal(year, month, day, 13, 4, 72, 0)
	now := utils.Timestamp()
	st := utils.Time2Stamp(startTime)
	et := utils.Time2Stamp(endTime)
	if now > st && now < et {
		return false
	}
	return true
}

//处理
func (r *RobotServer) run() {
	defer func() {
		glog.Infof("Robots closed online -> %d\n", len(r.online))
		glog.Infof("Robots closed offline -> %d\n", len(r.offline))
		glog.Infof("Robots closed phone -> %s\n", r.phone)
	}()
	glog.Infof("Robots started -> %s", r.phone)
	tick := time.Tick(time.Minute)
	for {
		select {
		case m, ok := <-r.msgCh:
			if !ok {
				glog.Errorf("Robots msgCh closed phone -> %s\n", r.phone)
				return
			}
			switch m.(type) {
			case *pb.RobotMsg:
				//启动机器人
				msg := m.(*pb.RobotMsg)
				glog.Infof("run msg -> %v", msg)
				r.run2(msg)
			case *pb.RobotReLogin:
				//重新尝试登录
				msg := m.(*pb.RobotReLogin)
				glog.Infof("ReLogin -> %#v", msg)
				go r.RunRobot(msg.Roomid, msg.Phone, msg.Code, msg.Rtype, msg.Ltype, msg.EnvBet, false)
			case *pb.RobotLogin:
				//登录成功
				msg := m.(*pb.RobotLogin)
				glog.Infof("login -> %#v", msg)
				delete(r.offline, msg.Phone)
				r.online[msg.Phone] = true
			case *pb.RobotLogout:
				//登出断开
				msg := m.(*pb.RobotLogout)
				glog.Infof("logout -> %#v", msg)
				if _, ok := r.online[msg.Phone]; ok {
					delete(r.online, msg.Phone)
				}
				if msg.Chip < 20000 {
					//TODO 自动充值
					//r.unused[msg.Phone] = msg.Chip
				} else {
					//TODO 暂时不重复
					//r.offline[msg.Phone] = true
				}
				if v, ok := r.rooms[msg.Roomid]; ok && v > 0 {
					r.rooms[msg.Roomid]--
				}
			case *pb.RobotStop:
				msg := m.(*pb.RobotStop)
				glog.Infof("robot stop %v", msg)
				//r.mutexConns.Lock()
				//for conn := range r.conns {
				//	conn.Close()
				//}
				//r.conns = nil
				//r.mutexConns.Unlock()
			case *pb.RobotRoomList:
				msg := m.(*pb.RobotRoomList)
				if v, ok := r.ltypes[msg.Ltype]; ok {
					var have bool
					for _, val := range v {
						if val == msg.Roomid {
							have = true
							break
						}
					}
					if !have {
						v = append(v, msg.Roomid)
						r.ltypes[msg.Ltype] = v
					}
				} else {
					v := make([]string, 0)
					v = append(v, msg.Roomid)
					r.ltypes[msg.Ltype] = v
				}
				glog.Debugf("room list %s ltypes %#v", msg.Roomid, r.ltypes)
			case *pb.RobotEnterRoom:
				msg := m.(*pb.RobotEnterRoom)
				glog.Debugf("RobotEnterRoom -> %#v", msg)
				r.rooms[msg.Roomid]++
				glog.Debugf("rooms -> %#v", r.rooms)
			case closeFlag:
				//停止发送消息
				close(r.stopCh)
				return
			}
		case <-tick:
			//逻辑处理
		}
	}
}

//启动机器人
func (r *RobotServer) run2(msg *pb.RobotMsg) {
	var code string = msg.Code
	var rtype uint32 = msg.Rtype
	var ltype uint32 = msg.Ltype
	var envBet int32 = msg.EnvBet
	var phone string
	//选择一个房间
	glog.Debugf("ltypes %#v", r.ltypes)
	var roomid string
	if len(msg.Roomid) != 0 {
		roomid = msg.Roomid
	} else {
		if s, ok := r.ltypes[ltype]; ok {
			for _, v := range s {
				//每个房间5个人
				if r.rooms[v] < 8 {
					roomid = v
					break
				}
			}
		}
	}
	//房间已经存在列表
	if roomid == "" && len(r.ltypes[ltype]) != 0 {
		return
	}
	glog.Debugf("roomid %s, ltype %d", roomid, ltype)
	for k, v := range r.offline {
		if v { //已经断开
			phone = k
			r.offline[k] = false //登录中
			glog.Infof("run offline robot -> %s, %s", roomid, phone)
			go r.RunRobot(roomid, phone, code, rtype, ltype, envBet, false)
			break
		}
	}
	glog.Infof("offline phone -> %s", phone)
	if len(phone) == 0 {
		phone = r.phone
		r.phone = utils.StringAdd(r.phone)
		if _, ok := r.unused[phone]; ok {
			//TODO 自动充值
			//Msg2Robots(msg, 1)
			//TODO 会出现死循环
		} else if _, ok := r.online[phone]; ok {
			//重复
			Msg2Robots(msg, 1)
		} else {
			//新机器人不用注册
			glog.Infof("run new robot -> %s, %s", roomid, phone)
			go r.RunRobot(roomid, phone, code, rtype, ltype, envBet, false)
		}
	}
	glog.Infof("new phone -> %s", phone)
	//重置
	phone1 := cfg.Section("robot").Key("phone").Value()
	if r.phone > utils.StringAdd2(phone1, "750") {
		r.phone = phone1
	}
}
