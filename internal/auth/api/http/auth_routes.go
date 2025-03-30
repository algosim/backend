package http

import (
	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes configures all authentication routes
func SetupAuthRoutes(r *gin.Engine, h *AuthHandler) {
	auth := r.Group("/api/v1/auth")
	{
		// OAuth routes
		auth.GET("/oauth/login", h.InitiateOAuthLogin)
		auth.GET("/oauth/callback", h.OAuthCallback)

		// Token management
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/logout", h.Logout)
	}
}
