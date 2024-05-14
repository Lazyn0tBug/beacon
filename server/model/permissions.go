// models/perrmissions.go
package model

type Permissions struct {
	GBN_MODEL
	// global.GVA_MODEL
	PermissionName string    `gorm:"comment:权限名"`       // 用户登录名
}

func (Permissions) TableName() string {
	return "Permissions"
}


