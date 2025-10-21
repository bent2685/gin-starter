package routes

import (
	"gin-starter/internal/application/services/rbac"
	"gin-starter/internal/middleware"
	"gin-starter/pkg/utils/res"

	"github.com/gin-gonic/gin"
)

type ProtectedRouter struct {
	rbacService *rbac.RBACService
}

// NewProtectedRouter 创建受保护路由实例
func NewProtectedRouter(rbacService *rbac.RBACService) *ProtectedRouter {
	return &ProtectedRouter{
		rbacService: rbacService,
	}
}

// RegisterRoutes 注册路由
func (pr *ProtectedRouter) RegisterRoutes(router *gin.RouterGroup) {
	// 受保护的路由组，需要认证
	protectedGroup := router.Group("/protected")
	protectedGroup.Use(middleware.AuthMiddleware()) // 使用认证中间件
	{
		protectedGroup.GET("/data", pr.getProtectedData)
		protectedGroup.POST("/data", pr.createProtectedData)
	}

	// 管理员路由组，需要认证和管理员权限
	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware())                        // 使用认证中间件
	adminGroup.Use(middleware.RoleMiddleware(pr.rbacService, "admin")) // 使用角色中间件
	{
		adminGroup.GET("/users", pr.getUsers)
		adminGroup.DELETE("/users/:id", pr.deleteUser)
	}

	// 部门路由组，需要认证和特定部门权限
	departmentGroup := router.Group("/department")
	departmentGroup.Use(middleware.AuthMiddleware())                           // 使用认证中间件
	departmentGroup.Use(middleware.DepartmentMiddleware(pr.rbacService, "IT")) // 使用部门中间件
	{
		departmentGroup.GET("/resources", pr.getDepartmentResources)
	}
}

// getProtectedData 获取受保护的数据
func (pr *ProtectedRouter) getProtectedData(c *gin.Context) {
	// 从上下文中获取用户信息
	userID := c.MustGet("user_id").(uint)
	username := c.MustGet("username").(string)

	data := map[string]interface{}{
		"message":  "这是受保护的数据",
		"user_id":  userID,
		"username": username,
	}

	res.Success(c, data)
}

// createProtectedData 创建受保护的数据
func (pr *ProtectedRouter) createProtectedData(c *gin.Context) {
	// 从上下文中获取用户信息
	userID := c.MustGet("user_id").(uint)
	username := c.MustGet("username").(string)

	// 这里可以处理创建数据的逻辑

	data := map[string]interface{}{
		"message":  "数据创建成功",
		"user_id":  userID,
		"username": username,
	}

	res.Success(c, data)
}

// getUsers 获取用户列表（仅管理员）
func (pr *ProtectedRouter) getUsers(c *gin.Context) {
	// 这里可以实现获取用户列表的逻辑
	data := map[string]interface{}{
		"message": "管理员访问成功",
		"users":   []string{"user1", "user2", "user3"},
	}

	res.Success(c, data)
}

// deleteUser 删除用户（仅管理员）
func (pr *ProtectedRouter) deleteUser(c *gin.Context) {
	// 从URL参数中获取用户ID
	userID := c.Param("id")

	// 这里可以实现删除用户的逻辑

	data := map[string]interface{}{
		"message": "用户删除成功",
		"user_id": userID,
	}

	res.Success(c, data)
}

// getDepartmentResources 获取部门资源（仅IT部门）
func (pr *ProtectedRouter) getDepartmentResources(c *gin.Context) {
	// 从上下文中获取用户信息
	userID := c.MustGet("user_id").(uint)
	username := c.MustGet("username").(string)

	data := map[string]interface{}{
		"message":   "IT部门资源访问成功",
		"user_id":   userID,
		"username":  username,
		"resources": []string{"server1", "server2", "database1"},
	}

	res.Success(c, data)
}
