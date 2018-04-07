package main

import (
	"fmt"
	"testing"
	"time"

	"goplays/pb"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

/*
func TestNode(t *testing.T) {
	remote.Start("127.0.0.1:0")
	activate("hello")
	<-time.After(time.Minute)
	activate("hello")
	console.ReadLine()
}

func activate(name string) {
	timeout := 1 * time.Second
	pid, err := remote.SpawnNamed("127.0.0.1:7002", "remote1", name, timeout)
	if err != nil {
		//fmt.Println(err)
		//return
	}
	res, _ := pid.RequestFuture(new(pb.Request), timeout).Result()
	fmt.Println("res ", res)
	response := res.(*pb.Response)
	fmt.Println(response)
	//pid.Stop()
	//
	//pid, _ = remote.SpawnNamed("127.0.0.1:8080", "remote2", name, timeout)
	res, _ = pid.RequestFuture(new(pb.Request), timeout).Result()
	response = res.(*pb.Response)
	fmt.Println(response)
}
*/

func TestHall(t *testing.T) {
	remote.Start("127.0.0.1:8080")
	activate2()
	<-time.After(time.Minute)
	console.ReadLine()
}

func activate2() {
	name := "V2jyQR6wbNbQ+PA2VOxbcouZG3lBtElenMK2EtOvlrE"
	bind := "127.0.0.1:7012"
	timeout := 3 * time.Second
	//pid, _ := remote.SpawnNamed(bind, "remote", name, timeout)
	//pid := pidResp.Pid
	pid := actor.NewPID(bind, name)
	msg1 := new(pb.LoginHall)
	msg1.Sender = pid
	msg1.Userid = "1"
	res1, err1 := pid.RequestFuture(msg1, timeout).Result()
	if err1 != nil {
		fmt.Printf("LoginHall err: %v", err1)
	}
	fmt.Println("res1 ", res1, err1)
	//var s *remote.ActorPidResponse
}
