package routes

import (
	"gin-starter/internal/interfaces/handlers"

	"github.com/gin-gonic/gin"
)

type DepartmentRouter struct{}

func NewDepartmentRouter() *DepartmentRouter {
	return &DepartmentRouter{}
}

func (dr *DepartmentRouter) RegisterRoutes(router *gin.RouterGroup) {
	departmentGroup := router.Group("/departments")
	{
		departmentGroup.POST("", handlers.CreateDepartment)
		departmentGroup.GET("", handlers.GetAllDepartments)
		departmentGroup.GET("/:id", handlers.GetDepartment)
		departmentGroup.PUT("/:id", handlers.UpdateDepartment)
		departmentGroup.DELETE("/:id", handlers.DeleteDepartment)
		departmentGroup.GET("/tree", handlers.GetDepartmentTree)
	}
}
