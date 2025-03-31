package usecase

import (
	"testing"
	"time"

	"github.com/algosim/backend/configs"
	"github.com/algosim/backend/internal/auth/domain"
	"github.com/algosim/backend/internal/auth/infrastructure/oauth"
	"github.com/algosim/backend/internal/auth/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByOAuthProviderID(provider, providerID string) (*domain.User, error) {
	args := m.Called(provider, providerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockTokenRepository is a mock implementation of TokenRepository
type MockTokenRepository struct {
	mock.Mock
}

func (m *MockTokenRepository) Create(token *domain.Token) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockTokenRepository) FindByID(id uuid.UUID) (*domain.Token, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Token), args.Error(1)
}

func (m *MockTokenRepository) FindByUserID(userID uuid.UUID) ([]*domain.Token, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Token), args.Error(1)
}

func (m *MockTokenRepository) FindByRefreshToken(refreshToken string) (*domain.Token, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Token), args.Error(1)
}

func (m *MockTokenRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTokenRepository) DeleteByUserID(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

// MockGoogleOAuth is a mock implementation of Google OAuth
type MockGoogleOAuth struct {
	mock.Mock
}

func (m *MockGoogleOAuth) GetAuthURL(state string) string {
	args := m.Called(state)
	return args.String(0)
}

func (m *MockGoogleOAuth) ExchangeCodeForToken(code string) (*domain.Token, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Token), args.Error(1)
}

func (m *MockGoogleOAuth) GetUserInfo(accessToken string) (*oauth.GoogleUserInfo, error) {
	args := m.Called(accessToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*oauth.GoogleUserInfo), args.Error(1)
}

func (m *MockGoogleOAuth) CreateUserFromGoogleInfo(info *oauth.GoogleUserInfo) *domain.User {
	args := m.Called(info)
	return args.Get(0).(*domain.User)
}

func TestAuthUseCase(t *testing.T) {
	// Setup test configuration
	config := &configs.Config{
		Auth: struct {
			JWTSecret string `mapstructure:"jwt_secret" env:"AUTH_JWT_SECRET"`
			TokenTTL  int    `mapstructure:"token_ttl" env:"AUTH_TOKEN_TTL"`
		}{
			JWTSecret: "test-secret-key-123",
			TokenTTL:  3600, // 1 hour
		},
		GoogleOAuth: struct {
			ClientID     string   `mapstructure:"client_id" env:"GOOGLE_OAUTH_CLIENT_ID"`
			ClientSecret string   `mapstructure:"client_secret" env:"GOOGLE_OAUTH_CLIENT_SECRET"`
			RedirectURI  string   `mapstructure:"redirect_uri" env:"GOOGLE_OAUTH_REDIRECT_URI"`
			Scopes       []string `mapstructure:"scopes"`
		}{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURI:  "http://localhost:8080/callback",
		},
	}

	// Create test user
	testUser := &domain.User{
		ID:              uuid.New(),
		Email:           "test@example.com",
		OAuthProvider:   "google",
		OAuthProviderID: "test-google-id",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Create test token
	testToken := &domain.Token{
		ID:           uuid.New(),
		UserID:       testUser.ID,
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	t.Run("InitiateOAuthLogin", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockTokenRepository)
		mockGoogleOAuth := new(MockGoogleOAuth)
		authUseCase := usecase.NewAuthUseCase(mockUserRepo, mockTokenRepo, mockGoogleOAuth, config)

		state := "test-state"
		expectedURL := "https://accounts.google.com/o/oauth2/v2/auth?client_id=test-client-id&state=test-state"
		mockGoogleOAuth.On("GetAuthURL", state).Return(expectedURL)

		authURL := authUseCase.InitiateOAuthLogin(state)
		assert.Equal(t, expectedURL, authURL)
		mockGoogleOAuth.AssertExpectations(t)
	})

	t.Run("HandleOAuthCallback - New User", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockTokenRepository)
		mockGoogleOAuth := new(MockGoogleOAuth)
		authUseCase := usecase.NewAuthUseCase(mockUserRepo, mockTokenRepo, mockGoogleOAuth, config)

		// Setup expectations
		mockGoogleOAuth.On("ExchangeCodeForToken", "test-code").Return(testToken, nil)
		mockGoogleOAuth.On("GetUserInfo", testToken.AccessToken).Return(&oauth.GoogleUserInfo{
			ID:            "test-google-id",
			Email:         "test@example.com",
			VerifiedEmail: true,
			Name:          "Test User",
		}, nil)
		mockUserRepo.On("FindByOAuthProviderID", "google", "test-google-id").Return(nil, assert.AnError)
		mockUserRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)
		mockTokenRepo.On("Create", mock.AnythingOfType("*domain.Token")).Return(nil)

		// Test with mock Google OAuth response
		token, err := authUseCase.HandleOAuthCallback("test-code")
		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.NotEmpty(t, token.AccessToken)
		assert.NotEmpty(t, token.RefreshToken)

		mockUserRepo.AssertExpectations(t)
		mockTokenRepo.AssertExpectations(t)
		mockGoogleOAuth.AssertExpectations(t)
	})

	t.Run("HandleOAuthCallback - Existing User", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockTokenRepository)
		mockGoogleOAuth := new(MockGoogleOAuth)
		authUseCase := usecase.NewAuthUseCase(mockUserRepo, mockTokenRepo, mockGoogleOAuth, config)

		// Setup expectations
		mockGoogleOAuth.On("ExchangeCodeForToken", "test-code").Return(testToken, nil)
		mockGoogleOAuth.On("GetUserInfo", testToken.AccessToken).Return(&oauth.GoogleUserInfo{
			ID:            "test-google-id",
			Email:         "test@example.com",
			VerifiedEmail: true,
			Name:          "Test User",
		}, nil)
		mockUserRepo.On("FindByOAuthProviderID", "google", "test-google-id").Return(testUser, nil)
		mockTokenRepo.On("Create", mock.AnythingOfType("*domain.Token")).Return(nil)

		// Test with mock Google OAuth response
		token, err := authUseCase.HandleOAuthCallback("test-code")
		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.NotEmpty(t, token.AccessToken)
		assert.NotEmpty(t, token.RefreshToken)

		mockUserRepo.AssertExpectations(t)
		mockTokenRepo.AssertExpectations(t)
		mockGoogleOAuth.AssertExpectations(t)
	})

	// t.Run("RefreshToken", func(t *testing.T) {
	// 	mockUserRepo := new(MockUserRepository)
	// 	mockTokenRepo := new(MockTokenRepository)
	// 	googleOAuth := oauth.NewGoogleOAuth(config)
	// 	authUseCase := usecase.NewAuthUseCase(mockUserRepo, mockTokenRepo, googleOAuth, config)

	// 	// Setup expectations
	// 	mockTokenRepo.On("FindByRefreshToken", "test-refresh-token").Return(testToken, nil)
	// 	mockUserRepo.On("FindByID", testUser.ID).Return(testUser, nil)
	// 	mockTokenRepo.On("Create", mock.AnythingOfType("*domain.Token")).Return(nil)
	// 	mockTokenRepo.On("Delete", testToken.ID).Return(nil)

	// 	// Test token refresh
	// 	newToken, err := authUseCase.RefreshToken("test-refresh-token")
	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, newToken)
	// 	assert.NotEmpty(t, newToken.AccessToken)
	// 	assert.NotEmpty(t, newToken.RefreshToken)

	// 	mockUserRepo.AssertExpectations(t)
	// 	mockTokenRepo.AssertExpectations(t)
	// })

	// t.Run("Logout", func(t *testing.T) {
	// 	mockUserRepo := new(MockUserRepository)
	// 	mockTokenRepo := new(MockTokenRepository)
	// 	googleOAuth := oauth.NewGoogleOAuth(config)
	// 	authUseCase := usecase.NewAuthUseCase(mockUserRepo, mockTokenRepo, googleOAuth, config)

	// 	// Setup expectations
	// 	mockTokenRepo.On("FindByRefreshToken", "test-refresh-token").Return(testToken, nil)
	// 	mockTokenRepo.On("Delete", testToken.ID).Return(nil)

	// 	// Test logout
	// 	err := authUseCase.Logout("test-refresh-token")
	// 	assert.NoError(t, err)

	// 	mockTokenRepo.AssertExpectations(t)
	// })

	// t.Run("ValidateToken", func(t *testing.T) {
	// 	mockUserRepo := new(MockUserRepository)
	// 	mockTokenRepo := new(MockTokenRepository)
	// 	googleOAuth := oauth.NewGoogleOAuth(config)
	// 	authUseCase := usecase.NewAuthUseCase(mockUserRepo, mockTokenRepo, googleOAuth, config)

	// 	// Setup expectations
	// 	mockUserRepo.On("FindByID", testUser.ID).Return(testUser, nil)

	// 	// Test token validation
	// 	user, err := authUseCase.ValidateToken(testToken.AccessToken)
	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, user)
	// 	assert.Equal(t, testUser.ID, user.ID)
	// 	assert.Equal(t, testUser.Email, user.Email)

	// 	mockUserRepo.AssertExpectations(t)
	// })
}
