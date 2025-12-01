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
	Host          string
	Port          string
	User          string
	Password      string
	DBName        string
	SSLMode       string
	AdminUser     string
	AdminPassword string
}

func New() (Config, error) {
	dbUser := getEnv("DB_USER", "bananas_user")
	dbPassword := getEnv("DB_PASSWORD", "bananas_pass")

	config := Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DatabaseConfig: DatabaseConfig{
			Host:          getEnv("DB_HOST", "localhost"),
			Port:          getEnv("DB_PORT", "5432"),
			User:          dbUser,
			Password:      dbPassword,
			DBName:        getEnv("DB_NAME", "bananas_dev"),
			SSLMode:       getEnv("DB_SSL_MODE", "disable"),
			AdminUser:     getEnv("DB_ADMIN_USER", dbUser),
			AdminPassword: getEnv("DB_ADMIN_PASSWORD", dbPassword),
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