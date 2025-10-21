package handlers

import (
	"strconv"

	rbacService "gin-starter/internal/application/services/rbac"
	"gin-starter/pkg/utils/res"

	"github.com/gin-gonic/gin"
)

// RBACHandler RBAC处理器
type RBACHandler struct {
	rbacService *rbacService.RBACService
}

// NewRBACHandler 创建RBAC处理器实例
func NewRBACHandler(rbacService *rbacService.RBACService) *RBACHandler {
	return &RBACHandler{
		rbacService: rbacService,
	}
}

// AddPolicy 添加策略
func (h *RBACHandler) AddPolicy(c *gin.Context) {
	var req struct {
		Sub string `json:"sub" binding:"required"`
		Obj string `json:"obj" binding:"required"`
		Act string `json:"act" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		res.HandleValidationErrors(c, err)
		return
	}

	ok, err := h.rbacService.AddPolicy(req.Sub, req.Obj, req.Act)
	if err != nil {
		res.Error(c, 500001, "添加策略失败: "+err.Error())
		return
	}

	if !ok {
		res.Error(c, 400001, "策略已存在")
		return
	}

	res.SuccessWithMessage(c, "策略添加成功", nil)
}

// AddRoleForUser 为用户添加角色
func (h *RBACHandler) AddRoleForUser(c *gin.Context) {
	var req struct {
		UserID uint   `json:"user_id" binding:"required"`
		Role   string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		res.HandleValidationErrors(c, err)
		return
	}

	ok, err := h.rbacService.AddRoleForUser(strconv.Itoa(int(req.UserID)), req.Role)
	if err != nil {
		res.Error(c, 500001, "添加角色失败: "+err.Error())
		return
	}

	if !ok {
		res.Error(c, 400001, "角色已存在")
		return
	}

	res.SuccessWithMessage(c, "角色添加成功", nil)
}

// GetRolesForUser 获取用户角色
func (h *RBACHandler) GetRolesForUser(c *gin.Context) {
	userIDStr := c.Param("user_id")

	roles, err := h.rbacService.GetRolesForUser(userIDStr)
	if err != nil {
		res.Error(c, 500001, "获取角色失败: "+err.Error())
		return
	}

	res.Success(c, roles)
}

// Enforce 验证权限
func (h *RBACHandler) Enforce(c *gin.Context) {
	var req struct {
		Sub string `json:"sub" binding:"required"`
		Obj string `json:"obj" binding:"required"`
		Act string `json:"act" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		res.HandleValidationErrors(c, err)
		return
	}

	ok, err := h.rbacService.Enforce(req.Sub, req.Obj, req.Act)
	if err != nil {
		res.Error(c, 500001, "权限验证失败: "+err.Error())
		return
	}

	res.Success(c, map[string]interface{}{
		"allowed": ok,
	})
}
