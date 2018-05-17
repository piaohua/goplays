/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-11-19 13:11:36
 * Filename      : client.go
 * Description   : 机器人
 * *******************************************************/
package main

import (
	"bytes"
	"errors"
	"time"

	"goplays/glog"
	"goplays/pb"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second    // Time allowed to write a message to the peer.
	pongWait       = 60 * time.Second    // Time allowed to read the next pong message from the peer.
	pingPeriod     = (pongWait * 9) / 10 // Send pings to peer with this period. Must be less than pongWait.
	maxMessageSize = 1024                // Maximum message size allowed from peer.
	waitForLogin   = 20 * time.Second    // 连接建立后5秒内没有收到登陆请求,断开socket
)

type WebsocketConnSet map[*websocket.Conn]struct{}

// 机器人连接数据
type Robot struct {
	conn *websocket.Conn // websocket连接

	stopCh chan struct{}    // 关闭通道
	msgCh  chan interface{} // 消息通道

	maxMsgLen uint32 // 最大消息长度
	index     int    // 包序

	//玩家游戏数据
	data *user //数据
	//游戏桌子数据
	code   string //邀请码
	gtype  uint32 //游戏类型
	rtype  uint32 //房间类型
	ltype  uint32 //游戏类型
	roomid string //房间id
	envBet int32  //下注规则
	seat   uint32 //位置
	//临时数据
	round  uint32 //玩牌局数
	sits   uint32 //尝试坐下次数
	bits   uint32 //尝试下注次数
	bitNum uint32 //尝试下注数量
	regist bool   //注册标识
	timer  uint32 //在线时间
	//
	dealerSeat uint32 //庄家位置
	betSeat    uint32 //下注位置
}

// 基本数据
type user struct {
	Userid   string // 用户id
	Nickname string // 用户昵称
	Sex      uint32 // 用户性别,男1 女2 非男非女3
	Phone    string // 绑定的手机号码
	Coin     int64  // 金币
	Diamond  int64  // 钻石
	Card     int64  // 房卡
	Chip     int64  // 筹码
	Vip      uint32 // vip
}

type WSPING int

//通道关闭信号
type closeFlag int

//创建连接
func newRobot(conn *websocket.Conn, pendingWriteNum int, maxMsgLen uint32) *Robot {
	return &Robot{
		maxMsgLen: maxMsgLen,

		conn: conn,
		data: new(user),

		msgCh:  make(chan interface{}, pendingWriteNum),
		stopCh: make(chan struct{}),
	}
}

//断开连接
func (ws *Robot) Close() {
	select {
	case <-ws.stopCh:
		return
	default:
		//关闭消息通道
		ws.Sender(closeFlag(1))
		//关闭连接
		ws.conn.Close()
		//Logout message
		Logout(ws.roomid, ws.data.Phone, ws.code, ws.data.Chip)
	}
}

//接收
func (ws *Robot) Router(id uint32, body []byte) {
	//body = pbAesDe(body) //解密
	msg, err := pb.Runpack(id, body)
	if err != nil {
		glog.Error("protocol unpack err:", id, err)
		return
	}
	ws.receive(msg)
}

//发送消息
func (ws *Robot) Sender(msg interface{}) {
	if ws.msgCh == nil {
		glog.Errorf("WSConn msg channel closed %#v", msg)
		return
	}
	if len(ws.msgCh) == cap(ws.msgCh) {
		glog.Errorf("send msg channel full -> %d", len(ws.msgCh))
		return
	}
	select {
	case <-ws.stopCh:
		glog.Info("sender closed")
		return
	default: //防止阻塞
	}
	select {
	case <-ws.stopCh:
		return
	case ws.msgCh <- msg:
	}
}

//时钟
func (ws *Robot) ticker() {
	tick := time.Tick(10 * time.Second)
	for {
		select {
		case <-ws.stopCh:
			glog.Info("ticker closed")
			return
		default: //防止阻塞
		}
		select {
		case <-ws.stopCh:
			glog.Info("ticker closed")
			return
		case <-tick:
			ws.SendPing()
		}
	}
}

func (ws *Robot) readPump() {
	defer ws.Close()
	ws.conn.SetReadLimit(maxMessageSize)
	ws.conn.SetReadDeadline(time.Now().Add(pongWait))
	ws.conn.SetPongHandler(func(string) error { ws.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	// 消息缓冲
	msgbuf := bytes.NewBuffer(make([]byte, 0, 1024))
	// 消息长度
	var length int = 0
	// 包序长度
	var index int = 0
	// 协议编号
	var proto uint32 = 0
	for {
		n, message, err := ws.conn.ReadMessage()
		if err != nil {
			glog.Errorf("Read error: %s, %d\n", err, n)
			break
		}
		// 数据添加到消息缓冲
		m, err := msgbuf.Write(message)
		if err != nil {
			glog.Errorf("Buffer write error: %s, %d\n", err, m)
			return
		}
		// 消息分割循环
		for {
			// 消息头
			if length == 0 && msgbuf.Len() >= 9 {
				index = int(msgbuf.Next(1)[0])             //包序
				proto = decodeUint32(msgbuf.Next(4))       //协议号
				length = int(decodeUint32(msgbuf.Next(4))) //消息长度
				// 检查超长消息
				if length > 1024 {
					glog.Errorf("Message too length: %d\n", length)
					return
				}
			}
			//fmt.Printf("index: %d, proto: %d, length: %d, len: %d\n", index, proto, length, msgbuf.Len())
			// 消息体
			if length > 0 && msgbuf.Len() >= length {
				//fmt.Printf("Client messge: %s\n", string(msgbuf.Next(length)))
				//包序验证
				//fmt.Printf("Message index error: %d, %d\n", index, ws.index)
				if ws.index != index {
					//glog.Errorf("Message index error: %d, %d\n", index, ws.index)
					//return
				}
				//路由
				ws.Router(proto, msgbuf.Next(length))
				length = 0
			} else {
				break
			}
		}
	}
}

//消息写入
func (ws *Robot) writePump() {
	for {
		select {
		case message, ok := <-ws.msgCh:
			if !ok {
				ws.write(websocket.CloseMessage, []byte{})
				return
			}
			err := ws.write(websocket.BinaryMessage, message)
			if err != nil {
				//停止发送消息
				close(ws.stopCh)
				return
			}
		}
	}
}

//Send pings
func (ws *Robot) pingPump() {
	tick := time.Tick(pingPeriod)
	for {
		select {
		case <-tick:
			ws.Sender(WSPING(1))
		case <-ws.stopCh:
			return
		}
	}
}

//写入
func (ws *Robot) write(mt int, msg interface{}) error {
	var message []byte
	switch msg.(type) {
	case closeFlag:
		return errors.New("msg channel closed")
	case WSPING:
		mt = websocket.PingMessage
		message = nil
	case []byte:
		message = msg.([]byte)
	default:
		code, body, err := pb.Rpacket(msg)
		if err != nil {
			glog.Errorf("write msg err %#v", msg)
			return err
		}
		message = pack(code, body, ws.index)
		if ws.index >= 255 {
			ws.index = 0
		} else {
			ws.index += 1
		}
	}
	if uint32(len(message)) > ws.maxMsgLen {
		glog.Errorf("write msg too long -> %d", len(message))
		return errors.New("write msg too long")
	}
	ws.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return ws.conn.WriteMessage(mt, message)
}

func decodeUint32(b []byte) (i uint32) {
	i = uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
	return
}

func encodeUint32(i uint32) (b []byte) {
	b = append(b, byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
	return
}

//封包
func pack(code uint32, msg []byte, index int) []byte {
	//msg = pbAesEn(msg) //加密
	buff := make([]byte, 9+len(msg))
	msglen := uint32(len(msg))
	buff[0] = byte(index)
	copy(buff[1:5], encodeUint32(code))
	copy(buff[5:9], encodeUint32(msglen))
	copy(buff[9:], msg)
	return buff
}
