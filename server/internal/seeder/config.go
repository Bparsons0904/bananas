package seeder

// Config holds the target counts for seeding
type Config struct {
	// Master data
	Categories int
	Suppliers  int
	Warehouses int

	// Products and relationships
	Products                int
	ProductCategoriesPerProduct int // avg categories per product
	SupplierProductsPerProduct  int // avg suppliers per product
	ProductPricesPerProduct     int // price history records
	ProductCostsPerProduct      int // cost history records

	// Customers
	Customers              int
	CustomerAddressesPerCustomer int // avg addresses per customer

	// Sales
	SalesOrders           int
	SalesOrderItemsPerOrder int // avg items per order
	SalesOrderPaymentsPerOrder int // avg payments per order (some split payments)

	// Purchasing
	PurchaseOrders            int
	PurchaseOrderItemsPerOrder int // avg items per order
	PurchaseOrderReceiptsPerItem int // avg receipts per item (partial deliveries)

	// Inventory
	InventoryRecordsPerProduct   int // avg warehouses per product
	InventoryTransactionsPerProduct int // transaction history

	// Batch sizes for insertion
	BatchSize int
}

// DefaultConfig returns the default seeding configuration
// Targets: ~500K products, ~2M customers, ~20M sales orders
func DefaultConfig() *Config {
	return &Config{
		// Master data
		Categories: 1000,
		Suppliers:  5000,
		Warehouses: 500,

		// Products: 500,000 products
		Products:                    500_000,
		ProductCategoriesPerProduct: 2,      // 1-3 categories
		SupplierProductsPerProduct:  2,      // 1-5 suppliers
		ProductPricesPerProduct:     3,      // 3 price history records
		ProductCostsPerProduct:      3,      // 3 cost history records

		// Customers: 2,000,000 customers
		Customers:                    2_000_000,
		CustomerAddressesPerCustomer: 2, // 1-3 addresses

		// Sales: 20,000,000 orders with ~60M line items
		SalesOrders:                20_000_000,
		SalesOrderItemsPerOrder:    3,  // 1-10 items
		SalesOrderPaymentsPerOrder: 1,  // 1-2 payments

		// Purchasing: 100,000 orders with ~5M line items
		PurchaseOrders:               100_000,
		PurchaseOrderItemsPerOrder:   50, // 10-100 items (bulk orders)
		PurchaseOrderReceiptsPerItem: 1,  // 1-2 receipts (partial deliveries)

		// Inventory: full coverage
		InventoryRecordsPerProduct:      2,  // 1-3 warehouses
		InventoryTransactionsPerProduct: 10, // transaction history

		// Performance tuning
		BatchSize: 5000,
	}
}

// SmallConfig returns a smaller config for testing
func SmallConfig() *Config {
	return &Config{
		Categories: 100,
		Suppliers:  500,
		Warehouses: 50,

		Products:                    5_000,
		ProductCategoriesPerProduct: 2,
		SupplierProductsPerProduct:  2,
		ProductPricesPerProduct:     2,
		ProductCostsPerProduct:      2,

		Customers:                    10_000,
		CustomerAddressesPerCustomer: 2,

		SalesOrders:                100_000,
		SalesOrderItemsPerOrder:    3,
		SalesOrderPaymentsPerOrder: 1,

		PurchaseOrders:               1_000,
		PurchaseOrderItemsPerOrder:   50,
		PurchaseOrderReceiptsPerItem: 1,

		InventoryRecordsPerProduct:      2,
		InventoryTransactionsPerProduct: 5,

		BatchSize: 1000,
	}
}

// TotalRecordsEstimate returns the estimated total number of records
func (c *Config) TotalRecordsEstimate() int {
	total := 0
	total += c.Categories
	total += c.Suppliers
	total += c.Warehouses
	total += c.Products
	total += c.Products * c.ProductCategoriesPerProduct
	total += c.Products * c.SupplierProductsPerProduct
	total += c.Products * c.ProductPricesPerProduct
	total += c.Products * c.ProductCostsPerProduct
	total += c.Customers
	total += c.Customers * c.CustomerAddressesPerCustomer
	total += c.SalesOrders
	total += c.SalesOrders * c.SalesOrderItemsPerOrder
	total += c.SalesOrders * c.SalesOrderPaymentsPerOrder
	total += c.PurchaseOrders
	total += c.PurchaseOrders * c.PurchaseOrderItemsPerOrder
	total += c.PurchaseOrders * c.PurchaseOrderItemsPerOrder * c.PurchaseOrderReceiptsPerItem
	total += c.Products * c.InventoryRecordsPerProduct
	total += c.Products * c.InventoryTransactionsPerProduct
	return total
}
