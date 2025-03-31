package main

import (
	"log"

	"github.com/algosim/backend/configs"
	"github.com/algosim/backend/internal/server"
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

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Load configuration
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create and setup server
	srv := server.NewServer(cfg)
	srv.SetupRoutes()

	// Run server
	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
