package main

import (
	"bananas/internal/config"
	"bananas/internal/database"
	"bananas/internal/logger"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/migration/main.go [create-db|up|down|seed]")
		os.Exit(1)
	}

	command := os.Args[1]
	log := logger.New("migration")

	cfg, err := config.New()
	if err != nil {
		log.Er("failed to initialize config", err)
		os.Exit(1)
	}

	switch command {
	case "create-db":
		err = createDatabase(cfg, log)
	case "up", "down", "seed":
		db, err := database.New(cfg)
		if err != nil {
			log.Er("failed to connect to database", err)
			os.Exit(1)
		}
		defer db.Close()

		switch command {
		case "up":
			err = migrateUp(db)
		case "down":
			err = migrateDown(db)
		case "seed":
			err = seed(db)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Available commands: create-db, up, down, seed")
		os.Exit(1)
	}

	if err != nil {
		log.Er("command failed", err)
		os.Exit(1)
	}

	log.Info("Command completed successfully")
}

func createDatabase(cfg config.Config, log logger.Logger) error {
	log.Info("Creating database if it doesn't exist")

	adminDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		cfg.DatabaseConfig.Host,
		cfg.DatabaseConfig.Port,
		cfg.DatabaseConfig.AdminUser,
		cfg.DatabaseConfig.AdminPassword,
		cfg.DatabaseConfig.SSLMode,
	)

	adminDB, err := sql.Open("postgres", adminDSN)
	if err != nil {
		log.Er("failed to connect to postgres database", err)
		return err
	}
	defer adminDB.Close()

	if err := adminDB.Ping(); err != nil {
		log.Er("failed to ping postgres database", err)
		return err
	}

	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", cfg.DatabaseConfig.DBName)
	err = adminDB.QueryRow(query).Scan(&exists)
	if err != nil {
		log.Er("failed to check if database exists", err)
		return err
	}

	if exists {
		log.Info("Database already exists", "database", cfg.DatabaseConfig.DBName)
		return nil
	}

	createDBQuery := fmt.Sprintf("CREATE DATABASE %s", cfg.DatabaseConfig.DBName)
	_, err = adminDB.Exec(createDBQuery)
	if err != nil {
		log.Er("failed to create database", err)
		return err
	}

	log.Info("Database created successfully", "database", cfg.DatabaseConfig.DBName)

	if cfg.DatabaseConfig.User != cfg.DatabaseConfig.AdminUser {
		var userExists bool
		userQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_roles WHERE rolname = '%s')", cfg.DatabaseConfig.User)
		err = adminDB.QueryRow(userQuery).Scan(&userExists)
		if err != nil {
			log.Er("failed to check if user exists", err)
			return err
		}

		if !userExists {
			createUserQuery := fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", cfg.DatabaseConfig.User, cfg.DatabaseConfig.Password)
			_, err = adminDB.Exec(createUserQuery)
			if err != nil {
				log.Er("failed to create user", err)
				return err
			}
			log.Info("Database user created", "user", cfg.DatabaseConfig.User)
		}

		grantQuery := fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s", cfg.DatabaseConfig.DBName, cfg.DatabaseConfig.User)
		_, err = adminDB.Exec(grantQuery)
		if err != nil {
			log.Er("failed to grant privileges", err)
			return err
		}
		log.Info("Database privileges granted", "user", cfg.DatabaseConfig.User)
	}

	return nil
}

func migrateUp(db *database.DB) error {
	log := db.Logger.Function("migrateUp")

	queries := []string{
		`CREATE TABLE IF NOT EXISTS frameworks (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			type VARCHAR(50) NOT NULL,
			description TEXT,
			enabled BOOLEAN DEFAULT true,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS test_results (
			id SERIAL PRIMARY KEY,
			framework VARCHAR(255) NOT NULL,
			test_type VARCHAR(255) NOT NULL,
			execution_ms INTEGER NOT NULL,
			success BOOLEAN DEFAULT true,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_test_results_framework ON test_results(framework)`,
		`CREATE INDEX IF NOT EXISTS idx_test_results_created_at ON test_results(created_at)`,
	}

	for _, query := range queries {
		_, err := db.SQL.Exec(query)
		if err != nil {
			log.Er("failed to execute migration query", err)
			return err
		}
	}

	log.Info("Database migration completed successfully")
	return nil
}

func migrateDown(db *database.DB) error {
	log := db.Logger.Function("migrateDown")

	queries := []string{
		`DROP TABLE IF EXISTS test_results`,
		`DROP TABLE IF EXISTS frameworks`,
	}

	for _, query := range queries {
		_, err := db.SQL.Exec(query)
		if err != nil {
			log.Er("failed to execute rollback query", err)
			return err
		}
	}

	log.Info("Database rollback completed successfully")
	return nil
}

func seed(db *database.DB) error {
	log := db.Logger.Function("seed")

	// Clear existing data first
	_, err := db.SQL.Exec("DELETE FROM test_results")
	if err != nil {
		log.Er("failed to clear test results", err)
		return err
	}

	_, err = db.SQL.Exec("DELETE FROM frameworks")
	if err != nil {
		log.Er("failed to clear frameworks", err)
		return err
	}

	// Insert backend frameworks
	backendFrameworks := []string{"standard", "gin", "fiber", "echo", "chi", "gorilla"}
	for _, name := range backendFrameworks {
		query := `INSERT INTO frameworks (name, type, description) VALUES ($1, 'backend', $2)`
		_, err := db.SQL.Exec(query, name, fmt.Sprintf("Go %s web framework", name))
		if err != nil {
			log.Er("failed to insert backend framework", err)
			return err
		}
	}

	// Insert frontend frameworks
	frontendFrameworks := []string{"react", "vue", "svelte", "solid", "angular", "htmx", "templ"}
	for _, name := range frontendFrameworks {
		query := `INSERT INTO frameworks (name, type, description) VALUES ($1, 'frontend', $2)`
		_, err := db.SQL.Exec(query, name, fmt.Sprintf("%s frontend framework", name))
		if err != nil {
			log.Er("failed to insert frontend framework", err)
			return err
		}
	}

	log.Info("Database seeded successfully")
	return nil
}