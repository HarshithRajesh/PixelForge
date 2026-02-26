// Package middleware for JWT authentication

package middleware

import (
	"net/http"
	"strings"

	"github.com/HarshithRajesh/PixelForge/internal/domain"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token from header
		// frontend sends: Authorization: Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort() // stop the request here
			return
		}

		// remove "Bearer " prefix
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// validate the token
		claims, err := domain.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// store user info in context so handlers can use it
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)

		c.Next() // continue to the actual handler
	}
}
