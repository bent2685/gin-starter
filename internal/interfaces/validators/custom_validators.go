package validators

import (
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// RegisterCustomValidators 注册自定义验证器
func RegisterCustomValidators() {
	// 获取验证器实例
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册自定义验证器
		v.RegisterValidation("username", validateUsername)
	}
}

// validateUsername 自定义用户名验证器
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	// 用户名不能为空
	if username == "" {
		return false
	}

	// 用户名长度应在3-50个字符之间
	if len(username) < 3 || len(username) > 50 {
		return false
	}

	// 用户名只能包含字母、数字、下划线和连字符
	for _, r := range username {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-') {
			return false
		}
	}

	// 不能以连字符或下划线开头或结尾
	if username[0] == '-' || username[0] == '_' || username[len(username)-1] == '-' || username[len(username)-1] == '_' {
		return false
	}

	return true
}

// ValidateStruct 验证结构体
func ValidateStruct(s interface{}) error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		return v.Struct(s)
	}
	return nil
}

// GetValidationError 获取验证错误信息
func GetValidationError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string
		for _, e := range validationErrors {
			errorMessages = append(errorMessages, formatValidationError(e))
		}
		return strings.Join(errorMessages, "; ")
	}
	return err.Error()
}

// formatValidationError 格式化验证错误信息
func formatValidationError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " 为必填字段"
	case "email":
		return fe.Field() + " 必须是有效的邮箱地址"
	case "min":
		return fe.Field() + " 长度不能少于 " + fe.Param() + " 个字符"
	case "max":
		return fe.Field() + " 长度不能超过 " + fe.Param() + " 个字符"
	case "username":
		return fe.Field() + " 必须是3-50个字符，只能包含字母、数字、下划线和连字符，且不能以下划线或连字符开头或结尾"
	default:
		return fe.Field() + " 格式不正确"
	}
}
