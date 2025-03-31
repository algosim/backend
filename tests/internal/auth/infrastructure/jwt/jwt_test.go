package tests

import (
	"testing"
	"time"

	"github.com/algosim/backend/configs"
	"github.com/algosim/backend/internal/auth/domain"
	jwtinfra "github.com/algosim/backend/internal/auth/infrastructure/jwt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJWTManager(t *testing.T) {
	// Setup test configuration
	config := &configs.Config{
		Auth: struct {
			JWTSecret string `mapstructure:"jwt_secret" env:"AUTH_JWT_SECRET"`
			TokenTTL  int    `mapstructure:"token_ttl" env:"AUTH_TOKEN_TTL"`
		}{
			JWTSecret: "test-secret-key-123",
			TokenTTL:  3600, // 1 hour
		},
	}

	jwtManager := jwtinfra.NewJWTManager(config)

	// Create test user
	testUser := &domain.User{
		ID:    uuid.New(),
		Email: "test@example.com",
	}

	t.Run("GenerateToken", func(t *testing.T) {
		token, err := jwtManager.GenerateToken(testUser)
		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.NotEmpty(t, token.AccessToken)
		assert.NotEmpty(t, token.RefreshToken)
		assert.Equal(t, testUser.ID, token.UserID)
		assert.True(t, time.Now().Before(token.ExpiresAt))
	})

	t.Run("ValidateAccessToken", func(t *testing.T) {
		// Generate a token
		token, err := jwtManager.GenerateToken(testUser)
		assert.NoError(t, err)

		// Test valid token
		claims, err := jwtManager.ValidateAccessToken(token.AccessToken)
		assert.NoError(t, err)
		assert.Equal(t, testUser.ID, claims.UserID)
		assert.Equal(t, testUser.Email, claims.Email)

		// Test invalid token
		invalidToken := "invalid.token.string"
		claims, err = jwtManager.ValidateAccessToken(invalidToken)
		assert.Error(t, err)
		assert.Nil(t, claims)

		// Test expired token
		expiredClaims := jwtinfra.Claims{
			UserID: testUser.ID,
			Email:  testUser.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			},
		}
		expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
		expiredTokenString, err := expiredToken.SignedString([]byte(config.Auth.JWTSecret))
		assert.NoError(t, err)

		claims, err = jwtManager.ValidateAccessToken(expiredTokenString)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("ValidateRefreshToken", func(t *testing.T) {
		// Generate a token
		token, err := jwtManager.GenerateToken(testUser)
		assert.NoError(t, err)

		// Test valid refresh token
		err = jwtManager.ValidateRefreshToken(token)
		assert.NoError(t, err)

		// Test expired refresh token
		expiredToken := &domain.Token{
			ID:        uuid.New(),
			UserID:    testUser.ID,
			ExpiresAt: time.Now().Add(-1 * time.Hour),
		}
		err = jwtManager.ValidateRefreshToken(expiredToken)
		assert.Error(t, err)
	})

	t.Run("ExtractUserIDFromToken", func(t *testing.T) {
		// Generate a token
		token, err := jwtManager.GenerateToken(testUser)
		assert.NoError(t, err)

		// Test valid token
		userID, err := jwtManager.ExtractUserIDFromToken(token.AccessToken)
		assert.NoError(t, err)
		assert.Equal(t, testUser.ID, userID)

		// Test invalid token
		userID, err = jwtManager.ExtractUserIDFromToken("invalid.token.string")
		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, userID)
	})

	t.Run("ExtractEmailFromToken", func(t *testing.T) {
		// Generate a token
		token, err := jwtManager.GenerateToken(testUser)
		assert.NoError(t, err)

		// Test valid token
		email, err := jwtManager.ExtractEmailFromToken(token.AccessToken)
		assert.NoError(t, err)
		assert.Equal(t, testUser.Email, email)

		// Test invalid token
		email, err = jwtManager.ExtractEmailFromToken("invalid.token.string")
		assert.Error(t, err)
		assert.Empty(t, email)
	})
}
