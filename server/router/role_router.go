package router

import (
	"github.com/Lazyn0tBug/beacon/server/controller"
	"github.com/gin-gonic/gin"
)

func InitRoleRouter(r *gin.Engine) {
	roleRouter := r.Group("/auth")
	userController := controller.UserController{}
	authController := controller.AuthController{}
	{
		roleRouter.POST("/register", func(ctx *gin.Context) {
			userController.Register(ctx)
		})
		roleRouter.POST("/login", func(ctx *gin.Context) {
			authController.Login(ctx)
		})
		roleRouter.POST("/logout", func(ctx *gin.Context) {
			authController.Logout(ctx)
		})
		roleRouter.GET("/user/:id", func(ctx *gin.Context) {
			userController.GetUser(ctx)
		})
	}
}
