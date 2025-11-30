package repositories

import (
	"bananas/internal/database"
	"bananas/internal/logger"
	"bananas/internal/models"
	"context"
	"time"
)

type RepositoryInterface interface {
	CreateTestResult(ctx context.Context, result *models.TestResult) error
	GetTestResults(ctx context.Context, limit int) ([]*models.TestResult, error)
	CreateFramework(ctx context.Context, framework *models.Framework) error
	GetFrameworks(ctx context.Context, frameworkType string) ([]*models.Framework, error)
}

type SQLRepository struct {
	DB     *database.DB
	Logger logger.Logger
}

func NewSQLRepository(db *database.DB) *SQLRepository {
	return &SQLRepository{
		DB:     db,
		Logger: logger.New("sql-repository"),
	}
}

func (r *SQLRepository) CreateTestResult(ctx context.Context, result *models.TestResult) error {
	query := `
		INSERT INTO test_results (framework, test_type, execution_ms, success, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	
	now := time.Now()
	result.CreatedAt = now
	result.UpdatedAt = now
	
	err := r.DB.SQL.QueryRowContext(ctx, query,
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

func (r *SQLRepository) GetTestResults(ctx context.Context, limit int) ([]*models.TestResult, error) {
	query := `
		SELECT id, framework, test_type, execution_ms, success, created_at, updated_at
		FROM test_results
		ORDER BY created_at DESC
		LIMIT $1
	`
	
	rows, err := r.DB.SQL.QueryContext(ctx, query, limit)
	if err != nil {
		r.Logger.Er("failed to query test results", err)
		return nil, err
	}
	defer rows.Close()
	
	var results []*models.TestResult
	for rows.Next() {
		result := &models.TestResult{}
		err := rows.Scan(
			&result.ID,
			&result.Framework,
			&result.TestType,
			&result.ExecutionMs,
			&result.Success,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			r.Logger.Er("failed to scan test result", err)
			return nil, err
		}
		results = append(results, result)
	}
	
	if err := rows.Err(); err != nil {
		r.Logger.Er("error iterating test results", err)
		return nil, err
	}
	
	return results, nil
}

func (r *SQLRepository) CreateFramework(ctx context.Context, framework *models.Framework) error {
	query := `
		INSERT INTO frameworks (name, type, description, enabled, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	
	now := time.Now()
	framework.CreatedAt = now
	framework.UpdatedAt = now
	
	err := r.DB.SQL.QueryRowContext(ctx, query,
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

func (r *SQLRepository) GetFrameworks(ctx context.Context, frameworkType string) ([]*models.Framework, error) {
	query := `
		SELECT id, name, type, description, enabled, created_at, updated_at
		FROM frameworks
		WHERE $1 = '' OR type = $1
		ORDER BY name
	`
	
	rows, err := r.DB.SQL.QueryContext(ctx, query, frameworkType)
	if err != nil {
		r.Logger.Er("failed to query frameworks", err)
		return nil, err
	}
	defer rows.Close()
	
	var frameworks []*models.Framework
	for rows.Next() {
		framework := &models.Framework{}
		err := rows.Scan(
			&framework.ID,
			&framework.Name,
			&framework.Type,
			&framework.Description,
			&framework.Enabled,
			&framework.CreatedAt,
			&framework.UpdatedAt,
		)
		if err != nil {
			r.Logger.Er("failed to scan framework", err)
			return nil, err
		}
		frameworks = append(frameworks, framework)
	}
	
	if err := rows.Err(); err != nil {
		r.Logger.Er("error iterating frameworks", err)
		return nil, err
	}
	
	return frameworks, nil
}