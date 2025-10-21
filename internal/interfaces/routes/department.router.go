package routes

import (
	"gin-starter/internal/application/services"
	"gin-starter/internal/interfaces/handlers"

	"github.com/gin-gonic/gin"
)

type DepartmentRouter struct {
	departmentHandler *handlers.DepartmentHandler
}

// NewDepartmentRouter 创建部门路由实例
func NewDepartmentRouter() *DepartmentRouter {
	// 创建部门服务
	departmentService := services.NewDepartmentService()

	// 创建部门处理器
	departmentHandler := handlers.NewDepartmentHandler(departmentService)

	return &DepartmentRouter{
		departmentHandler: departmentHandler,
	}
}

// RegisterRoutes 注册路由
func (dr *DepartmentRouter) RegisterRoutes(router *gin.RouterGroup) {
	departmentGroup := router.Group("/departments")
	{
		departmentGroup.POST("", dr.departmentHandler.CreateDepartment)
		departmentGroup.GET("", dr.departmentHandler.GetAllDepartments)
		departmentGroup.GET("/:id", dr.departmentHandler.GetDepartment)
		departmentGroup.PUT("/:id", dr.departmentHandler.UpdateDepartment)
		departmentGroup.DELETE("/:id", dr.departmentHandler.DeleteDepartment)
		departmentGroup.GET("/tree", dr.departmentHandler.GetDepartmentTree)
	}
}
