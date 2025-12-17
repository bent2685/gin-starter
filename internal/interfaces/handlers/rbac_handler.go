package handlers

import (
	"gin-starter/internal/application/services/rbac"
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

func AddPolicy(c *gin.Context) {
	var req AddPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, err.Error())
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

func AddRoleForUser(c *gin.Context) {
	var req AddRoleForUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, err.Error())
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

func GetRolesForUser(c *gin.Context) {
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

func AddDepartmentForUser(c *gin.Context) {
	var req AddDepartmentForUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, err.Error())
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

func GetDepartmentsForUser(c *gin.Context) {
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

func RBACEnforce(c *gin.Context) {
	var req EnforceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, err.Error())
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
