package domain_test

import (
	"os"
	"testing"
	"time"

	"github.com/HarshithRajesh/PixelForge/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken_Sucess(t *testing.T) {
	token, err := domain.GenerateToken("alice@example.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken_Sucess(t *testing.T) {
	token, _ := domain.GenerateToken("alice@example.com")

	claims, err := domain.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "alice@example.com", claims.Email)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	claims, err := domain.ValidateToken("sending.invalid.token")
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	claims := &domain.Claims{
		Email: "alice@example.com",
	}
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-1 * time.Hour))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	result, err := domain.ValidateToken(tokenString)

	assert.Error(t, err)
	assert.Nil(t, result)
}
