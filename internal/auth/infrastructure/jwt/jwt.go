package jwt

import (
	"fmt"
	"time"

	"github.com/algosim/backend/configs"
	"github.com/algosim/backend/internal/auth/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents the JWT claims structure
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT operations
type JWTManager struct {
	secretKey       []byte
	tokenTTL        time.Duration
	refreshTokenTTL time.Duration
}

// NewJWTManager creates a new JWT manager instance
func NewJWTManager(config *configs.Config) *JWTManager {
	return &JWTManager{
		secretKey:       []byte(config.Auth.JWTSecret),
		tokenTTL:        time.Duration(config.Auth.TokenTTL) * time.Second,
		refreshTokenTTL: 24 * time.Hour, // Refresh tokens last 24 hours
	}
}

// GenerateAccessToken creates a new access token
func (m *JWTManager) GenerateToken(user *domain.User) (*domain.Token, error) {
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString(m.secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTokenString := uuid.New().String()
	expiresAt := time.Now().Add(m.refreshTokenTTL)

	token := domain.NewToken(user.ID, accessTokenString, refreshTokenString, expiresAt)
	return token, nil
}

// ValidateAccessToken validates an access token and returns the claims
func (m *JWTManager) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ValidateRefreshToken validates a refresh token
func (m *JWTManager) ValidateRefreshToken(token *domain.Token) error {
	if time.Now().After(token.ExpiresAt) {
		return fmt.Errorf("refresh token expired")
	}
	return nil
}

// ExtractUserIDFromToken extracts the user ID from a valid access token
func (m *JWTManager) ExtractUserIDFromToken(tokenString string) (uuid.UUID, error) {
	claims, err := m.ValidateAccessToken(tokenString)
	if err != nil {
		return uuid.Nil, err
	}
	return claims.UserID, nil
}

// ExtractEmailFromToken extracts the email from a valid access token
func (m *JWTManager) ExtractEmailFromToken(tokenString string) (string, error) {
	claims, err := m.ValidateAccessToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Email, nil
}

// ValidateToken validates an access token and returns the user information
func (m *JWTManager) ValidateToken(tokenString string) (*domain.User, error) {
	claims, err := m.ValidateAccessToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return &domain.User{
		ID:    claims.UserID,
		Email: claims.Email,
	}, nil
}
