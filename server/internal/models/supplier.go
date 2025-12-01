package models

import (
	"time"

	"github.com/google/uuid"
)

type Supplier struct {
	ID          uuid.UUID  `json:"id" db:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string     `json:"name" db:"name" gorm:"type:varchar(255);not null"`
	ContactName *string    `json:"contact_name,omitempty" db:"contact_name" gorm:"type:varchar(255)"`
	Email       *string    `json:"email,omitempty" db:"email" gorm:"type:varchar(255)"`
	Phone       *string    `json:"phone,omitempty" db:"phone" gorm:"type:varchar(50)"`
	Address     *string    `json:"address,omitempty" db:"address" gorm:"type:text"`
	City        *string    `json:"city,omitempty" db:"city" gorm:"type:varchar(100)"`
	State       *string    `json:"state,omitempty" db:"state" gorm:"type:varchar(100)"`
	PostalCode  *string    `json:"postal_code,omitempty" db:"postal_code" gorm:"type:varchar(20)"`
	Country     *string    `json:"country,omitempty" db:"country" gorm:"type:varchar(100)"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at" gorm:"type:timestamptz;default:now()"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at" gorm:"type:timestamptz;default:now()"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at" gorm:"type:timestamptz;index"`
}
