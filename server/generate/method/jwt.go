package method

import (
	"github.com/Lazyn0tBug/beacon/server/model/system"
)

type JwtBlacklist interface {
	// 用户登出时，创建一条 jwt_blacklist 记录
	//
	// INSERT INTO jwt_blacklist (user_id, jti, jwt, expiration, is_active)
	// VALUES (@user_id, @jti, @jwt, @expiration, 1)
	Create(jwtBlacklist *system.JwtBlacklist) error

	// 当token过期时或用户再次登录时，将 jwt_blacklist 记录设置为有效
	//
	// Update @@table SET is_active = 0 WHERE jti = @jti AND user_id = @user_id
	DeActive(userID uint64, jti string) error

	// 当用户退出时，将 jwt_blacklist 记录设置为有效
	//
	// Update @@table SET is_active = 1 WHERE jti = @jti AND user_id = @user_id
	ReActive(userID uint64, jti string) error

	// List 获取 jwt_blacklist 记录列表
	//
	// SELECT * FROM jwt_blacklist WHERE expiration > @now ORDER BY created_at DESC LIMIT @limit OFFSET @offset
	// @now 表示当前时间戳，用于过滤已过期的记录
	List(now uint64, limit, offset int) ([]*system.JwtBlacklist, error)

	// Sync 将 Redis 中的 jwt_blacklist 记录同步到数据库中
	//
	// 该方法在系统崩溃时和定时任务中调用，以确保数据库中的 jwt_blacklist 记录与 Redis 中的记录保持一致
	Sync() error

	// Load 从数据库中加载 jwt_blacklist 记录到 Redis 中
	//
	// 该方法在系统启动时调用，以将数据库中的 jwt_blacklist 记录加载到 Redis 中，并在加载过程中自动忽略已经过期的记录
	Load() error

	// BatchInvalidateTokens 批量使过期的 jwt_blacklist 记录失效，即将其中的 is_active 字段设置为 false
	//
	// 该方法可以在定时任务中调用，以确保 Redis 中的 jwt_blacklist 记录与数据库中的记录保持一致
	//
	// UPDATE jwt_blacklist SET is_active = false WHERE jti IN (@jtis) AND expiration < @now
	BatchInvalidateTokens() error
}

type JwtInActiveDAO interface {
	// CreateOrUpdate 创建或更新一条 jwt_inactive 记录
	//
	// INSERT INTO jwt_inactive (user_id, jti, token, ip, expiration, active)
	// VALUES (@user_id, @jti, @token, @ip, @expiration, @active)
	// ON CONFLICT (user_id) DO UPDATE SET jti=@jti, token=@token, ip=@ip, expiration=@expiration, active=@active
	CreateOrUpdate(jwtInActive *system.JwtInActive) error

	// List 获取 jwt_inactive 记录列表
	//
	// SELECT * FROM jwt_inactive WHERE deleted_at IS NULL AND active = true ORDER BY created_at DESC LIMIT @limit OFFSET @offset
	List(limit, offset int) ([]*system.JwtInActive, error)

	// Sync 将 Redis 中的 jwt_inactive 记录同步到数据库中
	//
	// 该方法在系统崩溃时和定时任务中调用，以确保数据库中的 jwt_inactive 记录与 Redis 中的记录保持一致
	Sync() error

	// Load 从数据库中加载 jwt_inactive 记录到 Redis 中
	//
	// 该方法在系统启动时调用，以将数据库中的 jwt_inactive 记录加载到 Redis 中，并在加载过程中自动忽略已经过期的记录
	Load() error

	// CheckExpiredTokens 定期检查过期的 jwt_inactive 记录，并将其中的 active 字段设置为 false
	//
	// 该方法可以在定时任务中调用，以确保数据库中的 jwt_inactive 记录的 active 字段与实际情况保持一致
	CheckExpiredTokens() error
}
