package system

import (
	"github.com/Lazyn0tBug/beacon/server/global"
)

type JwtBlacklist struct {
	global.GBN_MODEL
	Jwt string `gorm:"type:text;comment:jwt"`
}

