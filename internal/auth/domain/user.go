package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents the core user entity in the domain
type User struct {
	ID               uuid.UUID
	Email            string
	CodeforcesHandle string
	AtcoderHandle    string
	OAuthProvider    string
	OAuthProviderID  string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NewUser creates a new User instance with default values
func NewUser(email, oauthProvider, oauthProviderID string) *User {
	now := time.Now()
	return &User{
		ID:              uuid.New(),
		Email:           email,
		OAuthProvider:   oauthProvider,
		OAuthProviderID: oauthProviderID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

}
