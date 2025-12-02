package seeder

import (
	"bananas/internal/database"
	"bananas/internal/logger"
	"bananas/internal/seeder/generators"
	"bananas/internal/seeder/inserters"
	"context"
	"fmt"
)

// Seeder orchestrates the seeding process
type Seeder struct {
	db       *database.DB
	config   *Config
	inserter *inserters.PGXCopyInserter
	logger   logger.Logger
	progress *ProgressTracker
}

// New creates a new seeder
func New(db *database.DB, config *Config) *Seeder {
	return &Seeder{
		db:       db,
		config:   config,
		inserter: inserters.NewPGXCopyInserter(db.PGX),
		logger:   logger.New("seeder"),
		progress: NewProgressTracker(),
	}
}

// SeedAll seeds all tables with data
func (s *Seeder) SeedAll(ctx context.Context) error {
	s.logger.Info("Starting full database seeding")
	s.logger.Info(fmt.Sprintf("Estimated total records: %d", s.config.TotalRecordsEstimate()))

	// Phase 1: Master Data
	s.logger.Info("=== Phase 1: Master Data ===")
	idMap, err := s.seedMasterData(ctx)
	if err != nil {
		return fmt.Errorf("failed to seed master data: %w", err)
	}

	// Phase 2: Products and Relationships
	s.logger.Info("=== Phase 2: Products and Relationships ===")
	if err := s.seedProducts(ctx, idMap); err != nil {
		return fmt.Errorf("failed to seed products: %w", err)
	}

	// Phase 3: Customers
	s.logger.Info("=== Phase 3: Customers ===")
	if err := s.seedCustomers(ctx, idMap); err != nil {
		return fmt.Errorf("failed to seed customers: %w", err)
	}

	// Phase 4: Sales Orders
	s.logger.Info("=== Phase 4: Sales Orders ===")
	if err := s.seedSalesOrders(ctx, idMap); err != nil {
		return fmt.Errorf("failed to seed sales orders: %w", err)
	}

	// Phase 5: Purchase Orders
	s.logger.Info("=== Phase 5: Purchase Orders ===")
	if err := s.seedPurchaseOrders(ctx, idMap); err != nil {
		return fmt.Errorf("failed to seed purchase orders: %w", err)
	}

	// Phase 6: Inventory
	s.logger.Info("=== Phase 6: Inventory ===")
	if err := s.seedInventory(ctx, idMap); err != nil {
		return fmt.Errorf("failed to seed inventory: %w", err)
	}

	elapsed := s.progress.Elapsed()
	s.logger.Info(fmt.Sprintf("Seeding completed in %s", FormatDuration(elapsed)))

	return nil
}

func (s *Seeder) seedMasterData(ctx context.Context) (*generators.IDMap, error) {
	// Categories
	s.logger.Info(fmt.Sprintf("Generating %d categories...", s.config.Categories))
	categoryRows, idMap := generators.GenerateCategories(s.config.Categories)
	s.logger.Info(fmt.Sprintf("Generated %d categories, inserting...", len(categoryRows)))
	s.progress.StartTable("Categories", s.config.Categories)
	if err := s.inserter.BulkInsertBatched(ctx, "categories", generators.CategoryColumns(), categoryRows, s.config.BatchSize); err != nil {
		return nil, err
	}
	s.progress.Finish()
	s.logger.Info(fmt.Sprintf("Seeded %d categories", len(categoryRows)))

	// Suppliers
	s.logger.Info(fmt.Sprintf("Generating %d suppliers...", s.config.Suppliers))
	supplierRows := generators.GenerateSuppliers(s.config.Suppliers, idMap)
	s.logger.Info(fmt.Sprintf("Generated %d suppliers, inserting...", len(supplierRows)))
	s.progress.StartTable("Suppliers", s.config.Suppliers)
	if err := s.inserter.BulkInsertBatched(ctx, "suppliers", generators.SupplierColumns(), supplierRows, s.config.BatchSize); err != nil {
		return nil, err
	}
	s.progress.Finish()
	s.logger.Info(fmt.Sprintf("Seeded %d suppliers", len(supplierRows)))

	// Warehouses
	s.logger.Info(fmt.Sprintf("Generating %d warehouses...", s.config.Warehouses))
	warehouseRows := generators.GenerateWarehouses(s.config.Warehouses, idMap)
	s.logger.Info(fmt.Sprintf("Generated %d warehouses, inserting...", len(warehouseRows)))
	s.progress.StartTable("Warehouses", s.config.Warehouses)
	if err := s.inserter.BulkInsertBatched(ctx, "warehouses", generators.WarehouseColumns(), warehouseRows, s.config.BatchSize); err != nil {
		return nil, err
	}
	s.progress.Finish()
	s.logger.Info(fmt.Sprintf("Seeded %d warehouses", len(warehouseRows)))

	return idMap, nil
}

func (s *Seeder) seedProducts(ctx context.Context, idMap *generators.IDMap) error {
	// Products
	s.logger.Info(fmt.Sprintf("Generating %d products...", s.config.Products))
	productRows := generators.GenerateProducts(s.config.Products, idMap)
	s.logger.Info(fmt.Sprintf("Generated %d products, inserting...", len(productRows)))
	s.progress.StartTable("Products", len(productRows))
	if err := s.inserter.BulkInsertBatched(ctx, "products", generators.ProductColumns(), productRows, s.config.BatchSize); err != nil {
		return err
	}
	s.progress.Finish()
	s.logger.Info(fmt.Sprintf("Seeded %d products", len(productRows)))

	// Product Categories
	s.logger.Info("Generating product-category relationships...")
	pcRows := generators.GenerateProductCategories(idMap, s.config.ProductCategoriesPerProduct)
	s.logger.Info(fmt.Sprintf("Generated %d product-category relationships, inserting...", len(pcRows)))
	s.progress.StartTable("Product Categories", len(pcRows))
	if err := s.inserter.BulkInsertBatched(ctx, "product_categories", generators.ProductCategoryColumns(), pcRows, s.config.BatchSize); err != nil {
		return err
	}
	s.progress.Finish()
	s.logger.Info(fmt.Sprintf("Seeded %d product-category relationships", len(pcRows)))

	// Product Prices
	s.logger.Info("Generating product prices...")
	priceRows := generators.GenerateProductPrices(idMap, s.config.ProductPricesPerProduct)
	s.logger.Info(fmt.Sprintf("Generated %d product prices, inserting...", len(priceRows)))
	s.progress.StartTable("Product Prices", len(priceRows))
	if err := s.inserter.BulkInsertBatched(ctx, "product_prices", generators.ProductPriceColumns(), priceRows, s.config.BatchSize); err != nil {
		return err
	}
	s.progress.Finish()
	s.logger.Info(fmt.Sprintf("Seeded %d product prices", len(priceRows)))

	// Product Costs
	s.logger.Info("Generating product costs...")
	costRows := generators.GenerateProductCosts(idMap, s.config.ProductCostsPerProduct)
	s.logger.Info(fmt.Sprintf("Generated %d product costs, inserting...", len(costRows)))
	s.progress.StartTable("Product Costs", len(costRows))
	if err := s.inserter.BulkInsertBatched(ctx, "product_costs", generators.ProductCostColumns(), costRows, s.config.BatchSize); err != nil {
		return err
	}
	s.progress.Finish()
	s.logger.Info(fmt.Sprintf("Seeded %d product costs", len(costRows)))

	// Supplier Products
	s.logger.Info("Generating supplier-product relationships...")
	spRows := generators.GenerateSupplierProducts(idMap, s.config.SupplierProductsPerProduct)
	s.logger.Info(fmt.Sprintf("Generated %d supplier-product relationships, inserting...", len(spRows)))
	s.progress.StartTable("Supplier Products", len(spRows))
	if err := s.inserter.BulkInsertBatched(ctx, "supplier_products", generators.SupplierProductColumns(), spRows, s.config.BatchSize); err != nil {
		return err
	}
	s.progress.Finish()
	s.logger.Info(fmt.Sprintf("Seeded %d supplier-product relationships", len(spRows)))

	return nil
}

func (s *Seeder) seedCustomers(ctx context.Context, idMap *generators.IDMap) error {
	// Customers
	s.logger.Info(fmt.Sprintf("Generating %d customers...", s.config.Customers))
	customerRows := generators.GenerateCustomers(s.config.Customers, idMap)
	s.logger.Info(fmt.Sprintf("Generated %d customers, inserting...", len(customerRows)))
	s.progress.StartTable("Customers", len(customerRows))
	if err := s.inserter.BulkInsertBatched(ctx, "customers", generators.CustomerColumns(), customerRows, s.config.BatchSize); err != nil {
		return err
	}
	s.progress.Finish()
	s.logger.Info(fmt.Sprintf("Seeded %d customers", len(customerRows)))

	// Customer Addresses
	s.logger.Info("Generating customer addresses...")
	addressRows := generators.GenerateCustomerAddresses(idMap, s.config.CustomerAddressesPerCustomer)
	s.logger.Info(fmt.Sprintf("Generated %d customer addresses, inserting...", len(addressRows)))
	s.progress.StartTable("Customer Addresses", len(addressRows))
	if err := s.inserter.BulkInsertBatched(ctx, "customer_addresses", generators.CustomerAddressColumns(), addressRows, s.config.BatchSize); err != nil {
		return err
	}
	s.progress.Finish()
	s.logger.Info(fmt.Sprintf("Seeded %d customer addresses", len(addressRows)))

	return nil
}

func (s *Seeder) seedSalesOrders(ctx context.Context, idMap *generators.IDMap) error {
	// Sales Orders
	s.logger.Info(fmt.Sprintf("Generating %d sales orders...", s.config.SalesOrders))
	orderRows, salesData := generators.GenerateSalesOrders(s.config.SalesOrders, idMap)
	s.logger.Info(fmt.Sprintf("Generated %d sales orders, inserting...", len(orderRows)))
	s.progress.StartTable("Sales Orders", len(orderRows))
	if err := s.inserter.BulkInsertBatched(ctx, "sales_orders", generators.SalesOrderColumns(), orderRows, s.config.BatchSize); err != nil {
		return err
	}
	s.progress.Finish()
	s.logger.Info(fmt.Sprintf("Seeded %d sales orders", len(orderRows)))

	// Sales Order Items
	s.logger.Info("Generating sales order items...")
	itemRows := generators.GenerateSalesOrderItems(salesData, idMap, s.config.SalesOrderItemsPerOrder)
	s.logger.Info(fmt.Sprintf("Generated %d sales order items, inserting...", len(itemRows)))
	s.progress.StartTable("Sales Order Items", len(itemRows))
	if err := s.inserter.BulkInsertBatched(ctx, "sales_order_items", generators.SalesOrderItemColumns(), itemRows, s.config.BatchSize); err != nil {
		return err
	}
	s.progress.Finish()
	s.logger.Info(fmt.Sprintf("Seeded %d sales order items", len(itemRows)))

	// Sales Order Payments
	s.logger.Info("Generating sales order payments...")
	paymentRows := generators.GenerateSalesOrderPayments(salesData, s.config.SalesOrderPaymentsPerOrder)
	s.logger.Info(fmt.Sprintf("Generated %d sales order payments, inserting...", len(paymentRows)))
	s.progress.StartTable("Sales Order Payments", len(paymentRows))
	if err := s.inserter.BulkInsertBatched(ctx, "sales_order_payments", generators.SalesOrderPaymentColumns(), paymentRows, s.config.BatchSize); err != nil {
		return err
	}
	s.progress.Finish()
	s.logger.Info(fmt.Sprintf("Seeded %d sales order payments", len(paymentRows)))

	return nil
}

func (s *Seeder) seedPurchaseOrders(ctx context.Context, idMap *generators.IDMap) error {
	// Purchase Orders
	s.logger.Info(fmt.Sprintf("Generating %d purchase orders...", s.config.PurchaseOrders))
	poRows, purchaseData := generators.GeneratePurchaseOrders(s.config.PurchaseOrders, idMap)
	s.logger.Info(fmt.Sprintf("Generated %d purchase orders, inserting...", len(poRows)))
	s.progress.StartTable("Purchase Orders", len(poRows))
	if err := s.inserter.BulkInsertBatched(ctx, "purchase_orders", generators.PurchaseOrderColumns(), poRows, s.config.BatchSize); err != nil {
		return err
	}
	s.progress.Finish()
	s.logger.Info(fmt.Sprintf("Seeded %d purchase orders", len(poRows)))

	// Purchase Order Items
	s.logger.Info("Generating purchase order items...")
	poItemRows := generators.GeneratePurchaseOrderItems(purchaseData, idMap, s.config.PurchaseOrderItemsPerOrder)
	s.logger.Info(fmt.Sprintf("Generated %d purchase order items, inserting...", len(poItemRows)))
	s.progress.StartTable("Purchase Order Items", len(poItemRows))
	if err := s.inserter.BulkInsertBatched(ctx, "purchase_order_items", generators.PurchaseOrderItemColumns(), poItemRows, s.config.BatchSize); err != nil {
		return err
	}
	s.progress.Finish()
	s.logger.Info(fmt.Sprintf("Seeded %d purchase order items", len(poItemRows)))

	// Purchase Order Receipts
	s.logger.Info("Generating purchase order receipts...")
	receiptRows := generators.GeneratePurchaseOrderReceipts(purchaseData, s.config.PurchaseOrderReceiptsPerItem)
	s.logger.Info(fmt.Sprintf("Generated %d purchase order receipts, inserting...", len(receiptRows)))
	s.progress.StartTable("Purchase Order Receipts", len(receiptRows))
	if err := s.inserter.BulkInsertBatched(ctx, "purchase_order_receipts", generators.PurchaseOrderReceiptColumns(), receiptRows, s.config.BatchSize); err != nil {
		return err
	}
	s.progress.Finish()
	s.logger.Info(fmt.Sprintf("Seeded %d purchase order receipts", len(receiptRows)))

	return nil
}

func (s *Seeder) seedInventory(ctx context.Context, idMap *generators.IDMap) error {
	// Inventory
	s.logger.Info("Generating inventory records...")
	invRows := generators.GenerateInventory(idMap, s.config.InventoryRecordsPerProduct)
	s.logger.Info(fmt.Sprintf("Generated %d inventory records, inserting...", len(invRows)))
	s.progress.StartTable("Inventory", len(invRows))
	if err := s.inserter.BulkInsertBatched(ctx, "inventory", generators.InventoryColumns(), invRows, s.config.BatchSize); err != nil {
		return err
	}
	s.progress.Finish()
	s.logger.Info(fmt.Sprintf("Seeded %d inventory records", len(invRows)))

	// Inventory Transactions
	s.logger.Info("Generating inventory transactions...")
	txnRows := generators.GenerateInventoryTransactions(idMap, s.config.InventoryTransactionsPerProduct)
	s.logger.Info(fmt.Sprintf("Generated %d inventory transactions, inserting...", len(txnRows)))
	s.progress.StartTable("Inventory Transactions", len(txnRows))
	if err := s.inserter.BulkInsertBatched(ctx, "inventory_transactions", generators.InventoryTransactionColumns(), txnRows, s.config.BatchSize); err != nil {
		return err
	}
	s.progress.Finish()
	s.logger.Info(fmt.Sprintf("Seeded %d inventory transactions", len(txnRows)))

	return nil
}

// CheckNeedsSeeding checks if seeding is needed based on customer count
func CheckNeedsSeeding(db *database.DB, targetCount int) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM customers"
	err := db.SQL.QueryRow(query).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check customer count: %w", err)
	}

	return count < targetCount, nil
}
