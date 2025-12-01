package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID  `json:"id" db:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SKU         string     `json:"sku" db:"sku" gorm:"type:varchar(100);unique;not null;index"`
	Name        string     `json:"name" db:"name" gorm:"type:varchar(255);not null"`
	Description *string    `json:"description,omitempty" db:"description" gorm:"type:text"`
	Weight      *float64   `json:"weight,omitempty" db:"weight" gorm:"type:decimal(10,2)"`
	Dimensions  *string    `json:"dimensions,omitempty" db:"dimensions" gorm:"type:varchar(100)"`
	IsActive    bool       `json:"is_active" db:"is_active" gorm:"default:true;index"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at" gorm:"type:timestamptz;default:now()"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at" gorm:"type:timestamptz;default:now()"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at" gorm:"type:timestamptz;index"`
}

type ProductCategory struct {
	ID         uuid.UUID `json:"id" db:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID  uuid.UUID `json:"product_id" db:"product_id" gorm:"type:uuid;not null;index"`
	CategoryID uuid.UUID `json:"category_id" db:"category_id" gorm:"type:uuid;not null;index"`
	CreatedAt  time.Time `json:"created_at" db:"created_at" gorm:"type:timestamptz;default:now()"`
}

type ProductPrice struct {
	ID            uuid.UUID  `json:"id" db:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID     uuid.UUID  `json:"product_id" db:"product_id" gorm:"type:uuid;not null;index"`
	Price         float64    `json:"price" db:"price" gorm:"type:decimal(10,2);not null"`
	Currency      string     `json:"currency" db:"currency" gorm:"type:varchar(3);default:'USD'"`
	EffectiveDate time.Time  `json:"effective_date" db:"effective_date" gorm:"type:timestamptz;default:now();index"`
	EndDate       *time.Time `json:"end_date,omitempty" db:"end_date" gorm:"type:timestamptz"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at" gorm:"type:timestamptz;default:now()"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at" gorm:"type:timestamptz;default:now()"`
}

type ProductCost struct {
	ID            uuid.UUID  `json:"id" db:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID     uuid.UUID  `json:"product_id" db:"product_id" gorm:"type:uuid;not null;index"`
	Cost          float64    `json:"cost" db:"cost" gorm:"type:decimal(10,2);not null"`
	Currency      string     `json:"currency" db:"currency" gorm:"type:varchar(3);default:'USD'"`
	EffectiveDate time.Time  `json:"effective_date" db:"effective_date" gorm:"type:timestamptz;default:now();index"`
	EndDate       *time.Time `json:"end_date,omitempty" db:"end_date" gorm:"type:timestamptz"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at" gorm:"type:timestamptz;default:now()"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at" gorm:"type:timestamptz;default:now()"`
}

type SupplierProduct struct {
	ID                   uuid.UUID  `json:"id" db:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SupplierID           uuid.UUID  `json:"supplier_id" db:"supplier_id" gorm:"type:uuid;not null;index"`
	ProductID            uuid.UUID  `json:"product_id" db:"product_id" gorm:"type:uuid;not null;index"`
	SupplierSKU          *string    `json:"supplier_sku,omitempty" db:"supplier_sku" gorm:"type:varchar(100)"`
	Cost                 *float64   `json:"cost,omitempty" db:"cost" gorm:"type:decimal(10,2)"`
	Currency             string     `json:"currency" db:"currency" gorm:"type:varchar(3);default:'USD'"`
	LeadTimeDays         *int       `json:"lead_time_days,omitempty" db:"lead_time_days" gorm:"type:integer"`
	MinimumOrderQuantity *int       `json:"minimum_order_quantity,omitempty" db:"minimum_order_quantity" gorm:"type:integer"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at" gorm:"type:timestamptz;default:now()"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at" gorm:"type:timestamptz;default:now()"`
}
