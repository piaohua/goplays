package main

import (
	"time"

	"goplays/glog"
	"goplays/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/gogo/protobuf/proto"
)

var (
	nodePid *actor.PID
)

//桌子服务
type DeskActor struct {
	Name string
	//大厅服务
	hallPid *actor.PID
	//数据中心服务
	dbmsPid *actor.PID
	//房间服务
	roomPid *actor.PID
	//节点角色进程
	roles map[string]*actor.PID
	//所有桌子roomid-deskPid
	desks map[string]*actor.PID
	//房间人数roomid-numbers
	count map[string]uint32
	//配置映射unique-roomid
	rules map[string]string
	//关闭通道
	stopCh chan struct{}
	//更新状态
	status bool
	//计时
	timer int64
	//房间状态,0准备中,1游戏中,2封盘,3结算
	state uint32
	//下一个状态时间点
	nexttime int64
	//当天开始结束时间
	startTime time.Time
	endTime   time.Time
	//获取到的开奖结果
	expect        string
	opencode      string
	opentime      string
	opentimestamp int64
	//上一期
	lastexpect   string
	lastopencode string
	//code 种类
	code string
	//开奖时间间隔
	interval int64
	//封盘时间
	sealTime int64
	//抓取结果时间,延长封盘时间
	grabTime int64
	//结算时间
	overTime int64
}

func (a *DeskActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *pb.Request:
		ctx.Respond(&pb.Response{})
	case *actor.Started:
		glog.Notice("Starting, initialize actor here")
	case *actor.Stopping:
		glog.Notice("Stopping, actor is about to shut down")
	case *actor.Stopped:
		glog.Notice("Stopped, actor and its children are stopped")
	case *actor.Restarting:
		glog.Notice("Restarting, actor is about to restart")
	case *actor.ReceiveTimeout:
		glog.Infof("ReceiveTimeout: %v", ctx.Self().String())
	case proto.Message:
		a.Handler(msg, ctx)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

func newDeskActor() actor.Actor {
	a := new(DeskActor)
	a.Name = cfg.Section(nodeName).Name()
	//roles key=userid
	a.roles = make(map[string]*actor.PID)
	a.desks = make(map[string]*actor.PID)
	a.count = make(map[string]uint32)
	a.rules = make(map[string]string)
	a.stopCh = make(chan struct{})
	a.code = cfg.Section(nodeName).Key("code").Value()
	//
	a.interval = 300 //开奖时间间隔
	a.sealTime = 42  //封盘时间
	a.grabTime = 108 //150 最长抓取时间 60 + 48
	a.overTime = 30  //180 结算动画时间 60 剩余下注时间 120
	return a
}

func NewRemote(bind, name string) {
	remote.Start(bind)
	huiyinProps := actor.FromProducer(newDeskActor)
	remote.Register(name, huiyinProps)
	nodePid, err = actor.SpawnNamed(huiyinProps, name)
	if err != nil {
		glog.Fatalf("nodePid err %v", err)
	}
	glog.Infof("nodePid %s", nodePid.String())
	nodePid.Tell(new(pb.ServeStart))
}

//关闭
func Stop() {
	timeout := 10 * time.Second
	msg := new(pb.ServeStop)
	if nodePid != nil {
		res1, err1 := nodePid.RequestFuture(msg, timeout).Result()
		if err1 != nil {
			//TODO future: timeout
			glog.Errorf("nodePid Stop err: %v", err1)
		}
		if response1, ok := res1.(*pb.ServeStoped); ok {
			glog.Debugf("response1: %#v", response1)
		}
		nodePid.Stop()
	}
}
