// models/case.go
package model

import "github.com/Lazyn0tBug/beacon/server/global"

type Case struct {
	global.GBN_MODEL
	Title       string `gorm:"column:title;type:varchar(100);not null" json:"title"`
	Description string `gorm:"column:description;type:text" json:"description"`
	Status      string `gorm:"column:status;type:varchar(20);not null;default:'open'" json:"status"`
}

func (Case) TableName() string {
	return "Case"
}
