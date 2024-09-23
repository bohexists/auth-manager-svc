package config

import (
	"log"
	"os"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	DBDSN      string
	JWTSecret  string
}

// LoadConfig loads environment variables from Docker environment
func LoadConfig() *Config {
	// Если файл .env не найден, используем переменные среды из контейнера
	config := &Config{
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", ""),
		DBHost:     getEnv("DB_HOST", ""),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBDSN:      getEnv("DB_DSN", ""),
		JWTSecret:  getEnv("JWT_SECRET", "secret-key"),
	}

	return config
}

// Helper function to retrieve environment variables or set a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("Environment variable %s not found, using default value: %s", key, defaultValue)
		return defaultValue
	}
	log.Printf("Environment variable %s loaded with value: %s", key, value)
	return value
}
