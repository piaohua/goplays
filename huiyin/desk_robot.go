package main

import (
	"goplays/data"
	"goplays/game/config"
	"goplays/glog"
	"goplays/pb"
)

//r真实人数,f机器人
func (a *Desk) realRoles() (f, r uint32) {
	for _, v := range a.players {
		if v.GetRobot() {
			f++
		} else {
			r++
		}
	}
	return
}

//是否有真实玩家下注
func (a *Desk) realRoleBet() bool {
	for k, v := range a.players {
		if v.GetRobot() {
		} else {
			if _, ok := a.HuiYinDeskData.roleBets[k]; ok {
				return true
			}
		}
	}
	return false
}

//var num int32 = config.GetEnv(data.ENV12)
//var allot1 int32 = config.GetEnv(data.ENV13)
//var allot2 int32 = config.GetEnv(data.ENV14)
//var bet int32 = config.GetEnv(data.ENV15)
//
//robot
//bind := cfg.Section("robot").Key("bind").Value()
//name := cfg.Section("cookie").Key("name").Value()
//robotPid := actor.NewPID(bind, name)

func (a *DeskActor) robotAllot(now int64) {
	if a.state != data.STATE_BET {
		return
	}
	if (a.nexttime - now) == 200 {
		var bet int32 = config.GetEnv(data.ENV15)
		var allot1 int32 = config.GetEnv(data.ENV13)
		if allot1 > 0 {
			msg1 := &pb.RobotAllot{
				Type:    1,
				EnvBet:  bet,
				HallPid: a.hallPid,
			}
			a.broadcast(msg1)
		}
		var allot2 int32 = config.GetEnv(data.ENV14)
		glog.Debugf("allot1 %d, allot2 %d\n", allot1, allot2)
		if allot2 > 0 {
			msg1 := &pb.RobotAllot{
				Type:    2,
				EnvBet:  bet,
				HallPid: a.hallPid,
			}
			a.broadcast(msg1)
		}
		//TODO 暂时用这个协议控制机器人下注
		msg1 := &pb.SPing{
			Time: 200,
		}
		a.broadcast(msg1)
	} else if (a.nexttime - now) == 10 {
		//最后10秒
		//TODO 暂时用这个协议控制机器人下注
		msg1 := &pb.SPing{
			Time: 10,
		}
		a.broadcast(msg1)
	}
}
