package vo

import "time"

// UserVO 用户视图对象
type UserVO struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserListVO 用户列表视图对象
type UserListVO struct {
	Users []UserVO `json:"users"`
	Total int64    `json:"total"`
}

// LoginResponse 登录响应VO
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	User      UserVO `json:"user"`
}
