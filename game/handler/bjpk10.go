package handler

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"goplays/data"
	"goplays/glog"

	"github.com/PuerkitoBio/goquery"
	jsoniter "github.com/json-iterator/go"
)

//免费api获取(请求间隔大于3秒,存在延迟,备选方案)
func GetPk10Api2(bjpk10Api string) (d []data.Bjpk10, err error) {
	result, err := doHttpGet(bjpk10Api)
	if err != nil {
		glog.Errorf("doHttpGet err %v", err)
		return
	}
	b := new(data.Bjpk10Info)
	err = jsoniter.Unmarshal(result, &b)
	if err != nil {
		return
	}
	if len(b.Data) == 0 || b.Rows == 0 {
		err = fmt.Errorf("data empty")
		return
	}
	d = b.Data
	return
}

func doHttpGet(targetUrl string) ([]byte, error) {
	req, err := http.NewRequest("GET", targetUrl, bytes.NewBuffer([]byte{}))
	if err != nil {
		return []byte(""), err
	}
	req.Header.Add("Content-type", "application/x-www-form-urlencoded;charset=UTF-8")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	client := &http.Client{Transport: tr}

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

//免费api获取(请求间隔大于3秒,存在延迟,备选方案)
func GetPk10Api(bjpk10Api string) (index, codes string,
	opentimestamp int64, err error) {
	//bjpk10Api := "http://f.apiplus.net/bjpk10.json"
	res, err := doRequest(bjpk10Api)
	if err != nil {
		glog.Errorf("doRequest err %v", err)
		return
	}
	result, err := doRequest2(res)
	if err != nil {
		return
	}
	b := new(data.Bjpk10Info)
	err = jsoniter.Unmarshal(result, &b)
	if err != nil {
		return
	}
	if len(b.Data) == 0 || b.Rows == 0 {
		err = fmt.Errorf("data empty")
		return
	}
	index = b.Data[0].Expect
	codes = b.Data[0].Opencode
	opentimestamp = b.Data[0].Opentimestamp
	return
}

// 解析
func doRequest2(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}
	return respData, nil
}

// doRequest get the order in json format with a sign
func doRequest(targetUrl string) (*http.Response, error) {
	req, err := http.NewRequest("GET", targetUrl, bytes.NewBuffer([]byte{}))
	if err != nil {
		glog.Errorf("Request url %s err %v", targetUrl, err)
		return nil, err
	}
	req.Header.Add("Content-type", "application/x-www-form-urlencoded;charset=UTF-8")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   3 * time.Second,
	}

	return client.Do(req)
}

//官网页面上抓取(请求间隔5分钟内)
func GetPk10() (index, codes string, err error) {
	bjpk10Url := "http://www.bwlc.net/"
	res, err := doRequest(bjpk10Url)
	if err != nil {
		glog.Errorf("doRequest err %v", err)
		return
	}
	//doc, err := goquery.NewDocument(bjpk10Url)
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		glog.Errorf("NewDocument err %v", err)
		return
	}

	//期号
	index = doc.Find("div.game_list").Has("div.pk10_bg").First().Find("span.ml10").Text()
	glog.Debugf("index %s", index)

	//codes
	codes = doc.Find("div.game_list").Has("div.pk10_bg").First().Find("li").Text()
	glog.Debugf("codes %s", codes)

	if len(index) == 0 || len(codes) != 20 {
		err = fmt.Errorf("index or codes error")
		return
	}

	var codes2 string
	for k, v := range codes {
		if ((k % 2) == 1) && k != 19 {
			codes2 += string(v) + ","
		} else {
			codes2 += string(v)
		}
	}
	codes = codes2
	glog.Debugf("codes %s", codes)
	return
}
