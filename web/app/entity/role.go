package entity

import "time"

// 角色自增id
type RoleIDGen struct {
	Id         string `bson:"_id"`
	LastRoleId string `bson:"last_role_id"`
}

// 角色, 分组
type Role struct {
	Id          string    `bson:"_id"`         // AUTO_INCREMENT, PRIMARY KEY (`id`)
	RoleName    string    `bson:"role_name"`   // 角色名称
	ProjectIds  string    `bson:"project_ids"` // 项目权限
	Description string    `bson:"description"` // 说明
	CreateTime  time.Time `bson:"create_time"` // 创建时间
	UpdateTime  time.Time `bson:"update_time"` // 更新时间
	PermList    []Perm    `bson:"perm_list"`   // 权限列表
	UserList    []User    `bson:"user_list"`   // 用户列表
}

// 角色权限, 分组权限设置
type RolePerm struct {
	Id     string `bson:"_id"`     // PRIMARY KEY (`role_id`,`perm`)
	RoleId string `bson:"role_id"` // 角色id
	Perm   string `bson:"perm"`    // 权限
}
