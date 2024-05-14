package database

import (
	"gorm.io/gorm"
)

type DBInterface interface {
	GetDB() *gorm.DB
}
