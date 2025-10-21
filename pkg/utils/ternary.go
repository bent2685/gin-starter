package utils

// Ternary 三元表达式模拟
func Ternary[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

// DefaultString 字符串默认值
func DefaultString(primary, fallback string) string {
	return Ternary(primary != "", primary, fallback)
}

// DefaultInt 整数默认值
func DefaultInt(primary, fallback int) int {
	return Ternary(primary != 0, primary, fallback)
}
