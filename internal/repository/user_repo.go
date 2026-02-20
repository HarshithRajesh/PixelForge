// Package repository is to define the user signup logic
package repository

import (
	"errors"

	"github.com/HarshithRajesh/PixelForge/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUser(email string) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetUser(email string) (bool, error) {
	var user models.User
	err := r.db.Where("email=?", email).First(&user).Error
	if err == nil {
		return true, nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else {
		return false, err
	}
}
