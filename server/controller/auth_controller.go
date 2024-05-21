package controller

import (
	"net/http"
	"time"

	"github.com/Lazyn0tBug/beacon/server/global"
	"github.com/Lazyn0tBug/beacon/server/service"
	system_service "github.com/Lazyn0tBug/beacon/server/service/system"
	"github.com/Lazyn0tBug/beacon/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type AuthController struct {
	userService *service.UserService
}

func NewAuthController(userService *service.UserService) *AuthController {
	return &AuthController{userService}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (authController *AuthController) Login(c *gin.Context) {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := authController.userService.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})
	tokenString, err := token.SignedString([]byte("secret")) // You should use a more secure secret in production

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (authController *AuthController) Logout(c *gin.Context) {
	// Nothing to do here as JWT tokens are stateless
	token := utils.GetToken(c)
	jwt := system_service.JwtBlacklist{Jwt: token}
	error := jwtService.JsonInBlacklist(jwt)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	if err != nil {
		global.GVA_LOG.Error("jwt作废失败!", zap.Error(err))
		response.FailWithMessage("jwt作废失败", c)
		return
	}
	utils.ClearToken(c)
	response.OkWithMessage("jwt作废成功", c)
}
