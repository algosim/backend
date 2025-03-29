package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/algosim/backend/configs"
	"github.com/algosim/backend/internal/auth/domain"
	"github.com/google/uuid"
)

const (
	googleAuthURL     = "https://accounts.google.com/o/oauth2/v2/auth"
	googleTokenURL    = "https://oauth2.googleapis.com/token"
	googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
)

// GoogleUserInfo represents the user information returned by Google
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// GoogleOAuth handles Google OAuth authentication
type GoogleOAuth struct {
	config *configs.Config
}

// NewGoogleOAuth creates a new Google OAuth handler
func NewGoogleOAuth(config *configs.Config) *GoogleOAuth {
	return &GoogleOAuth{
		config: config,
	}
}

// GetAuthURL generates the Google OAuth authorization URL
func (g *GoogleOAuth) GetAuthURL(state string) string {
	params := url.Values{}
	params.Add("client_id", g.config.GoogleOAuth.ClientID)
	params.Add("redirect_uri", g.config.GoogleOAuth.RedirectURI)
	params.Add("response_type", "code")
	params.Add("scope", "openid email profile")
	params.Add("state", state)
	params.Add("access_type", "offline")
	params.Add("prompt", "consent")

	return fmt.Sprintf("%s?%s", googleAuthURL, params.Encode())
}

// ExchangeCodeForToken exchanges the authorization code for access and refresh tokens
func (g *GoogleOAuth) ExchangeCodeForToken(code string) (*domain.Token, error) {
	params := url.Values{}
	params.Add("client_id", g.config.GoogleOAuth.ClientID)
	params.Add("client_secret", g.config.GoogleOAuth.ClientSecret)
	params.Add("code", code)
	params.Add("grant_type", "authorization_code")
	params.Add("redirect_uri", g.config.GoogleOAuth.RedirectURI)

	resp, err := http.PostForm(googleTokenURL, params)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to exchange code for token: %s", string(body))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	// Create domain token
	token := &domain.Token{
		ID:           uuid.New(),
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
	}

	return token, nil
}

// GetUserInfo retrieves user information from Google
func (g *GoogleOAuth) GetUserInfo(accessToken string) (*GoogleUserInfo, error) {
	req, err := http.NewRequest("GET", googleUserInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: %s", string(body))
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}

// CreateUserFromGoogleInfo creates a domain User from Google user info
func (g *GoogleOAuth) CreateUserFromGoogleInfo(info *GoogleUserInfo) *domain.User {
	return &domain.User{
		ID:              uuid.New(),
		Email:           info.Email,
		OAuthProvider:   "google",
		OAuthProviderID: info.ID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}
