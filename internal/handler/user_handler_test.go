package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HarshithRajesh/PixelForge/internal/handler"
	"github.com/HarshithRajesh/PixelForge/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) SignUp(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) Login(user *models.Login) error {
	args := m.Called(user)
	return args.Error(0)
}

func setupRouter(handler *handler.UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/signup", handler.SignUp)
	r.POST("/login", handler.Login)
	return r
}

func TestSignUpHandler_InvalidBody(t *testing.T) {
	mockService := new(MockUserService)
	handler := handler.NewUserHandler(mockService)
	router := setupRouter(handler)

	body := `{"email":"invalid`
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockService.AssertNotCalled(t, "SignUp")
}

func TestSignUpHandler_ServiceFails(t *testing.T) {
	// ARRANGE
	mockService := new(MockUserService)
	handler := handler.NewUserHandler(mockService)
	router := setupRouter(handler)

	// valid JSON body
	body := models.Login{
		Email:    "alice@example.com",
		Password: "abc123",
	}
	jsonBody, _ := json.Marshal(body)

	mockService.On("SignUp", mock.Anything).Return(errors.New("user exists"))

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// check the response body contains the error
	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, "user exists", response["error"])

	mockService.AssertExpectations(t)
}

func TestSignUpHandler_Success(t *testing.T) {
	// ARRANGE
	mockService := new(MockUserService)
	handler := handler.NewUserHandler(mockService)
	router := setupRouter(handler)

	body := models.Login{
		Email:    "alice@example.com",
		Password: "abc123",
	}
	jsonBody, _ := json.Marshal(body)

	// tell mock: SignUp succeeds
	mockService.On("SignUp", mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// ACT
	router.ServeHTTP(w, req)

	// ASSERT
	assert.Equal(t, http.StatusOK, w.Code)

	// check the success message
	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, "User Create Successfully", response["message"])

	mockService.AssertExpectations(t)
}

func TestLoginHandler_InvalidBody(t *testing.T) {
	mockService := new(MockUserService)
	handler := handler.NewUserHandler(mockService)
	router := setupRouter(handler)

	body := `{"email":"invalid`
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockService.AssertNotCalled(t, "Login")
}

func TestLoginHandler_Fails(t *testing.T) {
	mockService := new(MockUserService)
	handler := handler.NewUserHandler(mockService)
	router := setupRouter(handler)

	// valid JSON body
	body := models.User{
		Email:    "alice@example.com",
		Password: "abc123",
	}
	jsonBody, _ := json.Marshal(body)

	mockService.On("Login", mock.Anything).Return(errors.New("user doesnt exist"))

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, "user doesnt exist", response["error"])

	mockService.AssertExpectations(t)
}

func TestLoginHandler_Success(t *testing.T) {
	mockService := new(MockUserService)
	handler := handler.NewUserHandler(mockService)
	router := setupRouter(handler)

	body := models.Login{
		Email:    "alice@example.com",
		Password: "123",
	}
	jsonBody, _ := json.Marshal(body)

	mockService.On("Login", mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, "Login Successfull", response["message"])
}
