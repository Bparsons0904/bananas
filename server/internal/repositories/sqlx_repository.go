package repositories

import (
	"bananas/internal/database"
	"bananas/internal/logger"
	"bananas/internal/models"
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type SQLxRepository struct {
	DB     *sqlx.DB
	Logger logger.Logger
}

func NewSQLxRepository(db *database.DB) *SQLxRepository {
	return &SQLxRepository{
		DB:     db.SQLx,
		Logger: logger.New("sqlx-repository"),
	}
}

func (r *SQLxRepository) CreateTestResult(ctx context.Context, result *models.TestResult) error {
	query := `
		INSERT INTO test_results (framework, test_type, execution_ms, success, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	now := time.Now()
	result.CreatedAt = now
	result.UpdatedAt = now

	err := r.DB.QueryRowContext(ctx, query,
		result.Framework,
		result.TestType,
		result.ExecutionMs,
		result.Success,
		now,
		now,
	).Scan(&result.ID)

	if err != nil {
		r.Logger.Er("failed to create test result", err)
		return err
	}

	return nil
}

func (r *SQLxRepository) GetTestResults(ctx context.Context, limit int) ([]*models.TestResult, error) {
	query := `
		SELECT id, framework, test_type, execution_ms, success, created_at, updated_at
		FROM test_results
		ORDER BY created_at DESC
		LIMIT $1
	`

	var results []*models.TestResult
	err := r.DB.SelectContext(ctx, &results, query, limit)
	if err != nil {
		r.Logger.Er("failed to query test results", err)
		return nil, err
	}

	return results, nil
}

func (r *SQLxRepository) CreateFramework(ctx context.Context, framework *models.Framework) error {
	query := `
		INSERT INTO frameworks (name, type, description, enabled, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	now := time.Now()
	framework.CreatedAt = now
	framework.UpdatedAt = now

	err := r.DB.QueryRowContext(ctx, query,
		framework.Name,
		framework.Type,
		framework.Description,
		framework.Enabled,
		now,
		now,
	).Scan(&framework.ID)

	if err != nil {
		r.Logger.Er("failed to create framework", err)
		return err
	}

	return nil
}

func (r *SQLxRepository) GetFrameworks(ctx context.Context, frameworkType string) ([]*models.Framework, error) {
	query := `
		SELECT id, name, type, description, enabled, created_at, updated_at
		FROM frameworks
		WHERE $1 = '' OR type = $1
		ORDER BY name
	`

	var frameworks []*models.Framework
	err := r.DB.SelectContext(ctx, &frameworks, query, frameworkType)
	if err != nil {
		r.Logger.Er("failed to query frameworks", err)
		return nil, err
	}

	return frameworks, nil
}
