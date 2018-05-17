package main

import (
	"bytes"
	"errors"
	"strings"
	"time"

	"goplays/glog"
	"goplays/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second    // Time allowed to write a message to the peer.
	pongWait       = 60 * time.Second    // Time allowed to read the next pong message from the peer.
	pingPeriod     = (pongWait * 9) / 10 // Send pings to peer with this period. Must be less than pongWait.
	maxMessageSize = 1024                // Maximum message size allowed from peer.
	waitForLogin   = 20 * time.Second    // 连接建立后5秒内没有收到登陆请求,断开socket
)

type WSPING int

//通道关闭信号
type closeFlag int

type WebsocketConnSet map[*websocket.Conn]struct{}

type WSConn struct {
	conn *websocket.Conn // websocket连接

	maxMsgLen uint32 // 最大消息长度
	index     int    // 包序

	stopCh chan struct{}    // 关闭通道
	msgCh  chan interface{} // 消息通道

	pid     *actor.PID // ws进程ID,登录成功后切换为rs进程
	rolePid *actor.PID // 角色服务

	online bool //在线状态
}

//创建连接
func newWSConn(conn *websocket.Conn, pendingWriteNum int, maxMsgLen uint32) *WSConn {
	return &WSConn{
		conn:      conn,
		maxMsgLen: maxMsgLen,
		msgCh:     make(chan interface{}, pendingWriteNum),
		stopCh:    make(chan struct{}),
	}
}

//连接地址
func (ws *WSConn) localAddr() string {
	return ws.conn.LocalAddr().String()
}

func (ws *WSConn) remoteAddr() string {
	return ws.conn.RemoteAddr().String()
}

func (ws *WSConn) GetIPAddr() string {
	return strings.Split(ws.remoteAddr(), ":")[0]
}

//断开连接
func (ws *WSConn) Close() {
	select {
	case <-ws.stopCh:
		return
	default:
		//glog.Debugf("ws closed closeFlag %d", len(ws.msgCh))
		//关闭消息通道
		ws.Send(closeFlag(1))
		//停止发送消息
		close(ws.stopCh)
		//关闭连接
		ws.conn.Close()
	}
}

//index(1byte) + proto(4byte) + msgLen(4byte) + msg
func (ws *WSConn) readPump() {
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
					//fmt.Printf("Message index error: %d, %d\n", index, ws.index)
					//glog.Errorf("Message index error: %d, %d\n", index, ws.index)
					//return
				}
				if ws.index >= 255 {
					ws.index = 0
				} else {
					ws.index += 1
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
func (ws *WSConn) writePump() {
	for {
		select {
		case msg, ok := <-ws.msgCh:
			if !ok {
				ws.write(websocket.CloseMessage, []byte{})
				return
			}
			err := ws.write(websocket.BinaryMessage, msg)
			if err != nil {
				return
			}
		}
	}
}

//Send pings
func (ws *WSConn) pingPump() {
	tick := time.Tick(pingPeriod)
	for {
		select {
		case <-tick:
			ws.Send(WSPING(1))
		case <-ws.stopCh:
			return
		}
	}
}

//写入
func (ws *WSConn) write(mt int, msg interface{}) error {
	var message []byte
	switch msg.(type) {
	case closeFlag:
		return errors.New("msg channel closed")
	case WSPING:
		mt = websocket.PingMessage
	case []byte:
		message = msg.([]byte)
	default:
		code, body, err := pb.Packet(msg)
		if err != nil {
			glog.Errorf("write msg err %#v", msg)
			return err
		}
		message = pack(code, body, ws.index)
	}
	if uint32(len(message)) > ws.maxMsgLen {
		glog.Errorf("write msg too long -> %d", len(message))
		//return errors.New("write msg too long")
	}
	ws.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return ws.conn.WriteMessage(mt, message)
}
