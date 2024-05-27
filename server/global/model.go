package global

import (
	"gorm.io/gorm"
)

type GBN_MODEL struct {
	ID        uint64         `json:"id" gorm:"primaryKey;index;autoIncrement:true"` // 主键ID
	CreatedAt gorm.CreatedAt // 创建时间
	UpdatedAt gorm.UpdatedAt // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}
