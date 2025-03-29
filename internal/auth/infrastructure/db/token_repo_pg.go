package db

import (
	"sync"

	"github.com/algosim/backend/internal/auth/domain"
	"github.com/algosim/backend/internal/auth/repository"
	"github.com/google/uuid"
)

// TokenRepoPG implements TokenRepository interface with in-memory storage
type TokenRepoPG struct {
	mu     sync.RWMutex
	tokens map[uuid.UUID]*domain.Token
}

// NewTokenRepoPG creates a new instance of TokenRepoPG
func NewTokenRepoPG() *TokenRepoPG {
	return &TokenRepoPG{
		tokens: make(map[uuid.UUID]*domain.Token),
	}
}

func (r *TokenRepoPG) Create(token *domain.Token) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tokens[token.ID] = token
	return nil
}

func (r *TokenRepoPG) FindByID(id uuid.UUID) (*domain.Token, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	token, exists := r.tokens[id]
	if !exists {
		return nil, domain.ErrTokenNotFound
	}

	return token, nil
}

func (r *TokenRepoPG) FindByUserID(userID uuid.UUID) ([]*domain.Token, error) {
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

func (r *TokenRepoPG) FindByRefreshToken(refreshToken string) (*domain.Token, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, token := range r.tokens {
		if token.RefreshToken == refreshToken {
			return token, nil
		}
	}

	return nil, domain.ErrTokenNotFound
}

func (r *TokenRepoPG) Delete(id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tokens[id]; !exists {
		return domain.ErrTokenNotFound
	}

	delete(r.tokens, id)
	return nil
}

func (r *TokenRepoPG) DeleteByUserID(userID uuid.UUID) error {
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

// Ensure TokenRepoPG implements TokenRepository interface
var _ repository.TokenRepository = (*TokenRepoPG)(nil)
