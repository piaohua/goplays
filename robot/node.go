/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-11-19 13:12:54
 * Filename      : node.go
 * Description   : 机器人
 * *******************************************************/
package main

import (
	"goplays/glog"
	"goplays/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

func (server *RobotServer) NewRemote(bind, name string) {
	if bind == "" {
		glog.Panic("bind empty")
	}
	// Start the remote server
	remote.Start(bind)
	server.remoteRecv(name) //接收远程消息
}

//接收远程消息
func (server *RobotServer) remoteRecv(name string) {
	//create the channel
	server.channel = make(chan interface{}, 100) //protos中定义

	//create an actor receiving messages and pushing them onto the channel
	props := actor.FromFunc(func(context actor.Context) {
		select {
		case <-server.stopCh:
			return
		default:
		}
		server.channel <- context.Message()
	})
	nodePid, err = actor.SpawnNamed(props, name)
	if err != nil {
		glog.Panic(err)
	}
	server.remoteHall()

	//consume the channel just like you use to
	go func() {
		for msg := range server.channel {
			server.remoteHandler(msg)
		}
		//channel closed
		server.disconnectHall()
	}()
}

//处理
func (server *RobotServer) remoteHandler(message interface{}) {
	switch message.(type) {
	case *pb.RobotMsg:
		msg := message.(*pb.RobotMsg)
		//分配机器人
		go Msg2Robots(msg, msg.Num)
		glog.Infof("node msg -> %#v", msg)
	}
}

func (server *RobotServer) remoteHall() {
	//hall
	name := cfg.Section("cookie").Key("name").Value()
	bind := cfg.Section("hall").Key("bind").Value()
	hallPid = actor.NewPID(bind, name)
	//name
	server.Name = cfg.Section("robot").Name()
	connect := &pb.Connect{
		Name: server.Name,
	}
	hallPid.Request(connect, nodePid)
}

func (server *RobotServer) disconnectHall() {
	//hall
	disconnect := &pb.Disconnect{
		Name: server.Name,
	}
	if hallPid != nil {
		hallPid.Tell(disconnect)
	}
	if nodePid != nil {
		nodePid.Stop()
	}
}
