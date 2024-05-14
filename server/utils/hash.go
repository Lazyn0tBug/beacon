package utils

import (
	"crypto/sha256"
	"encoding/hex"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// BcryptHash 使用 bcrypt 对密码进行加密
func BcryptHash(password string) (string, error) {
	logger := GetLogger()
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("密码哈希失败", zap.Error(err))
		return "", err
	}

	hashPassword := string(bytes)
	logger.Info("密码哈希完成", zap.String("hashedPassword", hashPassword))
	return hashPassword, nil
}

// BcryptCheck 对比明文密码和数据库的哈希值
func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// @author: [Lazyn0tBug](https://github.com/Lazyn0tBug)
// @function: SHA256V
// @description: sha256加密
// @param: str []byte
// @return: string
func SHA256V(str []byte) string {
	hash := sha256.Sum256(str)
	return hex.EncodeToString(hash[:])
}
