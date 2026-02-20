package user_test

import (
	"errors"
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

func TestSignUp_PasswordDontMatch(t *testing.T) {
	mockrepo := new(MockUserRepository)
	svc := service.NewUserService(mockrepo)

	user := &models.User{
		Email:           "example@test.com",
		Password:        "123",
		ConfirmPassword: "wrong123",
	}

	err := svc.SignUp(user)
	assert.EqualError(t, err, "password are not matching")
	mockrepo.AssertNotCalled(t, "GetUser")
	mockrepo.AssertNotCalled(t, "CreateUser")
}

func TestSignUp_UserAlreadyExists(t *testing.T) {
	mockrepo := new(MockUserRepository)
	svc := service.NewUserService(mockrepo)

	user := &models.User{
		Email:           "example@test.com",
		Password:        "123",
		ConfirmPassword: "123",
	}

	mockrepo.On("GetUser", "example@test.com").Return(true, nil)
	err := svc.SignUp(user)
	assert.EqualError(t, err, "user exists")
	mockrepo.AssertNotCalled(t, "CreateUser")
	mockrepo.AssertExpectations(t)
}

func TestSignUp_GetUserDBError(t *testing.T) {
	mockrepo := new(MockUserRepository)
	svc := service.NewUserService(mockrepo)

	user := &models.User{
		Email:           "example@test.com",
		Password:        "123",
		ConfirmPassword: "123",
	}

	mockrepo.On("GetUser", "example@test.com").Return(false, errors.New("db connection lost"))
	err := svc.SignUp(user)

	assert.EqualError(t, err, "db connection lost")
	mockrepo.AssertNotCalled(t, "CreateUser")
	mockrepo.AssertExpectations(t)
}

func TestSignUp_CreateUserFails(t *testing.T) {
	mockrepo := new(MockUserRepository)
	svc := service.NewUserService(mockrepo)

	user := &models.User{
		Email:           "example@test.com",
		Password:        "123",
		ConfirmPassword: "123",
	}
	mockrepo.On("GetUser", "example@test.com").Return(false, nil)
	mockrepo.On("CreateUser", user).Return(errors.New("insert failed"))

	err := svc.SignUp(user)
	assert.EqualError(t, err, "insert failed")
	mockrepo.AssertExpectations(t)
}

func TestSignUp_Success(t *testing.T) {
	mockrepo := new(MockUserRepository)
	svc := service.NewUserService(mockrepo)

	user := &models.User{
		Email:           "example@test.com",
		Password:        "123",
		ConfirmPassword: "123",
	}
	mockrepo.On("GetUser", "example@test.com").Return(false, nil)
	mockrepo.On("CreateUser", user).Return(nil)

	err := svc.SignUp(user)
	assert.NoError(t, err)
	mockrepo.AssertExpectations(t)
}
