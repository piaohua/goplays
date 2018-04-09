package main

import (
	"fmt"
	"strings"
	"time"

	"goplays/game/handler"
	"goplays/glog"
	"goplays/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//web请求处理
func (a *HallActor) HandlerWeb(arg *pb.WebRequest,
	rsp *pb.WebResponse, ctx actor.Context) {
	switch arg.Code {
	case pb.WebOnline:
		msg1 := make([]string, 0)
		err1 := json.Unmarshal(arg.Data, &msg1)
		if err1 != nil {
			rsp.ErrMsg = fmt.Sprintf("msg err: %v", err1)
			return
		}
		//响应
		resp := make(map[string]int)
		for _, v := range msg1 {
			if _, ok := a.roles[v]; ok {
				resp[v] = 1
			} else {
				resp[v] = 0
			}
		}
		result, err2 := json.Marshal(resp)
		if err2 != nil {
			rsp.ErrMsg = fmt.Sprintf("msg err: %v", err2)
			return
		}
		rsp.Result = result
	case pb.WebShop:
		//更新配置
		msg2 := handler.SyncConfig2(pb.CONFIG_SHOP, arg.Atype, arg.Data)
		err1 := handler.SyncConfig(msg2)
		if err1 != nil {
			rsp.ErrMsg = fmt.Sprintf("msg err: %v", err1)
			return
		}
		//广播所有节点,主动通知同步配置,只同步修改数据
		a.broadcast(msg2, ctx)
	case pb.WebEnv:
		//更新配置
		msg2 := handler.SyncConfig2(pb.CONFIG_ENV, arg.Atype, arg.Data)
		err1 := handler.SyncConfig(msg2)
		if err1 != nil {
			rsp.ErrMsg = fmt.Sprintf("msg err: %v", err1)
			return
		}
		//广播所有节点,主动通知同步配置,只同步修改数据
		a.broadcast(msg2, ctx)
	case pb.WebNotice:
		//更新配置
		msg2 := handler.SyncConfig2(pb.CONFIG_NOTICE, arg.Atype, arg.Data)
		err1 := handler.SyncConfig(msg2)
		if err1 != nil {
			rsp.ErrMsg = fmt.Sprintf("msg err: %v", err1)
			return
		}
		//广播所有节点,主动通知同步配置,只同步修改数据
		a.broadcast(msg2, ctx)
	case pb.WebGame:
		//更新配置
		msg2 := handler.SyncConfig2(pb.CONFIG_GAMES, arg.Atype, arg.Data)
		err1 := handler.SyncConfig(msg2)
		if err1 != nil {
			rsp.ErrMsg = fmt.Sprintf("msg err: %v", err1)
			return
		}
		//广播所有节点,主动通知同步配置,只同步修改数据
		a.broadcast(msg2, ctx)
	case pb.WebVip:
		//更新配置
		msg2 := handler.SyncConfig2(pb.CONFIG_VIP, arg.Atype, arg.Data)
		err1 := handler.SyncConfig(msg2)
		if err1 != nil {
			rsp.ErrMsg = fmt.Sprintf("msg err: %v", err1)
			return
		}
		//广播所有节点,主动通知同步配置,只同步修改数据
		a.broadcast(msg2, ctx)
	case pb.WebBuild:
		//TODO
	case pb.WebGive:
		//后台货币赠送同步到game房间
		msg2 := new(pb.PayCurrency)
		err1 := msg2.Unmarshal(arg.Data)
		if err1 != nil {
			rsp.ErrMsg = fmt.Sprintf("msg err: %v", err1)
			return
		}
		//消息
		a.msg2role(msg2.Userid, msg2)
	case pb.WebNumber:
		result, err2 := a.getNumber(ctx)
		if err2 != nil {
			rsp.ErrMsg = fmt.Sprintf("msg err: %v", err2)
			return
		}
		rsp.Result = result
	default:
		glog.Errorf("unknown message %v", arg)
	}
}

//广播所有节点,游戏逻辑服,dbms
func (a *HallActor) broadcast(msg interface{}, ctx actor.Context) {
	for k, v := range a.serve {
		if strings.Contains(k, "dbms") ||
			strings.Contains(k, "gate.") ||
			strings.Contains(k, "game.") {
			v.Tell(msg)
		}
	}
}

//消息通知到玩家
func (a *HallActor) msg2role(userid string, msg interface{}) {
	gate := a.roles[userid]
	//存活在节点中
	if v, ok := a.serve[gate]; ok {
		v.Tell(msg)
		return
	}
	//不在节点中直接同步到数据库
	role := cfg.Section("role").Name()
	if v, ok := a.serve[role]; ok {
		v.Tell(msg)
		return
	}
}

//
func (a *HallActor) getNumber(ctx actor.Context) ([]byte, error) {
	msg2 := new(pb.GetNumber)
	res2 := a.reqRole(msg2, ctx)
	var response2 *pb.GotNumber
	var ok bool
	if response2, ok = res2.(*pb.GotNumber); !ok {
		glog.Error("response msg error")
		return nil, fmt.Errorf("response msg error")
	}
	//响应1 机器人,2 玩家
	resp := make(map[int]int)
	for _, v := range response2.Robot {
		if _, ok := a.roles[v]; ok {
			resp[1] += 1
		}
	}
	for _, v := range response2.Role {
		if _, ok := a.roles[v]; ok {
			resp[2] += 1
		}
	}
	result, err2 := json.Marshal(resp)
	if err2 != nil {
		glog.Errorf("msg err: %v", err2)
		return nil, fmt.Errorf("msg err: %v", err2)
	}
	return result, nil
}

//数据处理
func (a *HallActor) reqRole(msg interface{}, ctx actor.Context) interface{} {
	role := cfg.Section("role").Name()
	if v, ok := a.serve[role]; ok {
		timeout := 3 * time.Second
		res1, err1 := v.RequestFuture(msg, timeout).Result()
		if err1 != nil {
			glog.Errorf("reqRole err: %v, msg %#v", err1, msg)
			return nil
		}
		return res1
	}
	return nil
}
