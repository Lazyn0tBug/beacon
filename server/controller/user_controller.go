package controller

import (
	"net/http"
	"strconv"

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

// GetUserByusername 根据用户名获取用户信息
func (userController *UserController) GetUserByName(c *gin.Context) {
	username := c.Param("username")
	user, err := userController.userService.GetUserByName(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetUserByID 根据用户ID获取用户信息
func (userController *UserController) GetUserByID(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	user, err := userController.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser 更新用户信息
func (userController *UserController) UpdateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user.ID = userID
	if err := userController.userService.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteUser 删除用户
func (userController *UserController) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	if err := userController.userService.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func (userController *UserController) ActivateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := userController.userService.ActivateUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Activate user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User Activated"})
}

func (userController *UserController) SuspendUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := userController.userService.SuspendUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Suspend user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User Suspended"})
}

func (userController *UserController) GetUsersWithPagination(c *gin.Context) {
	pageNumber, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pagesize", "10"))
	if pageSize <= 0 {
		pageSize = 10
	}

	users, err := userController.userService.GetUsersWithPagination(pageNumber, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
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

func (userController *UserController) ChangePassword(c *gin.Context) {
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

	if err := userController.userService.ChangePassword(user, changePasswordRequest.OldPassword, changePasswordRequest.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}
