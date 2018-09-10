package entity

import "time"

// 表结构
type MailTpl struct {
	Id         string    `bson:"_id"`
	UserId     string    `bson:"user_id"`
	Name       string    `bson:"name"`
	Subject    string    `bson:"subject"`
	Content    string    `bson:"content"`
	MailTo     string    `bson:"mail_to"`
	MailCc     string    `bson:"mail_cc"`
	CreateTime time.Time `bson:"create_time"`
	UpdateTime time.Time `bson:"update_time"`
}
