package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"

	"github.com/bohexists/auth-manager-svc/config"
	"github.com/bohexists/auth-manager-svc/database/migration"
	"github.com/bohexists/auth-manager-svc/internal/repositorys"
	"github.com/bohexists/auth-manager-svc/internal/services"
	"github.com/bohexists/auth-manager-svc/transport/http/handlers"
	"github.com/bohexists/auth-manager-svc/transport/http/middleware"
	"github.com/bohexists/auth-manager-svc/transport/http/router"
)

func main() {
	// Load configuration from environment variables
	cfg := config.LoadConfig()

	// Initialize database connection using DSN from config
	db, err := gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Log successful connection
	log.Println("Successfully connected to the database.")

	// Apply migrations
	log.Println("Starting to apply migrations...")
	migration.ApplyMigrations(cfg)

	// Initialize repositories and services
	log.Println("Initializing repositories and services...")
	userRepo := repositorys.NewUserRepository(db)
	JWTService := services.NewJWTService(cfg)
	authService := services.NewAuthService(userRepo, JWTService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, JWTService)

	// Initialize JWT middleware
	jwtMiddleware := middleware.JWTAuthMiddleware(JWTService)

	// Setup routes using external router setup
	router := routes.SetupRouter(authHandler, jwtMiddleware)

	// Run the server
	log.Println("Starting the server on port 8080...")
	router.Run(":8080")
}
