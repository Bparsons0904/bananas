package generators

import (
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

// PurchaseOrderData holds generated purchase order IDs and item IDs
type PurchaseOrderData struct {
	OrderIDs []uuid.UUID
	ItemIDs  []uuid.UUID
}

// GeneratePurchaseOrders creates purchase order data
func GeneratePurchaseOrders(count int, idMap *IDMap) ([][]interface{}, *PurchaseOrderData) {
	rows := make([][]interface{}, 0, count)
	data := &PurchaseOrderData{
		OrderIDs: make([]uuid.UUID, 0, count),
	}
	now := time.Now()

	statuses := []string{"pending", "confirmed", "partially_received", "received", "cancelled"}

	for i := 0; i < count; i++ {
		id := uuid.New()
		data.OrderIDs = append(data.OrderIDs, id)

		// Random supplier
		supplierID := idMap.SupplierIDs[gofakeit.Number(0, len(idMap.SupplierIDs)-1)]

		// Random warehouse
		warehouseID := idMap.WarehouseIDs[gofakeit.Number(0, len(idMap.WarehouseIDs)-1)]

		// Random order date in the past year
		orderDate := now.AddDate(0, 0, -gofakeit.Number(0, 365))

		// Expected date 7-60 days after order
		expectedDays := gofakeit.Number(7, 60)
		expectedDate := orderDate.AddDate(0, 0, expectedDays)

		status := statuses[gofakeit.Number(0, len(statuses)-1)]
		notes := gofakeit.Sentence(15)

		// Totals will be calculated from items
		rows = append(rows, []interface{}{
			id,                               // id
			fmt.Sprintf("PO-%010d", i+1),     // po_number
			supplierID,                       // supplier_id
			warehouseID,                      // warehouse_id
			orderDate,                        // order_date
			&expectedDate,                    // expected_date
			status,                           // status
			0.0,                              // subtotal
			0.0,                              // tax
			0.0,                              // shipping
			0.0,                              // total
			&notes,                           // notes
			now,                              // created_at
			now,                              // updated_at
			nil,                              // deleted_at
		})
	}

	return rows, data
}

// GeneratePurchaseOrderItems creates purchase order line items
func GeneratePurchaseOrderItems(purchaseData *PurchaseOrderData, idMap *IDMap, itemsPerOrder int) [][]interface{} {
	rows := make([][]interface{}, 0, len(purchaseData.OrderIDs)*itemsPerOrder)
	now := time.Now()

	for _, orderID := range purchaseData.OrderIDs {
		// Generate 10-100 items per order (bulk purchases)
		numItems := gofakeit.Number(itemsPerOrder-40, itemsPerOrder+50)

		usedProducts := make(map[uuid.UUID]bool)

		for j := 0; j < numItems; j++ {
			productID := idMap.ProductIDs[gofakeit.Number(0, len(idMap.ProductIDs)-1)]

			// Avoid duplicate products in same order
			if usedProducts[productID] {
				continue
			}
			usedProducts[productID] = true

			itemID := uuid.New()
			purchaseData.ItemIDs = append(purchaseData.ItemIDs, itemID)

			quantity := gofakeit.Number(10, 1000) // Bulk quantities
			unitCost := float64(gofakeit.Number(200, 30000)) / 100.0
			tax := (unitCost * float64(quantity)) * 0.05 // 5% tax
			total := (unitCost * float64(quantity)) + tax

			// Some items partially received
			receivedQty := 0
			if gofakeit.Bool() {
				receivedQty = gofakeit.Number(0, quantity)
			}

			rows = append(rows, []interface{}{
				itemID,         // id
				orderID,        // purchase_order_id
				productID,      // product_id
				quantity,       // quantity
				unitCost,       // unit_cost
				tax,            // tax
				total,          // total
				receivedQty,    // received_quantity
				now,            // created_at
				now,            // updated_at
			})
		}
	}

	return rows
}

// GeneratePurchaseOrderReceipts creates receipt records for items
func GeneratePurchaseOrderReceipts(purchaseData *PurchaseOrderData, receiptsPerItem int) [][]interface{} {
	rows := make([][]interface{}, 0, len(purchaseData.ItemIDs)*receiptsPerItem)
	now := time.Now()

	for i, itemID := range purchaseData.ItemIDs {
		// Get the order ID for this item (use modulo to map back)
		orderIdx := i / 50 // Approximate, assuming ~50 items per order
		if orderIdx >= len(purchaseData.OrderIDs) {
			orderIdx = len(purchaseData.OrderIDs) - 1
		}
		orderID := purchaseData.OrderIDs[orderIdx]

		// Usually 1 receipt, sometimes 2 (partial deliveries)
		numReceipts := gofakeit.Number(1, receiptsPerItem+1)

		for j := 0; j < numReceipts; j++ {
			qtyReceived := gofakeit.Number(10, 100)
			receivedDate := now.AddDate(0, 0, -gofakeit.Number(0, 365))
			receivedBy := gofakeit.Name()
			notes := gofakeit.Sentence(10)

			rows = append(rows, []interface{}{
				uuid.New(),          // id
				orderID,             // purchase_order_id
				itemID,              // purchase_order_item_id
				qtyReceived,         // quantity_received
				receivedDate,        // received_date
				&receivedBy,         // received_by
				&notes,              // notes
				now,                 // created_at
			})
		}
	}

	return rows
}

// Purchase-related column functions
func PurchaseOrderColumns() []string {
	return []string{"id", "po_number", "supplier_id", "warehouse_id", "order_date",
		"expected_date", "status", "subtotal", "tax", "shipping", "total", "notes",
		"created_at", "updated_at", "deleted_at"}
}

func PurchaseOrderItemColumns() []string {
	return []string{"id", "purchase_order_id", "product_id", "quantity", "unit_cost",
		"tax", "total", "received_quantity", "created_at", "updated_at"}
}

func PurchaseOrderReceiptColumns() []string {
	return []string{"id", "purchase_order_id", "purchase_order_item_id", "quantity_received",
		"received_date", "received_by", "notes", "created_at"}
}
