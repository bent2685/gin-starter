package converter

import (
	"time"

	"github.com/jinzhu/copier"
	"github.com/mitchellh/mapstructure"
)

// Convert 结构体转换函数
func Convert(toValue interface{}, fromValue interface{}) error {
	return copier.Copy(toValue, fromValue)
}

// ConvertSlice 切片转换函数
func ConvertSlice(toSlice interface{}, fromSlice interface{}) error {
	return copier.Copy(toSlice, fromSlice)
}

// SafeConvert 安全转换函数（忽略错误）
func SafeConvert(toValue interface{}, fromValue interface{}) {
	_ = copier.Copy(toValue, fromValue)
}

// SafeConvertSlice 安全切片转换函数（忽略错误）
func SafeConvertSlice(toSlice interface{}, fromSlice interface{}) {
	_ = copier.Copy(toSlice, fromSlice)
}

// ConvertWithOption 带选项的转换函数
func ConvertWithOption(toValue interface{}, fromValue interface{}, opt copier.Option) error {
	return copier.CopyWithOption(toValue, fromValue, opt)
}

// MapStructureConvert 使用mapstructure进行转换
func MapStructureConvert(toValue interface{}, fromValue interface{}) error {
	return mapstructure.Decode(fromValue, toValue)
}

// MapStructureConvertWithConfig 使用mapstructure和自定义配置进行转换
func MapStructureConvertWithConfig(toValue interface{}, fromValue interface{}, config *mapstructure.DecoderConfig) error {
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(fromValue)
}

// FormatTime 时间格式化函数
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatTimePtr 时间指针格式化函数
func FormatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
