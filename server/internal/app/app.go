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

func (a *App) Close() error {
	if a.Database != nil {
		return a.Database.Close()
	}
	return nil
}