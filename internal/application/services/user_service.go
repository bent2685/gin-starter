package services

import (
	"errors"
	"gin-starter/internal/domain/models"
	"gin-starter/internal/infra/database"

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
func (s *UserService) CreateUser(user *models.User) error {
	// 验证用户数据
	if err := user.Validate(); err != nil {
		return err
	}

	// 检查用户名是否已存在
	if _, err := s.GetUserByUsername(user.Username); err == nil {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if _, err := s.GetUserByEmail(user.Email); err == nil {
		return errors.New("邮箱已存在")
	}

	// 创建用户
	return s.db.Create(user).Error
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := s.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAllUsers 获取所有用户
func (s *UserService) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := s.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(user *models.User) error {
	// 验证用户数据
	if err := user.Validate(); err != nil {
		return err
	}

	return s.db.Save(user).Error
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id uint) error {
	return s.db.Delete(&models.User{}, id).Error
}

// GetUserByEmail 根据邮箱获取用户
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ActivateUser 激活用户
func (s *UserService) ActivateUser(id uint) error {
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	user.Activate()
	return s.UpdateUser(user)
}

// DeactivateUser 停用用户
func (s *UserService) DeactivateUser(id uint) error {
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	user.Deactivate()
	return s.UpdateUser(user)
}

// ChangeUserEmail 更改用户邮箱
func (s *UserService) ChangeUserEmail(id uint, email string) error {
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	if err := user.ChangeEmail(email); err != nil {
		return err
	}

	// 检查邮箱是否已存在
	existingUser, err := s.GetUserByEmail(email)
	if err == nil && existingUser.ID != id {
		return errors.New("邮箱已存在")
	}

	return s.UpdateUser(user)
}
