package usecase

import (
	"time"

	"github.com/algosim/backend/internal/auth/domain"
	"github.com/algosim/backend/internal/auth/infrastructure/oauth"
	"github.com/algosim/backend/internal/auth/repository"
)

// AuthUseCase handles authentication business logic
type AuthUseCase struct {
	userRepo    repository.UserRepository
	tokenRepo   repository.TokenRepository
	googleOAuth *oauth.GoogleOAuth
}

// NewAuthUseCase creates a new AuthUseCase instance
func NewAuthUseCase(userRepo repository.UserRepository, tokenRepo repository.TokenRepository, googleOAuth *oauth.GoogleOAuth) *AuthUseCase {
	return &AuthUseCase{
		userRepo:    userRepo,
		tokenRepo:   tokenRepo,
		googleOAuth: googleOAuth,
	}
}

// InitiateOAuthLogin generates the OAuth login URL
func (u *AuthUseCase) InitiateOAuthLogin(state string) string {
	return u.googleOAuth.GetAuthURL(state)
}

// HandleOAuthCallback processes the OAuth callback
func (u *AuthUseCase) HandleOAuthCallback(code string) (*domain.Token, error) {
	// Exchange code for token
	token, err := u.googleOAuth.ExchangeCodeForToken(code)
	if err != nil {
		return nil, err
	}

	// Get user info from Google
	userInfo, err := u.googleOAuth.GetUserInfo(token.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Check if user exists
	existingUser, err := u.userRepo.FindByOAuthProviderID("google", userInfo.ID)
	if err != nil {
		// Create new user if not found
		newUser := u.googleOAuth.CreateUserFromGoogleInfo(userInfo)
		if err := u.userRepo.Create(newUser); err != nil {
			return nil, err
		}
		token.UserID = newUser.ID
	} else {
		token.UserID = existingUser.ID
	}

	// Save token
	if err := u.tokenRepo.Create(token); err != nil {
		return nil, err
	}

	return token, nil
}

// RefreshToken generates a new access token using refresh token
func (u *AuthUseCase) RefreshToken(refreshToken string) (*domain.Token, error) {
	// Find token
	token, err := u.tokenRepo.FindByRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Check if token is expired
	if time.Now().After(token.ExpiresAt) {
		return nil, domain.ErrTokenExpired
	}

	// TODO: Implement token refresh logic with Google OAuth
	// For now, just return the existing token
	return token, nil
}

// Logout invalidates the refresh token
func (u *AuthUseCase) Logout(refreshToken string) error {
	// Find token
	token, err := u.tokenRepo.FindByRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	// Delete token
	return u.tokenRepo.Delete(token.ID)
}
