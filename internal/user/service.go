// Package user is to define authentication functions
package user

import (
	"errors"

	"github.com/HarshithRajesh/PixelForge/internal/models"
	"github.com/HarshithRajesh/PixelForge/internal/repository"
)

type UserService interface {
	SignUp(user *models.User) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) SignUp(user *models.User) error {
	if user.ConfirmPassword != user.Password {
		return errors.New("password are not matching")
	}
	existing, err := s.repo.GetUser(user.Email)
	if existing {
		return errors.New("user exists")
	} else if err != nil {
		return err
	}

	err = s.repo.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}
