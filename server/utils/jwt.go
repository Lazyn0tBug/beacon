package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenExpiryDuration = time.Hour * 24
)

func GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(TokenExpiryDuration).Unix(),
	})

	// 使用密钥签名token
	// ...

	return tokenString, nil
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	// 解析和验证token
	// ...

	return token, nil
}
