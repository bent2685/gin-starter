package models

import (
	"errors"
	"time"
	"unicode"

	"gorm.io/gorm"
)

// User 用户领域模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Username string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username" binding:"required,min=3,max=50"`
	Email    string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email" binding:"required,email"`
	Password string `gorm:"type:varchar(255);not null" json:"-" binding:"required,min=6,max=100"`
	FullName string `gorm:"type:varchar(100)" json:"full_name" binding:"max=100"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// 可以在这里添加创建前的逻辑
	return
}

// BeforeUpdate 更新前钩子
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	// 可以在这里添加更新前的逻辑
	return
}

// Validate 验证用户数据
func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("用户名不能为空")
	}

	if len(u.Username) < 3 {
		return errors.New("用户名长度不能少于3个字符")
	}

	if u.Email == "" {
		return errors.New("邮箱不能为空")
	}

	if !isValidEmail(u.Email) {
		return errors.New("邮箱格式不正确")
	}

	if u.Password == "" {
		return errors.New("密码不能为空")
	}

	if len(u.Password) < 6 {
		return errors.New("密码长度不能少于6个字符")
	}

	return nil
}

// SetPassword 设置密码（包含加密逻辑）
func (u *User) SetPassword(password string) error {
	if len(password) < 6 {
		return errors.New("密码长度不能少于6个字符")
	}

	// 在实际项目中，这里应该使用 bcrypt 等加密算法
	// 示例中为了简化，直接存储密码（实际项目中不要这样做）
	u.Password = password
	return nil
}

// CheckPassword 检查密码
func (u *User) CheckPassword(password string) bool {
	// 在实际项目中，这里应该比较加密后的密码
	return u.Password == password
}

// Activate 激活用户
func (u *User) Activate() {
	u.IsActive = true
}

// Deactivate 停用用户
func (u *User) Deactivate() {
	u.IsActive = false
}

// UpdateProfile 更新用户资料
func (u *User) UpdateProfile(fullName string) {
	u.FullName = fullName
}

// ChangeEmail 更改邮箱
func (u *User) ChangeEmail(email string) error {
	if email == "" {
		return errors.New("邮箱不能为空")
	}

	if !isValidEmail(email) {
		return errors.New("邮箱格式不正确")
	}

	u.Email = email
	return nil
}

// isValidEmail 验证邮箱格式
func isValidEmail(email string) bool {
	// 简单的邮箱格式验证
	for _, r := range email {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '@' && r != '.' && r != '_' && r != '-' {
			return false
		}
	}
	return len(email) > 0 && len(email) <= 100 && email[0] != '@' && email[len(email)-1] != '@'
}
