package services

import (
	"bananas/internal/logger"
	"bananas/internal/models"
	"bananas/internal/repositories"
	"context"
	"time"
)

type Service struct {
	RepoManager *repositories.Manager
	Logger      logger.Logger
}

func New(repoManager *repositories.Manager) (*Service, error) {
	log := logger.New("service")

	return &Service{
		RepoManager: repoManager,
		Logger:      log,
	}, nil
}

func (s *Service) RunPerformanceTest(ctx context.Context, framework, testType, ormType string) (*models.TestResult, error) {
	log := s.Logger.Function("RunPerformanceTest")
	log.Info("Running performance test for framework: %s, test: %s, orm: %s", framework, testType, ormType)

	executionMs := s.simulateWork(framework, testType)

	result := &models.TestResult{
		Framework:   framework,
		TestType:    testType,
		ExecutionMs: executionMs,
		Success:     true,
	}

	repo := s.RepoManager.GetRepository(ormType)
	err := repo.CreateTestResult(ctx, result)
	if err != nil {
		log.Er("failed to save test result", err)
		return nil, err
	}

	log.Info("Performance test completed in %dms using %s ORM", executionMs, ormType)
	return result, nil
}

func (s *Service) simulateWork(framework, testType string) int {
	// Simulate different performance characteristics
	baseTime := 50 // base time in ms
	
	frameworkMultipliers := map[string]int{
		"standard": 100,
		"gin":       80,
		"fiber":     60,
		"echo":      85,
		"chi":       90,
		"gorilla":   95,
	}
	
	testMultipliers := map[string]int{
		"simple_request":   1,
		"database_query":   3,
		"json_response":    2,
		"file_upload":      5,
		"authentication":   4,
	}
	
	frameworkMultiplier := frameworkMultipliers[framework]
	if frameworkMultiplier == 0 {
		frameworkMultiplier = 100
	}
	
	testMultiplier := testMultipliers[testType]
	if testMultiplier == 0 {
		testMultiplier = 1
	}
	
	return baseTime * frameworkMultiplier * testMultiplier / 100
}

func (s *Service) GetTestResults(ctx context.Context, ormType string, limit int) ([]*models.TestResult, error) {
	repo := s.RepoManager.GetRepository(ormType)
	return repo.GetTestResults(ctx, limit)
}

func (s *Service) GetFrameworks(ctx context.Context, ormType string, frameworkType string) ([]*models.Framework, error) {
	repo := s.RepoManager.GetRepository(ormType)
	return repo.GetFrameworks(ctx, frameworkType)
}

func (s *Service) GetRecentOrders(ctx context.Context, ormType string, limit int) ([]*models.OrderWithDetails, int64, error) {
	start := time.Now()
	repo := s.RepoManager.GetRepository(ormType)
	orders, err := repo.GetRecentOrders(ctx, limit)
	dbTime := time.Since(start).Milliseconds()
	return orders, dbTime, err
}