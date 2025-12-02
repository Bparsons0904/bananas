package generators

import (
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

// GenerateProducts creates product data
func GenerateProducts(count int, idMap *IDMap) [][]interface{} {
	rows := make([][]interface{}, 0, count)
	now := time.Now()

	for i := 0; i < count; i++ {
		id := uuid.New()
		idMap.ProductIDs = append(idMap.ProductIDs, id)

		description := gofakeit.Paragraph(2, 3, 10, " ")
		weight := float64(gofakeit.Number(1, 10000)) / 100.0
		dimensions := fmt.Sprintf("%dx%dx%d cm",
			gofakeit.Number(1, 100),
			gofakeit.Number(1, 100),
			gofakeit.Number(1, 100))

		rows = append(rows, []interface{}{
			id,                                 // id
			fmt.Sprintf("SKU-%09d", i+1),       // sku
			gofakeit.ProductName(),             // name
			&description,                       // description
			&weight,                            // weight
			&dimensions,                        // dimensions
			true,                               // is_active
			now,                                // created_at
			now,                                // updated_at
			nil,                                // deleted_at
		})
	}

	return rows
}

// GenerateProductCategories creates product-category relationships
func GenerateProductCategories(idMap *IDMap, categoriesPerProduct int) [][]interface{} {
	rows := make([][]interface{}, 0, len(idMap.ProductIDs)*categoriesPerProduct)
	now := time.Now()

	for _, productID := range idMap.ProductIDs {
		// Randomly assign 1-3 categories per product
		numCategories := gofakeit.Number(1, categoriesPerProduct+1)
		usedCategories := make(map[uuid.UUID]bool)

		for j := 0; j < numCategories; j++ {
			categoryID := idMap.CategoryIDs[gofakeit.Number(0, len(idMap.CategoryIDs)-1)]

			// Avoid duplicate category assignments
			if usedCategories[categoryID] {
				continue
			}
			usedCategories[categoryID] = true

			rows = append(rows, []interface{}{
				uuid.New(),  // id
				productID,   // product_id
				categoryID,  // category_id
				now,         // created_at
			})
		}
	}

	return rows
}

// GenerateProductPrices creates pricing history
func GenerateProductPrices(idMap *IDMap, pricesPerProduct int) [][]interface{} {
	rows := make([][]interface{}, 0, len(idMap.ProductIDs)*pricesPerProduct)
	now := time.Now()

	for _, productID := range idMap.ProductIDs {
		basePrice := float64(gofakeit.Number(500, 50000)) / 100.0

		for j := 0; j < pricesPerProduct; j++ {
			// Create price history going back in time
			effectiveDate := now.AddDate(0, 0, -j*30)

			// Price varies by ±20% from base
			priceVariation := float64(gofakeit.Number(80, 120)) / 100.0
			price := basePrice * priceVariation

			var endDate *time.Time
			if j > 0 {
				end := now.AddDate(0, 0, -(j-1)*30).Add(-time.Second)
				endDate = &end
			}

			rows = append(rows, []interface{}{
				uuid.New(),      // id
				productID,       // product_id
				price,           // price
				"USD",           // currency
				effectiveDate,   // effective_date
				endDate,         // end_date
				now,             // created_at
				now,             // updated_at
			})
		}
	}

	return rows
}

// GenerateProductCosts creates cost history
func GenerateProductCosts(idMap *IDMap, costsPerProduct int) [][]interface{} {
	rows := make([][]interface{}, 0, len(idMap.ProductIDs)*costsPerProduct)
	now := time.Now()

	for _, productID := range idMap.ProductIDs {
		baseCost := float64(gofakeit.Number(200, 30000)) / 100.0

		for j := 0; j < costsPerProduct; j++ {
			// Create cost history going back in time
			effectiveDate := now.AddDate(0, 0, -j*30)

			// Cost varies by ±15% from base
			costVariation := float64(gofakeit.Number(85, 115)) / 100.0
			cost := baseCost * costVariation

			var endDate *time.Time
			if j > 0 {
				end := now.AddDate(0, 0, -(j-1)*30).Add(-time.Second)
				endDate = &end
			}

			rows = append(rows, []interface{}{
				uuid.New(),      // id
				productID,       // product_id
				cost,            // cost
				"USD",           // currency
				effectiveDate,   // effective_date
				endDate,         // end_date
				now,             // created_at
				now,             // updated_at
			})
		}
	}

	return rows
}

// GenerateSupplierProducts creates supplier-product relationships
func GenerateSupplierProducts(idMap *IDMap, suppliersPerProduct int) [][]interface{} {
	rows := make([][]interface{}, 0, len(idMap.ProductIDs)*suppliersPerProduct)
	now := time.Now()

	for _, productID := range idMap.ProductIDs {
		// Randomly assign 1-5 suppliers per product
		numSuppliers := gofakeit.Number(1, suppliersPerProduct+2)
		if numSuppliers > len(idMap.SupplierIDs) {
			numSuppliers = len(idMap.SupplierIDs)
		}

		usedSuppliers := make(map[uuid.UUID]bool)

		for j := 0; j < numSuppliers; j++ {
			supplierID := idMap.SupplierIDs[gofakeit.Number(0, len(idMap.SupplierIDs)-1)]

			// Avoid duplicate supplier assignments
			if usedSuppliers[supplierID] {
				continue
			}
			usedSuppliers[supplierID] = true

			supplierSKU := fmt.Sprintf("SUP-%s-%d", supplierID.String()[:8], gofakeit.Number(1000, 9999))
			cost := float64(gofakeit.Number(200, 30000)) / 100.0
			leadTime := gofakeit.Number(1, 60)
			minOrder := gofakeit.Number(1, 100)

			rows = append(rows, []interface{}{
				uuid.New(),     // id
				supplierID,     // supplier_id
				productID,      // product_id
				&supplierSKU,   // supplier_sku
				&cost,          // cost
				"USD",          // currency
				&leadTime,      // lead_time_days
				&minOrder,      // minimum_order_quantity
				now,            // created_at
				now,            // updated_at
			})
		}
	}

	return rows
}

// Product-related column functions
func ProductColumns() []string {
	return []string{"id", "sku", "name", "description", "weight", "dimensions",
		"is_active", "created_at", "updated_at", "deleted_at"}
}

func ProductCategoryColumns() []string {
	return []string{"id", "product_id", "category_id", "created_at"}
}

func ProductPriceColumns() []string {
	return []string{"id", "product_id", "price", "currency", "effective_date",
		"end_date", "created_at", "updated_at"}
}

func ProductCostColumns() []string {
	return []string{"id", "product_id", "cost", "currency", "effective_date",
		"end_date", "created_at", "updated_at"}
}

func SupplierProductColumns() []string {
	return []string{"id", "supplier_id", "product_id", "supplier_sku", "cost",
		"currency", "lead_time_days", "minimum_order_quantity", "created_at", "updated_at"}
}
