package usecase

import (
	"fmt"

	"github.com/algosim/backend/configs"
	"github.com/algosim/backend/internal/auth/domain"
	"github.com/algosim/backend/internal/auth/infrastructure/jwt"
	"github.com/algosim/backend/internal/auth/infrastructure/oauth"
	"github.com/algosim/backend/internal/auth/repository"
)

// AuthUseCase handles authentication business logic
type AuthUseCase struct {
	userRepo    repository.UserRepository
	tokenRepo   repository.TokenRepository
	googleOAuth oauth.GoogleOAuth
	jwtManager  *jwt.JWTManager
}

// NewAuthUseCase creates a new AuthUseCase instance
func NewAuthUseCase(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	googleOAuth oauth.GoogleOAuth,
	config *configs.Config,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:    userRepo,
		tokenRepo:   tokenRepo,
		googleOAuth: googleOAuth,
		jwtManager:  jwt.NewJWTManager(config),
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
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get user info from Google
	userInfo, err := u.googleOAuth.GetUserInfo(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Check if user exists
	user, err := u.userRepo.FindByOAuthProviderID("google", userInfo.ID)
	if err != nil {
		// Create new user if not found
		newUser := u.googleOAuth.CreateUserFromGoogleInfo(userInfo)
		if err := u.userRepo.Create(newUser); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		user = newUser
	}

	token, err = u.jwtManager.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Store refresh token
	if err := u.tokenRepo.Create(token); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return token, nil
}

// RefreshToken generates a new access token using a refresh token
func (u *AuthUseCase) RefreshToken(refreshTokenString string) (*domain.Token, error) {
	// Find refresh token
	token, err := u.tokenRepo.FindByRefreshToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to find refresh token: %w", err)
	}

	// Validate refresh token
	if err := u.jwtManager.ValidateRefreshToken(token); err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Get user
	user, err := u.userRepo.FindByID(token.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Generate new refresh token
	newToken, err := u.jwtManager.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	// Store new refresh token
	if err := u.tokenRepo.Create(newToken); err != nil {
		return nil, fmt.Errorf("failed to store new refresh token: %w", err)
	}

	// Delete old refresh token
	if err := u.tokenRepo.Delete(token.ID); err != nil {
		return nil, fmt.Errorf("failed to delete old refresh token: %w", err)
	}

	return newToken, nil
}

// Logout invalidates the refresh token
func (u *AuthUseCase) Logout(refreshTokenString string) error {
	// Find refresh token
	token, err := u.tokenRepo.FindByRefreshToken(refreshTokenString)
	if err != nil {
		return fmt.Errorf("failed to find refresh token: %w", err)
	}

	// Delete refresh token
	if err := u.tokenRepo.Delete(token.ID); err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	return nil
}

// ValidateToken validates an access token and returns the user information
func (u *AuthUseCase) ValidateToken(tokenString string) (*domain.User, error) {
	user, err := u.jwtManager.ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	// Verify user exists in database
	dbUser, err := u.userRepo.FindByID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return dbUser, nil
}
