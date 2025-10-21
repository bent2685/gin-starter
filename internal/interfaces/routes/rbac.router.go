package routes

import (
	"gin-starter/internal/application/services/rbac"
	"gin-starter/internal/interfaces/handlers"

	"github.com/gin-gonic/gin"
)

type RBACRouter struct {
	rbacHandler *handlers.RBACHandler
}

// NewRBACRouter 创建RBAC路由实例
func NewRBACRouter(rbacService *rbac.RBACService) *RBACRouter {
	return &RBACRouter{
		rbacHandler: handlers.NewRBACHandler(rbacService),
	}
}

// RegisterRoutes 注册路由
func (rr *RBACRouter) RegisterRoutes(router *gin.RouterGroup) {
	rbacGroup := router.Group("/rbac")
	{
		rbacGroup.POST("/policy", rr.rbacHandler.AddPolicy)
		rbacGroup.POST("/role", rr.rbacHandler.AddRoleForUser)
		rbacGroup.GET("/roles/:user_id", rr.rbacHandler.GetRolesForUser)
		rbacGroup.POST("/enforce", rr.rbacHandler.Enforce)
	}
}
