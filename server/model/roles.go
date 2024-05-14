// models/roles.go
package model

type Roles struct {
	GBN_MODEL
	// global.GVA_MODEL
	RoleName string    `gorm:"comment:角色名"`       // 用户登录名
    Permissions []uint `gorm:"one2many:role_permissions;"`
}

func (Roles) TableName() string {
	return "Roles"
}

