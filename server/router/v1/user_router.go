package v1

import (
	"github.com/Lazyn0tBug/beacon/server/controller"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(r *gin.Engine) {
	userRouter := r.Group("v1/user")
	userController := controller.UserController{}
	{
		userRouter.POST("/register", func(ctx *gin.Context) {
			userController.Register(ctx)
		})
		userRouter.GET("/user/list", func(ctx *gin.Context) {
			userController.GetUsersWithPagination(ctx)
		})
		userRouter.GET("/user/:id", func(ctx *gin.Context) {
			userController.GetUserByID(ctx)
		})
		userRouter.PUT("/user/:id", func(ctx *gin.Context) {
			userController.UpdateUser(ctx)
		})
		userRouter.DELETE("/user/:id", func(ctx *gin.Context) {
			userController.DeleteUser(ctx)
		})
	}
}
