package routes

import (
	"gin-starter/internal/interfaces/dto"
	"gin-starter/internal/interfaces/handlers"
	"gin-starter/internal/middleware"

	"github.com/gin-gonic/gin"
)

// UserRouter 使用中间件简化处理的路由
type UserRouter struct {
	userHandler handlers.UserHandler
}

func NewUserRouter() *UserRouter {
	return &UserRouter{
		userHandler: *handlers.NewUserHandler(),
	}
}

func (ur *UserRouter) RegisterRoutes(router *gin.RouterGroup) {
	userGroup := router.Group("/users")
	{
		// 使用中间件自动绑定和验证请求
		userGroup.POST("",
			middleware.BindRequest(&dto.CreateUserRequest{}),
			ur.userHandler.CreateUser,
		)

		userGroup.GET("", ur.userHandler.GetAllUsers)

		userGroup.GET("/:id", ur.userHandler.GetUser)

		userGroup.PUT("/:id",
			middleware.BindRequest(&dto.UpdateUserRequest{}),
			ur.userHandler.UpdateUser,
		)

		userGroup.DELETE("/:id", ur.userHandler.DeleteUser)

		userGroup.POST("/:id/activate", ur.userHandler.ActivateUser)

		userGroup.POST("/:id/deactivate", ur.userHandler.DeactivateUser)

		userGroup.PUT("/:id/email",
			middleware.BindRequest(&dto.ChangeEmailRequest{}),
			ur.userHandler.ChangeEmail,
		)

		userGroup.POST("/login",
			middleware.BindRequest(&dto.LoginRequest{}),
			ur.userHandler.Login,
		)
	}
}
