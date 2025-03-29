package http

import (
	"fmt"
	"net/http"

	"github.com/algosim/backend/internal/auth/usecase"
	"github.com/gin-gonic/gin"
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
	fmt.Println("authURL", authURL)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// OAuthCallback handles the OAuth callback
// @Summary OAuth Callback
// @Description Handles the callback from OAuth provider
// @Tags auth
// @Accept json
// @Produce json
// @Param code body OAuthCallbackRequest true "OAuth callback data"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} map[string]string
// @Router /auth/oauth/callback [post]
func (h *AuthHandler) OAuthCallback(c *gin.Context) {
	var req OAuthCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authUseCase.HandleOAuthCallback(req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  token.RefreshToken, // TODO: Generate proper access token
		RefreshToken: token.RefreshToken,
	})
}

// RefreshToken handles token refresh
// @Summary Refresh Token
// @Description Generates a new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
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
		AccessToken: token.RefreshToken, // TODO: Generate proper access token
	})
}

// Logout handles user logout
// @Summary Logout
// @Description Invalidates the refresh token
// @Tags auth
// @Accept json
// @Produce json
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
