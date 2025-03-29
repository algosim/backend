package db

import (
	"sync"

	"github.com/algosim/backend/internal/auth/domain"
	"github.com/algosim/backend/internal/auth/repository"
	"github.com/google/uuid"
)

// UserRepoPG implements UserRepository interface with in-memory storage
type UserRepoPG struct {
	mu    sync.RWMutex
	users map[uuid.UUID]*domain.User
}

// NewUserRepoPG creates a new instance of UserRepoPG
func NewUserRepoPG() *UserRepoPG {
	return &UserRepoPG{
		users: make(map[uuid.UUID]*domain.User),
	}
}

func (r *UserRepoPG) Create(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if user already exists
	if _, exists := r.users[user.ID]; exists {
		return domain.ErrUserAlreadyExists
	}

	// Store user
	r.users[user.ID] = user
	return nil
}

func (r *UserRepoPG) FindByID(id uuid.UUID) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

func (r *UserRepoPG) FindByEmail(email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, domain.ErrUserNotFound
}

func (r *UserRepoPG) FindByOAuthProviderID(provider, providerID string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.OAuthProvider == provider && user.OAuthProviderID == providerID {
			return user, nil
		}
	}

	return nil, domain.ErrUserNotFound
}

func (r *UserRepoPG) Update(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return domain.ErrUserNotFound
	}

	r.users[user.ID] = user
	return nil
}

func (r *UserRepoPG) Delete(id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return domain.ErrUserNotFound
	}

	delete(r.users, id)
	return nil
}

// Ensure UserRepoPG implements UserRepository interface
var _ repository.UserRepository = (*UserRepoPG)(nil)
