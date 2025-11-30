package repositories

import (
	"bananas/internal/database"
	"bananas/internal/logger"
	"bananas/internal/models"
	"context"
	"time"

	"gorm.io/gorm"
)

type GORMRepository struct {
	DB     *gorm.DB
	Logger logger.Logger
}

func NewGORMRepository(db *database.DB) *GORMRepository {
	return &GORMRepository{
		DB:     db.GORM,
		Logger: logger.New("gorm-repository"),
	}
}

func (r *GORMRepository) CreateTestResult(ctx context.Context, result *models.TestResult) error {
	now := time.Now()
	result.CreatedAt = now
	result.UpdatedAt = now

	err := r.DB.WithContext(ctx).Table("test_results").Create(result).Error
	if err != nil {
		r.Logger.Er("failed to create test result", err)
		return err
	}

	return nil
}

func (r *GORMRepository) GetTestResults(ctx context.Context, limit int) ([]*models.TestResult, error) {
	var results []*models.TestResult

	err := r.DB.WithContext(ctx).
		Table("test_results").
		Order("created_at DESC").
		Limit(limit).
		Find(&results).Error

	if err != nil {
		r.Logger.Er("failed to query test results", err)
		return nil, err
	}

	return results, nil
}

func (r *GORMRepository) CreateFramework(ctx context.Context, framework *models.Framework) error {
	now := time.Now()
	framework.CreatedAt = now
	framework.UpdatedAt = now

	err := r.DB.WithContext(ctx).Table("frameworks").Create(framework).Error
	if err != nil {
		r.Logger.Er("failed to create framework", err)
		return err
	}

	return nil
}

func (r *GORMRepository) GetFrameworks(ctx context.Context, frameworkType string) ([]*models.Framework, error) {
	var frameworks []*models.Framework

	query := r.DB.WithContext(ctx).Table("frameworks")
	if frameworkType != "" {
		query = query.Where("type = ?", frameworkType)
	}

	err := query.Order("name").Find(&frameworks).Error
	if err != nil {
		r.Logger.Er("failed to query frameworks", err)
		return nil, err
	}

	return frameworks, nil
}
