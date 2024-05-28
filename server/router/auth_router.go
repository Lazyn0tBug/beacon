package router

import (
	"github.com/Lazyn0tBug/beacon/server/controller"
	"github.com/gin-gonic/gin"
)

func InitAuthRouter(r *gin.Engine) {
	authRouter := r.Group("/v1/auth")
	authController := controller.AuthController{}
	{
		authRouter.POST("/login", authController.Login)
		authRouter.POST("/logout", authController.Logout)
	}
}
