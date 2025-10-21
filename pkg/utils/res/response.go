package res

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    20000,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 带自定义消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
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
	ErrInvalidParam   = NewBusinessError(400010, "参数错误")
	ErrUnauthorized   = NewBusinessError(400001, "未授权")
	ErrForbidden      = NewBusinessError(400003, "禁止访问")
	ErrNotFound       = NewBusinessError(400004, "资源不存在")
	ErrWps            = NewBusinessError(500010, "WPS异常")
	ErrInternalServer = NewBusinessError(500001, "服务器内部错误")
)
