package model

import "github.com/Lazyn0tBug/beacon/server/global"

type CaseInfo struct {
	global.GBN_MODEL
	Title       string `gorm:"column:title;not null" json:"title"`
	Description string `gorm:"column:description" json:"description"`
	Status      string `gorm:"column:status;not null;default:'open'" json:"status"`
}

func (CaseInfo) TableName() string {
	return "CaseInfo"
}
