package handlers

import (
	"gin-starter/internal/application/services"
	"gin-starter/internal/interfaces/dto"
	"gin-starter/pkg/utils/jwt"
	"gin-starter/pkg/utils/res"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, err.Error())
		return
	}

	user, err := h.userService.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}

	res.Success(c, user)
}

// GetAllUsers 获取所有用户
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}

	res.Success(c, users)
}

// GetUser 获取用户详情
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, "无效的用户ID")
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		res.ErrNotFound.ThrowWithMessage(c, "用户不存在")
		return
	}

	res.Success(c, user)
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, "无效的用户ID")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, err.Error())
		return
	}

	user, err := h.userService.UpdateUser(uint(id), req.Username, req.Email)
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}

	res.Success(c, user)
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, "无效的用户ID")
		return
	}

	err = h.userService.DeleteUser(uint(id))
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}

	res.SuccessWithMessage(c, "用户删除成功", nil)
}

// ActivateUser 激活用户
func (h *UserHandler) ActivateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, "无效的用户ID")
		return
	}

	err = h.userService.ActivateUser(uint(id))
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}

	res.SuccessWithMessage(c, "用户激活成功", nil)
}

// DeactivateUser 停用用户
func (h *UserHandler) DeactivateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, "无效的用户ID")
		return
	}

	err = h.userService.DeactivateUser(uint(id))
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}

	res.SuccessWithMessage(c, "用户停用成功", nil)
}

// ChangeEmail 修改邮箱
func (h *UserHandler) ChangeEmail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, "无效的用户ID")
		return
	}

	var req dto.ChangeEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, err.Error())
		return
	}

	err = h.userService.ChangeEmail(uint(id), req.Email)
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}

	res.SuccessWithMessage(c, "邮箱修改成功", nil)
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, err.Error())
		return
	}

	// 验证用户凭据
	user, err := h.userService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		res.ErrInvalidCredentials.ThrowWithMessage(c, "用户名或密码错误")
		return
	}

	// 生成JWT Token
	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, "Token生成失败")
		return
	}

	// 返回Token和用户信息
	response := dto.LoginResponse{
		Token: token,
		User:  user,
	}

	c.JSON(http.StatusOK, res.Response{
		Code:    20000,
		Message: "登录成功",
		Data:    response,
	})
}
