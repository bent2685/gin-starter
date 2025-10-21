package routes

import (
	"gin-starter/internal/interfaces/handlers"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	userHandler *handlers.UserHandler
}

// NewUserRouter 创建用户路由实例
func NewUserRouter() *UserRouter {
	return &UserRouter{
		userHandler: handlers.NewUserHandler(),
	}
}

// RegisterRoutes 注册路由
func (ur *UserRouter) RegisterRoutes(router *gin.RouterGroup) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("", ur.userHandler.CreateUser)
		userGroup.GET("", ur.userHandler.GetAllUsers)
		userGroup.GET("/:id", ur.userHandler.GetUser)
		userGroup.PUT("/:id", ur.userHandler.UpdateUser)
		userGroup.DELETE("/:id", ur.userHandler.DeleteUser)

		// 用户状态管理
		userGroup.POST("/:id/activate", ur.userHandler.ActivateUser)
		userGroup.POST("/:id/deactivate", ur.userHandler.DeactivateUser)

		// 用户信息更新
		userGroup.PUT("/:id/email", ur.userHandler.ChangeEmail)

		// 用户认证
		userGroup.POST("/login", ur.userHandler.Login)
	}
}
