/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2018-01-27 11:33:32
 * Filename      : node.go
 * Description   : 连接
 * *******************************************************/
package node

import (
	"time"

	"goplays/glog"
	"goplays/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

type Node interface {
	Receive(interface{})
}

type NodeServer struct {
	Name     string
	HallBind string
	HallName string
	HallPID  *actor.PID
	NodePID  *actor.PID
	msgCh    chan interface{} //消息通道
}

func NewNodeServer(serverName, HallBind, HallName,
	NodeBind, NodeName string) *NodeServer {
	server := new(NodeServer)
	server.Name = serverName
	server.HallBind = HallBind
	server.HallName = HallName
	server.NewRemote(NodeBind, NodeName)
	return server
}

func (server *NodeServer) NewRemote(bind, name string) {
	if bind == "" {
		glog.Panic("bind empty")
	}
	// Start the remote server
	remote.Start(bind)
	server.remoteRecv(name) //接收远程消息
}

//接收远程消息
func (server *NodeServer) remoteRecv(name string) {
	//create the channel
	server.msgCh = make(chan interface{}, 20) //protos中定义

	//create an actor receiving messages and pushing them onto the channel
	props := actor.FromFunc(func(context actor.Context) {
		msg := context.Message()
		server.msgCh <- msg
	})
	var err1 error
	server.NodePID, err1 = actor.SpawnNamed(props, name)
	if err1 != nil {
		glog.Panic(err1)
	}
	server.remoteHall()

	//consume the channel just like you use to
	go func() {
		for msg := range server.msgCh {
			server.Receive(msg)
		}
		//channel closed
		server.disconnectHall()
	}()
}

func (server *NodeServer) remoteHall() {
	//hall
	server.HallPID = actor.NewPID(server.HallBind, server.HallName)
	//name
	connect := &pb.Connect{
		Name: server.Name,
	}
	server.Request(connect)
}

func (server *NodeServer) disconnectHall() {
	//hall
	disconnect := &pb.Disconnect{
		Name: server.Name,
	}
	if server.HallPID != nil {
		server.HallPID.Tell(disconnect)
	}
	if server.NodePID != nil {
		server.NodePID.Stop()
	}
}

func (server *NodeServer) Send(msg interface{}) {
	if server.HallPID == nil || server.NodePID == nil {
		return
	}
	server.HallPID.Request(msg, server.NodePID)
}

func (server *NodeServer) Request(msg interface{}) interface{} {
	if server.HallPID == nil || server.NodePID == nil {
		return nil
	}
	timeout := 3 * time.Second
	res1, err1 := server.NodePID.RequestFuture(msg, timeout).Result()
	if err1 != nil {
		glog.Errorf("Request err: %v", err1)
		return nil
	}
	return res1
}

func (server *NodeServer) Receive(msg interface{}) {
	switch msg.(type) {
	case *pb.ServeClose:
		server.NodePID.Stop()
		server.HallPID = nil
		server.NodePID = nil
	default:
		glog.Errorf("unknow msg -> %v", msg)
	}
}
