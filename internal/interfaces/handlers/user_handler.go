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

var (
	userService = services.User
)

func CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, err.Error())
		return
	}
	user, err := userService.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}
	res.Success(c, user)
}

func GetAllUsers(c *gin.Context) {
	users, err := userService.GetAllUsers()
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}
	res.Success(c, users)
}

func GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, "无效的用户ID")
		return
	}
	user, err := userService.GetUserByID(uint(id))
	if err != nil {
		res.ErrNotFound.ThrowWithMessage(c, "用户不存在")
		return
	}
	res.Success(c, user)
}

func UpdateUser(c *gin.Context) {
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
	user, err := userService.UpdateUser(uint(id), req.Username, req.Email)
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}
	res.Success(c, user)
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, "无效的用户ID")
		return
	}
	err = userService.DeleteUser(uint(id))
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}
	res.SuccessWithMessage(c, "用户删除成功", nil)
}

func ActivateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, "无效的用户ID")
		return
	}
	err = userService.ActivateUser(uint(id))
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}
	res.SuccessWithMessage(c, "用户激活成功", nil)
}

func DeactivateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, "无效的用户ID")
		return
	}
	err = userService.DeactivateUser(uint(id))
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}
	res.SuccessWithMessage(c, "用户停用成功", nil)
}

func ChangeEmail(c *gin.Context) {
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
	err = userService.ChangeEmail(uint(id), req.Email)
	if err != nil {
		res.ErrInternalServer.ThrowWithMessage(c, err.Error())
		return
	}
	res.SuccessWithMessage(c, "邮箱修改成功", nil)
}

func Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ErrInvalidParam.ThrowWithMessage(c, err.Error())
		return
	}
	user, err := userService.AuthenticateUser(req.Username, req.Password)
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
