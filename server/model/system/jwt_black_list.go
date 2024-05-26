package system

import (
	"github.com/Lazyn0tBug/beacon/server/global"
)

type JwtBlacklist struct {
	global.GBN_MODEL
	UserID     uint64 `gorm:"type:bigint;comment:user id" json:"user_id"`
	Jti        string `gorm:"type:text;comment:jwt token id"`
	Jwt        string `gorm:"type:text;comment:jwt token"`
	Expiration uint64 `gorm:"type:bigint;comment:过期时间"`
}
