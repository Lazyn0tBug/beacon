package dao

import (
	"errors"

	"github.com/Lazyn0tBug/beacon/server/database"
	"github.com/Lazyn0tBug/beacon/server/model"
	"github.com/Lazyn0tBug/beacon/server/utils"
	"gorm.io/gorm"
)

// Updated UserDAOInterface with new methods
type UserDAOInterface interface {
	GetByUsername(username string) (*model.User, error)
	GetUserByUserID(userID uint64) (*model.User, error) // 获取指定ID的用户
	Create(user *model.User) error
	UpdatePassword(user *model.User, oldPassword string, newPassword string) error
	Update(user *model.User) error                                         // 更新用户信息
	Delete(userID uint) error                                              // 删除用户
	Activate(userID uint) error                                            // 激活用户
	Suspend(userID uint) error                                             // 冻结用户
	GetUsersWithPagination(pageNumber, pageSize int) ([]model.User, error) // 分页读取用户列表
	VerifyUser(username, password string) (*model.User, error)             // 验证用户
}

type UserDAO struct {
	dbInterface database.DBInterface
}

var ErrUserNotFound = errors.New("user not found")

func NewUserDao(dbInterface database.DBInterface) *UserDAO {
	db := dbInterface.GetDB()
	// 注册BeforeSave钩子
	db.Callback().Create().Before("gorm:save_hook").Register("hash_password", hashPasswordBeforeSave)
	db.Callback().Update().Before("gorm:save_hook").Register("hash_password", hashPasswordBeforeSave)

	return &UserDAO{dbInterface}
}

func (userDao *UserDAO) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
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
func (userDAO *UserDAO) Create(user *model.User) error {
	user.IsActive = 1 // Set the user as active after registration
	result := userDAO.dbInterface.GetDB().Create(user)
	return result.Error
}

// Update updates an existing user.
func (userDAO *UserDAO) UpdatePassword2(userID int64, newPassword string) error {
	result := userDAO.dbInterface.GetDB().Model(&model.User{}).Where("id = ?", userID).Update("password", newPassword)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Update updates an existing user.
func (userDAO *UserDAO) Update(user *model.User) error {
	// 创建一个map来存储要更新的字段，排除password字段
	updates := make(map[string]interface{})
	// 假设User结构体中有Email, Username等字段需要更新，但不包括Password
	updates["email"] = user.Email
	updates["username"] = user.Username
	// 使用Updates方法更新指定字段，排除password字段
	result := userDAO.dbInterface.GetDB().Model(user).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Delete deletes a user by ID.
func (userDAO *UserDAO) Delete(userID uint64) error {
	user := &model.User{}
	user.ID = userID
	result := userDAO.dbInterface.GetDB().Delete(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Activate sets a user's status to active.
func (userDAO *UserDAO) Activate(userID uint64) error {
	user := &model.User{}
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
	user := &model.User{}
	user.ID = userID  // Fixed: assign the ID after creating the struct instance
	user.IsActive = 0 // Assuming 0 represents suspended status
	result := userDAO.dbInterface.GetDB().Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (userDAO *UserDAO) GetUserByUserID(userID uint) (*model.User, error) {
	user := &model.User{}
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
func (userDAO *UserDAO) GetUsersWithPagination(pageNumber, pageSize int) ([]model.User, error) {
	var users []model.User
	offset := (pageNumber - 1) * pageSize
	result := userDAO.dbInterface.GetDB().Limit(pageSize).Offset(offset).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// ... (other methods and NewUserDao function omitted for brevity)

// VerifyUser verifies the username and password, returning the matching user (if exists and the password is correct).
func (userDAO *UserDAO) VerifyUser(username, password string) (*model.User, error) {
	// Query the database to find the user
	var user model.User
	result := userDAO.dbInterface.GetDB().Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// User does not exist
			return nil, errors.New("user not found")
		}
		// Other database errors
		return nil, result.Error
	}

	// Verify if the password is correct
	if !user.CheckPassword(password) {
		// Incorrect password
		return nil, errors.New("incorrect password")
	}

	// User exists and password is correct, return the user pointer
	return &user, nil
}

// hashPasswordBeforeSave 是在保存用户之前调用的钩子，用于哈希密码
func hashPasswordBeforeSave(db *gorm.DB) {
	if db.Error != nil {
		return
	}

	if db.Statement.Changed("password") {
		newPassword, res := db.Statement.Get("password")
		if res == false {
			return
		}

		hashedPassword, err := utils.BcryptHash(newPassword.(string))
		if err != nil {
			db.AddError(err)
			return
		}
		db.Statement.SetColumn("password", string(hashedPassword))
	}
}

// UpdatePassword 更新用户密码，验证旧密码并自动应用bcrypt哈希（钩子已处理哈希）
func (userDAO *UserDAO) UpdatePassword(userID uint64, oldPassword string, newPassword string) error {
	// 开启事务
	tx := userDAO.dbInterface.GetDB().Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 获取用户以验证旧密码
	var user model.User
	if err := tx.First(&user, "id = ?", userID).Error; err != nil {
		tx.Rollback() // 如果出现其他错误，也回滚事务
		return err
	}

	// 验证旧密码是否正确
	if !user.CheckPassword(oldPassword) {
		tx.Rollback() // 如果旧密码不正确，回滚事务
		return errors.New("incorrect old password")
	}

	// 更新数据库中的密码字段，由于钩子会自动哈希新密码，所以直接保存即可
	if err := tx.Model(&user).Update("password", newPassword).Error; err != nil {
		tx.Rollback() // 如果更新失败，回滚事务
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return err // 如果提交事务失败，返回错误
	}

	return nil
}
