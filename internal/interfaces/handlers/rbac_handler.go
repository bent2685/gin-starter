package handlers

import (
	"gin-starter/internal/application/services/rbac"
	"gin-starter/internal/interfaces/validators"
	"gin-starter/pkg/utils/res"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddPolicyRequest struct {
	Sub string `json:"sub" binding:"required"`
	Obj string `json:"obj" binding:"required"`
	Act string `json:"act" binding:"required"`
}

type AddRoleForUserRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required"`
}

type AddDepartmentForUserRequest struct {
	UserID     uint   `json:"user_id" binding:"required"`
	Department string `json:"department" binding:"required"`
}

type EnforceRequest struct {
	Sub string `json:"sub" binding:"required"`
	Obj string `json:"obj" binding:"required"`
	Act string `json:"act" binding:"required"`
}

type RBACHandler struct{}

func NewRBACHandler() *RBACHandler {
	return &RBACHandler{}
}

// AddPolicy godoc
// @Summary 添加权限策略
// @Description 添加RBAC权限策略
// @Tags RBAC权限管理
// @Accept json
// @Produce json
// @Param request body AddPolicyRequest true "添加策略请求"
// @Success 200 {object} res.Response "添加成功"
// @Router /rbac/policy [post]
// @Security Bearer
func (h *RBACHandler) AddPolicy(c *gin.Context) {
	var req AddPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, validators.GetValidationError(err))
		return
	}
	ok, err := rbac.AddPolicy(req.Sub, req.Obj, req.Act)
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

// AddRoleForUser godoc
// @Summary 为用户分配角色
// @Description 为指定用户分配角色
// @Tags RBAC权限管理
// @Accept json
// @Produce json
// @Param request body AddRoleForUserRequest true "分配角色请求"
// @Success 200 {object} res.Response "分配成功"
// @Router /rbac/role [post]
// @Security Bearer
func (h *RBACHandler) AddRoleForUser(c *gin.Context) {
	var req AddRoleForUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, validators.GetValidationError(err))
		return
	}
	ok, err := rbac.AddRoleForUser(rbac.GetUserID(req.UserID), req.Role)
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

// GetRolesForUser godoc
// @Summary 获取用户角色
// @Description 获取指定用户的所有角色
// @Tags RBAC权限管理
// @Produce json
// @Param user_id path string true "用户ID"
// @Success 200 {object} res.Response{data=[]string} "获取成功"
// @Router /rbac/role/{user_id} [get]
// @Security Bearer
func (h *RBACHandler) GetRolesForUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		res.ErrInvalidParam.ThrowWithMessage(c, "用户ID不能为空")
		return
	}
	roles, err := rbac.GetRolesForUser(userID)
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}
	res.Success(c, roles)
}

// AddDepartmentForUser godoc
// @Summary 为用户分配部门
// @Description 为指定用户分配部门
// @Tags RBAC权限管理
// @Accept json
// @Produce json
// @Param request body AddDepartmentForUserRequest true "分配部门请求"
// @Success 200 {object} res.Response "分配成功"
// @Router /rbac/department [post]
// @Security Bearer
func (h *RBACHandler) AddDepartmentForUser(c *gin.Context) {
	var req AddDepartmentForUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, validators.GetValidationError(err))
		return
	}
	ok, err := rbac.AddDepartmentForUser(rbac.GetUserID(req.UserID), req.Department)
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

// GetDepartmentsForUser godoc
// @Summary 获取用户部门
// @Description 获取指定用户的所有部门
// @Tags RBAC权限管理
// @Produce json
// @Param user_id path string true "用户ID"
// @Success 200 {object} res.Response{data=[]string} "获取成功"
// @Router /rbac/department/{user_id} [get]
// @Security Bearer
func (h *RBACHandler) GetDepartmentsForUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		res.ErrInvalidParam.ThrowWithMessage(c, "用户ID不能为空")
		return
	}
	departments, err := rbac.GetDepartmentsForUser(userID)
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}
	res.Success(c, departments)
}

// RBACEnforce godoc
// @Summary 权限验证
// @Description 验证用户是否有权限执行指定操作
// @Tags RBAC权限管理
// @Accept json
// @Produce json
// @Param request body EnforceRequest true "权限验证请求"
// @Success 200 {object} res.Response{data=map[string]bool} "验证完成"
// @Router /rbac/enforce [post]
// @Security Bearer
func (h *RBACHandler) RBACEnforce(c *gin.Context) {
	var req EnforceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, validators.GetValidationError(err))
		return
	}
	ok, err := rbac.Enforce(req.Sub, req.Obj, req.Act)
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
