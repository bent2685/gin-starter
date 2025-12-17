package handlers

import (
	"gin-starter/internal/application/services"
	"gin-starter/internal/interfaces/dto"
	"gin-starter/pkg/utils/res"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	departmentService = services.Department
)

func CreateDepartment(c *gin.Context) {
	var req dto.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, err.Error())
		return
	}
	department, err := departmentService.CreateDepartment(req.Name, req.Description, req.ParentID)
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}
	res.SuccessWithMessage(c, "部门创建成功", department)
}

func GetAllDepartments(c *gin.Context) {
	departments, err := departmentService.GetAllDepartments()
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}
	res.Success(c, departments)
}

func GetDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, "无效的部门ID")
		return
	}
	department, err := departmentService.GetDepartmentByID(uint(id))
	if err != nil {
		res.ErrNotFound.ThrowWithMessage(c, "部门不存在")
		return
	}
	res.Success(c, department)
}

func UpdateDepartment(c *gin.Context) {
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
	department, err := departmentService.UpdateDepartment(uint(id), req.Name, req.Description, req.ParentID)
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

func DeleteDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, "无效的部门ID")
		return
	}
	err = departmentService.DeleteDepartment(uint(id))
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}
	res.SuccessWithMessage(c, "部门删除成功", nil)
}

func GetDepartmentTree(c *gin.Context) {
	departments, err := departmentService.GetDepartmentTree()
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}
	res.Success(c, departments)
}
