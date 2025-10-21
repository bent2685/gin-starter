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

// RoleMiddleware 角色中间件（基于角色的访问控制）
func RoleMiddleware(rbacService *rbac.RBACService, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户ID
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}

		// 获取用户的角色
		roles, err := rbacService.GetRolesForUser(userID)
		if err != nil {
			res.ErrForbidden.ThrowWithMessage(c, "角色验证失败: "+err.Error())
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
			res.ErrForbidden.ThrowWithMessage(c, "角色验证失败，需要角色: "+requiredRole)

			c.Abort()
			return
		}

		c.Next()
	}
}
