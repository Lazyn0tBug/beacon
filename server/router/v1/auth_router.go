package v1

import (
	"github.com/Lazyn0tBug/beacon/server/controller"
	"github.com/gin-gonic/gin"
)

func InitAuthRouter(r *gin.Engine) {
	authRouter := r.Group("/v1/auth")
	authController := controller.AuthController{}
	{
		authRouter.POST("/login", func(ctx *gin.Context) {
			authController.Login(ctx)
		})
		authRouter.POST("/logout", func(ctx *gin.Context) {
			authController.Logout(ctx)
		})
	}
}
