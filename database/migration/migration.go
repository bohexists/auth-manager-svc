package migration

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"

	"github.com/bohexists/auth-manager-svc/config"
)

// ApplyMigrations applies migrations to the given database connection
func ApplyMigrations(cfg *config.Config) {
	// Используем строку подключения к базе данных из конфигурации
	dbConnStr := cfg.DBDSN

	// Используем путь к миграциям и строку подключения
	m, err := migrate.New(
		"file:///app/database/migrations",
		dbConnStr)

	if err != nil {
		log.Fatal("Migration failed: ", err)
	}

	// Применяем миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Could not apply migrations: ", err)
	}
}
