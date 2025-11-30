package main

import (
	"bananas/internal/config"
	"bananas/internal/database"
	"bananas/internal/logger"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/migration/main.go [up|down|seed]")
		os.Exit(1)
	}

	command := os.Args[1]
	log := logger.New("migration")

	cfg, err := config.New()
	if err != nil {
		log.Er("failed to initialize config", err)
		os.Exit(1)
	}

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
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}

	if err != nil {
		log.Er("migration failed", err)
		os.Exit(1)
	}

	log.Info("Migration completed successfully")
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