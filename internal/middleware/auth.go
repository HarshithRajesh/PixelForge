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
		tokenString, err := c.Cookie("token")
		if err != nil {
			authHeader := c.GetHeader("Authorization")

			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
				c.Abort()
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
				c.Abort()
				return
			}
		}
		claims, err := domain.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Next()
	}
}
