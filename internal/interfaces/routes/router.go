package routes

import (
	"github.com/gin-gonic/gin"
)

// Router 路由管理接口
type Router interface {
	RegisterRoutes(rg *gin.RouterGroup)
}

// RouterManager 路由管理器
type RouterManager struct {
	routers []Router
}

// NewRouterManager 创建路由管理器
func NewRouterManager() *RouterManager {
	return &RouterManager{
		routers: make([]Router, 0),
	}
}

// RegisterRouter 注册路由
func (rm *RouterManager) RegisterRouter(router Router) {
	rm.routers = append(rm.routers, router)
}

// SetupRoutes 设置所有路由
func (rm *RouterManager) SetupRoutes(engine *gin.Engine) {
	// API版本分组
	v1 := engine.Group("")

	// 注册所有路由
	for _, router := range rm.routers {
		router.RegisterRoutes(v1)
	}
}
