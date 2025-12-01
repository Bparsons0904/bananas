package models

import (
	"time"

	"github.com/google/uuid"
)

type SalesOrder struct {
	ID          uuid.UUID  `json:"id" db:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderNumber string     `json:"order_number" db:"order_number" gorm:"type:varchar(100);unique;not null;index"`
	CustomerID  uuid.UUID  `json:"customer_id" db:"customer_id" gorm:"type:uuid;not null;index"`
	OrderDate   time.Time  `json:"order_date" db:"order_date" gorm:"type:timestamptz;default:now();index"`
	Status      string     `json:"status" db:"status" gorm:"type:varchar(50);default:'pending';index"`
	Subtotal    float64    `json:"subtotal" db:"subtotal" gorm:"type:decimal(10,2);not null;default:0"`
	Tax         float64    `json:"tax" db:"tax" gorm:"type:decimal(10,2);not null;default:0"`
	Shipping    float64    `json:"shipping" db:"shipping" gorm:"type:decimal(10,2);not null;default:0"`
	Total       float64    `json:"total" db:"total" gorm:"type:decimal(10,2);not null;default:0"`
	Notes       *string    `json:"notes,omitempty" db:"notes" gorm:"type:text"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at" gorm:"type:timestamptz;default:now()"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at" gorm:"type:timestamptz;default:now()"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at" gorm:"type:timestamptz;index"`
}

type SalesOrderItem struct {
	ID            uuid.UUID `json:"id" db:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SalesOrderID  uuid.UUID `json:"sales_order_id" db:"sales_order_id" gorm:"type:uuid;not null;index"`
	ProductID     uuid.UUID `json:"product_id" db:"product_id" gorm:"type:uuid;not null;index"`
	Quantity      int       `json:"quantity" db:"quantity" gorm:"not null"`
	UnitPrice     float64   `json:"unit_price" db:"unit_price" gorm:"type:decimal(10,2);not null"`
	Discount      float64   `json:"discount" db:"discount" gorm:"type:decimal(10,2);default:0"`
	Tax           float64   `json:"tax" db:"tax" gorm:"type:decimal(10,2);default:0"`
	Total         float64   `json:"total" db:"total" gorm:"type:decimal(10,2);not null"`
	CreatedAt     time.Time `json:"created_at" db:"created_at" gorm:"type:timestamptz;default:now()"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at" gorm:"type:timestamptz;default:now()"`
}

type SalesOrderPayment struct {
	ID            uuid.UUID `json:"id" db:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SalesOrderID  uuid.UUID `json:"sales_order_id" db:"sales_order_id" gorm:"type:uuid;not null;index"`
	PaymentMethod string    `json:"payment_method" db:"payment_method" gorm:"type:varchar(50);not null"`
	Amount        float64   `json:"amount" db:"amount" gorm:"type:decimal(10,2);not null"`
	TransactionID *string   `json:"transaction_id,omitempty" db:"transaction_id" gorm:"type:varchar(255)"`
	Status        string    `json:"status" db:"status" gorm:"type:varchar(50);default:'pending';index"`
	PaymentDate   time.Time `json:"payment_date" db:"payment_date" gorm:"type:timestamptz;default:now();index"`
	CreatedAt     time.Time `json:"created_at" db:"created_at" gorm:"type:timestamptz;default:now()"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at" gorm:"type:timestamptz;default:now()"`
}
