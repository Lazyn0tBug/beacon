package service

import (
	"errors"
	"fmt"

	"github.com/Lazyn0tBug/beacon/server/dao"
	"github.com/Lazyn0tBug/beacon/server/model"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type UserService struct {
	userDAO dao.UserDAOInterface
}

func NewUserService(userDAO dao.UserDAOInterface) *UserService {
	return &UserService{
		userDAO: userDAO,
	}
}

func (s *UserService) Register(userData model.Users) error {
	// 检查用户名是否存在
	_, err := s.userDAO.GetByUsername(userData.Username)
	if err == nil {
		return ErrUserAlreadyExists
	}
	if !errors.Is(err, dao.ErrUserNotFound) {
		return err // 返回非“用户不存在”的错误
	}

	// 调用DAO层创建用户
	if err := s.userDAO.Create(&userData); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
