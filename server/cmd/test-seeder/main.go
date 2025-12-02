package main

import (
	"bananas/internal/config"
	"bananas/internal/database"
	"bananas/internal/logger"
	"bananas/internal/seeder"
	"context"
	"fmt"
	"os"
)

func main() {
	log := logger.New("test-seeder")

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

	// Use small config for testing
	seedCfg := seeder.SmallConfig()
	log.Info(fmt.Sprintf("Starting seeding with config: %d products, %d customers", seedCfg.Products, seedCfg.Customers))

	s := seeder.New(db, seedCfg)
	ctx := context.Background()

	if err := s.SeedAll(ctx); err != nil {
		log.Er("seeding failed", err)
		os.Exit(1)
	}

	log.Info("Seeding completed successfully!")
}
