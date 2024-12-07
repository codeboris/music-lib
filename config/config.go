package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port          string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPass        string
	DBName        string
	SSLMode       string
	ConfigMigrate ConfigMigrate
}

type ConfigMigrate struct {
	MigrationPath string
	DatabaseURL   string
}

func LoadConfig() Config {
	cfg := Config{
		Port:    getEnv("APP_PORT", "8000"),
		DBHost:  getEnv("DB_HOST", "postgres"),
		DBPort:  getEnv("DB_PORT", "5432"),
		DBUser:  getEnv("DB_USER", "postgres"),
		DBName:  getEnv("DB_NAME", "postgres"),
		DBPass:  getEnv("DB_PASSWORD", "12345"),
		SSLMode: getEnv("DB_SSL_MODE", "disable"),
	}

	cfg.ConfigMigrate.MigrationPath = getEnv("MIGRATE_PATH", "file://migrations")

	cfg.ConfigMigrate.DatabaseURL = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.SSLMode,
	)

	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
