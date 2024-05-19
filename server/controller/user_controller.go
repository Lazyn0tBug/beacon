package controller

import (
	"net/http"

	"github.com/Lazyn0tBug/beacon/server/model"
	"github.com/Lazyn0tBug/beacon/server/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	NickName string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"email"`
}

func (userController *UserController) Register(c *gin.Context) {
	var registerRequest RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 假设model.User结构体的字段与RegisterRequest的字段一致或可以通过registerRequest构建
	user := model.User{
		Username: registerRequest.Username,
		Password: registerRequest.Password, // 可能需要进行加密处理
		NickName: registerRequest.NickName,
		Email:    registerRequest.Email,
		RoleID:   uint(1),
		IsActive: int(1),
	}

	if err := userController.userService.Register(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
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

// user_controller.go

func (ctrl *UserController) ChangePassword(c *gin.Context) {
	session := sessions.Default(c)
	user, exists := session.Get("user").(*model.User)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var changePasswordRequest struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&changePasswordRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.userService.ChangePassword(user, changePasswordRequest.OldPassword, changePasswordRequest.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}
