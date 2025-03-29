package repository

import (
	"github.com/algosim/backend/internal/auth/domain"
	"github.com/google/uuid"
)

// UserRepository defines the interface for user persistence operations
type UserRepository interface {
	// Create creates a new user
	Create(user *domain.User) error
	// FindByID finds a user by their ID
	FindByID(id uuid.UUID) (*domain.User, error)
	// FindByEmail finds a user by their email
	FindByEmail(email string) (*domain.User, error)
	// FindByOAuthProviderID finds a user by their OAuth provider ID
	FindByOAuthProviderID(provider, providerID string) (*domain.User, error)
	// Update updates an existing user
	Update(user *domain.User) error
	// Delete deletes a user by their ID
	Delete(id uuid.UUID) error
}
