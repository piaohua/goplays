package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"goplays/pb"
	"goplays/web/app/entity"
	"utils"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

var gmHost string
var gmPort string
var gmPath string
var gmKey string
var gmUrl string

func init() {
	gmHost = beego.AppConfig.String("gm.host")
	gmPort = beego.AppConfig.String("gm.port")
	gmPath = beego.AppConfig.String("gm.path")
	gmKey = beego.AppConfig.String("gm.key")
	gmUrl = beego.AppConfig.String("gm.url")
}

var cstDialer = websocket.Dialer{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//Gm操作,角色ID,操作类型,操作物品,操作数量
func Gm(msgName, msg string) (string, error) {
	addr := gmHost + ":" + gmPort
	//fmt.Println("addr -> ", addr)
	u := url.URL{Scheme: "wss", Host: addr, Path: gmPath}
	TimeStr := GmTime()
	Token := GmToke(TimeStr)
	//c, _, err := websocket.DefaultDialer.Dial(u.String(),
	//	http.Header{"Token": {Token}})
	d := cstDialer
	d.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c, _, err := d.Dial(u.String(),
		http.Header{"Token": {Token}})
	if err != nil {
		return "", errors.New(fmt.Sprintf("dial err -> %v", err))
	}
	sign := GmSign(msg, TimeStr)
	msg2 := GmMsg(msg, msgName, sign, TimeStr)
	if c != nil {
		c.WriteMessage(websocket.TextMessage, []byte(msg2))
		defer c.Close()
		_, message, err := c.ReadMessage()
		if err != nil {
			return "", errors.New(fmt.Sprintf("read err -> %v", err))
		}
		resp := new(entity.RespErr)
		err = json.Unmarshal(message, resp)
		if err != nil {
			return "", errors.New(fmt.Sprintf("Unmarshal err -> %v", err))
		}
		if resp.ErrCode != 0 {
			return "", errors.New(fmt.Sprintf("resp.ErrMsgi %s", resp.ErrMsg))
		}
		return resp.Result, nil
	}
	return "", errors.New(fmt.Sprintf("c empty err -> %v", err))
}

//字符串时间
func GmTime() string {
	Time := utils.Timestamp()
	TimeStr := utils.String(Time)
	return TimeStr
}

// Sign := utils.Md5(Key+Now)
// Token := Sign+Now+RandNum
func GmToke(TimeStr string) string {
	Sign := utils.Md5(gmKey + TimeStr)
	Token := Sign + TimeStr + utils.RandStr(6)
	return Token
}

// Sign := TimeStr + Key + Md5(msg)
func GmSign(msg, TimeStr string) string {
	return utils.Md5(TimeStr + gmKey + utils.Md5(msg))
}

// Timestr|Sign|msg_name|msg
func GmMsg(msg, msgName, sign, TimeStr string) string {
	return TimeStr + "|" + sign + "|" + msgName + "|" + msg
}

// http request

func GmRequest(code pb.WebCode, atype pb.ConfigAtype,
	b interface{}) (interface{}, error) {
	//pack
	body, err1 := gmPack(code, atype, b)
	if err1 != nil {
		return nil, err1
	}
	//request
	result, err3 := doHttpPost(gmUrl, body)
	if err3 != nil {
		return nil, err3
	}
	//unpack
	return gmUnpack(code, result)
}

func gmPack(code pb.WebCode, atype pb.ConfigAtype,
	b interface{}) ([]byte, error) {
	msg := new(pb.WebRequest)
	msg.Code = code
	msg.Atype = atype
	switch b.(type) {
	case []byte:
		msg.Data = b.([]byte)
	case *pb.ChangeCurrency:
		msg2 := b.(*pb.ChangeCurrency)
		result, err2 := msg2.Marshal()
		if err2 != nil {
			return []byte{}, err2
		}
		msg.Data = result
	case *pb.PayCurrency:
		msg2 := b.(*pb.PayCurrency)
		result, err2 := msg2.Marshal()
		if err2 != nil {
			return []byte{}, err2
		}
		msg.Data = result
	default:
		result, err2 := json.Marshal(b)
		if err2 != nil {
			return []byte{}, err2
		}
		msg.Data = result
	}
	body, err1 := msg.Marshal()
	if err1 != nil {
		return []byte{}, err1
	}
	return body, nil
}

func gmUnpack(code pb.WebCode, body []byte) (interface{}, error) {
	resp := new(pb.WebResponse)
	err1 := resp.Unmarshal(body)
	if err1 != nil {
		fmt.Println("gmUnpack code ", code, err1)
		return nil, err1
	}
	if resp.ErrCode != 0 || resp.ErrMsg != "" {
		fmt.Println("gmUnpack code ", resp.ErrCode, resp.ErrMsg)
		return nil, fmt.Errorf("ErrCode %d, ErrMsg %s", resp.ErrCode, resp.ErrMsg)
	}
	switch code {
	case pb.WebOnline:
		b := make(map[string]int)
		err2 := json.Unmarshal(resp.Result, &b)
		if err2 != nil {
			fmt.Println("gmUnpack code ", code, err2)
			return nil, err2
		}
		return b, nil
	case pb.WebNumber:
		b := make(map[int]int)
		err2 := json.Unmarshal(resp.Result, &b)
		if err2 != nil {
			fmt.Println("gmUnpack code ", code, err2)
			return nil, err2
		}
		return b, nil
	case pb.WebShop:
	case pb.WebEnv:
		return nil, nil
	case pb.WebNotice:
	case pb.WebGame:
		return nil, nil
	case pb.WebVip:
	case pb.WebBuild:
	case pb.WebGive:
		return nil, nil
	}
	return nil, fmt.Errorf("unknown code %d", code)
}

//http post
func doHttpPost(targetUrl string, body []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", targetUrl, bytes.NewBuffer(body))
	if err != nil {
		return []byte(""), err
	}
	req.Header.Add("Content-type", "text/plain;charset=UTF-8")

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   0,
			KeepAlive: 0,
		}).Dial,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: false},
		TLSHandshakeTimeout: 10 * time.Second,
	}

	client := &http.Client{Transport: transport}

	resp, err := client.Do(req)
	if err != nil {
		return []byte(""), err
	}

	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}

	return respData, nil
}
