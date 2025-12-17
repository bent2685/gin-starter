package handlers

import (
	"gin-starter/internal/application/services"
	"gin-starter/internal/interfaces/dto"
	"gin-starter/internal/interfaces/validators"
	"gin-starter/pkg/utils/res"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	departmentService services.DepartmentService
}

func NewDepartmentHandler() *DepartmentHandler {
	return &DepartmentHandler{
		departmentService: services.DepartmentService{},
	}
}

// CreateDepartment godoc
// @Summary 创建部门
// @Description 创建新的部门
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param request body dto.CreateDepartmentRequest true "创建部门请求"
// @Success 200 {object} res.Response{data=rbac.Department} "创建成功"
// @Router /departments [post]
// @Security Bearer
func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	var req dto.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, validators.GetValidationError(err))
		return
	}
	department, err := h.departmentService.CreateDepartment(req.Name, req.Description, req.ParentID)
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}
	res.SuccessWithMessage(c, "部门创建成功", department)
}

// GetAllDepartments godoc
// @Summary 获取所有部门
// @Description 获取所有部门列表
// @Tags 部门管理
// @Produce json
// @Success 200 {object} res.Response{data=[]rbac.Department} "获取成功"
// @Router /departments [get]
// @Security Bearer
func (h *DepartmentHandler) GetAllDepartments(c *gin.Context) {
	departments, err := h.departmentService.GetAllDepartments()
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}
	res.Success(c, departments)
}

// GetDepartment godoc
// @Summary 获取部门详情
// @Description 根据ID获取部门详情
// @Tags 部门管理
// @Produce json
// @Param id path int true "部门ID"
// @Success 200 {object} res.Response{data=rbac.Department} "获取成功"
// @Failure 404 {object} res.Response "部门不存在"
// @Router /departments/{id} [get]
// @Security Bearer
func (h *DepartmentHandler) GetDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, "无效的部门ID")
		return
	}
	department, err := h.departmentService.GetDepartmentByID(uint(id))
	if err != nil {
		res.ErrNotFound.ThrowWithMessage(c, "部门不存在")
		return
	}
	res.Success(c, department)
}

// UpdateDepartment godoc
// @Summary 更新部门
// @Description 根据ID更新部门信息
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param id path int true "部门ID"
// @Param request body dto.UpdateDepartmentRequest true "更新部门请求"
// @Success 200 {object} res.Response{data=rbac.Department} "更新成功"
// @Router /departments/{id} [put]
// @Security Bearer
func (h *DepartmentHandler) UpdateDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, "无效的部门ID")
		return
	}
	var req dto.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, validators.GetValidationError(err))
		return
	}
	department, err := h.departmentService.UpdateDepartment(uint(id), req.Name, req.Description, req.ParentID)
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, res.Response{
		Code:    20000,
		Message: "部门更新成功",
		Data:    department,
	})
}

// DeleteDepartment godoc
// @Summary 删除部门
// @Description 根据ID删除部门
// @Tags 部门管理
// @Produce json
// @Param id path int true "部门ID"
// @Success 200 {object} res.Response "删除成功"
// @Router /departments/{id} [delete]
// @Security Bearer
func (h *DepartmentHandler) DeleteDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, "无效的部门ID")
		return
	}
	err = h.departmentService.DeleteDepartment(uint(id))
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}
	res.SuccessWithMessage(c, "部门删除成功", nil)
}

// GetDepartmentTree godoc
// @Summary 获取部门树
// @Description 获取部门树形结构
// @Tags 部门管理
// @Produce json
// @Success 200 {object} res.Response{data=[]rbac.Department} "获取成功"
// @Router /departments/tree [get]
// @Security Bearer
func (h *DepartmentHandler) GetDepartmentTree(c *gin.Context) {
	departments, err := h.departmentService.GetDepartmentTree()
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}
	res.Success(c, departments)
}
