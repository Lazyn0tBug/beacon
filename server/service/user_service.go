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

// GetByUsername 通过用户名获取用户
func (userService *UserService) GetUserByName(username string) (*model.User, error) {
	return userService.userDAO.GetUserByName(username)
}

// GetUserByUserID 通过用户ID获取用户
func (userService *UserService) GetUserByID(userID uint64) (*model.User, error) {
	return userService.userDAO.GetUserByID(userID)
}

// UpdateUser 更新用户信息
func (userService *UserService) UpdateUser(user *model.User) error {
	return userService.userDAO.Update(user)
}

// DeleteUser 删除用户
func (userService *UserService) DeleteUser(userID uint64) error {
	return userService.userDAO.Delete(userID)
}

// ActivateUser 激活用户
func (userService *UserService) ActivateUser(userID uint64) error {
	return userService.userDAO.Activate(userID)
}

// SuspendUser 冻结用户
func (userService *UserService) SuspendUser(userID uint64) error {
	return userService.userDAO.Suspend(userID)
}

// GetUsersWithPagination 分页读取用户列表
func (userService *UserService) GetUsersWithPagination(pageNumber, pageSize int) ([]model.User, error) {
	return userService.userDAO.GetUsersWithPagination(pageNumber, pageSize)
}

// VerifyUser 验证用户
func (userService *UserService) VerifyUser(username, password string) (*model.User, error) {
	user, err := userService.userDAO.VerifyUser(username, password)
	if err != nil {
		return nil, fmt.Errorf("user verification failed: %w", err)
	}
	return user, nil
}

func (userService *UserService) Register(userData *model.User) error {
	// 检查用户名是否存在
	_, err := userService.userDAO.GetUserByName(userData.Username)
	if err == nil {
		return ErrUserAlreadyExists
	}
	if !errors.Is(err, dao.ErrUserNotFound) {
		return err // 返回非“用户不存在”的错误
	}

	// 调用DAO层创建用户
	if err := userService.userDAO.Create(userData); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (userService *UserService) Login(username, password string) (*model.User, error) {
	user, err := userService.VerifyUser(username, password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (userService *UserService) Logout(session sessions.Session) error {
	// 清除用户会话信息
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
