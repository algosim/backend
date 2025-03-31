package http

import (
	"net/http"
	"time"

	"github.com/algosim/backend/internal/auth/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

// InitiateOAuthLogin handles the OAuth login initiation
// @Summary Initiate OAuth Login
// @Description Redirects to Google OAuth login page
// @Tags auth
// @Accept json
// @Produce json
// @Param provider query string true "OAuth provider (e.g., google)"
// @Success 302 {string} string "Redirect to OAuth provider"
// @Failure 400 {object} map[string]string
// @Router /auth/oauth/login [get]
func (h *AuthHandler) InitiateOAuthLogin(c *gin.Context) {
	provider := c.Query("provider")
	if provider != "google" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported provider"})
		return
	}

	// Generate state for CSRF protection
	state := "random-state" // TODO: Generate proper random state
	authURL := h.authUseCase.InitiateOAuthLogin(state)

	// Instead of redirecting, return the URL to the frontend
	c.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
	})
}

// OAuthCallback handles the OAuth callback
// @Summary OAuth Callback
// @Description Handles the callback from OAuth provider
// @Tags auth
// @Accept json
// @Produce json
// @Param code query string true "Authorization code from OAuth provider"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} map[string]string
// @Router /auth/oauth/callback [get]
func (h *AuthHandler) OAuthCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorization code is required"})
		return
	}

	token, err := h.authUseCase.HandleOAuthCallback(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return tokens to the frontend
	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	})
}

// RefreshToken handles token refresh
// @Summary Refresh Token
// @Description Generates a new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param token body RefreshTokenRequest true "Refresh token data"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} map[string]string
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authUseCase.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	})
}

// Logout handles user logout
// @Summary Logout
// @Description Invalidates the refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param token body LogoutRequest true "Logout data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authUseCase.Logout(req.RefreshToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// ValidateToken handles token validation
// @Summary Validate Token
// @Description Validates the access token and returns user information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserResponse
// @Failure 401 {object} map[string]string
// @Router /auth/validate [get]
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	// Get token from Authorization header
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no token provided"})
		return
	}

	// Remove "Bearer " prefix if present
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	user, err := h.authUseCase.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:               user.ID,
		Email:            user.Email,
		CodeforcesHandle: user.CodeforcesHandle,
		AtcoderHandle:    user.AtcoderHandle,
		OAuthProvider:    user.OAuthProvider,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
	})
}

// UserResponse represents the user information response
type UserResponse struct {
	ID               uuid.UUID `json:"id"`
	Email            string    `json:"email"`
	CodeforcesHandle string    `json:"codeforces_handle,omitempty"`
	AtcoderHandle    string    `json:"atcoder_handle,omitempty"`
	OAuthProvider    string    `json:"oauth_provider"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Request types
type OAuthCallbackRequest struct {
	Code string `json:"code" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}
