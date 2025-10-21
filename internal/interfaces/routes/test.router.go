package routes

import (
	"gin-starter/pkg/utils/res"

	"github.com/gin-gonic/gin"
)

type TestRouter struct{}

// NewUserRouter 创建用户路由实例
func NewTestRouter() *TestRouter {
	return &TestRouter{}
}

// RegisterRoutes 注册路由
func (tr *TestRouter) RegisterRoutes(router *gin.RouterGroup) {
	testGroup := router.Group("/test")
	{
		testGroup.GET("/health", tr.Health)
	}
}

// Health 健康检查
func (tr *TestRouter) Health(c *gin.Context) {
	res.SuccessWithMessage(c, "PONG", true)
}
