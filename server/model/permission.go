// models/perrmissions.go
package model

import "github.com/Lazyn0tBug/beacon/server/global"

type Permission struct {
	global.GBN_MODEL
	// global.GVA_MODEL
	PermissionName string `gorm:"comment:权限名"` // 用户登录名
}

func (Permission) TableName() string {
	return "Permission"
}
