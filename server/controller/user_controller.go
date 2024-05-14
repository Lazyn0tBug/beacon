package controller

import (
	"net/http"

	"net/http"

	"github.com/Lazyn0tBug/beacon/server/model"
	"github.com/Lazyn0tBug/beacon/server/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func RegisterUser(c *gin.Context) {
	var user model.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.RegisterUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func LoginUser(c *gin.Context) {
	// 省略验证逻辑和JWT生成
	// ...
}

func (uc *UserController) CreateUser(user *model.User) error {
	err := uc.DB.Create(user).Error
	return err
}

func (uc *UserController) GetUser(id uint) (*model.User, error) {
	var user model.User
	err := uc.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (uc *UserController) UpdateUser(user *model.User) error {
	err := uc.DB.Save(user).Error
	return err
}

func (uc *UserController) DeleteUser(id uint) error {
	user := model.User{}
	err := uc.DB.Delete(&user, id).Error
	return err
}
