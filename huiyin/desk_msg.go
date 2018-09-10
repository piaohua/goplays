package main

import (
	"time"

	"goplays/data"
	"goplays/glog"
	"goplays/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func (a *Desk) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.ServeStop:
		//关闭服务
		glog.Infof("ServeStop: %#v", msg)
		a.handlerStop(ctx)
		//响应登录
		//rsp := new(pb.ServeStoped)
		//ctx.Respond(rsp)
	case *pb.ServeStart:
		a.start(ctx)
		//响应
		//rsp := new(pb.ServeStarted)
		//ctx.Respond(rsp)
	case *pb.Tick:
		a.ding(ctx)
	default:
		//glog.Errorf("unknown message %v", msg)
		a.HandlerLogic(msg, ctx)
	}
}

//启动服务
func (a *Desk) start(ctx actor.Context) {
	glog.Infof("desk start: %v", ctx.Self().String())
	//启动
	//go a.ticker(ctx)
}

//时钟
func (a *Desk) ticker(ctx actor.Context) {
	tick := time.Tick(time.Second)
	msg := new(pb.Tick)
	for {
		select {
		case <-a.stopCh:
			glog.Info("desk ticker closed")
			return
		default: //防止阻塞
		}
		select {
		case <-a.stopCh:
			glog.Info("desk ticker closed")
			return
		case <-tick:
			ctx.Self().Tell(msg)
		}
	}
}

//钟声
func (a *Desk) ding(ctx actor.Context) {
	//glog.Debugf("ding: %v", ctx.Self().String())
	//逻辑处理
	//a.tickerHandler()
	//TODO
}

//钟声
func (a *Desk) tickerHandler(ctx actor.Context) {
	//逻辑处理
	//t.timer += 1
	//TODO
}

//关闭时钟
func (a *Desk) closeTick() {
	select {
	case <-a.stopCh:
		return
	default:
		//停止发送消息
		close(a.stopCh)
	}
	//离开消息
	for k, p := range a.pids {
		msg2 := &pb.LeaveDesk{
			Roomid: a.id,
			Userid: k,
		}
		nodePid.Tell(msg2)
		if p == nil {
			continue
		}
		p.Tell(msg2)
	}
	//逻辑处理
	//a.close()
	//TODO
}

func (a *Desk) handlerStop(ctx actor.Context) {
	glog.Infof("handlerStop: %s", a.Name)
	//关闭
	a.closeTick()
	//断开处理
	//TODO 玩家结算退出
	if a.state == data.STATE_BET {
		a.state = data.STATE_OVER
		//a.gameOver()
		//直接退款
	}
	//直接退款
	a.state = data.STATE_OVER
	a.closeRefund()
}
