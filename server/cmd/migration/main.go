package main

import (
	"bananas/internal/config"
	"bananas/internal/database"
	"bananas/internal/logger"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/migration/main.go [create-db|up|down|seed]")
		os.Exit(1)
	}

	command := os.Args[1]
	log := logger.New("migration")

	cfg, err := config.New()
	if err != nil {
		log.Er("failed to initialize config", err)
		os.Exit(1)
	}

	switch command {
	case "create-db":
		err = createDatabase(cfg, log)
	case "up", "down", "seed":
		db, err := database.New(cfg)
		if err != nil {
			log.Er("failed to connect to database", err)
			os.Exit(1)
		}
		defer db.Close()

		switch command {
		case "up":
			err = migrateUp(db)
		case "down":
			err = migrateDown(db)
		case "seed":
			err = seed(db)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Available commands: create-db, up, down, seed")
		os.Exit(1)
	}

	if err != nil {
		log.Er("command failed", err)
		os.Exit(1)
	}

	log.Info("Command completed successfully")
}

func createDatabase(cfg config.Config, log logger.Logger) error {
	log.Info("Creating database if it doesn't exist")

	adminDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		cfg.DatabaseConfig.Host,
		cfg.DatabaseConfig.Port,
		cfg.DatabaseConfig.AdminUser,
		cfg.DatabaseConfig.AdminPassword,
		cfg.DatabaseConfig.SSLMode,
	)

	adminDB, err := sql.Open("postgres", adminDSN)
	if err != nil {
		log.Er("failed to connect to postgres database", err)
		return err
	}
	defer adminDB.Close()

	if err := adminDB.Ping(); err != nil {
		log.Er("failed to ping postgres database", err)
		return err
	}

	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", cfg.DatabaseConfig.DBName)
	err = adminDB.QueryRow(query).Scan(&exists)
	if err != nil {
		log.Er("failed to check if database exists", err)
		return err
	}

	if exists {
		log.Info("Database already exists", "database", cfg.DatabaseConfig.DBName)
		return nil
	}

	createDBQuery := fmt.Sprintf("CREATE DATABASE %s", cfg.DatabaseConfig.DBName)
	_, err = adminDB.Exec(createDBQuery)
	if err != nil {
		log.Er("failed to create database", err)
		return err
	}

	log.Info("Database created successfully", "database", cfg.DatabaseConfig.DBName)

	if cfg.DatabaseConfig.User != cfg.DatabaseConfig.AdminUser {
		var userExists bool
		userQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_roles WHERE rolname = '%s')", cfg.DatabaseConfig.User)
		err = adminDB.QueryRow(userQuery).Scan(&userExists)
		if err != nil {
			log.Er("failed to check if user exists", err)
			return err
		}

		if !userExists {
			createUserQuery := fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", cfg.DatabaseConfig.User, cfg.DatabaseConfig.Password)
			_, err = adminDB.Exec(createUserQuery)
			if err != nil {
				log.Er("failed to create user", err)
				return err
			}
			log.Info("Database user created", "user", cfg.DatabaseConfig.User)
		}

		grantQuery := fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s", cfg.DatabaseConfig.DBName, cfg.DatabaseConfig.User)
		_, err = adminDB.Exec(grantQuery)
		if err != nil {
			log.Er("failed to grant privileges", err)
			return err
		}
		log.Info("Database privileges granted", "user", cfg.DatabaseConfig.User)
	}

	return nil
}

func migrateUp(db *database.DB) error {
	log := db.Logger.Function("migrateUp")

	queries := []string{
		// Master Data Tables
		`CREATE TABLE IF NOT EXISTS categories (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			parent_id UUID REFERENCES categories(id) ON DELETE SET NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_categories_parent_id ON categories(parent_id)`,
		`CREATE INDEX IF NOT EXISTS idx_categories_deleted_at ON categories(deleted_at)`,

		`CREATE TABLE IF NOT EXISTS suppliers (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			contact_name VARCHAR(255),
			email VARCHAR(255),
			phone VARCHAR(50),
			address TEXT,
			city VARCHAR(100),
			state VARCHAR(100),
			postal_code VARCHAR(20),
			country VARCHAR(100),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_suppliers_deleted_at ON suppliers(deleted_at)`,

		`CREATE TABLE IF NOT EXISTS customers (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			first_name VARCHAR(255) NOT NULL,
			last_name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			phone VARCHAR(50),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_customers_email ON customers(email)`,
		`CREATE INDEX IF NOT EXISTS idx_customers_deleted_at ON customers(deleted_at)`,

		`CREATE TABLE IF NOT EXISTS warehouses (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			code VARCHAR(50) UNIQUE NOT NULL,
			address TEXT,
			city VARCHAR(100),
			state VARCHAR(100),
			postal_code VARCHAR(20),
			country VARCHAR(100),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_warehouses_code ON warehouses(code)`,
		`CREATE INDEX IF NOT EXISTS idx_warehouses_deleted_at ON warehouses(deleted_at)`,

		`CREATE TABLE IF NOT EXISTS products (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			sku VARCHAR(100) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			weight DECIMAL(10, 2),
			dimensions VARCHAR(100),
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku)`,
		`CREATE INDEX IF NOT EXISTS idx_products_is_active ON products(is_active)`,
		`CREATE INDEX IF NOT EXISTS idx_products_deleted_at ON products(deleted_at)`,

		// Inventory Tables
		`CREATE TABLE IF NOT EXISTS inventory (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
			warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
			quantity INTEGER NOT NULL DEFAULT 0,
			reserved_quantity INTEGER NOT NULL DEFAULT 0,
			reorder_point INTEGER DEFAULT 0,
			reorder_quantity INTEGER DEFAULT 0,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			UNIQUE(product_id, warehouse_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_inventory_product_id ON inventory(product_id)`,
		`CREATE INDEX IF NOT EXISTS idx_inventory_warehouse_id ON inventory(warehouse_id)`,

		`CREATE TABLE IF NOT EXISTS inventory_transactions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
			warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
			transaction_type VARCHAR(50) NOT NULL,
			quantity INTEGER NOT NULL,
			reference_id UUID,
			reference_type VARCHAR(50),
			notes TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_inventory_transactions_product_id ON inventory_transactions(product_id)`,
		`CREATE INDEX IF NOT EXISTS idx_inventory_transactions_warehouse_id ON inventory_transactions(warehouse_id)`,
		`CREATE INDEX IF NOT EXISTS idx_inventory_transactions_reference ON inventory_transactions(reference_type, reference_id)`,
		`CREATE INDEX IF NOT EXISTS idx_inventory_transactions_created_at ON inventory_transactions(created_at)`,

		// Pricing Tables
		`CREATE TABLE IF NOT EXISTS product_prices (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
			price DECIMAL(10, 2) NOT NULL,
			currency VARCHAR(3) DEFAULT 'USD',
			effective_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			end_date TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_product_prices_product_id ON product_prices(product_id)`,
		`CREATE INDEX IF NOT EXISTS idx_product_prices_effective_date ON product_prices(effective_date)`,

		`CREATE TABLE IF NOT EXISTS product_costs (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
			cost DECIMAL(10, 2) NOT NULL,
			currency VARCHAR(3) DEFAULT 'USD',
			effective_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			end_date TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_product_costs_product_id ON product_costs(product_id)`,
		`CREATE INDEX IF NOT EXISTS idx_product_costs_effective_date ON product_costs(effective_date)`,

		`CREATE TABLE IF NOT EXISTS supplier_products (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			supplier_id UUID NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
			product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
			supplier_sku VARCHAR(100),
			cost DECIMAL(10, 2),
			currency VARCHAR(3) DEFAULT 'USD',
			lead_time_days INTEGER,
			minimum_order_quantity INTEGER,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			UNIQUE(supplier_id, product_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_supplier_products_supplier_id ON supplier_products(supplier_id)`,
		`CREATE INDEX IF NOT EXISTS idx_supplier_products_product_id ON supplier_products(product_id)`,

		// Sales Order Tables
		`CREATE TABLE IF NOT EXISTS sales_orders (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			order_number VARCHAR(100) UNIQUE NOT NULL,
			customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
			order_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			status VARCHAR(50) DEFAULT 'pending',
			subtotal DECIMAL(10, 2) NOT NULL DEFAULT 0,
			tax DECIMAL(10, 2) NOT NULL DEFAULT 0,
			shipping DECIMAL(10, 2) NOT NULL DEFAULT 0,
			total DECIMAL(10, 2) NOT NULL DEFAULT 0,
			notes TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_sales_orders_order_number ON sales_orders(order_number)`,
		`CREATE INDEX IF NOT EXISTS idx_sales_orders_customer_id ON sales_orders(customer_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sales_orders_order_date ON sales_orders(order_date)`,
		`CREATE INDEX IF NOT EXISTS idx_sales_orders_status ON sales_orders(status)`,
		`CREATE INDEX IF NOT EXISTS idx_sales_orders_deleted_at ON sales_orders(deleted_at)`,

		`CREATE TABLE IF NOT EXISTS sales_order_items (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			sales_order_id UUID NOT NULL REFERENCES sales_orders(id) ON DELETE CASCADE,
			product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
			quantity INTEGER NOT NULL,
			unit_price DECIMAL(10, 2) NOT NULL,
			discount DECIMAL(10, 2) DEFAULT 0,
			tax DECIMAL(10, 2) DEFAULT 0,
			total DECIMAL(10, 2) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_sales_order_items_sales_order_id ON sales_order_items(sales_order_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sales_order_items_product_id ON sales_order_items(product_id)`,

		`CREATE TABLE IF NOT EXISTS sales_order_payments (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			sales_order_id UUID NOT NULL REFERENCES sales_orders(id) ON DELETE CASCADE,
			payment_method VARCHAR(50) NOT NULL,
			amount DECIMAL(10, 2) NOT NULL,
			transaction_id VARCHAR(255),
			status VARCHAR(50) DEFAULT 'pending',
			payment_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_sales_order_payments_sales_order_id ON sales_order_payments(sales_order_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sales_order_payments_status ON sales_order_payments(status)`,
		`CREATE INDEX IF NOT EXISTS idx_sales_order_payments_payment_date ON sales_order_payments(payment_date)`,

		// Purchase Order Tables
		`CREATE TABLE IF NOT EXISTS purchase_orders (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			po_number VARCHAR(100) UNIQUE NOT NULL,
			supplier_id UUID NOT NULL REFERENCES suppliers(id) ON DELETE RESTRICT,
			warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE RESTRICT,
			order_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			expected_date TIMESTAMP WITH TIME ZONE,
			status VARCHAR(50) DEFAULT 'pending',
			subtotal DECIMAL(10, 2) NOT NULL DEFAULT 0,
			tax DECIMAL(10, 2) NOT NULL DEFAULT 0,
			shipping DECIMAL(10, 2) NOT NULL DEFAULT 0,
			total DECIMAL(10, 2) NOT NULL DEFAULT 0,
			notes TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_purchase_orders_po_number ON purchase_orders(po_number)`,
		`CREATE INDEX IF NOT EXISTS idx_purchase_orders_supplier_id ON purchase_orders(supplier_id)`,
		`CREATE INDEX IF NOT EXISTS idx_purchase_orders_warehouse_id ON purchase_orders(warehouse_id)`,
		`CREATE INDEX IF NOT EXISTS idx_purchase_orders_order_date ON purchase_orders(order_date)`,
		`CREATE INDEX IF NOT EXISTS idx_purchase_orders_status ON purchase_orders(status)`,
		`CREATE INDEX IF NOT EXISTS idx_purchase_orders_deleted_at ON purchase_orders(deleted_at)`,

		`CREATE TABLE IF NOT EXISTS purchase_order_items (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			purchase_order_id UUID NOT NULL REFERENCES purchase_orders(id) ON DELETE CASCADE,
			product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
			quantity INTEGER NOT NULL,
			unit_cost DECIMAL(10, 2) NOT NULL,
			tax DECIMAL(10, 2) DEFAULT 0,
			total DECIMAL(10, 2) NOT NULL,
			received_quantity INTEGER DEFAULT 0,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_purchase_order_items_purchase_order_id ON purchase_order_items(purchase_order_id)`,
		`CREATE INDEX IF NOT EXISTS idx_purchase_order_items_product_id ON purchase_order_items(product_id)`,

		`CREATE TABLE IF NOT EXISTS purchase_order_receipts (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			purchase_order_id UUID NOT NULL REFERENCES purchase_orders(id) ON DELETE CASCADE,
			purchase_order_item_id UUID NOT NULL REFERENCES purchase_order_items(id) ON DELETE CASCADE,
			quantity_received INTEGER NOT NULL,
			received_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			received_by VARCHAR(255),
			notes TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_purchase_order_receipts_po_id ON purchase_order_receipts(purchase_order_id)`,
		`CREATE INDEX IF NOT EXISTS idx_purchase_order_receipts_po_item_id ON purchase_order_receipts(purchase_order_item_id)`,
		`CREATE INDEX IF NOT EXISTS idx_purchase_order_receipts_received_date ON purchase_order_receipts(received_date)`,

		// Junction/Supporting Tables
		`CREATE TABLE IF NOT EXISTS product_categories (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
			category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			UNIQUE(product_id, category_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_product_categories_product_id ON product_categories(product_id)`,
		`CREATE INDEX IF NOT EXISTS idx_product_categories_category_id ON product_categories(category_id)`,

		`CREATE TABLE IF NOT EXISTS customer_addresses (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
			address_type VARCHAR(50) NOT NULL,
			address_line1 VARCHAR(255) NOT NULL,
			address_line2 VARCHAR(255),
			city VARCHAR(100) NOT NULL,
			state VARCHAR(100),
			postal_code VARCHAR(20),
			country VARCHAR(100) NOT NULL,
			is_default BOOLEAN DEFAULT false,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_customer_addresses_customer_id ON customer_addresses(customer_id)`,
		`CREATE INDEX IF NOT EXISTS idx_customer_addresses_is_default ON customer_addresses(is_default)`,
	}

	for _, query := range queries {
		_, err := db.SQL.Exec(query)
		if err != nil {
			log.Er("failed to execute migration query", err)
			return err
		}
	}

	log.Info("Database migration completed successfully")
	return nil
}

func migrateDown(db *database.DB) error {
	log := db.Logger.Function("migrateDown")

	queries := []string{
		// Drop tables in reverse order of dependencies
		`DROP TABLE IF EXISTS customer_addresses CASCADE`,
		`DROP TABLE IF EXISTS product_categories CASCADE`,
		`DROP TABLE IF EXISTS purchase_order_receipts CASCADE`,
		`DROP TABLE IF EXISTS purchase_order_items CASCADE`,
		`DROP TABLE IF EXISTS purchase_orders CASCADE`,
		`DROP TABLE IF EXISTS sales_order_payments CASCADE`,
		`DROP TABLE IF EXISTS sales_order_items CASCADE`,
		`DROP TABLE IF EXISTS sales_orders CASCADE`,
		`DROP TABLE IF EXISTS supplier_products CASCADE`,
		`DROP TABLE IF EXISTS product_costs CASCADE`,
		`DROP TABLE IF EXISTS product_prices CASCADE`,
		`DROP TABLE IF EXISTS inventory_transactions CASCADE`,
		`DROP TABLE IF EXISTS inventory CASCADE`,
		`DROP TABLE IF EXISTS products CASCADE`,
		`DROP TABLE IF EXISTS warehouses CASCADE`,
		`DROP TABLE IF EXISTS customers CASCADE`,
		`DROP TABLE IF EXISTS suppliers CASCADE`,
		`DROP TABLE IF EXISTS categories CASCADE`,
	}

	for _, query := range queries {
		_, err := db.SQL.Exec(query)
		if err != nil {
			log.Er("failed to execute rollback query", err)
			return err
		}
	}

	log.Info("Database rollback completed successfully")
	return nil
}

func seed(db *database.DB) error {
	log := db.Logger.Function("seed")

	// Clear existing data in reverse dependency order
	tables := []string{
		"customer_addresses", "product_categories", "purchase_order_receipts",
		"purchase_order_items", "purchase_orders", "sales_order_payments",
		"sales_order_items", "sales_orders", "supplier_products",
		"product_costs", "product_prices", "inventory_transactions",
		"inventory", "products", "warehouses", "customers", "suppliers", "categories",
	}

	for _, table := range tables {
		_, err := db.SQL.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			log.Er(fmt.Sprintf("failed to clear %s", table), err)
			return err
		}
	}

	// Seed Categories (hierarchical)
	categoryIDs := make(map[string]string)
	categories := []struct{ name, parent, desc string }{
		{"Electronics", "", "Electronic devices and accessories"},
		{"Computers", "Electronics", "Desktop and laptop computers"},
		{"Laptops", "Computers", "Portable computers"},
		{"Desktops", "Computers", "Desktop computers"},
		{"Peripherals", "Computers", "Computer accessories"},
		{"Home & Garden", "", "Home and garden products"},
		{"Furniture", "Home & Garden", "Home furniture"},
		{"Tools", "Home & Garden", "Hand and power tools"},
		{"Clothing", "", "Apparel and accessories"},
		{"Men's Clothing", "Clothing", "Men's apparel"},
		{"Women's Clothing", "Clothing", "Women's apparel"},
	}

	for _, cat := range categories {
		var parentID interface{}
		if cat.parent != "" {
			parentID = categoryIDs[cat.parent]
		}
		var id string
		err := db.SQL.QueryRow(
			`INSERT INTO categories (name, description, parent_id) VALUES ($1, $2, $3) RETURNING id`,
			cat.name, cat.desc, parentID,
		).Scan(&id)
		if err != nil {
			log.Er("failed to insert category", err)
			return err
		}
		categoryIDs[cat.name] = id
	}

	// Seed Suppliers
	supplierIDs := make([]string, 0, 10)
	suppliers := []struct{ name, contact, email, city, country string }{
		{"Tech Distributors Inc", "John Smith", "john@techdist.com", "San Francisco", "USA"},
		{"Global Electronics Supply", "Mary Johnson", "mary@globalelec.com", "New York", "USA"},
		{"Furniture Wholesale Co", "Bob Williams", "bob@furnwhole.com", "Chicago", "USA"},
		{"Apparel Source Ltd", "Sarah Davis", "sarah@apparelsrc.com", "Los Angeles", "USA"},
		{"Hardware Imports", "Mike Wilson", "mike@hardimports.com", "Seattle", "USA"},
	}

	for _, sup := range suppliers {
		var id string
		err := db.SQL.QueryRow(
			`INSERT INTO suppliers (name, contact_name, email, city, country) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			sup.name, sup.contact, sup.email, sup.city, sup.country,
		).Scan(&id)
		if err != nil {
			log.Er("failed to insert supplier", err)
			return err
		}
		supplierIDs = append(supplierIDs, id)
	}

	// Seed Warehouses
	warehouseIDs := make([]string, 0, 5)
	warehouses := []struct{ name, code, city, country string }{
		{"Main Distribution Center", "WH-MAIN", "Dallas", "USA"},
		{"West Coast Warehouse", "WH-WEST", "Los Angeles", "USA"},
		{"East Coast Warehouse", "WH-EAST", "New York", "USA"},
	}

	for _, wh := range warehouses {
		var id string
		err := db.SQL.QueryRow(
			`INSERT INTO warehouses (name, code, city, country) VALUES ($1, $2, $3, $4) RETURNING id`,
			wh.name, wh.code, wh.city, wh.country,
		).Scan(&id)
		if err != nil {
			log.Er("failed to insert warehouse", err)
			return err
		}
		warehouseIDs = append(warehouseIDs, id)
	}

	// Seed Products (50 products)
	productIDs := make([]string, 0, 50)
	products := []struct{ sku, name, desc string; weight float64 }{
		{"LAPTOP-001", "Dell XPS 13", "13-inch ultrabook laptop", 1.2},
		{"LAPTOP-002", "MacBook Pro 14", "Apple 14-inch laptop", 1.6},
		{"LAPTOP-003", "ThinkPad X1", "Lenovo business laptop", 1.4},
		{"DESKTOP-001", "Gaming PC Pro", "High-end gaming desktop", 8.5},
		{"DESKTOP-002", "Office Workstation", "Business desktop computer", 7.2},
		{"MOUSE-001", "Wireless Mouse", "Ergonomic wireless mouse", 0.1},
		{"KEYBOARD-001", "Mechanical Keyboard", "RGB mechanical keyboard", 0.8},
		{"MONITOR-001", "27\" 4K Monitor", "4K UHD display", 5.5},
		{"MONITOR-002", "34\" Ultrawide", "Curved ultrawide monitor", 7.8},
		{"CHAIR-001", "Ergonomic Office Chair", "Adjustable office chair", 15.0},
	}

	for i, prod := range products {
		var id string
		err := db.SQL.QueryRow(
			`INSERT INTO products (sku, name, description, weight, is_active) VALUES ($1, $2, $3, $4, true) RETURNING id`,
			prod.sku, prod.name, prod.desc, prod.weight,
		).Scan(&id)
		if err != nil {
			log.Er("failed to insert product", err)
			return err
		}
		productIDs = append(productIDs, id)

		// Link products to categories
		catID := categoryIDs["Laptops"]
		if i >= 3 && i < 5 {
			catID = categoryIDs["Desktops"]
		} else if i >= 5 && i < 9 {
			catID = categoryIDs["Peripherals"]
		} else if i >= 9 {
			catID = categoryIDs["Furniture"]
		}

		_, err = db.SQL.Exec(
			`INSERT INTO product_categories (product_id, category_id) VALUES ($1, $2)`,
			id, catID,
		)
		if err != nil {
			log.Er("failed to link product to category", err)
			return err
		}

		// Add product prices
		price := float64(299 + (i * 50))
		_, err = db.SQL.Exec(
			`INSERT INTO product_prices (product_id, price, currency) VALUES ($1, $2, 'USD')`,
			id, price,
		)
		if err != nil {
			log.Er("failed to insert product price", err)
			return err
		}

		// Add product costs (70% of price)
		cost := price * 0.7
		_, err = db.SQL.Exec(
			`INSERT INTO product_costs (product_id, cost, currency) VALUES ($1, $2, 'USD')`,
			id, cost,
		)
		if err != nil {
			log.Er("failed to insert product cost", err)
			return err
		}

		// Link to suppliers
		supplierID := supplierIDs[i%len(supplierIDs)]
		_, err = db.SQL.Exec(
			`INSERT INTO supplier_products (supplier_id, product_id, cost, lead_time_days, minimum_order_quantity) VALUES ($1, $2, $3, $4, $5)`,
			supplierID, id, cost, 7+i%14, 10+i*5,
		)
		if err != nil {
			log.Er("failed to link product to supplier", err)
			return err
		}

		// Add inventory for each warehouse
		for _, whID := range warehouseIDs {
			qty := 100 + (i * 10)
			_, err = db.SQL.Exec(
				`INSERT INTO inventory (product_id, warehouse_id, quantity, reserved_quantity, reorder_point, reorder_quantity) VALUES ($1, $2, $3, 0, 50, 100)`,
				id, whID, qty,
			)
			if err != nil {
				log.Er("failed to insert inventory", err)
				return err
			}
		}
	}

	// Seed Customers (20 customers)
	customerIDs := make([]string, 0, 20)
	firstNames := []string{"John", "Jane", "Bob", "Alice", "Charlie", "Diana", "Frank", "Grace", "Henry", "Iris"}
	lastNames := []string{"Smith", "Doe", "Johnson", "Williams", "Brown", "Davis", "Miller", "Wilson", "Moore", "Taylor"}

	for i := 0; i < 20; i++ {
		firstName := firstNames[i%len(firstNames)]
		lastName := lastNames[i%len(lastNames)]
		email := fmt.Sprintf("%s.%s%d@example.com", firstName, lastName, i+1)
		phone := fmt.Sprintf("+1-555-%04d", 1000+i)

		var id string
		err := db.SQL.QueryRow(
			`INSERT INTO customers (first_name, last_name, email, phone) VALUES ($1, $2, $3, $4) RETURNING id`,
			firstName, lastName, email, phone,
		).Scan(&id)
		if err != nil {
			log.Er("failed to insert customer", err)
			return err
		}
		customerIDs = append(customerIDs, id)

		// Add customer addresses
		_, err = db.SQL.Exec(
			`INSERT INTO customer_addresses (customer_id, address_type, address_line1, city, state, postal_code, country, is_default) VALUES ($1, 'shipping', $2, 'Austin', 'TX', '78701', 'USA', true)`,
			id, fmt.Sprintf("%d Main St", 100+i),
		)
		if err != nil {
			log.Er("failed to insert customer address", err)
			return err
		}
	}

	// Seed Sales Orders (30 orders)
	for i := 0; i < 30; i++ {
		customerID := customerIDs[i%len(customerIDs)]
		orderNum := fmt.Sprintf("SO-%06d", 1000+i)

		var orderID string
		err := db.SQL.QueryRow(
			`INSERT INTO sales_orders (order_number, customer_id, status, subtotal, tax, shipping, total) VALUES ($1, $2, $3, 0, 0, 0, 0) RETURNING id`,
			orderNum, customerID, "completed",
		).Scan(&orderID)
		if err != nil {
			log.Er("failed to insert sales order", err)
			return err
		}

		// Add 1-5 items per order
		itemCount := 1 + (i % 5)
		orderTotal := 0.0

		for j := 0; j < itemCount; j++ {
			productID := productIDs[(i+j)%len(productIDs)]
			qty := 1 + (j % 3)
			unitPrice := 299.0 + float64((i+j)*50)
			total := float64(qty) * unitPrice
			orderTotal += total

			_, err = db.SQL.Exec(
				`INSERT INTO sales_order_items (sales_order_id, product_id, quantity, unit_price, discount, tax, total) VALUES ($1, $2, $3, $4, 0, $5, $6)`,
				orderID, productID, qty, unitPrice, total*0.08, total,
			)
			if err != nil {
				log.Er("failed to insert sales order item", err)
				return err
			}
		}

		// Update order totals
		tax := orderTotal * 0.08
		shipping := 15.0
		finalTotal := orderTotal + tax + shipping

		_, err = db.SQL.Exec(
			`UPDATE sales_orders SET subtotal = $1, tax = $2, shipping = $3, total = $4 WHERE id = $5`,
			orderTotal, tax, shipping, finalTotal, orderID,
		)
		if err != nil {
			log.Er("failed to update sales order totals", err)
			return err
		}

		// Add payment
		_, err = db.SQL.Exec(
			`INSERT INTO sales_order_payments (sales_order_id, payment_method, amount, status, transaction_id) VALUES ($1, 'credit_card', $2, 'completed', $3)`,
			orderID, finalTotal, fmt.Sprintf("TXN-%d", 10000+i),
		)
		if err != nil {
			log.Er("failed to insert sales order payment", err)
			return err
		}
	}

	// Seed Purchase Orders (10 large orders)
	for i := 0; i < 10; i++ {
		supplierID := supplierIDs[i%len(supplierIDs)]
		warehouseID := warehouseIDs[i%len(warehouseIDs)]
		poNum := fmt.Sprintf("PO-%06d", 2000+i)

		var poID string
		err := db.SQL.QueryRow(
			`INSERT INTO purchase_orders (po_number, supplier_id, warehouse_id, status, subtotal, tax, shipping, total) VALUES ($1, $2, $3, $4, 0, 0, 0, 0) RETURNING id`,
			poNum, supplierID, warehouseID, "received",
		).Scan(&poID)
		if err != nil {
			log.Er("failed to insert purchase order", err)
			return err
		}

		// Add many items (simulate large delivery)
		itemCount := 5 + (i * 2)
		poTotal := 0.0

		for j := 0; j < itemCount && j < len(productIDs); j++ {
			productID := productIDs[j]
			qty := 50 + (j * 10)
			unitCost := 200.0 + float64(j*30)
			total := float64(qty) * unitCost
			poTotal += total

			var poItemID string
			err = db.SQL.QueryRow(
				`INSERT INTO purchase_order_items (purchase_order_id, product_id, quantity, unit_cost, tax, total, received_quantity) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
				poID, productID, qty, unitCost, total*0.06, total, qty,
			).Scan(&poItemID)
			if err != nil {
				log.Er("failed to insert purchase order item", err)
				return err
			}

			// Add receipt record
			_, err = db.SQL.Exec(
				`INSERT INTO purchase_order_receipts (purchase_order_id, purchase_order_item_id, quantity_received, received_by) VALUES ($1, $2, $3, 'Warehouse Staff')`,
				poID, poItemID, qty,
			)
			if err != nil {
				log.Er("failed to insert purchase order receipt", err)
				return err
			}
		}

		// Update PO totals
		tax := poTotal * 0.06
		shipping := 50.0
		finalTotal := poTotal + tax + shipping

		_, err = db.SQL.Exec(
			`UPDATE purchase_orders SET subtotal = $1, tax = $2, shipping = $3, total = $4 WHERE id = $5`,
			poTotal, tax, shipping, finalTotal, poID,
		)
		if err != nil {
			log.Er("failed to update purchase order totals", err)
			return err
		}
	}

	log.Info("Database seeded successfully with sales/inventory data")
	return nil
}