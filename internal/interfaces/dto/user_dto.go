package dto

// CreateUserRequest 创建用户请求DTO
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,username"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	FullName string `json:"full_name" binding:"max=100"`
}

// UpdateUserRequest 更新用户请求DTO
type UpdateUserRequest struct {
	FullName string `json:"full_name" binding:"max=100"`
	IsActive bool   `json:"is_active"`
}

// ChangeEmailRequest 更改邮箱请求DTO
type ChangeEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// LoginRequest 登录请求DTO
type LoginRequest struct {
	Username string `json:"username" binding:"required,username"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}
