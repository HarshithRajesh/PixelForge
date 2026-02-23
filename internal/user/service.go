// Package user is to define authentication functions
package user

import (
	"errors"
	"log"

	"github.com/HarshithRajesh/PixelForge/internal/domain"
	"github.com/HarshithRajesh/PixelForge/internal/models"
	"github.com/HarshithRajesh/PixelForge/internal/repository"
)

type UserService interface {
	SignUp(user *models.User) error
	Login(user *models.Login) error
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
	existingUser, err := s.repo.GetUser(user.Email)
	if err != nil {
		log.Printf("database error: %v", err)
		return err
	}

	if existingUser != nil {
		return errors.New("user exists")
	}
	user.Password, err = domain.HashPassword(user.Password)
	if err != nil {
		return err
	}
	err = s.repo.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) Login(user *models.Login) error {
	var existingUser *models.User
	existingUser, err := s.repo.GetUser(user.Email)
	if existingUser == nil {
		return errors.New("user doesnt exist")
	} else if err != nil {
		return err
	}
	if !domain.CheckPasswordHash(user.Password, existingUser.Password) {
		return errors.New("invalid password")
	}
	return nil
}
