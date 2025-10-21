package services

import (
	"gin-starter/internal/domain/models"
	"gin-starter/internal/infra/database"
	"gin-starter/pkg/utils/res"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	db *gorm.DB
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{
		db: database.GetDB(),
	}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(username, email, password string) (*models.User, error) {
	// 检查用户名是否已存在
	var existingUser models.User
	if err := s.db.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return nil, res.ErrUsernameTaken
	}

	// 检查邮箱是否已存在
	if err := s.db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, res.ErrEmailAlreadyUsed
	}

	// 对密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户对象
	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		IsActive: true, // 默认激活状态
	}

	// 保存到数据库
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	// 清除密码字段再返回
	user.Password = ""

	return user, nil
}

// GetAllUsers 获取所有用户
func (s *UserService) GetAllUsers() ([]*models.User, error) {
	var users []*models.User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}

	// 清除密码字段
	for _, user := range users {
		user.Password = ""
	}

	return users, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, res.ErrUserNotFound
	}

	// 清除密码字段
	user.Password = ""

	return &user, nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(id uint, username, email string) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, res.ErrUserNotFound
	}

	// 检查用户名是否已被其他用户使用
	var existingUser models.User
	if err := s.db.Where("username = ? AND id != ?", username, id).First(&existingUser).Error; err == nil {
		return nil, res.ErrUsernameTaken
	}

	// 检查邮箱是否已被其他用户使用
	if err := s.db.Where("email = ? AND id != ?", email, id).First(&existingUser).Error; err == nil {
		return nil, res.ErrEmailAlreadyUsed
	}

	// 更新用户信息
	user.Username = username
	user.Email = email

	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	// 清除密码字段再返回
	user.Password = ""

	return &user, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id uint) error {
	// 软删除用户
	if err := s.db.Delete(&models.User{}, id).Error; err != nil {
		return err
	}

	return nil
}

// ActivateUser 激活用户
func (s *UserService) ActivateUser(id uint) error {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return res.ErrUserNotFound
	}

	user.IsActive = true

	if err := s.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

// DeactivateUser 停用用户
func (s *UserService) DeactivateUser(id uint) error {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return res.ErrUserNotFound
	}

	user.IsActive = false

	if err := s.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

// ChangeEmail 修改邮箱
func (s *UserService) ChangeEmail(id uint, email string) error {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return res.ErrUserNotFound
	}

	// 检查邮箱是否已被其他用户使用
	var existingUser models.User
	if err := s.db.Where("email = ? AND id != ?", email, id).First(&existingUser).Error; err == nil {
		return res.ErrEmailAlreadyUsed
	}

	user.Email = email

	if err := s.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

// AuthenticateUser 验证用户凭据
func (s *UserService) AuthenticateUser(username, password string) (*models.User, error) {
	var user models.User
	// 根据用户名查找用户
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, res.ErrInvalidCredentials
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, res.ErrInvalidCredentials
	}

	// 清除密码字段再返回
	user.Password = ""

	return &user, nil
}
