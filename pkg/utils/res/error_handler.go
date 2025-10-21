package res

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ValidationError 验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// HandleValidationErrors 处理验证错误
func HandleValidationErrors(c *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var errorsList []ValidationError
		for _, fieldError := range validationErrors {
			errorsList = append(errorsList, ValidationError{
				Field:   fieldError.Field(),
				Message: getValidationMessage(fieldError),
			})
		}

		ErrorWithDetails(c, 400001, "参数验证失败", errorsList)
		return
	}

	// 处理其他绑定错误
	Error(c, 400001, "参数错误: "+err.Error())
}

// ErrorWithDetails 带详细信息的错误响应
func ErrorWithDetails(c *gin.Context, code int, message string, details interface{}) {
	c.JSON(400, Response{
		Code:    code,
		Message: message,
		Data:    details,
	})
	c.Abort()
}

// getValidationMessage 获取验证错误信息
func getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "该字段为必填项"
	case "email":
		return "请输入有效的邮箱地址"
	case "min":
		return "长度不能少于 " + fe.Param() + " 个字符"
	case "max":
		return "长度不能超过 " + fe.Param() + " 个字符"
	case "username":
		return "用户名必须是3-50个字符，只能包含字母、数字、下划线和连字符，且不能以下划线或连字符开头或结尾"
	default:
		return "格式不正确"
	}
}

// ParseError 解析错误信息
func ParseError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string
		for _, e := range validationErrors {
			errorMessages = append(errorMessages, e.Field()+": "+getValidationMessage(e))
		}
		return strings.Join(errorMessages, "; ")
	}
	return err.Error()
}
