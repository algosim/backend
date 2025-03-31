package memory

import (
	"fmt"
	"sync"

	"github.com/algosim/backend/internal/auth/domain"
	"github.com/algosim/backend/internal/auth/repository"
	"github.com/google/uuid"
)

// UserRepoMemo implements UserRepository interface using in-memory storage
type UserRepoMemo struct {
	users map[uuid.UUID]*domain.User
	mu    sync.RWMutex
}

// NewUserRepoMemo creates a new in-memory user repository
func NewUserRepoMemo() *UserRepoMemo {
	return &UserRepoMemo{
		users: make(map[uuid.UUID]*domain.User),
	}
}

// Create stores a new user
func (r *UserRepoMemo) Create(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return fmt.Errorf("user already exists")
	}

	r.users[user.ID] = user
	return nil
}

// FindByID retrieves a user by ID
func (r *UserRepoMemo) FindByID(id uuid.UUID) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

// FindByEmail retrieves a user by email
func (r *UserRepoMemo) FindByEmail(email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, fmt.Errorf("user not found")
}

// FindByOAuthProviderID retrieves a user by OAuth provider ID
func (r *UserRepoMemo) FindByOAuthProviderID(provider, providerID string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.OAuthProvider == provider && user.OAuthProviderID == providerID {
			return user, nil
		}
	}

	return nil, fmt.Errorf("user not found")
}

// Update updates an existing user
func (r *UserRepoMemo) Update(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return fmt.Errorf("user not found")
	}

	r.users[user.ID] = user
	return nil
}

// Delete removes a user
func (r *UserRepoMemo) Delete(id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return fmt.Errorf("user not found")
	}

	delete(r.users, id)
	return nil
}

// Ensure UserRepoMemo implements UserRepository interface
var _ repository.UserRepository = (*UserRepoMemo)(nil)
