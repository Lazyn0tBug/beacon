package system

import (
	"github.com/Lazyn0tBug/beacon/server/global"
)

type JwtBlacklist struct {
	global.GBN_MODEL
	UserID     uint64 `gorm:"type:bigint;index;column:user_id;comment:user id" json:"user_id"`
	Jti        string `gorm:"type:text;column:jti;comment:jwt token id" json:"jti"`
	Jwt        string `gorm:"type:text;column:jwt;comment:jwt token" json:"jwt"`
	Expiration uint64 `gorm:"type:bigint;column:expiration;comment:过期时间" json:"expiration"`
	IsActive   bool   `gorm:"column:is_active; size:1;not null;default:1;comment: 是否为活跃令牌" json:"active"`
}

type JwtInActive struct {
	global.GBN_MODEL
	UserID     uint64 `gorm:"type:bigint;column:user_id;comment:user id" json:"user_id"`
	Jti        string `gorm:"type:text;column:jti;comment:jwt token id" json:"jti"`
	Jwt        string `gorm:"type:text;column:jwt;comment:jwt token" json:"jwt"`
	IP         string `gorm:"type:text;column:jti;comment:ip address" json:"ip"`
	Expiration uint64 `gorm:"type:bigint;index;column:expiration;comment:过期时间" json:"expiration"`
	IsActive   bool   `gorm:"column:is_active; size:1;not null;default:1;comment: 是否为活跃令牌" json:"active"`
}
