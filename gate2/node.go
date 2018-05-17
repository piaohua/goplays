package main

import (
	"time"

	"goplays/glog"
	"goplays/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/mailbox"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/gogo/protobuf/proto"
)

var (
	nodePid *actor.PID
)

//网关服务
type GateActor struct {
	Name string
	//大厅服务
	hallPid *actor.PID
	//数据中心服务
	dbmsPid *actor.PID
	//房间服务
	roomPid *actor.PID
	//角色服务
	rolePid *actor.PID
	//节点角色进程 userid - pid
	roles map[string]*actor.PID
	//节点离线角色进程 userid - pid
	offline map[string]*actor.PID
	//节点离线角色进程关闭时间
	offtime map[string]int
	//关闭通道
	stopCh chan struct{}
	//更新状态
	status bool
	//计时
	timer int
}

func (a *GateActor) Receive(ctx actor.Context) {
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

func newGateActor() actor.Producer {
	return func() actor.Actor {
		a := new(GateActor)
		a.Name = cfg.Section(nodeName).Name()
		glog.Debugf("Name %s", a.Name)
		//roles key=userid
		a.roles = make(map[string]*actor.PID)
		a.offline = make(map[string]*actor.PID)
		a.offtime = make(map[string]int)
		a.stopCh = make(chan struct{})
		return a
	}
}

func NewRemote(bind, name string) {
	remote.Start(bind)
	props := actor.
		FromProducer(newGateActor()).
		WithMailbox(mailbox.Bounded(20000))
	remote.Register(name, props)
	nodePid, err = actor.SpawnNamed(props, name)
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
			glog.Errorf("nodePid Stop err: %v", err1)
		}
		response1 := res1.(*pb.ServeStoped)
		glog.Debugf("response1: %#v", response1)
		nodePid.Stop()
	}
}
