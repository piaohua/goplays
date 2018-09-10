package main

import (
	"time"

	"goplays/data"
	"goplays/game/handler"
	"goplays/glog"
	"goplays/pb"
	"utils"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func (a *DeskActor) HandlerMsg(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.Connected:
		//连接成功
		arg := msg.(*pb.Connected)
		glog.Infof("Connected %s", arg.Name)
	case *pb.Disconnected:
		//成功断开
		arg := msg.(*pb.Disconnected)
		glog.Infof("Disconnected %s", arg.Name)
	case *pb.CloseDesk:
		arg := msg.(*pb.CloseDesk)
		glog.Debugf("CloseDesk %#v", arg)
		//移除
		delete(a.desks, arg.Roomid)
		delete(a.count, arg.Roomid)
		//响应
		//rsp := new(pb.ClosedDesk)
		//ctx.Respond(rsp)
	case *pb.LeaveDesk:
		arg := msg.(*pb.LeaveDesk)
		glog.Debugf("LeaveDesk %#v", arg)
		//移除
		delete(a.roles, arg.Userid)
		if n, ok := a.count[arg.Roomid]; ok && n > 0 {
			a.count[arg.Roomid] = n - 1
		}
		if arg.Type == 1 {
			a.roomPid.Request(arg, ctx.Self())
			a.hallPid.Request(arg, ctx.Self())
		}
		//响应
		//rsp := new(pb.LeftDesk)
		//ctx.Respond(rsp)
	case *pb.JoinDesk:
		arg := msg.(*pb.JoinDesk)
		glog.Debugf("JoinDesk %#v", arg)
		//房间数据变更
		if _, ok := a.roles[arg.Userid]; !ok {
			a.count[arg.Roomid] += 1
		}
		a.roles[arg.Userid] = arg.Sender
		a.roomPid.Request(arg, ctx.Self())
		a.hallPid.Request(arg, ctx.Self())
		//响应
		//rsp := new(pb.EnteredRoom)
		//ctx.Respond(rsp)
	case *pb.SyncConfig:
		//同步配置
		arg := msg.(*pb.SyncConfig)
		glog.Debugf("SyncConfig %#v", arg)
		a.syncDesk(arg, ctx)
	case *pb.ChangeCurrency:
		arg := msg.(*pb.ChangeCurrency)
		//离开房间更新
		a.changeCurrency(arg)
	case *pb.OfflineCurrency:
		arg := msg.(*pb.OfflineCurrency)
		msg2 := &pb.ChangeCurrency{
			Userid:  arg.Userid,
			Type:    arg.Type,
			Coin:    arg.Coin,
			Diamond: arg.Diamond,
			Chip:    arg.Chip,
			Card:    arg.Card,
		}
		if v, ok := a.roles[arg.Userid]; ok {
			glog.Infof("OfflineCurrency %#v", arg)
			v.Tell(msg2)
		} else {
			glog.Infof("OfflineCurrency %#v", arg)
			a.hallPid.Tell(msg2)
		}
	case *pb.HuiYinOpenTime:
		arg := msg.(*pb.HuiYinOpenTime)
		glog.Debugf("HuiYinOpenTime %#v", arg)
		rsp := a.leftOpen()
		arg.Sender.Tell(rsp)
	case *pb.Pk10TrendLog:
		a.dbmsPid.Tell(msg)
	case *pb.Pk10UseridLog:
		a.dbmsPid.Tell(msg)
	case *pb.Pk10GameLog:
		a.dbmsPid.Tell(msg)
	case *pb.LogChip:
		a.dbmsPid.Tell(msg)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

//在这里触发创建房间
func (a *DeskActor) syncDesk(arg *pb.SyncConfig, ctx actor.Context) {
	switch arg.Type {
	case pb.CONFIG_GAMES:
		b := make(map[string]data.Game)
		err = json.Unmarshal(arg.Data, &b)
		if err != nil {
			glog.Errorf("syncConfig Unmarshal err %v", err)
			return
		}
		for _, v := range b {
			//不是当前节点房间
			if v.Node != nodeName {
				continue
			}
			switch arg.Atype {
			case pb.CONFIG_DELETE:
				//关闭房间
				v2 := v
				a.closeDesk(&v2, ctx)
			case pb.CONFIG_UPSERT:
				//创建房间
				v2 := v
				//已经存在
				if id, ok := a.rules[v2.Id]; ok {
					if p, ok2 := a.desks[id]; ok2 {
						p.Tell(arg)
						continue
					}
				}
				a.spawnDesk(&v2, ctx)
			}
		}
	default:
	}
	handler.SyncConfig(arg)
}

//关闭一张桌子
func (a *DeskActor) closeDesk(gameData *data.Game, ctx actor.Context) {
	glog.Debugf("close Desk %#v", gameData)
	//可以去room服务中取
	if k, ok := a.rules[gameData.Id]; ok {
		glog.Debugf("close Desk %s", k)
		if v, ok := a.desks[k]; ok {
			glog.Debugf("close Desk %s", v.String())
			msg1 := new(pb.ServeStop)
			v.Request(msg1, ctx.Self())
			//停掉服务
			msg2 := new(pb.CloseDesk)
			msg2.Roomid = k
			//TODO 添加类型,桌子中关闭
			a.roomPid.Request(msg2, ctx.Self())
			a.hallPid.Request(msg2, ctx.Self())
			v.Stop()
			delete(a.desks, k)
		}
		delete(a.rules, gameData.Id)
	}
}

//启动新服务,新开的房间同步状态
func (a *DeskActor) spawnDesk(gameData *data.Game, ctx actor.Context) {
	glog.Debugf("spawn Desk %#v", gameData)
	deskData := handler.NewDeskData(gameData)
	msg1 := new(pb.GenDesk)
	msg1.Rtype = deskData.Rtype
	msg1.Gtype = deskData.Gtype
	res1 := a.reqRoom(msg1, ctx)
	var response1 *pb.GenedDesk
	var ok bool
	if response1, ok = res1.(*pb.GenedDesk); !ok {
		glog.Error("spawn desk failed")
		return
	}
	glog.Debugf("response1: %#v", response1)
	deskData.Rid = response1.Roomid
	//新桌子
	newDesk1 := NewDesk(deskData)
	//spawn desk
	deskPid := newDesk1.newDesk()
	//添加新桌子
	a.desks[deskData.Rid] = deskPid
	//规则暂时只闭时用到
	a.rules[deskData.Unique] = deskData.Rid
	glog.Debugf("deskPid: %#v", deskPid.String())
	//启动
	deskPid.Tell(new(pb.ServeStart))
	//添加桌子
	msg2 := new(pb.AddDesk)
	msg2.Desk = deskPid
	msg2.Roomid = deskData.Rid
	msg2.Rtype = deskData.Rtype
	msg2.Gtype = deskData.Gtype
	msg2.Data = handler.Desk2Data(deskData) //打包
	if len(msg2.Data) == 0 {
		glog.Error("add desk failed")
		a.closeDesk(gameData, ctx)
		return
	}
	res2 := a.reqRoom(msg2, ctx)
	var response2 *pb.AddedDesk
	if response2, ok = res2.(*pb.AddedDesk); !ok {
		glog.Error("spawn desk failed")
		a.closeDesk(gameData, ctx)
		return
	}
	if response2.Error != pb.OK {
		glog.Error("spawn desk failed")
		a.closeDesk(gameData, ctx)
		return
	}
	glog.Debugf("spawn Desk successfully %s, %s", deskData.Rid, deskPid.String())
	a.hallPid.Request(msg2, ctx.Self())
}

//登录成功数据处理
func (a *DeskActor) reqRoom(msg interface{}, ctx actor.Context) interface{} {
	timeout := 3 * time.Second
	res1, err1 := a.roomPid.RequestFuture(msg, timeout).Result()
	if err1 != nil {
		glog.Errorf("reqRoom err: %v, msg %#v", err1, msg)
		return nil
	}
	return res1
}

//更新货币
func (a *DeskActor) changeCurrency(arg *pb.ChangeCurrency) {
	if v, ok := a.roles[arg.Userid]; ok {
		v.Tell(arg)
	}
}

//广播更新状态
func (a *DeskActor) pushDeskState() {
	msg1 := new(pb.PushDeskState)
	msg1.State = a.state
	msg1.Expect = a.expect
	msg1.Opencode = a.opencode
	msg1.Opentime = a.opentime
	msg1.Opentimestamp = a.opentimestamp
	msg1.Nexttime = a.nexttime
	//封盘重启设置上期期号
	if a.state == data.STATE_SEAL && a.expect == "" {
		msg1.Expect = a.lastexpect
		msg1.Opencode = a.lastopencode
	}
	a.broadcast(msg1)
}

//广播消息
func (a *DeskActor) broadcast(msg interface{}) {
	for _, v := range a.desks {
		v.Tell(msg)
	}
}

//开奖时间
func (a *DeskActor) leftOpen() (msg *pb.HuiYinOpenedTime) {
	now := utils.Timestamp()
	st := utils.Time2Stamp(a.startTime)
	et := utils.Time2Stamp(a.endTime)
	nextopentime, left := leftCount(now, st, et)
	msg = new(pb.HuiYinOpenedTime)
	msg2 := &pb.HuiYinGame{
		Gtype: data.GAME_JIU,
		State: a.state,
		Timer: nextopentime,
		Left:  left,
	}
	msg3 := &pb.HuiYinGame{
		Gtype: data.GAME_NIU,
		State: a.state,
		Timer: nextopentime,
		Left:  left,
	}
	msg4 := &pb.HuiYinGame{
		Gtype: data.GAME_SAN,
		State: a.state,
		Timer: nextopentime,
		Left:  left,
	}
	msg.List = append(msg.List, msg2)
	msg.List = append(msg.List, msg3)
	msg.List = append(msg.List, msg4)
	return
}
