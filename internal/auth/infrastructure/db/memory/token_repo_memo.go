package memory

import (
	"fmt"
	"sync"

	"github.com/algosim/backend/internal/auth/domain"
	"github.com/algosim/backend/internal/auth/repository"
	"github.com/google/uuid"
)

// TokenRepoMemo implements TokenRepository interface using in-memory storage
type TokenRepoMemo struct {
	tokens map[uuid.UUID]*domain.Token
	mu     sync.RWMutex
}

// NewTokenRepoMemo creates a new in-memory token repository
func NewTokenRepoMemo() *TokenRepoMemo {
	return &TokenRepoMemo{
		tokens: make(map[uuid.UUID]*domain.Token),
	}
}

// Create stores a new token
func (r *TokenRepoMemo) Create(token *domain.Token) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tokens[token.ID]; exists {
		return fmt.Errorf("token already exists")
	}

	r.tokens[token.ID] = token
	return nil
}

// FindByID retrieves a token by ID
func (r *TokenRepoMemo) FindByID(id uuid.UUID) (*domain.Token, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	token, exists := r.tokens[id]
	if !exists {
		return nil, fmt.Errorf("token not found")
	}

	return token, nil
}

// FindByRefreshToken retrieves a token by refresh token string
func (r *TokenRepoMemo) FindByRefreshToken(refreshToken string) (*domain.Token, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, token := range r.tokens {
		if token.RefreshToken == refreshToken {
			return token, nil
		}
	}

	return nil, fmt.Errorf("token not found")
}

// Delete removes a token
func (r *TokenRepoMemo) Delete(id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tokens[id]; !exists {
		return fmt.Errorf("token not found")
	}

	delete(r.tokens, id)
	return nil
}

// FindByUserID retrieves all tokens for a specific user
func (r *TokenRepoMemo) FindByUserID(userID uuid.UUID) ([]*domain.Token, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userTokens []*domain.Token
	for _, token := range r.tokens {
		if token.UserID == userID {
			userTokens = append(userTokens, token)
		}
	}

	return userTokens, nil
}

// DeleteByUserID removes all tokens for a specific user
func (r *TokenRepoMemo) DeleteByUserID(userID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Find all tokens for the user
	var tokensToDelete []uuid.UUID
	for id, token := range r.tokens {
		if token.UserID == userID {
			tokensToDelete = append(tokensToDelete, id)
		}
	}

	// Delete all found tokens
	for _, id := range tokensToDelete {
		delete(r.tokens, id)
	}

	return nil
}

// Ensure TokenRepoMemo implements TokenRepository interface
var _ repository.TokenRepository = (*TokenRepoMemo)(nil)
