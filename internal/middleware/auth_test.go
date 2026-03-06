package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HarshithRajesh/PixelForge/internal/domain"
	"github.com/HarshithRajesh/PixelForge/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// helper to setup a test router with middleware and a dummy protected route
func setupMiddlewareRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// dummy protected route just returns 200 if middleware passes
		protected.GET("/profile", func(c *gin.Context) {
			email := c.MustGet("email").(string)
			c.JSON(http.StatusOK, gin.H{
				"email": email,
			})
		})
	}

	return r
}

func TestAuthMiddleware(t *testing.T) {
	// generate a valid token to use in tests
	validToken, _ := domain.GenerateToken("alice@example.com")

	tests := []struct {
		name           string
		setupRequest   func(req *http.Request)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "no token at all",
			setupRequest: func(req *http.Request) {
				// don't set anything
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   nil,
		},
		{
			name: "invalid token in header",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer invalidtoken")
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   nil,
		},
		{
			name: "missing Bearer prefix",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", validToken) // no "Bearer " prefix
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   nil,
		},
		// {
		// 	name: "valid token in header",
		// 	setupRequest: func(req *http.Request) {
		// 		req.Header.Set("Authorization", "Bearer "+validToken)
		// 	},
		// 	expectedStatus: http.StatusOK,
		// 	expectedBody:   nil,
		// },
		{
			name: "valid token in cookie",
			setupRequest: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  "token",
					Value: validToken,
				})
			},
			expectedStatus: http.StatusOK,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupMiddlewareRouter()

			req := httptest.NewRequest(http.MethodGet, "/profile", nil)
			tt.setupRequest(req) // set up headers or cookies for this test
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
