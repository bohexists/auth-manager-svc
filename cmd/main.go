package main

import (
	"github.com/bohexists/auth-manager-svc/database/migration"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"

	"github.com/bohexists/auth-manager-svc/config"
	"github.com/bohexists/auth-manager-svc/internal/auth"
	"github.com/bohexists/auth-manager-svc/internal/user"
	"github.com/bohexists/auth-manager-svc/transport/http/handlers"
	"github.com/bohexists/auth-manager-svc/transport/http/router"
)

func main() {
	// Load configuration from environment variables
	cfg := config.LoadConfig()

	// Log the database connection string
	log.Printf("Connecting to the database using DSN: %s", cfg.DBDSN)

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
	userRepo := user.NewUserRepository(db)
	authService := auth.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// Setup routes using external router setup
	log.Println("Setting up routes...")
	router := routes.SetupRouter(authHandler)

	// Run the server
	log.Println("Starting the server on port 8080...")
	router.Run(":8080")
}
