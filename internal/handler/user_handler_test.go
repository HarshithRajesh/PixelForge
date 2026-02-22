package handler_test

import (
	"bytes"
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

func setupRouter(handler *handler.UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/signup", handler.SignUp)
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

func TestSignUpHandler (t *testing.T){
	tests := []struct{
		name string,
		body string,
		mockService func(m *MockUserService){
	},
	expectedStatus: http.StatusBadRequest,

	}
}
