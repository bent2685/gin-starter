package middleware

import (
	"gin-starter/internal/application/services/rbac"
	"gin-starter/pkg/utils/jwt"
	"gin-starter/pkg/utils/res"
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

// AuthorizationMiddleware 基于策略的权限中间件
func AuthorizationMiddleware(rbacService *rbac.RBACService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户信息
		userID, exists := c.Get("user_id")
		if !exists {
			res.ErrUnauthorized.ThrowWithMessage(c, "用户未认证")
			return
		}

		// 构造资源和操作
		resource := c.Request.URL.Path
		action := c.Request.Method

		// 首先检查用户是否具有super_admin权限（通过Casbin策略验证）
		allowed, err := rbacService.Enforce(rbac.GetUserID(userID.(uint)), "*", "*")
		if err != nil {
			res.ErrInternalServer.ThrowWithMessage(c, "权限验证失败")
			return
		}

		// 如果用户具有super_admin权限，允许访问所有资源
		if allowed {
			c.Next()
			return
		}

		// 获取用户的角色
		roles, err := rbacService.GetRolesForUser(rbac.GetUserID(userID.(uint)))
		if err != nil {
			res.ErrInternalServer.ThrowWithMessage(c, "角色获取失败")
			return
		}

		// 检查用户的角色权限
		allowed = false
		for _, role := range roles {
			// 验证权限
			roleAllowed, err := rbacService.Enforce(role, resource, action)
			if err != nil {
				res.ErrInternalServer.ThrowWithMessage(c, "权限验证失败")
				return
			}
			if roleAllowed {
				allowed = true
				break
			}
		}

		// 如果角色没有权限，检查用户个人权限
		if !allowed {
			allowed, err = rbacService.Enforce(rbac.GetUserID(userID.(uint)), resource, action)
			if err != nil {
				res.ErrInternalServer.ThrowWithMessage(c, "权限验证失败")
				return
			}
		}

		if !allowed {
			res.ErrForbidden.ThrowWithMessage(c, "权限不足")
			return
		}

		c.Next()
	}
}

// RoleMiddleware 基于角色的权限中间件
func RoleMiddleware(rbacService *rbac.RBACService, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户信息
		userID, exists := c.Get("user_id")
		if !exists {
			res.ErrUnauthorized.ThrowWithMessage(c, "用户未认证")
			return
		}

		// 首先检查用户是否具有super_admin权限（通过Casbin策略验证）
		allowed, err := rbacService.Enforce(rbac.GetUserID(userID.(uint)), "*", "*")
		if err != nil {
			res.ErrInternalServer.ThrowWithMessage(c, "权限验证失败")
			return
		}

		// 如果用户具有super_admin权限，允许访问
		if allowed {
			c.Next()
			return
		}

		// 获取用户的角色
		roles, err := rbacService.GetRolesForUser(rbac.GetUserID(userID.(uint)))
		if err != nil {
			res.ErrInternalServer.ThrowWithMessage(c, "角色获取失败")
			return
		}

		// 检查用户是否具有所需角色
		hasRole := false
		for _, role := range roles {
			if role == requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			res.ErrForbidden.ThrowWithMessage(c, "权限不足")
			return
		}

		c.Next()
	}
}

// DepartmentMiddleware 基于部门的权限中间件
func DepartmentMiddleware(rbacService *rbac.RBACService, requiredDepartment string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户信息
		userID, exists := c.Get("user_id")
		if !exists {
			res.ErrUnauthorized.ThrowWithMessage(c, "用户未认证")
			return
		}

		// 首先检查用户是否具有super_admin权限（通过Casbin策略验证）
		allowed, err := rbacService.Enforce(rbac.GetUserID(userID.(uint)), "*", "*")
		if err != nil {
			res.ErrInternalServer.ThrowWithMessage(c, "权限验证失败")
			return
		}

		// 如果用户具有super_admin权限，允许访问
		if allowed {
			c.Next()
			return
		}

		// 获取用户的部门
		departments, err := rbacService.GetDepartmentsForUser(rbac.GetUserID(userID.(uint)))
		if err != nil {
			res.ErrInternalServer.ThrowWithMessage(c, "部门获取失败")
			return
		}

		// 检查用户是否属于所需部门
		inDepartment := false
		for _, department := range departments {
			if department == requiredDepartment {
				inDepartment = true
				break
			}
		}

		if !inDepartment {
			res.ErrForbidden.ThrowWithMessage(c, "权限不足")
			return
		}

		c.Next()
	}
}
