// Package user is to define authentication functions
package user

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/HarshithRajesh/PixelForge/internal/config"
	"github.com/HarshithRajesh/PixelForge/internal/domain"
	"github.com/HarshithRajesh/PixelForge/internal/middleware"
	"github.com/HarshithRajesh/PixelForge/internal/models"
	"github.com/HarshithRajesh/PixelForge/internal/repository"
)

type UserService interface {
	SignUp(user *models.User) error
	Login(ctx context.Context, user *models.Login) (*middleware.Tokens, error)
}

type userService struct {
	repo repository.UserRepository
	rds  *config.Redis
}

func NewUserService(repo repository.UserRepository, rds *config.Redis) UserService {
	return &userService{
		repo: repo,
		rds:  rds,
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

func (s *userService) Login(ctx context.Context, user *models.Login) (*middleware.Tokens, error) {
	var existingUser *models.User
	existingUser, err := s.repo.GetUser(user.Email)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, errors.New("user doesnt exist")
	}
	if !domain.CheckPasswordHash(user.Password, existingUser.Password) {
		return nil, errors.New("invalid password")
	}
	userID := strconv.FormatUint(uint64(existingUser.ID), 10)
	token, err := middleware.IssueTokens(userID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	if err := middleware.Persist(ctx, s.rds, token); err != nil {
		return nil, err
	}
	return token, nil
}
