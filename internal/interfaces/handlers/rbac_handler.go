package handlers

import (
	"gin-starter/internal/application/services/rbac"
	"gin-starter/pkg/utils/res"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RBACHandler RBAC处理器
type RBACHandler struct {
	rbacService *rbac.RBACService
}

// NewRBACHandler 创建RBAC处理器实例
func NewRBACHandler(rbacService *rbac.RBACService) *RBACHandler {
	return &RBACHandler{
		rbacService: rbacService,
	}
}

// AddPolicyRequest 添加策略请求
type AddPolicyRequest struct {
	Sub string `json:"sub" binding:"required"` // 角色或部门名称
	Obj string `json:"obj" binding:"required"` // 资源路径
	Act string `json:"act" binding:"required"` // 操作方法
}

// AddRoleForUserRequest 为用户添加角色请求
type AddRoleForUserRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required"`
}

// AddDepartmentForUserRequest 为用户添加部门请求
type AddDepartmentForUserRequest struct {
	UserID     uint   `json:"user_id" binding:"required"`
	Department string `json:"department" binding:"required"`
}

// EnforceRequest 权限验证请求
type EnforceRequest struct {
	Sub string `json:"sub" binding:"required"` // 用户ID
	Obj string `json:"obj" binding:"required"` // 资源路径
	Act string `json:"act" binding:"required"` // 操作方法
}

// AddPolicy 添加策略
func (h *RBACHandler) AddPolicy(c *gin.Context) {
	var req AddPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, err.Error())
		return
	}

	ok, err := h.rbacService.AddPolicy(req.Sub, req.Obj, req.Act)
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}

	if ok {
		res.SuccessWithMessage(c, "策略添加成功", nil)
	} else {
		res.SuccessWithMessage(c, "策略已存在", nil)
	}
}

// AddRoleForUser 为用户添加角色
func (h *RBACHandler) AddRoleForUser(c *gin.Context) {
	var req AddRoleForUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, err.Error())
		return
	}

	ok, err := h.rbacService.AddRoleForUser(rbac.GetUserID(req.UserID), req.Role)
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}

	if ok {
		res.SuccessWithMessage(c, "角色分配成功", nil)
	} else {
		res.SuccessWithMessage(c, "角色已分配", nil)
	}
}

// GetRolesForUser 获取用户的角色
func (h *RBACHandler) GetRolesForUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		res.ErrInvalidParam.ThrowWithMessage(c, "用户ID不能为空")
		return
	}

	roles, err := h.rbacService.GetRolesForUser(userID)
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}

	res.Success(c, roles)
}

// AddDepartmentForUser 为用户添加部门
func (h *RBACHandler) AddDepartmentForUser(c *gin.Context) {
	var req AddDepartmentForUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, err.Error())
		return
	}

	ok, err := h.rbacService.AddDepartmentForUser(rbac.GetUserID(req.UserID), req.Department)
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}

	if ok {
		res.SuccessWithMessage(c, "部门分配成功", nil)
	} else {
		res.SuccessWithMessage(c, "部门已分配", nil)
	}
}

// GetDepartmentsForUser 获取用户的部门
func (h *RBACHandler) GetDepartmentsForUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		res.ErrInvalidParam.ThrowWithMessage(c, "用户ID不能为空")
		return
	}

	departments, err := h.rbacService.GetDepartmentsForUser(userID)
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}

	res.Success(c, departments)
}

// Enforce 验证权限
func (h *RBACHandler) Enforce(c *gin.Context) {
	var req EnforceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, err.Error())
		return
	}

	ok, err := h.rbacService.Enforce(req.Sub, req.Obj, req.Act)
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, res.Response{
		Code:    20000,
		Message: "权限验证完成",
		Data: map[string]bool{
			"allowed": ok,
		},
	})
}
