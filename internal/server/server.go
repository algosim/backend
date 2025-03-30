package server

import (
	"fmt"
	"log"

	"github.com/algosim/backend/configs"
	"github.com/algosim/backend/docs"
	"github.com/algosim/backend/internal/auth/api/http"
	"github.com/algosim/backend/internal/auth/infrastructure/db"
	"github.com/algosim/backend/internal/auth/infrastructure/oauth"
	"github.com/algosim/backend/internal/auth/usecase"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server represents the HTTP server
type Server struct {
	router *gin.Engine
	config *configs.Config
}

// NewServer creates a new Server instance
func NewServer(config *configs.Config) *Server {
	return &Server{
		router: gin.Default(),
		config: config,
	}
}

// SetupRoutes configures all routes for the server
func (s *Server) SetupRoutes() {
	// Swagger
	docs.SwaggerInfo.BasePath = "/api/v1"
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Initialize repositories
	userRepo := db.NewUserRepoPG()
	tokenRepo := db.NewTokenRepoPG()

	// Initialize Google OAuth
	googleOAuth := oauth.NewGoogleOAuth(s.config)

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, tokenRepo, googleOAuth, s.config)

	// Initialize handlers
	authHandler := http.NewAuthHandler(authUseCase)

	// Setup auth routes
	http.SetupAuthRoutes(s.router, authHandler)
}

// Run starts the server
func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	log.Printf("Server starting on %s", addr)
	return s.router.Run(addr)
}
