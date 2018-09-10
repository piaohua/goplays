package entity

import "time"

// 用户动作
type Action struct {
	Id         string    `bson:"_id"`
	Action     string    `bson:"action"`      // 动作类型
	Actor      string    `bson:"actor"`       // 操作角色
	ObjectType string    `bson:"object_type"` // 操作对象类型
	ObjectId   string    `bson:"object_id"`   // 操作对象id
	Extra      string    `bson:"extra"`       // 额外信息
	CreateTime time.Time `bson:"create_time"` // 更新时间
	Message    string    `bson:"message"`     // 格式化后的消息
}
