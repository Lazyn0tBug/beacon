package utils

import (
	"errors"
	"time"

	"github.com/Lazyn0tBug/beacon/server/global"
	"github.com/Lazyn0tBug/beacon/server/model/system/request"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

const (
	TokenExpiryDuration = time.Hour * 24
)

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
	Logger           = GetLogger()
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.GVA_CONFIG.JWT.SigningKey),
	}
}

func (jwtUtility *JWT) CreateClaims(baseClaims request.BaseClaims) request.CustomClaims {
	bf, _ := ParseDuration(global.GVA_CONFIG.JWT.BufferTime)
	ep, _ := ParseDuration(global.GVA_CONFIG.JWT.ExpiresTime)
	claims := request.CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: int64(bf / time.Second), // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"GVA"},                   // 受众
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)), // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ep)),    // 过期时间 7天  配置文件
			Issuer:    global.GVA_CONFIG.JWT.Issuer,              // 签名的发行者
		},
	}
	return claims
}

// 创建一个token
func (jwtUtility *JWT) CreateToken(claims request.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtUtility.SigningKey)
}

// CreateTokenByOldToken 旧token 换新token 使用归并回源避免并发问题
func (jwtUtility *JWT) CreateTokenByOldToken(oldToken string, claims request.CustomClaims) (string, error) {
	v, err, _ := global.GVA_Concurrency_Control.Do("JWT:"+oldToken, func() (interface{}, error) {
		return jwtUtility.CreateToken(claims)
	})
	return v.(string), err
}

// 解析 token
func (jwtUtility *JWT) ParseToken(tokenString string) (*request.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &request.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return jwtUtility.SigningKey, nil
	})

	if err != nil {
		Logger.Fatal("请检查请求头是否存在x-token且claims是否为规定结构")
	}

	if token != nil {
		if claims, ok := token.Claims.(*request.CustomClaims); ok && token.Valid {
			issuer := zap.String("issuer", claims.RegisteredClaims.Issuer)
			Logger.Info("来自 {} 的token验证成功", issuer)
			return claims, nil
		}
		Logger.Error("未知的claim类型，token验证失败")
		return nil, TokenInvalid
	} else {
		Logger.Error("token为空，验证失败")
		return nil, TokenInvalid
	}
}

func (jwtUtility *JWT) ValidateJWT(tokenString string) bool {

	result := true
	if _, err := jwtUtility.ParseToken(tokenString); err != nil {
		Logger.Error("token验证失败")
		result = false
	}
	Logger.Error("token验证成功")
	return result
}
