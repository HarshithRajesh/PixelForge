package user_test

import (
	"testing"

	"github.com/HarshithRajesh/PixelForge/internal/domain"
	"github.com/HarshithRajesh/PixelForge/internal/models"
	service "github.com/HarshithRajesh/PixelForge/internal/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUser(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestSignUp(t *testing.T) {
	tests := []struct {
		name          string
		input         *models.User
		mockSetup     func(mock *MockUserRepository)
		expectedError string
	}{
		{
			name: "passwords do not match",
			input: &models.User{
				Password:        "abc123",
				ConfirmPassword: "wrong",
			},
			mockSetup:     func(m *MockUserRepository) {}, // no mock needed
			expectedError: "password are not matching",
		},
		{
			name: "user already exists",
			input: &models.User{
				Email:           "alice@example.com",
				Password:        "abc123",
				ConfirmPassword: "abc123",
			},
			mockSetup: func(m *MockUserRepository) {
				existingUser := &models.User{
					ID:              1,
					Name:            "sample",
					Email:           "alice@example.com",
					Password:        "123",
					ConfirmPassword: "123",
				}
				m.On("GetUser", "alice@example.com").Return(existingUser, nil)
			},
			expectedError: "user exists",
		},
		{
			name: "success",
			input: &models.User{
				Email:           "alice@example.com",
				Password:        "abc123",
				ConfirmPassword: "abc123",
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUser", "alice@example.com").Return(nil, nil)
				m.On("CreateUser", mock.Anything).Return(nil)
			},
			expectedError: "", // empty means we expect no error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo) // set up mock for this specific case
			svc := service.NewUserService(mockRepo)

			err := svc.SignUp(tt.input)

			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name          string
		input         *models.Login
		mockSetup     func(mock *MockUserRepository)
		expectedError string
	}{
		{
			name: "success",
			input: &models.Login{
				Email:    "alice@example.com",
				Password: "correctpassword",
			},
			mockSetup: func(m *MockUserRepository) {
				// hash it exactly like your service does when signing up
				hashedPassword, _ := domain.HashPassword("correctpassword")

				existingUser := &models.User{
					ID:       1,
					Email:    "alice@example.com",
					Password: hashedPassword, // ✅ hashed password
				}
				m.On("GetUser", "alice@example.com").Return(existingUser, nil)
			},
			expectedError: "",
		},
		{
			name: "wrong password",
			input: &models.Login{
				Email:    "alice@example.com",
				Password: "wrongpassword", // ← different from hashed one
			},
			mockSetup: func(m *MockUserRepository) {
				hashedPassword, _ := domain.HashPassword("correctpassword")

				existingUser := &models.User{
					ID:       1,
					Email:    "alice@example.com",
					Password: hashedPassword,
				}
				m.On("GetUser", "alice@example.com").Return(existingUser, nil)
			},
			expectedError: "invalid password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)
			svc := service.NewUserService(mockRepo)

			err := svc.Login(tt.input)

			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError) // ✅ fixed
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
