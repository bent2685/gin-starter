package routes

import (
	"gin-starter/internal/interfaces/handlers"

	"github.com/gin-gonic/gin"
)

type RBACRouter struct {
	rabcHandler handlers.RBACHandler
}

func NewRBACRouter() *RBACRouter {
	return &RBACRouter{
		rabcHandler: *handlers.NewRBACHandler(),
	}
}

func (rr *RBACRouter) RegisterRoutes(router *gin.RouterGroup) {
	rbacGroup := router.Group("/rbac")
	{
		rbacGroup.POST("/policy", rr.rabcHandler.AddPolicy)
		rbacGroup.POST("/role", rr.rabcHandler.AddRoleForUser)
		rbacGroup.GET("/roles/:user_id", rr.rabcHandler.GetRolesForUser)
		rbacGroup.POST("/department", rr.rabcHandler.AddDepartmentForUser)
		rbacGroup.GET("/departments/:user_id", rr.rabcHandler.GetDepartmentsForUser)
		rbacGroup.POST("/enforce", rr.rabcHandler.RBACEnforce)
	}
}
