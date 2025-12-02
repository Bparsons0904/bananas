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

func (r *GORMRepository) GetRecentOrders(ctx context.Context, limit int) ([]*models.OrderWithDetails, error) {
	type OrderPreload struct {
		models.SalesOrder
		Customer models.Customer         `gorm:"foreignKey:CustomerID"`
		Items    []models.SalesOrderItem `gorm:"foreignKey:SalesOrderID"`
	}

	var ordersPreload []OrderPreload
	err := r.DB.WithContext(ctx).
		Table("sales_orders").
		Preload("Customer").
		Preload("Items").
		Where("sales_orders.deleted_at IS NULL").
		Order("sales_orders.order_date DESC, sales_orders.created_at DESC").
		Limit(limit).
		Find(&ordersPreload).Error

	if err != nil {
		r.Logger.Er("failed to query orders with preload", err)
		return nil, err
	}

	if len(ordersPreload) == 0 {
		return []*models.OrderWithDetails{}, nil
	}

	productIDs := make([]string, 0)
	for _, order := range ordersPreload {
		for _, item := range order.Items {
			productIDs = append(productIDs, item.ProductID.String())
		}
	}

	var products []models.Product
	productMap := make(map[string]models.Product)
	if len(productIDs) > 0 {
		err = r.DB.WithContext(ctx).
			Table("products").
			Where("id IN ?", productIDs).
			Find(&products).Error

		if err != nil {
			r.Logger.Er("failed to query products", err)
			return nil, err
		}

		for _, product := range products {
			productMap[product.ID.String()] = product
		}
	}

	results := make([]*models.OrderWithDetails, len(ordersPreload))
	for i, order := range ordersPreload {
		items := make([]models.OrderItemWithProduct, len(order.Items))
		for j, item := range order.Items {
			items[j] = models.OrderItemWithProduct{
				Item:    item,
				Product: productMap[item.ProductID.String()],
			}
		}

		results[i] = &models.OrderWithDetails{
			Order:    order.SalesOrder,
			Customer: order.Customer,
			Items:    items,
		}
	}

	return results, nil
}
