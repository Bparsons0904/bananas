package generators

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

// GenerateInventory creates inventory records (product-warehouse combinations)
func GenerateInventory(idMap *IDMap, warehousesPerProduct int) [][]interface{} {
	rows := make([][]interface{}, 0, len(idMap.ProductIDs)*warehousesPerProduct)
	now := time.Now()

	for _, productID := range idMap.ProductIDs {
		// Each product exists in 1-3 warehouses
		numWarehouses := gofakeit.Number(1, warehousesPerProduct+1)
		if numWarehouses > len(idMap.WarehouseIDs) {
			numWarehouses = len(idMap.WarehouseIDs)
		}

		usedWarehouses := make(map[uuid.UUID]bool)

		for j := 0; j < numWarehouses; j++ {
			warehouseID := idMap.WarehouseIDs[gofakeit.Number(0, len(idMap.WarehouseIDs)-1)]

			// Avoid duplicate warehouse assignments
			if usedWarehouses[warehouseID] {
				continue
			}
			usedWarehouses[warehouseID] = true

			quantity := gofakeit.Number(0, 10000)
			reservedQuantity := gofakeit.Number(0, quantity/10) // Reserve up to 10%
			reorderPoint := gofakeit.Number(50, 500)
			reorderQuantity := gofakeit.Number(100, 1000)

			rows = append(rows, []interface{}{
				uuid.New(),        // id
				productID,         // product_id
				warehouseID,       // warehouse_id
				quantity,          // quantity
				reservedQuantity,  // reserved_quantity
				reorderPoint,      // reorder_point
				reorderQuantity,   // reorder_quantity
				now,               // created_at
				now,               // updated_at
			})
		}
	}

	return rows
}

// GenerateInventoryTransactions creates transaction history
func GenerateInventoryTransactions(idMap *IDMap, transactionsPerProduct int) [][]interface{} {
	rows := make([][]interface{}, 0, len(idMap.ProductIDs)*transactionsPerProduct)
	now := time.Now()

	transactionTypes := []string{"purchase", "sale", "adjustment", "return", "transfer", "damage"}

	for _, productID := range idMap.ProductIDs {
		// Random warehouse for transactions
		warehouseID := idMap.WarehouseIDs[gofakeit.Number(0, len(idMap.WarehouseIDs)-1)]

		for j := 0; j < transactionsPerProduct; j++ {
			txnType := transactionTypes[gofakeit.Number(0, len(transactionTypes)-1)]

			// Quantity is positive for additions (purchase, return) and negative for subtractions (sale, damage)
			quantity := gofakeit.Number(1, 100)
			if txnType == "sale" || txnType == "damage" {
				quantity = -quantity
			}

			// Random reference ID (could be sales order, purchase order, etc.)
			referenceID := uuid.New()
			referenceType := "sales_order"
			if txnType == "purchase" {
				referenceType = "purchase_order"
			}

			notes := gofakeit.Sentence(8)
			createdAt := now.AddDate(0, 0, -gofakeit.Number(0, 365))

			rows = append(rows, []interface{}{
				uuid.New(),        // id
				productID,         // product_id
				warehouseID,       // warehouse_id
				txnType,           // transaction_type
				quantity,          // quantity
				&referenceID,      // reference_id
				&referenceType,    // reference_type
				&notes,            // notes
				createdAt,         // created_at
			})
		}
	}

	return rows
}

// Inventory-related column functions
func InventoryColumns() []string {
	return []string{"id", "product_id", "warehouse_id", "quantity", "reserved_quantity",
		"reorder_point", "reorder_quantity", "created_at", "updated_at"}
}

func InventoryTransactionColumns() []string {
	return []string{"id", "product_id", "warehouse_id", "transaction_type", "quantity",
		"reference_id", "reference_type", "notes", "created_at"}
}
