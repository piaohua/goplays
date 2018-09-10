package service

import (
	"errors"
	"strings"

	"goplays/web/app/entity"
	util "utils"

	"gopkg.in/mgo.v2/bson"
)

type roleService struct{}

// 生成ID
func (this *roleService) GetRoleIDGen() (string, error) {
	gen := &entity.RoleIDGen{}
	gen.Id = "last_role_id"
	Get(GenIDs, gen.Id, gen)
	if gen.LastRoleId == "" {
		gen.LastRoleId = "2"
	}
	id := gen.LastRoleId
	gen.LastRoleId = util.StringAdd(id)
	if Upsert(GenIDs, bson.M{"_id": gen.Id}, gen) {
		return id, nil
	}
	return id, errors.New("生成错误")
}

// 根据id获取角色信息
func (this *roleService) GetRole(roleId string) (*entity.Role, error) {
	role := new(entity.Role)
	Get(Roles, roleId, role)
	if role.Id == "" {
		return nil, errors.New("角色不存在")
	}
	this.loadRoleExtra(role)
	return role, nil
}

// 根据名称获取角色
func (this *roleService) GetRoleByName(roleName string) (*entity.Role, error) {
	role := &entity.Role{
		RoleName: roleName,
	}
	GetByQ(Roles, bson.M{"role_name": role.RoleName}, role)
	if role.Id == "" {
		return nil, errors.New("角色不存在")
	}
	this.loadRoleExtra(role)
	return role, nil
}

func (this *roleService) loadRoleExtra(role *entity.Role) {
	var rolePerm []entity.RolePerm
	ListByQ(RolePerms, bson.M{"role_id": role.Id}, &rolePerm)
	for _, v := range rolePerm {
		str := strings.Split(v.Perm, ".")
		if len(str) < 2 {
			continue
		}
		perm := entity.Perm{
			Module: str[0],
			Action: str[1],
			Key:    v.Perm,
		}
		role.PermList = append(role.PermList, perm)
	}
}

// 添加角色
func (this *roleService) AddRole(role *entity.Role) error {
	if Has(Roles, bson.M{"role_name": role.RoleName}) {
		return errors.New("角色已存在")
	}
	var err error
	role.Id, err = this.GetRoleIDGen()
	if err != nil {
		return errors.New("创建失败")
	}
	if Insert(Roles, role) {
		return nil
	}
	return errors.New("创建失败")
}

// 获取所有角色列表
func (this *roleService) GetAllRoles() ([]entity.Role, error) {
	var (
		roles []entity.Role // 角色列表
	)
	ListByQ(Roles, nil, &roles)
	return roles, nil
}

// 更新角色信息
func (this *roleService) UpdateRole(role *entity.Role, fields bson.M) error {
	if Has(Roles, bson.M{"role_name": role.RoleName}) {
		roleHad := new(entity.Role)
		GetByQ(Roles, bson.M{"role_name": role.RoleName}, roleHad)
		if roleHad.Id != role.Id {
			return errors.New("角色名称已存在")
		}
	}
	if Update(Roles, bson.M{"_id": role.Id}, bson.M{"$set": fields}) {
		return nil
	}
	return errors.New("更新失败")
}

// 设置角色权限
func (this *roleService) SetPerm(roleId string, perms []string) error {
	role := new(entity.Role)
	Get(Roles, roleId, role)
	if role.Id == "" {
		return errors.New("角色不存在")
	}
	all := SystemService.GetPermList()
	pmmap := make(map[string]bool)
	for _, list := range all {
		for _, perm := range list {
			pmmap[perm.Key] = true
		}
	}
	for _, v := range perms {
		if _, ok := pmmap[v]; !ok {
			return errors.New("权限名称无效:" + v)
		}
	}
	DeleteAll(RolePerms, bson.M{"role_id": roleId})
	rolePerm := &entity.RolePerm{RoleId: roleId}
	for _, v := range perms {
		rolePerm.Id = roleId + "." + v
		rolePerm.Perm = v
		Insert(RolePerms, rolePerm)
	}
	return nil
}

// 删除角色
func (this *roleService) DeleteRole(id string) error {
	role := new(entity.Role)
	Get(Roles, id, role)
	if role.Id == "" {
		return errors.New("角色不存在")
	}
	Delete(Roles, role)
	Delete(UserRoles, bson.M{"role_id": id})
	return nil
}
