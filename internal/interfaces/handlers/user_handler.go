package handlers

import (
	"strconv"

	"gin-starter/internal/application/services"
	"gin-starter/internal/domain/models"
	"gin-starter/internal/interfaces/dto"
	"gin-starter/internal/interfaces/vo"
	"gin-starter/pkg/utils/converter"
	"gin-starter/pkg/utils/res"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: services.NewUserService(),
	}
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.HandleValidationErrors(c, err)
		return
	}

	// 转换为领域模型
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		IsActive: true, // 新用户默认激活
	}

	// 设置密码
	if err := user.SetPassword(req.Password); err != nil {
		res.Error(c, 400002, "密码设置失败: "+err.Error())
		return
	}

	// 创建用户
	if err := h.userService.CreateUser(user); err != nil {
		res.Error(c, 500001, "创建用户失败: "+err.Error())
		return
	}

	// 使用转换工具转换为VO
	var userVO vo.UserVO
	converter.SafeConvert(&userVO, user)

	res.Success(c, userVO)
}

// GetUser 获取用户信息
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		res.Error(c, 400001, "参数错误: "+err.Error())
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		res.Error(c, 400004, "用户不存在")
		return
	}

	// 使用转换工具转换为VO
	var userVO vo.UserVO
	converter.SafeConvert(&userVO, user)

	res.Success(c, userVO)
}

// GetAllUsers 获取所有用户
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		res.Error(c, 500001, "获取用户列表失败: "+err.Error())
		return
	}

	// 使用转换工具转换为VO列表
	var userVOs []vo.UserVO
	converter.SafeConvertSlice(&userVOs, &users)

	// 构造列表VO
	listVO := vo.UserListVO{
		Users: userVOs,
		Total: int64(len(userVOs)),
	}

	res.Success(c, listVO)
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		res.Error(c, 400001, "参数错误: "+err.Error())
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.HandleValidationErrors(c, err)
		return
	}

	// 检查用户是否存在
	existingUser, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		res.Error(c, 400004, "用户不存在")
		return
	}

	// 更新用户信息
	existingUser.FullName = req.FullName
	existingUser.IsActive = req.IsActive

	if err := h.userService.UpdateUser(existingUser); err != nil {
		res.Error(c, 500001, "更新用户失败: "+err.Error())
		return
	}

	// 使用转换工具转换为VO
	var userVO vo.UserVO
	converter.SafeConvert(&userVO, existingUser)

	res.Success(c, userVO)
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		res.Error(c, 400001, "参数错误: "+err.Error())
		return
	}

	// 检查用户是否存在
	if _, err := h.userService.GetUserByID(uint(id)); err != nil {
		res.Error(c, 400004, "用户不存在")
		return
	}

	// 删除用户
	if err := h.userService.DeleteUser(uint(id)); err != nil {
		res.Error(c, 500001, "删除用户失败: "+err.Error())
		return
	}

	res.SuccessWithMessage(c, "用户删除成功", nil)
}

// ActivateUser 激活用户
func (h *UserHandler) ActivateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		res.Error(c, 400001, "参数错误: "+err.Error())
		return
	}

	if err := h.userService.ActivateUser(uint(id)); err != nil {
		res.Error(c, 500001, "激活用户失败: "+err.Error())
		return
	}

	res.SuccessWithMessage(c, "用户激活成功", nil)
}

// DeactivateUser 停用用户
func (h *UserHandler) DeactivateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		res.Error(c, 400001, "参数错误: "+err.Error())
		return
	}

	if err := h.userService.DeactivateUser(uint(id)); err != nil {
		res.Error(c, 500001, "停用用户失败: "+err.Error())
		return
	}

	res.SuccessWithMessage(c, "用户停用成功", nil)
}

// ChangeEmail 更改用户邮箱
func (h *UserHandler) ChangeEmail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		res.Error(c, 400001, "参数错误: "+err.Error())
		return
	}

	var req dto.ChangeEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.HandleValidationErrors(c, err)
		return
	}

	if err := h.userService.ChangeUserEmail(uint(id), req.Email); err != nil {
		res.Error(c, 500001, "更改邮箱失败: "+err.Error())
		return
	}

	res.SuccessWithMessage(c, "邮箱更改成功", nil)
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.HandleValidationErrors(c, err)
		return
	}

	// 根据用户名获取用户
	user, err := h.userService.GetUserByUsername(req.Username)
	if err != nil {
		res.Error(c, 400004, "用户不存在")
		return
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		res.Error(c, 400005, "密码错误")
		return
	}

	// 检查用户是否激活
	if !user.IsActive {
		res.Error(c, 400006, "用户未激活")
		return
	}

	// 使用转换工具转换为VO
	var userVO vo.UserVO
	converter.SafeConvert(&userVO, user)

	// 构造登录响应VO
	loginResponse := vo.LoginResponse{
		Token:     "fake-jwt-token", // 实际项目中应生成真实的JWT token
		ExpiresAt: 3600,             // 实际项目中应设置真实的过期时间
		User:      userVO,
	}

	res.Success(c, loginResponse)
}
