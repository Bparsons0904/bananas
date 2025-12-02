package repositories

import (
	"bananas/internal/database"
	"bananas/internal/logger"
	"bananas/internal/models"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PGXRepository struct {
	Pool   *pgxpool.Pool
	Logger logger.Logger
}

func NewPGXRepository(db *database.DB) *PGXRepository {
	return &PGXRepository{
		Pool:   db.PGX,
		Logger: logger.New("pgx-repository"),
	}
}

func (r *PGXRepository) CreateTestResult(ctx context.Context, result *models.TestResult) error {
	query := `
		INSERT INTO test_results (framework, test_type, execution_ms, success, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	now := time.Now()
	result.CreatedAt = now
	result.UpdatedAt = now

	err := r.Pool.QueryRow(ctx, query,
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

func (r *PGXRepository) GetTestResults(ctx context.Context, limit int) ([]*models.TestResult, error) {
	query := `
		SELECT id, framework, test_type, execution_ms, success, created_at, updated_at
		FROM test_results
		ORDER BY created_at DESC
		LIMIT $1
	`

	rows, err := r.Pool.Query(ctx, query, limit)
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

func (r *PGXRepository) CreateFramework(ctx context.Context, framework *models.Framework) error {
	query := `
		INSERT INTO frameworks (name, type, description, enabled, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	now := time.Now()
	framework.CreatedAt = now
	framework.UpdatedAt = now

	err := r.Pool.QueryRow(ctx, query,
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

func (r *PGXRepository) GetFrameworks(ctx context.Context, frameworkType string) ([]*models.Framework, error) {
	query := `
		SELECT id, name, type, description, enabled, created_at, updated_at
		FROM frameworks
		WHERE $1 = '' OR type = $1
		ORDER BY name
	`

	rows, err := r.Pool.Query(ctx, query, frameworkType)
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

func (r *PGXRepository) GetRecentOrders(ctx context.Context, limit int) ([]*models.OrderWithDetails, error) {
	orderQuery := `
		SELECT
			so.id, so.order_number, so.customer_id, so.order_date, so.status,
			so.subtotal, so.tax, so.shipping, so.total, so.notes,
			so.created_at, so.updated_at, so.deleted_at,
			c.id, c.first_name, c.last_name, c.email, c.phone,
			c.created_at, c.updated_at, c.deleted_at
		FROM sales_orders so
		JOIN customers c ON so.customer_id = c.id
		WHERE so.deleted_at IS NULL
		ORDER BY so.order_date DESC, so.created_at DESC
		LIMIT $1
	`

	rows, err := r.Pool.Query(ctx, orderQuery, limit)
	if err != nil {
		r.Logger.Er("failed to query orders", err)
		return nil, err
	}
	defer rows.Close()

	var results []*models.OrderWithDetails
	orderMap := make(map[string]*models.OrderWithDetails)
	orderIDs := make([]string, 0)

	for rows.Next() {
		order := &models.SalesOrder{}
		customer := &models.Customer{}

		err := rows.Scan(
			&order.ID, &order.OrderNumber, &order.CustomerID, &order.OrderDate, &order.Status,
			&order.Subtotal, &order.Tax, &order.Shipping, &order.Total, &order.Notes,
			&order.CreatedAt, &order.UpdatedAt, &order.DeletedAt,
			&customer.ID, &customer.FirstName, &customer.LastName, &customer.Email, &customer.Phone,
			&customer.CreatedAt, &customer.UpdatedAt, &customer.DeletedAt,
		)
		if err != nil {
			r.Logger.Er("failed to scan order", err)
			return nil, err
		}

		orderWithDetails := &models.OrderWithDetails{
			Order:    *order,
			Customer: *customer,
			Items:    []models.OrderItemWithProduct{},
		}
		orderMap[order.ID.String()] = orderWithDetails
		orderIDs = append(orderIDs, order.ID.String())
		results = append(results, orderWithDetails)
	}

	if err := rows.Err(); err != nil {
		r.Logger.Er("error iterating orders", err)
		return nil, err
	}

	if len(orderIDs) == 0 {
		return results, nil
	}

	itemQuery := `
		SELECT
			soi.id, soi.sales_order_id, soi.product_id, soi.quantity,
			soi.unit_price, soi.discount, soi.tax, soi.total,
			soi.created_at, soi.updated_at,
			p.id, p.sku, p.name, p.description, p.weight, p.dimensions,
			p.is_active, p.created_at, p.updated_at, p.deleted_at
		FROM sales_order_items soi
		JOIN products p ON soi.product_id = p.id
		WHERE soi.sales_order_id = ANY($1)
		ORDER BY soi.created_at
	`

	itemRows, err := r.Pool.Query(ctx, itemQuery, orderIDs)
	if err != nil {
		r.Logger.Er("failed to query order items", err)
		return nil, err
	}
	defer itemRows.Close()

	for itemRows.Next() {
		item := &models.SalesOrderItem{}
		product := &models.Product{}

		err := itemRows.Scan(
			&item.ID, &item.SalesOrderID, &item.ProductID, &item.Quantity,
			&item.UnitPrice, &item.Discount, &item.Tax, &item.Total,
			&item.CreatedAt, &item.UpdatedAt,
			&product.ID, &product.SKU, &product.Name, &product.Description, &product.Weight,
			&product.Dimensions, &product.IsActive, &product.CreatedAt, &product.UpdatedAt, &product.DeletedAt,
		)
		if err != nil {
			r.Logger.Er("failed to scan order item", err)
			return nil, err
		}

		if orderDetails, ok := orderMap[item.SalesOrderID.String()]; ok {
			orderDetails.Items = append(orderDetails.Items, models.OrderItemWithProduct{
				Item:    *item,
				Product: *product,
			})
		}
	}

	if err := itemRows.Err(); err != nil {
		r.Logger.Er("error iterating order items", err)
		return nil, err
	}

	return results, nil
}
