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
	// 创建RBAC处理器
	rbacHandler := handlers.NewRBACHandler(rbacService)

	return &RBACRouter{
		rbacHandler: rbacHandler,
	}
}

// RegisterRoutes 注册路由
func (rr *RBACRouter) RegisterRoutes(router *gin.RouterGroup) {
	rbacGroup := router.Group("/rbac")
	{
		// 策略管理
		rbacGroup.POST("/policy", rr.rbacHandler.AddPolicy)

		// 角色管理
		rbacGroup.POST("/role", rr.rbacHandler.AddRoleForUser)
		rbacGroup.GET("/roles/:user_id", rr.rbacHandler.GetRolesForUser)

		// 部门管理
		rbacGroup.POST("/department", rr.rbacHandler.AddDepartmentForUser)
		rbacGroup.GET("/departments/:user_id", rr.rbacHandler.GetDepartmentsForUser)

		// 权限验证
		rbacGroup.POST("/enforce", rr.rbacHandler.Enforce)
	}
}
