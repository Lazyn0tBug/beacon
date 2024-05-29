package system

import (
	"context"
	// "errors"
	// "time"

	"go.uber.org/zap"

	"github.com/Lazyn0tBug/beacon/server/global"
	"github.com/Lazyn0tBug/beacon/server/model/system"
	"github.com/Lazyn0tBug/beacon/server/utils"
	// jwt "github.com/lestrrat-go/jwx"
)

type JwtService struct{}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: JsonInBlacklist
//@description: 拉黑jwt
//@param: jwtList model.JwtBlacklist
//@return: err error

func (jwtService *JwtService) JsonInBlacklist(jwtList system.JwtBlacklist) (err error) {
	err = global.GVA_DB.Create(&jwtList).Error
	if err != nil {
		return
	}
	global.BlackCache.SetDefault(jwtList.Jwt, struct{}{})
	return
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: IsBlacklist
//@description: 判断JWT是否在黑名单内部
//@param: jwt string
//@return: bool

func (jwtService *JwtService) IsBlacklist(jwt string) bool {
	_, ok := global.BlackCache.Get(jwt)
	return ok
	// err := global.GVA_DB.Where("jwt = ?", jwt).First(&system.JwtBlacklist{}).Error
	// isNotFound := errors.Is(err, gorm.ErrRecordNotFound)
	// return !isNotFound
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetRedisJWT
//@description: 从redis取jwt
//@param: userName string
//@return: redisJWT string, err error

func (jwtService *JwtService) GetRedisJWT(userName string) (redisJWT string, err error) {
	redisJWT, err = global.GVA_REDIS.Get(context.Background(), userName).Result()
	return redisJWT, err
}

// GetRedisJWT retrieves the JWT ID associated with the username from Redis
// func (js *JwtService) GetRedisJWT(username string, tokenString string) (string, error) {
// 	// Parse the JWT token to extract the JWT ID
// 	token, err := jwt.ParseString(tokenString)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Get the JWT ID from the token
// 	jwtID, ok := token.Get(jwt.JWTPresentersKey)
// 	if !ok {
// 		return "", errors.New("JWT ID not found in token")
// 	}

// 	// Convert jwtID to string
// 	jwtIDStr, ok := jwtID.(string)
// 	if !ok {
// 		return "", errors.New("JWT ID is not a string")
// 	}

// 	// Define the key for the Redis set
// 	key := "jwt_blacklist:" + username

// 	// Check if the JWT ID exists in the Redis set
// 	ctx := context.Background()
// 	exists, err := global.GVA_REDIS.SIsMember(ctx, key, jwtIDStr).Result()
// 	if err != nil {
// 		return "", err
// 	}
// 	if !exists {
// 		return "", errors.New("JWT ID not found for the user")
// 	}

// 	return jwtIDStr, nil
// }

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetRedisJWT
//@description: jwt存入redis并设置过期时间
//@param: jwt string, userName string
//@return: err error

func (jwtService *JwtService) SetRedisJWT(jwt string, userName string) (err error) {
	// 此处过期时间等于jwt过期时间
	dr, err := utils.ParseDuration(global.GVA_CONFIG.JWT.ExpiresTime)
	if err != nil {
		return err
	}
	timer := dr
	err = global.GVA_REDIS.Set(context.Background(), userName, jwt, timer).Err()
	return err
}

// SetRedisJWT sets the JWT associated with the username and JWT ID in Redis blacklist
// func (jwtService *JwtService) SetRedisJWT(tokenString string, userName string) error {
// 	// Parse the JWT token to extract the JWT ID and expiration
// 	token, err := jwt.ParseString(tokenString)
// 	if err != nil {
// 		return err
// 	}

// 	// Get the JWT ID from the token
// 	jwtID, ok := token.Get(jwt.JWTPresentersKey)
// 	if !ok {
// 		return errors.New("JWT ID not found in token")
// 	}

// 	// Convert jwtID to string
// 	jwtIDStr, ok := jwtID.(string)
// 	if !ok {
// 		return errors.New("JWT ID is not a string")
// 	}

// 	// Get the expiration time from the token
// 	exp, ok := token.Get(jwt.ExpirationKey)
// 	if !ok {
// 		return errors.New("Expiration time not found in token")
// 	}

// 	// Convert exp to time.Time
// 	expirationTime := time.Unix(int64(exp.(float64)), 0)

// 	// Define the key for the Redis set
// 	key := "jwt_blacklist:" + userName

// 	// Set the JWT ID in the Redis set
// 	ctx := context.Background()
// 	err = jwtService.redisClient.SAdd(ctx, key, jwtIDStr).Err()
// 	if err != nil {
// 		return err
// 	}

// 	// Set the expiration for the key
// 	jwtService.redisClient.Expire(ctx, key, expirationTime.Sub(time.Now()))

// 	return nil
// }

func LoadAll() {
	var data []string
	err := global.GVA_DB.Model(&system.JwtBlacklist{}).Select("jwt").Find(&data).Error
	if err != nil {
		global.GVA_LOG.Error("加载数据库jwt黑名单失败!", zap.Error(err))
		return
	}
	for i := 0; i < len(data); i++ {
		global.BlackCache.SetDefault(data[i], struct{}{})
	} // jwt黑名单 加入 BlackCache 中
}
