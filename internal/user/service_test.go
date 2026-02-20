package user_test

import (
	"testing"

	"github.com/HarshithRajesh/PixelForge/internal/models"
	service "github.com/HarshithRajesh/PixelForge/internal/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUser(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// func TestSignUp_PasswordDontMatch(t *testing.T) {
// 	mockrepo := new(MockUserRepository)
// 	svc := service.NewUserService(mockrepo)
//
// 	user := &models.User{
// 		Email:           "example@test.com",
// 		Password:        "123",
// 		ConfirmPassword: "wrong123",
// 	}
//
// 	err := svc.SignUp(user)
// 	assert.EqualError(t, err, "password are not matching")
// 	mockrepo.AssertNotCalled(t, "GetUser")
// 	mockrepo.AssertNotCalled(t, "CreateUser")
// }
//
// func TestSignUp_UserAlreadyExists(t *testing.T) {
// 	mockrepo := new(MockUserRepository)
// 	svc := service.NewUserService(mockrepo)
//
// 	user := &models.User{
// 		Email:           "example@test.com",
// 		Password:        "123",
// 		ConfirmPassword: "123",
// 	}
//
// 	mockrepo.On("GetUser", "example@test.com").Return(true, nil)
// 	err := svc.SignUp(user)
// 	assert.EqualError(t, err, "user exists")
// 	mockrepo.AssertNotCalled(t, "CreateUser")
// 	mockrepo.AssertExpectations(t)
// }
//
// func TestSignUp_GetUserDBError(t *testing.T) {
// 	mockrepo := new(MockUserRepository)
// 	svc := service.NewUserService(mockrepo)
//
// 	user := &models.User{
// 		Email:           "example@test.com",
// 		Password:        "123",
// 		ConfirmPassword: "123",
// 	}
//
// 	mockrepo.On("GetUser", "example@test.com").Return(false, errors.New("db connection lost"))
// 	err := svc.SignUp(user)
//
// 	assert.EqualError(t, err, "db connection lost")
// 	mockrepo.AssertNotCalled(t, "CreateUser")
// 	mockrepo.AssertExpectations(t)
// }
//
// func TestSignUp_CreateUserFails(t *testing.T) {
// 	mockrepo := new(MockUserRepository)
// 	svc := service.NewUserService(mockrepo)
//
// 	user := &models.User{
// 		Email:           "example@test.com",
// 		Password:        "123",
// 		ConfirmPassword: "123",
// 	}
// 	mockrepo.On("GetUser", "example@test.com").Return(false, nil)
// 	mockrepo.On("CreateUser", user).Return(errors.New("insert failed"))
//
// 	err := svc.SignUp(user)
// 	assert.EqualError(t, err, "insert failed")
// 	mockrepo.AssertExpectations(t)
// }
//
// func TestSignUp_Success(t *testing.T) {
// 	mockrepo := new(MockUserRepository)
// 	svc := service.NewUserService(mockrepo)
//
// 	user := &models.User{
// 		Email:           "example@test.com",
// 		Password:        "123",
// 		ConfirmPassword: "123",
// 	}
// 	mockrepo.On("GetUser", "example@test.com").Return(false, nil)
// 	mockrepo.On("CreateUser", user).Return(nil)
//
// 	err := svc.SignUp(user)
// 	assert.NoError(t, err)
// 	mockrepo.AssertExpectations(t)
// }

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
				m.On("GetUser", "alice@example.com").Return(true, nil)
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
				m.On("GetUser", "alice@example.com").Return(false, nil)
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
