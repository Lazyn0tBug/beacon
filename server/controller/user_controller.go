package controller

import (
	"net/http"

	"github.com/Lazyn0tBug/beacon/server/model"
	"github.com/Lazyn0tBug/beacon/server/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.UserService
}

func (uc *UserController) RegisterUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := uc.userService.Register(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (uc *UserController) LoginUser(c *gin.Context) {
	// 省略验证逻辑和JWT生成
	// ...
}

// func (uc *UserController) GetUser(id uint) (*model.User, error) {
// 	var user model.User
// 	err := uc.DB.First(&user, id).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }

// func (uc *UserController) UpdateUser(user *model.User) error {
// 	err := uc.DB.Save(user).Error
// 	return err
// }

// func (uc *UserController) DeleteUser(id uint) error {
// 	user := model.User{}
// 	err := uc.DB.Delete(&user, id).Error
// 	return err
// }
