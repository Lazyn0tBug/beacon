package model

import "github.com/Lazyn0tBug/beacon/server/global"

// ServiceItem 模型表示服务项目信息。
type ServiceItem struct {
	global.GBN_MODEL
	Name         string  `gorm:"comment:服务项目名称;column:name" json:"name"`
	Price        float64 `gorm:"comment:服务项目价格;column:price" json:"price"`
	Introduction string  `gorm:"comment:服务项目简介;column:introduction" json:"introduction"`
}

func (ServiceItem) TableName() string {
	return "ServiceItem"
}
