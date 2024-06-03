package model

import "github.com/Lazyn0tBug/beacon/server/global"

// Hospital 模型表示医院信息。
type Hospital struct {
	global.GBN_MODEL
	Name         string `gorm:"comment:医院名称;column:name" json:"name"`
	Address      string `gorm:"comment:医院地址;column:address" json:"address"`
	Introduction string `gorm:"comment:医院简介;column:introduction" json:"introduction"`
	Avatar       string `gorm:"comment:医院头像;column:avatar" json:"avatar"`
}

func (Hospital) TableName() string {
	return "Hospital"
}
