package entity

// 权限设置
type Perm struct {
	Id     string `bson:"_id"` // UNIQUE KEY `module` (`module`,`action`)
	Module string `bson:"module"`
	Action string `bson:"action"`
	Key    string `bson:"key"` // Module.Action
}

func (p *Perm) TableUnique() [][]string {
	return [][]string{
		[]string{"Module", "Action"},
	}
}
