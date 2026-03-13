// Package repository is to define the user signup logic
package repository

import (
	"errors"
	"strconv"

	"github.com/HarshithRajesh/PixelForge/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUser(email string) (*models.User, error)
	SaveImageDB(metadata *models.Image) error
	GetAllImageData(userID uint) ([]*models.Image, error)
	GetImage(imageID string, userID string) (*models.Image, error)
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

func (r *userRepository) GetUser(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email=?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) SaveImageDB(metadata *models.Image) error {
	return r.db.Create(metadata).Error
}

func (r *userRepository) GetAllImageData(userID uint) ([]*models.Image, error) {
	var images []*models.Image
	err := r.db.Where("user_id=?", userID).Find(&images).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return images, nil
}

func (r *userRepository) GetImage(imageID string, userID string) (*models.Image, error) {
	newuserID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return nil, err
	}
	nuserID := uint(newuserID)
	var image *models.Image
	result := r.db.Where("id = ? AND user_id = ?", imageID, nuserID).First(&image)

	if result.Error != nil {
		return nil, result.Error
	}

	return image, nil
}
