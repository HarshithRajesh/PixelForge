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

func (m *MockUserService) Login(user *models.Login) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

// func (m *MockUserService)Logout()error{
// 	return args.Error(0)
// }

func setupRouter(h *handler.UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/signup", h.SignUp)
	r.POST("/login", h.Login)
	return r
}

func buildBody(body interface{}) []byte {
	switch v := body.(type) {
	case string:
		return []byte(v)
	default:
		jsonBody, _ := json.Marshal(v)
		return jsonBody
	}
}

// SIGNUP TESTS
func TestSignUpHandler(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		mockSetup      func(m *MockUserService)
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name: "invalid json body",
			body: `{"email":"invalid`,
			mockSetup: func(m *MockUserService) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name: "user already exists",
			body: models.User{
				Email:           "alice@example.com",
				Password:        "abc123",
				ConfirmPassword: "abc123",
			},
			mockSetup: func(m *MockUserService) {
				m.On("SignUp", mock.Anything).Return(errors.New("user exists"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "user exists"},
		},
		{
			name: "success",
			body: models.User{
				Email:           "alice@example.com",
				Password:        "abc123",
				ConfirmPassword: "abc123",
			},
			mockSetup: func(m *MockUserService) {
				m.On("SignUp", mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"message": "User Create Successfully"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			tt.mockSetup(mockService)
			h := handler.NewUserHandler(mockService)
			router := setupRouter(h)

			req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(buildBody(tt.body)))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response map[string]string
				json.NewDecoder(w.Body).Decode(&response)
				assert.Equal(t, tt.expectedBody, response)
			} else {
				mockService.AssertNotCalled(t, "SignUp")
			}

			mockService.AssertExpectations(t)
		})
	}
}

// LOGIN TESTS
func TestLoginHandler(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		mockSetup      func(m *MockUserService)
		expectedStatus int
		expectedBody   map[string]string
		expectedCookie string
	}{
		{
			name: "invalid json body",
			body: `{"email":"invalid`,
			mockSetup: func(m *MockUserService) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
			expectedCookie: "",
		},
		{
			name: "user does not exist",
			body: models.Login{
				Email:    "alice@example.com",
				Password: "abc123",
			},
			mockSetup: func(m *MockUserService) {
				m.On("Login", mock.Anything).Return("", errors.New("user doesnt exist"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "user doesnt exist"},
			expectedCookie: "",
		},
		{
			name: "success",
			body: models.Login{
				Email:    "alice@example.com",
				Password: "123",
			},
			mockSetup: func(m *MockUserService) {
				m.On("Login", mock.Anything).Return("mocked.jwt.token", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"message": "Login Successfull"},
			expectedCookie: "token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			mockService := new(MockUserService)
			tt.mockSetup(mockService)
			h := handler.NewUserHandler(mockService)
			router := setupRouter(h)

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(buildBody(tt.body)))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != nil {
				var response map[string]string
				json.NewDecoder(w.Body).Decode(&response)
				assert.Equal(t, tt.expectedBody, response)
			}

			// check cookie independently
			if tt.expectedCookie != "" {
				cookies := w.Result().Cookies()
				found := false
				for _, cookie := range cookies {
					if cookie.Name == tt.expectedCookie {
						found = true
						assert.NotEmpty(t, cookie.Value)
					}
				}
				assert.True(t, found, "expected cookie '%s' not found", tt.expectedCookie)
			}

			// check service not called only for invalid body case
			if tt.expectedBody == nil && tt.expectedCookie == "" {
				mockService.AssertNotCalled(t, "Login")
			}

			mockService.AssertExpectations(t)
		})
	}
}
