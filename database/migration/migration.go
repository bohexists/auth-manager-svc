package migration

import (
	"github.com/bohexists/auth-manager-svc/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

// ApplyMigrations applies migrations to the given database connection
func ApplyMigrations(cfg *config.Config) {
	// Log database connection string
	log.Printf("Starting migration with DB connection string: %s", cfg.DBDSN)

	// Get database connection string
	dbConnStr := cfg.DBDSN

	// Log migration path
	migrationPath := "file:///app/database/migrations"
	log.Printf("Using migration path: %s", migrationPath)

	// Log migrations
	log.Println("Initializing migrations...")

	// Log connection
	log.Println("Connecting to database...")

	m, err := migrate.New(migrationPath, dbConnStr)
	if err != nil {
		log.Fatalf("Migration initialization failed: %v", err)
	}

	// Apply migrations with log
	log.Println("Applying migrations...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Could not apply migrations: %v", err)
	} else if err == migrate.ErrNoChange {
		log.Println("No new migrations to apply")
	} else {
		log.Println("Migrations applied successfully")
	}
}
