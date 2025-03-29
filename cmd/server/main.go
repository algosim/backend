package main

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

// @title           Auth Service API
// @version         1.0
// @description     Authentication service for Problem Tracking & Challenge System
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	r := gin.Default()

	// Swagger
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Setup auth routes
	userRepo := db.NewUserRepoPG()
	tokenRepo := db.NewTokenRepoPG()
	googleOAuth := oauth.NewGoogleOAuth(cfg)
	authUseCase := usecase.NewAuthUseCase(userRepo, tokenRepo, googleOAuth)
	authHandler := http.NewAuthHandler(authUseCase)

	http.SetupAuthRoutes(r, authHandler)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
