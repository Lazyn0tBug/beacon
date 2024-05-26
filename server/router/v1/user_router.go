package v1

import (
	"github.com/Lazyn0tBug/beacon/server/controller"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(r *gin.Engine) {
	userRouter := r.Group("v1/user")
	userController := controller.UserController{}
	{
		userRouter.GET("/register", userController.Register)
		userRouter.GET("/list", userController.GetUsersWithPagination)
		userRouter.GET("/:id", userController.GetUserByID)
		userRouter.PUT("/:id", userController.UpdateUser)
		userRouter.DELETE("/:id", userController.DeleteUser)
	}
}
