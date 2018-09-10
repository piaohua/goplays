package controllers

import (
	"encoding/json"

	"goplays/web/app/entity"
	"goplays/web/app/service"
	"utils"

	"gopkg.in/mgo.v2/bson"
)

type ChartsController struct {
	BaseController
}

// 在线统计
func (this *ChartsController) Online() {

	list, _ := service.ChartsService.GetOnlineList(1, this.pageSize)
	//var list []entity.LogOnline
	//for i := 0; i < 5; i++ {
	//	l := entity.LogOnline{
	//		Num:   i,
	//		Ctime: utils.LocalTime(),
	//	}
	//	list = append(list, l)
	//}

	data1 := make([]entity.ChartData, 0)
	for _, v := range list {
		d := entity.ChartData{
			Label: utils.Format("Y-m-d H:i:s", v.Ctime),
			Value: utils.String(v.Num),
		}
		data1 = append(data1, d)
	}
	data, _ := json.Marshal(data1)
	/*
			data := `[
		              {
		                "label": "Mon",
		                "value": "4123"
		              },
		              {
		                "label": "Tue",
		                "value": "4633"
		              },
		              {
		                "label": "Wed",
		                "value": "5507"
		              },
		              {
		                "label": "Thu",
		                "value": "4910"
		              },
		              {
		                "label": "Fri",
		                "value": "5529"
		              },
		              {
		                "label": "Sat",
		                "value": "5803"
		              },
		              {
		                "label": "Sun",
		                "value": "6202"
		              }
		            ]`
	*/

	this.Data["pageTitle"] = "在线统计"
	this.Data["data"] = string(data)
	this.display()
}

// 注册统计
func (this *ChartsController) Regist() {

	list, _ := service.LoggerService.GetRegTodayList(1, this.pageSize, bson.M{})

	data1 := make([]entity.ChartData, 0)
	for _, v := range list {
		d := entity.ChartData{
			Label: utils.Format("Y-m-d", v.DayStamp),
			Value: utils.String(v.Num),
		}
		data1 = append(data1, d)
	}
	data, _ := json.Marshal(data1)

	this.Data["pageTitle"] = "注册统计"
	this.Data["data"] = string(data)
	this.display()
}

// 充值统计
func (this *ChartsController) Pay() {

	list, _ := service.LoggerService.GetPayTodayList(1, this.pageSize, bson.M{})

	data1 := make([]entity.ChartLabel, 0) //时间
	data2 := make([]entity.ChartValue, 0) //充值人数
	data3 := make([]entity.ChartValue, 0) //充值金额
	data4 := make([]entity.ChartValue, 0) //钻石数量
	for _, v := range list {
		d1 := entity.ChartLabel{
			Label: utils.Format("Y-m-d", v.DayStamp),
		}
		d2 := entity.ChartValue{
			Value: utils.String(v.Count),
		}
		d3 := entity.ChartValue{
			Value: utils.String((v.Money / 100)),
		}
		d4 := entity.ChartValue{
			Value: utils.String(v.Diamond),
		}
		data1 = append(data1, d1)
		data2 = append(data2, d2)
		data3 = append(data3, d3)
		data4 = append(data4, d4)
	}
	s1, _ := json.Marshal(data1)
	s2, _ := json.Marshal(data2)
	s3, _ := json.Marshal(data3)
	s4, _ := json.Marshal(data4)

	this.Data["pageTitle"] = "充值统计"
	this.Data["data1"] = string(s1)
	this.Data["data2"] = string(s2)
	this.Data["data3"] = string(s3)
	this.Data["data4"] = string(s4)
	this.display()
}
