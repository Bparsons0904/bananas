package generators

import (
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

// SalesOrderData holds generated sales order IDs for item relationships
type SalesOrderData struct {
	OrderIDs []uuid.UUID
}

// GenerateSalesOrders creates sales order data
func GenerateSalesOrders(count int, idMap *IDMap) ([][]interface{}, *SalesOrderData) {
	rows := make([][]interface{}, 0, count)
	data := &SalesOrderData{
		OrderIDs: make([]uuid.UUID, 0, count),
	}
	now := time.Now()

	statuses := []string{"pending", "confirmed", "processing", "shipped", "delivered", "cancelled"}

	for i := 0; i < count; i++ {
		id := uuid.New()
		data.OrderIDs = append(data.OrderIDs, id)

		// Random customer
		customerID := idMap.CustomerIDs[gofakeit.Number(0, len(idMap.CustomerIDs)-1)]

		// Random order date in the past year
		orderDate := now.AddDate(0, 0, -gofakeit.Number(0, 365))

		status := statuses[gofakeit.Number(0, len(statuses)-1)]
		notes := gofakeit.Sentence(15)

		// Totals will be calculated after items are generated
		// For now, use placeholder values
		rows = append(rows, []interface{}{
			id,                               // id
			fmt.Sprintf("SO-%010d", i+1),     // order_number
			customerID,                       // customer_id
			orderDate,                        // order_date
			status,                           // status
			0.0,                              // subtotal (will update)
			0.0,                              // tax (will update)
			0.0,                              // shipping (will update)
			0.0,                              // total (will update)
			&notes,                           // notes
			now,                              // created_at
			now,                              // updated_at
			nil,                              // deleted_at
		})
	}

	return rows, data
}

// GenerateSalesOrderItems creates sales order line items
func GenerateSalesOrderItems(salesData *SalesOrderData, idMap *IDMap, itemsPerOrder int) [][]interface{} {
	rows := make([][]interface{}, 0, len(salesData.OrderIDs)*itemsPerOrder)
	now := time.Now()

	for _, orderID := range salesData.OrderIDs {
		// Generate 1-10 items per order
		numItems := gofakeit.Number(1, itemsPerOrder+7)

		usedProducts := make(map[uuid.UUID]bool)

		for j := 0; j < numItems; j++ {
			productID := idMap.ProductIDs[gofakeit.Number(0, len(idMap.ProductIDs)-1)]

			// Avoid duplicate products in same order
			if usedProducts[productID] {
				continue
			}
			usedProducts[productID] = true

			quantity := gofakeit.Number(1, 20)
			unitPrice := float64(gofakeit.Number(500, 50000)) / 100.0
			discount := 0.0
			if gofakeit.Bool() {
				discount = float64(gofakeit.Number(0, 2000)) / 100.0
			}
			tax := (unitPrice * float64(quantity) - discount) * 0.08 // 8% tax
			total := (unitPrice * float64(quantity)) - discount + tax

			rows = append(rows, []interface{}{
				uuid.New(),     // id
				orderID,        // sales_order_id
				productID,      // product_id
				quantity,       // quantity
				unitPrice,      // unit_price
				discount,       // discount
				tax,            // tax
				total,          // total
				now,            // created_at
				now,            // updated_at
			})
		}
	}

	return rows
}

// GenerateSalesOrderPayments creates payment records
func GenerateSalesOrderPayments(salesData *SalesOrderData, paymentsPerOrder int) [][]interface{} {
	rows := make([][]interface{}, 0, len(salesData.OrderIDs)*paymentsPerOrder)
	now := time.Now()

	paymentMethods := []string{"credit_card", "debit_card", "paypal", "bank_transfer", "cash"}
	paymentStatuses := []string{"pending", "completed", "failed", "refunded"}

	for _, orderID := range salesData.OrderIDs {
		// Usually 1 payment, sometimes 2 (split payments)
		numPayments := gofakeit.Number(1, paymentsPerOrder+1)

		for j := 0; j < numPayments; j++ {
			method := paymentMethods[gofakeit.Number(0, len(paymentMethods)-1)]
			status := paymentStatuses[gofakeit.Number(0, len(paymentStatuses)-1)]
			amount := float64(gofakeit.Number(1000, 100000)) / 100.0
			transactionID := fmt.Sprintf("TXN-%s", uuid.New().String()[:13])
			paymentDate := now.AddDate(0, 0, -gofakeit.Number(0, 365))

			rows = append(rows, []interface{}{
				uuid.New(),       // id
				orderID,          // sales_order_id
				method,           // payment_method
				amount,           // amount
				&transactionID,   // transaction_id
				status,           // status
				paymentDate,      // payment_date
				now,              // created_at
				now,              // updated_at
			})
		}
	}

	return rows
}

// Sales-related column functions
func SalesOrderColumns() []string {
	return []string{"id", "order_number", "customer_id", "order_date", "status",
		"subtotal", "tax", "shipping", "total", "notes", "created_at", "updated_at", "deleted_at"}
}

func SalesOrderItemColumns() []string {
	return []string{"id", "sales_order_id", "product_id", "quantity", "unit_price",
		"discount", "tax", "total", "created_at", "updated_at"}
}

func SalesOrderPaymentColumns() []string {
	return []string{"id", "sales_order_id", "payment_method", "amount", "transaction_id",
		"status", "payment_date", "created_at", "updated_at"}
}
