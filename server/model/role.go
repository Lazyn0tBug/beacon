// models/roles.go
package model

import "github.com/Lazyn0tBug/beacon/server/global"

type Role struct {
	global.GBN_MODEL
	// global.GVA_MODEL
	RoleName    string `gorm:"not null;unique;comment:角色名"` // 用户登录名
	Permissions []uint `gorm:"one2many:role_permissions;"`
}

func (Role) TableName() string {
	return "Role"
}
