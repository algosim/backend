package domain

import (
	"time"

	"github.com/google/uuid"
)

// Token represents the refresh token entity in the domain
type Token struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	RefreshToken string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

// NewToken creates a new Token instance with default values
func NewToken(userID uuid.UUID, refreshToken string, expiresAt time.Time) *Token {
	now := time.Now()
	return &Token{
		ID:           uuid.New(),
		UserID:       userID,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		CreatedAt:    now,
	}
}
