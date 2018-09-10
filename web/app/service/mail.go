package service

import (
	"errors"

	"goplays/web/app/entity"

	"gopkg.in/mgo.v2/bson"
)

type mailService struct{}

func (this *mailService) AddMailTpl(tpl *entity.MailTpl) error {
	tpl.Id = bson.NewObjectId().String()
	if Insert(MailTpls, tpl) {
		return nil
	}
	return errors.New("添加失败")
}

func (this *mailService) DelMailTpl(id string) error {
	if Delete(MailTpls, bson.M{"_id": id}) {
		return nil
	}
	return errors.New("删除失败")
}

func (this *mailService) SaveMailTpl(tpl *entity.MailTpl) error {
	if Update(MailTpls, bson.M{"_id": tpl.Id}, tpl) {
		return nil
	}
	return errors.New("更新失败")
}

func (this *mailService) GetMailTpl(id string) (*entity.MailTpl, error) {
	tpl := &entity.MailTpl{}
	Get(MailTpls, id, tpl)
	if tpl.Id != "" {
		return tpl, nil
	}
	return tpl, errors.New("获取失败")
}

// 获取邮件模板列表
func (this *mailService) GetMailTplList() ([]entity.MailTpl, error) {
	var list []entity.MailTpl
	ListByQ(MailTpls, nil, &list)
	return list, nil
}
