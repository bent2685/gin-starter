package routes

import (
	"gin-starter/internal/interfaces/handlers"

	"github.com/gin-gonic/gin"
)

type DepartmentRouter struct {
	deptHandler handlers.DepartmentHandler
}

func NewDepartmentRouter() *DepartmentRouter {
	return &DepartmentRouter{
		deptHandler: *handlers.NewDepartmentHandler(),
	}
}

func (dr *DepartmentRouter) RegisterRoutes(router *gin.RouterGroup) {
	departmentGroup := router.Group("/departments")
	{
		departmentGroup.POST("", dr.deptHandler.CreateDepartment)
		departmentGroup.GET("", dr.deptHandler.GetAllDepartments)
		departmentGroup.GET("/:id", dr.deptHandler.GetDepartment)
		departmentGroup.PUT("/:id", dr.deptHandler.UpdateDepartment)
		departmentGroup.DELETE("/:id", dr.deptHandler.DeleteDepartment)
		departmentGroup.GET("/tree", dr.deptHandler.GetDepartmentTree)
	}
}
