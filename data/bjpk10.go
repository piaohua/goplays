package data

//api官网
//付费
//http://www.opencai.net/apipaid/
/*
{"rows":5,"code":"bjpk10","remain":"2189hrs","data":
[{"expect":"664880","opencode":"01,06,04,07,03,05,08,10,02,09","opentime":"2018-02-01 18:02:37","opentimestamp":1517479357},
{"expect":"664879","opencode":"01,08,07,06,04,09,05,03,02,10","opentime":"2018-02-01 17:57:42","opentimestamp":1517479062},
{"expect":"664878","opencode":"08,05,01,03,06,02,04,07,10,09","opentime":"2018-02-01 17:52:36","opentimestamp":1517478756},
{"expect":"664877","opencode":"03,06,07,02,01,10,09,05,04,08","opentime":"2018-02-01 17:47:39","opentimestamp":1517478459},
{"expect":"664876","opencode":"02,05,03,10,04,01,06,08,07,09","opentime":"2018-02-01 17:42:39","opentimestamp":1517478159}]}
*/

//免费
//http://www.opencai.net/apifree/

//限制
//开彩网API开放平台(免费接口)的开奖结果发布有2分钟至6分钟的随机延迟。
//开彩网API开放平台(免费接口)在调用过程中请注意控制每种彩票访问间隔不小于3秒1次。
//格式
//免费接口网址	http://f.apiplus.net/[彩票代码]-[返回行数].[返回格式]
//不填时默认返回5行数据
//http://f.apiplus.net/bjpk10.json
/*
{"rows":5,"code":"bjpk10","info":"免费接口随机延迟3-6分钟，实时接口请访问www.opencai.net查询、购买或续费","data":
[{"expect":"661447","opencode":"05,06,08,09,03,10,01,02,07,04","opentime":"2018-01-13 15:22:20","opentimestamp":1515828140},
{"expect":"661446","opencode":"06,08,01,07,10,04,09,05,03,02","opentime":"2018-01-13 15:17:20","opentimestamp":1515827840},
{"expect":"661445","opencode":"08,09,01,04,10,06,07,05,03,02","opentime":"2018-01-13 15:12:20","opentimestamp":1515827540},
{"expect":"661444","opencode":"03,10,09,01,05,08,02,04,07,06","opentime":"2018-01-13 15:07:20","opentimestamp":1515827240},
{"expect":"661443","opencode":"06,03,07,02,09,01,10,04,05,08","opentime":"2018-01-13 15:02:20","opentimestamp":1515826940}]}
*/

type Bjpk10Info struct {
	Rows   int      `json:"rows"`
	Code   string   `json:"code"`
	Info   string   `json:"info"`
	Remain string   `json:"remain"`
	Data   []Bjpk10 `json:"data"`
}

type Bjpk10 struct {
	Expect        string `json:"expect"`
	Opencode      string `json:"opencode"`
	Opentime      string `json:"opentime"`
	Opentimestamp int64  `json:"opentimestamp"`
}

//官网
//http://www.bwlc.net/
