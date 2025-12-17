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

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: services.UserService{},
	}
}

// CreateUser godoc
// @Summary 创建用户
// @Description 创建新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "创建用户请求"
// @Success 200 {object} res.Response{data=models.User} "创建成功"
// @Router /users [post]
// @Security Bearer
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := Bind(c, &req); err != nil {
		return
	}

	user, err := h.userService.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		Error(c, err)
		return
	}

	Success(c, user)
}

// GetAllUsers godoc
// @Summary 获取所有用户
// @Description 获取所有用户列表
// @Tags 用户管理
// @Produce json
// @Success 200 {object} res.Response{data=[]models.User} "获取成功"
// @Router /users [get]
// @Security Bearer
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		Error(c, err)
		return
	}
	Success(c, users)
}

// GetUser godoc
// @Summary 获取用户详情
// @Description 根据ID获取用户详情
// @Tags 用户管理
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} res.Response{data=models.User} "获取成功"
// @Router /users/{id} [get]
// @Security Bearer
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

	Success(c, user)
}

// UpdateUser godoc
// @Summary 更新用户
// @Description 根据ID更新用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param request body dto.UpdateUserRequest true "更新用户请求"
// @Success 200 {object} res.Response{data=models.User} "更新成功"
// @Router /users/{id} [put]
// @Security Bearer
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, "无效的用户ID")
		return
	}

	var req dto.UpdateUserRequest
	if err := Bind(c, &req); err != nil {
		return
	}

	user, err := h.userService.UpdateUser(uint(id), req.Username, req.Email)
	if err != nil {
		Error(c, err)
		return
	}

	Success(c, user)
}

// DeleteUser godoc
// @Summary 删除用户
// @Description 根据ID删除用户
// @Tags 用户管理
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} res.Response "删除成功"
// @Router /users/{id} [delete]
// @Security Bearer
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, "无效的用户ID")
		return
	}

	err = h.userService.DeleteUser(uint(id))
	if err != nil {
		Error(c, err)
		return
	}

	SuccessWithMessage(c, "用户删除成功", nil)
}

// ActivateUser godoc
// @Summary 激活用户
// @Description 根据ID激活用户
// @Tags 用户管理
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} res.Response "激活成功"
// @Router /users/{id}/activate [post]
// @Security Bearer
func (h *UserHandler) ActivateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, "无效的用户ID")
		return
	}

	err = h.userService.ActivateUser(uint(id))
	if err != nil {
		Error(c, err)
		return
	}

	SuccessWithMessage(c, "用户激活成功", nil)
}

// DeactivateUser godoc
// @Summary 停用用户
// @Description 根据ID停用用户
// @Tags 用户管理
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} res.Response "停用成功"
// @Router /users/{id}/deactivate [post]
// @Security Bearer
func (h *UserHandler) DeactivateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, "无效的用户ID")
		return
	}

	err = h.userService.DeactivateUser(uint(id))
	if err != nil {
		Error(c, err)
		return
	}

	SuccessWithMessage(c, "用户停用成功", nil)
}

// ChangeEmail godoc
// @Summary 修改邮箱
// @Description 根据ID修改用户邮箱
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param request body dto.ChangeEmailRequest true "修改邮箱请求"
// @Success 200 {object} res.Response "修改成功"
// @Router /users/{id}/email [put]
// @Security Bearer
func (h *UserHandler) ChangeEmail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, "无效的用户ID")
		return
	}

	var req dto.ChangeEmailRequest
	if err := Bind(c, &req); err != nil {
		return
	}

	err = h.userService.ChangeEmail(uint(id), req.Email)
	if err != nil {
		Error(c, err)
		return
	}

	SuccessWithMessage(c, "邮箱修改成功", nil)
}

// Login godoc
// @Summary 用户登录
// @Description 用户登录获取Token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "登录请求"
// @Success 200 {object} res.Response{data=dto.LoginResponse} "登录成功"
// @Router /users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := Bind(c, &req); err != nil {
		return
	}

	user, err := h.userService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		res.ErrInvalidCredentials.ThrowWithMessage(c, "用户名或密码错误")
		return
	}

	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, "Token生成失败")
		return
	}

	c.JSON(http.StatusOK, res.Response{
		Code:    20000,
		Message: "登录成功",
		Data: dto.LoginResponse{
			Token: token,
			User:  user,
		},
	})
}
