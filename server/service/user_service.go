package service

import (
	"errors"
	"fmt"

	"github.com/Lazyn0tBug/beacon/server/dao"
	"github.com/Lazyn0tBug/beacon/server/model"
	"github.com/gin-contrib/sessions"
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

func (s *UserService) Register(userData *model.User) error {
	// 检查用户名是否存在
	_, err := s.userDAO.GetByUsername(userData.Username)
	if err == nil {
		return ErrUserAlreadyExists
	}
	if !errors.Is(err, dao.ErrUserNotFound) {
		return err // 返回非“用户不存在”的错误
	}

	// 调用DAO层创建用户
	if err := s.userDAO.Create(userData); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// Login 用户登录
// func (userService *UserService) Login(username, password string) (string, error) {
// 	user, err := s.userDAO.VerifyUser(username, password)
// 	if err != nil {
// 		return "", err
// 	}
// 	token, err := utils.GenerateJWT(user.ID)
// 	if err != nil {
// 		return "", err
// 	}
// 	return token, nil
// }

func (userService *UserService) Login(username, password string) (*model.User, error) {
	user, err := userService.userDAO.VerifyUser(username, password)
	if err != nil {
		return nil, err
	}
	if !user.CheckPassword(password) {
		return nil, errors.New("invalid username or password")
	}
	return user, nil
}

func (userService *UserService) Logout(session sessions.Session) error {
	// 清除用户会话信息
	session.Clear()
	err := session.Save(r.Context())
	if err != nil {
		return err
	}
	return nil
}

func (userService *UserService) ChangePassword(user *model.User, oldPassword string, newPassword string) error {
	if !user.CheckPassword(oldPassword) {
		return errors.New("old password does not match")
	}
	if err := userService.userDAO.UpdatePassword(user, oldPassword, newPassword); err != nil {
		return err
	}
	return nil
}
