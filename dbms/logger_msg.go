package main

import (
	"goplays/data"
	"goplays/game/handler"
	"goplays/glog"
	"goplays/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func (a *LoggerActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.LogRegist:
		arg := msg.(*pb.LogRegist)
		data.RegistRecord(arg.Userid, arg.Nickname, arg.Ip, arg.Atype)
	case *pb.LogLogin:
		arg := msg.(*pb.LogLogin)
		data.LoginRecord(arg.Userid, arg.Ip, arg.Atype)
	case *pb.LogLogout:
		arg := msg.(*pb.LogLogout)
		data.LogoutRecord(arg.Userid, int(arg.Event))
	case *pb.LogDiamond:
		arg := msg.(*pb.LogDiamond)
		data.DiamondRecord(arg.Userid, arg.Type, arg.Rest, arg.Num)
	case *pb.LogCoin:
		arg := msg.(*pb.LogCoin)
		data.CoinRecord(arg.Userid, arg.Type, arg.Rest, arg.Num)
	case *pb.LogCard:
		arg := msg.(*pb.LogCard)
		data.CardRecord(arg.Userid, arg.Type, arg.Rest, arg.Num)
	case *pb.LogChip:
		arg := msg.(*pb.LogChip)
		data.ChipRecord(arg.Userid, arg.Type, arg.Rest, arg.Num)
	case *pb.Pk10RecordLog:
		arg := msg.(*pb.Pk10RecordLog)
		data.Pk10RecordLog(arg.Expect, arg.Opencode, arg.Opentime,
			arg.Code, arg.Opentimestamp)
	case *pb.Pk10TrendLog:
		arg := msg.(*pb.Pk10TrendLog)
		handler.Pk10TrendLog(arg)
	case *pb.Pk10UseridLog:
		arg := msg.(*pb.Pk10UseridLog)
		handler.Pk10UseridLog(arg)
	case *pb.Pk10GameLog:
		arg := msg.(*pb.Pk10GameLog)
		handler.Pk10GameLog(arg)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}
