package routes

import (
	"gin-starter/internal/interfaces/handlers"

	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func NewUserRouter() *UserRouter {
	return &UserRouter{}
}

func (ur *UserRouter) RegisterRoutes(router *gin.RouterGroup) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("", handlers.CreateUser)
		userGroup.GET("", handlers.GetAllUsers)
		userGroup.GET("/:id", handlers.GetUser)
		userGroup.PUT("/:id", handlers.UpdateUser)
		userGroup.DELETE("/:id", handlers.DeleteUser)
		userGroup.POST("/:id/activate", handlers.ActivateUser)
		userGroup.POST("/:id/deactivate", handlers.DeactivateUser)
		userGroup.PUT("/:id/email", handlers.ChangeEmail)
		userGroup.POST("/login", handlers.Login)
	}
}
