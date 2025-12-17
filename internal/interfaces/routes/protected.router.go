package routes

import (
	"gin-starter/internal/application/services/rbac"
	"gin-starter/internal/middleware"
	"gin-starter/pkg/utils/res"

	"github.com/gin-gonic/gin"
)

type ProtectedRouter struct{}

func NewProtectedRouter() *ProtectedRouter {
	return &ProtectedRouter{}
}

func (pr *ProtectedRouter) RegisterRoutes(router *gin.RouterGroup) {
	protectedGroup := router.Group("/protected")
	protectedGroup.Use(middleware.AuthMiddleware())
	{
		protectedGroup.GET("/data", getProtectedData)
		protectedGroup.POST("/data", createProtectedData)
	}

	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware())
	adminGroup.Use(middleware.RoleMiddleware(rbac.GetRBAC(), "admin"))
	{
		adminGroup.GET("/users", getUsers)
		adminGroup.DELETE("/users/:id", deleteUser)
	}

	departmentGroup := router.Group("/department")
	departmentGroup.Use(middleware.AuthMiddleware())
	departmentGroup.Use(middleware.DepartmentMiddleware(rbac.GetRBAC(), "IT"))
	{
		departmentGroup.GET("/resources", getDepartmentResources)
	}
}

func getProtectedData(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	username := c.MustGet("username").(string)
	res.Success(c, map[string]interface{}{
		"message":  "这是受保护的数据",
		"user_id":  userID,
		"username": username,
	})
}

func createProtectedData(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	username := c.MustGet("username").(string)
	res.Success(c, map[string]interface{}{
		"message":  "数据创建成功",
		"user_id":  userID,
		"username": username,
	})
}

func getUsers(c *gin.Context) {
	res.Success(c, map[string]interface{}{
		"message": "管理员访问成功",
		"users":   []string{"user1", "user2", "user3"},
	})
}

func deleteUser(c *gin.Context) {
	userID := c.Param("id")
	res.Success(c, map[string]interface{}{
		"message": "用户删除成功",
		"user_id": userID,
	})
}

func getDepartmentResources(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	username := c.MustGet("username").(string)
	res.Success(c, map[string]interface{}{
		"message":   "IT部门资源访问成功",
		"user_id":   userID,
		"username":  username,
		"resources": []string{"server1", "server2", "database1"},
	})
}
