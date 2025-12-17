package services

import (
	"gin-starter/internal/domain/models"
	"gin-starter/internal/infra/database"
	"gin-starter/pkg/utils/res"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

var User = &UserService{}

func (s *UserService) CreateUser(username, email, password string) (*models.User, error) {
	db := database.GetDB()
	var existingUser models.User
	if err := db.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return nil, res.ErrUsernameTaken
	}
	if err := db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, res.ErrEmailAlreadyUsed
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		IsActive: true,
	}
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (s *UserService) GetAllUsers() ([]*models.User, error) {
	db := database.GetDB()
	var users []*models.User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	for _, user := range users {
		user.Password = ""
	}
	return users, nil
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	db := database.GetDB()
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return nil, res.ErrUserNotFound
	}
	user.Password = ""
	return &user, nil
}

func (s *UserService) UpdateUser(id uint, username, email string) (*models.User, error) {
	db := database.GetDB()
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return nil, res.ErrUserNotFound
	}
	var existingUser models.User
	if err := db.Where("username = ? AND id != ?", username, id).First(&existingUser).Error; err == nil {
		return nil, res.ErrUsernameTaken
	}
	if err := db.Where("email = ? AND id != ?", email, id).First(&existingUser).Error; err == nil {
		return nil, res.ErrEmailAlreadyUsed
	}
	user.Username = username
	user.Email = email
	if err := db.Save(&user).Error; err != nil {
		return nil, err
	}
	user.Password = ""
	return &user, nil
}

func (s *UserService) DeleteUser(id uint) error {
	db := database.GetDB()
	return db.Delete(&models.User{}, id).Error
}

func (s *UserService) ActivateUser(id uint) error {
	db := database.GetDB()
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return res.ErrUserNotFound
	}
	user.IsActive = true
	return db.Save(&user).Error
}

func (s *UserService) DeactivateUser(id uint) error {
	db := database.GetDB()
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return res.ErrUserNotFound
	}
	user.IsActive = false
	return db.Save(&user).Error
}

func (s *UserService) ChangeEmail(id uint, email string) error {
	db := database.GetDB()
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return res.ErrUserNotFound
	}
	var existingUser models.User
	if err := db.Where("email = ? AND id != ?", email, id).First(&existingUser).Error; err == nil {
		return res.ErrEmailAlreadyUsed
	}
	user.Email = email
	return db.Save(&user).Error
}

func (s *UserService) AuthenticateUser(username, password string) (*models.User, error) {
	db := database.GetDB()
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, res.ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, res.ErrInvalidCredentials
	}
	user.Password = ""
	return &user, nil
}
