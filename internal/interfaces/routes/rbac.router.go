package routes

import (
	"gin-starter/internal/interfaces/handlers"

	"github.com/gin-gonic/gin"
)

type RBACRouter struct{}

func NewRBACRouter() *RBACRouter {
	return &RBACRouter{}
}

func (rr *RBACRouter) RegisterRoutes(router *gin.RouterGroup) {
	rbacGroup := router.Group("/rbac")
	{
		rbacGroup.POST("/policy", handlers.AddPolicy)
		rbacGroup.POST("/role", handlers.AddRoleForUser)
		rbacGroup.GET("/roles/:user_id", handlers.GetRolesForUser)
		rbacGroup.POST("/department", handlers.AddDepartmentForUser)
		rbacGroup.GET("/departments/:user_id", handlers.GetDepartmentsForUser)
		rbacGroup.POST("/enforce", handlers.RBACEnforce)
	}
}
