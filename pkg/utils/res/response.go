package res

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 通用响应结构
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    20000,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 带自定义消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    20000,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
	c.Abort()
}

// ErrorWithHttpStatus 带HTTP状态码的错误响应
func ErrorWithHttpStatus(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
	})
	c.Abort()
}

// BusinessError 业务异常结构
type BusinessError struct {
	Code    int
	Message string
}

func (e *BusinessError) Error() string {
	return e.Message
}

// WithMessage 覆盖异常消息
func (e *BusinessError) WithMessage(message string) *BusinessError {
	return &BusinessError{
		Code:    e.Code,
		Message: message,
	}
}

// ThrowWithMessage 抛出带自定义消息的异常
func (e *BusinessError) ThrowWithMessage(c *gin.Context, message string) {
	Error(c, e.Code, message)
}

// NewBusinessError 创建业务异常
func NewBusinessError(code int, message string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
	}
}

// 常用业务异常
var (
	// 通用错误
	ErrInvalidParam   = NewBusinessError(400010, "参数错误")
	ErrUnauthorized   = NewBusinessError(400001, "未授权")
	ErrForbidden      = NewBusinessError(400003, "禁止访问")
	ErrNotFound       = NewBusinessError(400004, "资源不存在")
	ErrInternalServer = NewBusinessError(500001, "服务器内部错误")

	// 认证相关错误
	ErrInvalidCredentials = NewBusinessError(400002, "用户名或密码错误")
	ErrTokenExpired       = NewBusinessError(400005, "Token已过期")
	ErrInvalidToken       = NewBusinessError(400006, "无效的Token")
	ErrTokenRequired      = NewBusinessError(400007, "未提供认证信息")
	ErrTokenFormat        = NewBusinessError(400008, "认证信息格式错误")

	// 用户相关错误
	ErrUserNotFound      = NewBusinessError(400101, "用户不存在")
	ErrUserAlreadyExists = NewBusinessError(400102, "用户已存在")
	ErrInvalidPassword   = NewBusinessError(400103, "密码错误")
	ErrEmailAlreadyUsed  = NewBusinessError(400104, "邮箱已被使用")
	ErrUsernameTaken     = NewBusinessError(400105, "用户名已存在")

	// 权限相关错误
	ErrInsufficientPermissions = NewBusinessError(400009, "权限不足")
)
