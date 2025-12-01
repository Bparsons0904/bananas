package models

import (
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	ID               uuid.UUID `json:"id" db:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID        uuid.UUID `json:"product_id" db:"product_id" gorm:"type:uuid;not null;index"`
	WarehouseID      uuid.UUID `json:"warehouse_id" db:"warehouse_id" gorm:"type:uuid;not null;index"`
	Quantity         int       `json:"quantity" db:"quantity" gorm:"not null;default:0"`
	ReservedQuantity int       `json:"reserved_quantity" db:"reserved_quantity" gorm:"not null;default:0"`
	ReorderPoint     int       `json:"reorder_point" db:"reorder_point" gorm:"default:0"`
	ReorderQuantity  int       `json:"reorder_quantity" db:"reorder_quantity" gorm:"default:0"`
	CreatedAt        time.Time `json:"created_at" db:"created_at" gorm:"type:timestamptz;default:now()"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at" gorm:"type:timestamptz;default:now()"`
}

type InventoryTransaction struct {
	ID              uuid.UUID  `json:"id" db:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID       uuid.UUID  `json:"product_id" db:"product_id" gorm:"type:uuid;not null;index"`
	WarehouseID     uuid.UUID  `json:"warehouse_id" db:"warehouse_id" gorm:"type:uuid;not null;index"`
	TransactionType string     `json:"transaction_type" db:"transaction_type" gorm:"type:varchar(50);not null"`
	Quantity        int        `json:"quantity" db:"quantity" gorm:"not null"`
	ReferenceID     *uuid.UUID `json:"reference_id,omitempty" db:"reference_id" gorm:"type:uuid"`
	ReferenceType   *string    `json:"reference_type,omitempty" db:"reference_type" gorm:"type:varchar(50);index:idx_inventory_transactions_reference,composite:reference_type"`
	Notes           *string    `json:"notes,omitempty" db:"notes" gorm:"type:text"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at" gorm:"type:timestamptz;default:now();index"`
}
