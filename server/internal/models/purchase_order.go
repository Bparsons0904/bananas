package models

import (
	"time"

	"github.com/google/uuid"
)

type PurchaseOrder struct {
	ID           uuid.UUID  `json:"id" db:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PONumber     string     `json:"po_number" db:"po_number" gorm:"type:varchar(100);unique;not null;index"`
	SupplierID   uuid.UUID  `json:"supplier_id" db:"supplier_id" gorm:"type:uuid;not null;index"`
	WarehouseID  uuid.UUID  `json:"warehouse_id" db:"warehouse_id" gorm:"type:uuid;not null;index"`
	OrderDate    time.Time  `json:"order_date" db:"order_date" gorm:"type:timestamptz;default:now();index"`
	ExpectedDate *time.Time `json:"expected_date,omitempty" db:"expected_date" gorm:"type:timestamptz"`
	Status       string     `json:"status" db:"status" gorm:"type:varchar(50);default:'pending';index"`
	Subtotal     float64    `json:"subtotal" db:"subtotal" gorm:"type:decimal(10,2);not null;default:0"`
	Tax          float64    `json:"tax" db:"tax" gorm:"type:decimal(10,2);not null;default:0"`
	Shipping     float64    `json:"shipping" db:"shipping" gorm:"type:decimal(10,2);not null;default:0"`
	Total        float64    `json:"total" db:"total" gorm:"type:decimal(10,2);not null;default:0"`
	Notes        *string    `json:"notes,omitempty" db:"notes" gorm:"type:text"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at" gorm:"type:timestamptz;default:now()"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at" gorm:"type:timestamptz;default:now()"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at" gorm:"type:timestamptz;index"`
}

type PurchaseOrderItem struct {
	ID               uuid.UUID `json:"id" db:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PurchaseOrderID  uuid.UUID `json:"purchase_order_id" db:"purchase_order_id" gorm:"type:uuid;not null;index"`
	ProductID        uuid.UUID `json:"product_id" db:"product_id" gorm:"type:uuid;not null;index"`
	Quantity         int       `json:"quantity" db:"quantity" gorm:"not null"`
	UnitCost         float64   `json:"unit_cost" db:"unit_cost" gorm:"type:decimal(10,2);not null"`
	Tax              float64   `json:"tax" db:"tax" gorm:"type:decimal(10,2);default:0"`
	Total            float64   `json:"total" db:"total" gorm:"type:decimal(10,2);not null"`
	ReceivedQuantity int       `json:"received_quantity" db:"received_quantity" gorm:"default:0"`
	CreatedAt        time.Time `json:"created_at" db:"created_at" gorm:"type:timestamptz;default:now()"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at" gorm:"type:timestamptz;default:now()"`
}

type PurchaseOrderReceipt struct {
	ID                  uuid.UUID `json:"id" db:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PurchaseOrderID     uuid.UUID `json:"purchase_order_id" db:"purchase_order_id" gorm:"type:uuid;not null;index"`
	PurchaseOrderItemID uuid.UUID `json:"purchase_order_item_id" db:"purchase_order_item_id" gorm:"type:uuid;not null;index"`
	QuantityReceived    int       `json:"quantity_received" db:"quantity_received" gorm:"not null"`
	ReceivedDate        time.Time `json:"received_date" db:"received_date" gorm:"type:timestamptz;default:now();index"`
	ReceivedBy          *string   `json:"received_by,omitempty" db:"received_by" gorm:"type:varchar(255)"`
	Notes               *string   `json:"notes,omitempty" db:"notes" gorm:"type:text"`
	CreatedAt           time.Time `json:"created_at" db:"created_at" gorm:"type:timestamptz;default:now()"`
}
