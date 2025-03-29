package repository

import (
	"github.com/algosim/backend/internal/auth/domain"
	"github.com/google/uuid"
)

// TokenRepository defines the interface for token persistence operations
type TokenRepository interface {
	// Create creates a new refresh token
	Create(token *domain.Token) error
	// FindByID finds a token by its ID
	FindByID(id uuid.UUID) (*domain.Token, error)
	// FindByUserID finds all tokens for a user
	FindByUserID(userID uuid.UUID) ([]*domain.Token, error)
	// FindByRefreshToken finds a token by its refresh token string
	FindByRefreshToken(refreshToken string) (*domain.Token, error)
	// Delete deletes a token by its ID
	Delete(id uuid.UUID) error
	// DeleteByUserID deletes all tokens for a user
	DeleteByUserID(userID uuid.UUID) error
}
