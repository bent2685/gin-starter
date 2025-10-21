package handlers

import (
	"gin-starter/internal/application/services"
	"gin-starter/internal/interfaces/dto"
	"gin-starter/pkg/utils/res"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DepartmentHandler 部门处理器
type DepartmentHandler struct {
	departmentService *services.DepartmentService
}

// NewDepartmentHandler 创建部门处理器实例
func NewDepartmentHandler(departmentService *services.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{
		departmentService: departmentService,
	}
}

// CreateDepartment 创建部门
func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	var req dto.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, err.Error())
		return
	}

	department, err := h.departmentService.CreateDepartment(req.Name, req.Description, req.ParentID)
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, res.Response{
		Code:    20000,
		Message: "部门创建成功",
		Data:    department,
	})
}

// GetAllDepartments 获取所有部门
func (h *DepartmentHandler) GetAllDepartments(c *gin.Context) {
	departments, err := h.departmentService.GetAllDepartments()
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}

	res.Success(c, departments)
}

// GetDepartment 获取部门详情
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

// UpdateDepartment 更新部门
func (h *DepartmentHandler) UpdateDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, "无效的部门ID")
		return
	}

	var req dto.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, err.Error())
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

// DeleteDepartment 删除部门
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

// GetDepartmentTree 获取部门树形结构
func (h *DepartmentHandler) GetDepartmentTree(c *gin.Context) {
	departments, err := h.departmentService.GetDepartmentTree()
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}

	res.Success(c, departments)
}
