package dao

import (
	"errors"

	"github.com/Lazyn0tBug/beacon/server/database"
	"github.com/Lazyn0tBug/beacon/server/model"
	"gorm.io/gorm"
)

// Updated UserDAOInterface with new methods
type UserDAOInterface interface {
	GetByUsername(username string) (*model.Users, error)
	GetUserByUserID(userID uint64) (*model.Users, error) // 获取指定ID的用户
	Create(user *model.Users) error
	Update(user *model.Users) error                                         // 更新用户信息
	Delete(userID uint) error                                               // 删除用户
	Activate(userID uint) error                                             // 激活用户
	Suspend(userID uint) error                                              // 冻结用户
	GetUsersWithPagination(pageNumber, pageSize int) ([]model.Users, error) // 分页读取用户列表
}

type UserDAO struct {
	dbInterface database.DBInterface
}

var ErrUserNotFound = errors.New("user not found")

func NewUserDao(db database.DBInterface) *UserDAO {
	return &UserDAO{db}
}

func (userDao *UserDAO) GetUserByUsername(username string) (*model.Users, error) {
	var user model.Users
	result := userDao.dbInterface.GetDB().First(&user, "username = ?", username)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// Create creates a new user.
func (userDAO *UserDAO) Create(user *model.Users) error {
	user.IsActive = 1 // Set the user as active after registration
	result := userDAO.dbInterface.GetDB().Create(user)
	return result.Error
}

// Update updates an existing user.
func (userDAO *UserDAO) Update(user *model.Users) error {
	result := userDAO.dbInterface.GetDB().Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Delete deletes a user by ID.
func (userDAO *UserDAO) Delete(userID uint64) error {
	user := &model.Users{}
	user.ID = userID
	result := userDAO.dbInterface.GetDB().Delete(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Activate sets a user's status to active.
func (userDAO *UserDAO) Activate(userID uint64) error {
	user := &model.Users{}
	user.ID = userID
	user.IsActive = 1 // Assuming 1 represents active status
	result := userDAO.dbInterface.GetDB().Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Suspend sets a user's status to suspended.
func (userDAO *UserDAO) Suspend(userID uint64) error {
	user := &model.Users{}
	user.ID = userID  // Fixed: assign the ID after creating the struct instance
	user.IsActive = 0 // Assuming 0 represents suspended status
	result := userDAO.dbInterface.GetDB().Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (userDAO *UserDAO) GetUserByUserID(userID uint) (*model.Users, error) {
	user := &model.Users{}
	result := userDAO.dbInterface.GetDB().First(user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil if the user is not found
		}
		return nil, result.Error
	}
	return user, nil
}

// GetUsersWithPagination retrieves a paginated list of users.
func (userDAO *UserDAO) GetUsersWithPagination(pageNumber, pageSize int) ([]model.Users, error) {
	var users []model.Users
	offset := (pageNumber - 1) * pageSize
	result := userDAO.dbInterface.GetDB().Limit(pageSize).Offset(offset).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// ... (other methods and NewUserDao function omitted for brevity)
