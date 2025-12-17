package middleware

import (
	"gin-starter/internal/interfaces/validators"
	"gin-starter/pkg/utils/res"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

// BindRequest 绑定并验证请求的中间件
func BindRequest(req interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建请求结构体的新实例
		reqValue := reflect.New(reflect.TypeOf(req).Elem())
		reqInterface := reqValue.Interface()

		// 绑定请求体
		if err := c.ShouldBindJSON(reqInterface); err != nil {
			res.ErrInvalidParam.ThrowWithMessage(c, validators.GetValidationError(err))
			c.Abort()
			return
		}

		// 将绑定后的请求对象存储到上下文中
		c.Set("request", reqInterface)
		c.Next()
	}
}

// GetRequest 从上下文中获取请求对象
func GetRequest(c *gin.Context, req interface{}) bool {
	if requestObj, exists := c.Get("request"); exists {
		// 使用反射将上下文中的对象赋值给传入的指针
		srcValue := reflect.ValueOf(requestObj)
		dstValue := reflect.ValueOf(req)

		if dstValue.Kind() == reflect.Ptr && !dstValue.IsNil() {
			dstValue.Elem().Set(srcValue.Elem())
			return true
		}
	}
	return false
}

// BindAndProcess 绑定请求并直接处理的高阶函数
func BindAndProcess(req interface{}, processor func(*gin.Context, interface{}) interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建请求结构体的新实例
		reqValue := reflect.New(reflect.TypeOf(req).Elem())
		reqInterface := reqValue.Interface()

		// 绑定请求体
		if err := c.ShouldBindJSON(reqInterface); err != nil {
			res.ErrInvalidParam.ThrowWithMessage(c, validators.GetValidationError(err))
			return
		}

		// 处理请求
		result := processor(c, reqInterface)

		// 返回结果
		c.JSON(http.StatusOK, result)
	}
}

// BindQuery 绑定查询参数的中间件
func BindQuery(req interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建请求结构体的新实例
		reqValue := reflect.New(reflect.TypeOf(req).Elem())
		reqInterface := reqValue.Interface()

		// 绑定查询参数
		if err := c.ShouldBindQuery(reqInterface); err != nil {
			res.ErrInvalidParam.ThrowWithMessage(c, validators.GetValidationError(err))
			c.Abort()
			return
		}

		// 将绑定后的请求对象存储到上下文中
		c.Set("request", reqInterface)
		c.Next()
	}
}

// BindURI 绑定URI参数的中间件
func BindURI(req interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建请求结构体的新实例
		reqValue := reflect.New(reflect.TypeOf(req).Elem())
		reqInterface := reqValue.Interface()

		// 绑定URI参数
		if err := c.ShouldBindUri(reqInterface); err != nil {
			res.ErrInvalidParam.ThrowWithMessage(c, validators.GetValidationError(err))
			c.Abort()
			return
		}

		// 将绑定后的请求对象存储到上下文中
		c.Set("request", reqInterface)
		c.Next()
	}
}
