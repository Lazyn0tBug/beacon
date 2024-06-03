package model

import "github.com/Lazyn0tBug/beacon/server/global"

// Consumption 模型表示消费信息。
type Consumption struct {
	global.GBN_MODEL
	UserID       uint    `gorm:"comment:用户ID;column:user_id" json:"user_id"`
	ConsumeTime  string  `gorm:"comment:消费时间;column:consume_time" json:"consume_time"`
	TotalPrice   float64 `gorm:"comment:消费总价;column:total_price" json:"total_price"`
	ServiceItems string  `gorm:"comment:消费服务项目;column:service_items" json:"service_items"`
}

func (Consumption) TableName() string {
	return "Consumption"
}
