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
	// Load configuration from .env
	cfg := config.LoadConfig()

	// Initialize database connection using DSN from config
	db, err := gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Call raw database connection for migrations
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get raw database connection: ", err)
	}

	// Apply migrations
	migration.ApplyMigrations(sqlDB)

	// Initialize repositories and services
	userRepo := user.NewUserRepository(db)
	authService := auth.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// Setup routes using external router setup
	router := routes.SetupRouter(authHandler)

	// Run the server
	router.Run(":8080")
}
