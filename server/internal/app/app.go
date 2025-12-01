package app

import (
	"bananas/internal/config"
	"bananas/internal/controllers"
	"bananas/internal/database"
	"bananas/internal/logger"
	"bananas/internal/repositories"
	"bananas/internal/services"
)

type App struct {
	Database    *database.DB
	Config      config.Config
	Services    *services.Service
	RepoManager *repositories.Manager
	Controllers *controllers.BaseController
}

func New() (*App, error) {
	log := logger.New("app")

	cfg, err := config.New()
	if err != nil {
		log.Er("failed to initialize config", err)
		return &App{}, err
	}

	db, err := database.New(cfg)
	if err != nil {
		log.Er("failed to create database", err)
		return &App{}, err
	}

	if err := runMigrations(db); err != nil {
		log.Er("failed to run migrations", err)
		return &App{}, err
	}

	repoManager, err := repositories.NewManager(db)
	if err != nil {
		log.Er("failed to initialize repository manager", err)
		return &App{}, err
	}

	service, err := services.New(repoManager)
	if err != nil {
		log.Er("failed to initialize services", err)
		return &App{}, err
	}

	controllers := controllers.New(service)

	app := &App{
		Database:    db,
		Config:      cfg,
		Services:    service,
		RepoManager: repoManager,
		Controllers: controllers,
	}

	log.Info("Application initialized successfully with multi-ORM support")

	return app, nil
}

func runMigrations(db *database.DB) error {
	log := db.Logger.Function("runMigrations")

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

	log.Info("Database migrations completed successfully")
	return nil
}

func (a *App) Close() error {
	if a.Database != nil {
		return a.Database.Close()
	}
	return nil
}