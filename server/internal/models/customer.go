package models

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID        uuid.UUID  `json:"id" db:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FirstName string     `json:"first_name" db:"first_name" gorm:"type:varchar(255);not null"`
	LastName  string     `json:"last_name" db:"last_name" gorm:"type:varchar(255);not null"`
	Email     string     `json:"email" db:"email" gorm:"type:varchar(255);unique;not null;index"`
	Phone     *string    `json:"phone,omitempty" db:"phone" gorm:"type:varchar(50)"`
	CreatedAt time.Time  `json:"created_at" db:"created_at" gorm:"type:timestamptz;default:now()"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at" gorm:"type:timestamptz;default:now()"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at" gorm:"type:timestamptz;index"`
}

type CustomerAddress struct {
	ID           uuid.UUID `json:"id" db:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CustomerID   uuid.UUID `json:"customer_id" db:"customer_id" gorm:"type:uuid;not null;index"`
	AddressType  string    `json:"address_type" db:"address_type" gorm:"type:varchar(50);not null"`
	AddressLine1 string    `json:"address_line1" db:"address_line1" gorm:"type:varchar(255);not null"`
	AddressLine2 *string   `json:"address_line2,omitempty" db:"address_line2" gorm:"type:varchar(255)"`
	City         string    `json:"city" db:"city" gorm:"type:varchar(100);not null"`
	State        *string   `json:"state,omitempty" db:"state" gorm:"type:varchar(100)"`
	PostalCode   *string   `json:"postal_code,omitempty" db:"postal_code" gorm:"type:varchar(20)"`
	Country      string    `json:"country" db:"country" gorm:"type:varchar(100);not null"`
	IsDefault    bool      `json:"is_default" db:"is_default" gorm:"default:false;index"`
	CreatedAt    time.Time `json:"created_at" db:"created_at" gorm:"type:timestamptz;default:now()"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at" gorm:"type:timestamptz;default:now()"`
}
