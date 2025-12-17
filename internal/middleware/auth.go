package middleware

import (
	"gin-starter/internal/application/services/rbac"
	"gin-starter/pkg/utils/jwt"
	"gin-starter/pkg/utils/res"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取Authorization字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			res.ErrTokenRequired.ThrowWithMessage(c, "未提供认证信息")
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			res.ErrTokenFormat.ThrowWithMessage(c, "认证信息格式错误")
			return
		}

		// 提取Token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 解析Token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			res.ErrInvalidToken.ThrowWithMessage(c, "无效的认证令牌")
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		// 继续处理请求
		c.Next()
	}
}

// RoleMiddleware 验证用户身份中间件
func RoleMiddleware(rbacService *rbac.RBACService, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			res.ErrUnauthorized.ThrowWithMessage(c, "用户未认证")
			return
		}
		allowed, err := rbac.Enforce(rbac.GetUserID(userID.(uint)), "*", "*")
		if err != nil {
			res.ErrInternalServer.ThrowWithMessage(c, "权限验证失败")
			return
		}
		if allowed {
			c.Next()
			return
		}
		roles, err := rbac.GetRolesForUser(rbac.GetUserID(userID.(uint)))
		if err != nil {
			res.ErrInternalServer.ThrowWithMessage(c, "角色获取失败")
			return
		}
		if !slices.Contains(roles, requiredRole) {
			res.ErrForbidden.ThrowWithMessage(c, "权限不足")
			return
		}
		c.Next()
	}
}

// DepartmentMiddleware 验证用户部门中间件
func DepartmentMiddleware(rbacService *rbac.RBACService, requiredDepartment string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			res.ErrUnauthorized.ThrowWithMessage(c, "用户未认证")
			return
		}
		allowed, err := rbac.Enforce(rbac.GetUserID(userID.(uint)), "*", "*")
		if err != nil {
			res.ErrInternalServer.ThrowWithMessage(c, "权限验证失败")
			return
		}
		if allowed {
			c.Next()
			return
		}
		departments, err := rbac.GetDepartmentsForUser(rbac.GetUserID(userID.(uint)))
		if err != nil {
			res.ErrInternalServer.ThrowWithMessage(c, "部门获取失败")
			return
		}
		if !slices.Contains(departments, requiredDepartment) {
			res.ErrForbidden.ThrowWithMessage(c, "权限不足")
			return
		}
		c.Next()
	}
}
