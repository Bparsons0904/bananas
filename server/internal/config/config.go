package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerPort     string
	DatabaseConfig DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func New() (Config, error) {
	config := Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DatabaseConfig: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "bananas_user"),
			Password: getEnv("DB_PASSWORD", "bananas_pass"),
			DBName:   getEnv("DB_NAME", "bananas_dev"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
	}

	return config, nil
}

func (c Config) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DatabaseConfig.Host,
		c.DatabaseConfig.Port,
		c.DatabaseConfig.User,
		c.DatabaseConfig.Password,
		c.DatabaseConfig.DBName,
		c.DatabaseConfig.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}