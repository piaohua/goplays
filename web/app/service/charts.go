package service

import "goplays/web/app/entity"

type chartsService struct{}

// 获取绑定日志列表
func (this *chartsService) GetOnlineList(page, pageSize int) ([]entity.LogOnline, error) {
	var list []entity.LogOnline
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	LogOnlines.
		Find(nil).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, nil
}

// 获取注册日志总数
func (this *chartsService) GetOnlineTotal() (int64, error) {
	return int64(Count(LogOnlines, nil)), nil
}
