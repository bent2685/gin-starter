package validators

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// fieldNameMap 字段名称中英文映射
var fieldNameMap = map[string]string{
	"Username":    "用户名",
	"Password":    "密码",
	"Email":       "邮箱",
	"Name":        "名称",
	"Description": "描述",
	"ParentID":    "父级ID",
	"UserID":      "用户ID",
	"Role":        "角色",
	"Department":  "部门",
	"Sub":         "主体",
	"Obj":         "对象",
	"Act":         "操作",
}

// getFieldName 获取字段中文名称
func getFieldName(fieldName string) string {
	if chineseName, exists := fieldNameMap[fieldName]; exists {
		return chineseName
	}
	return fieldName
}

// RegisterCustomValidators 注册自定义验证器
func RegisterCustomValidators() {
	// 获取验证器实例
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册自定义验证器
		v.RegisterValidation("username", validateUsername)

		// 注册自定义验证错误消息翻译器
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("json")
			if name == "" {
				name = fld.Name
			}
			return name
		})
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
func ValidateStruct(s any) error {
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
	fieldName := getFieldName(fe.Field())

	switch fe.Tag() {
	case "required":
		return fieldName + "为必填字段"
	case "email":
		return fieldName + "必须是有效的邮箱地址"
	case "min":
		return fieldName + "长度不能少于" + fe.Param() + "个字符"
	case "max":
		return fieldName + "长度不能超过" + fe.Param() + "个字符"
	case "username":
		return fieldName + "必须是3-50个字符，只能包含字母、数字、下划线和连字符，且不能以下划线或连字符开头或结尾"
	default:
		return fieldName + "格式不正确"
	}
}
