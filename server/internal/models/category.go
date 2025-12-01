package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID  `json:"id" db:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string     `json:"name" db:"name" gorm:"type:varchar(255);not null"`
	Description *string    `json:"description,omitempty" db:"description" gorm:"type:text"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty" db:"parent_id" gorm:"type:uuid"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at" gorm:"type:timestamptz;default:now()"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at" gorm:"type:timestamptz;default:now()"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at" gorm:"type:timestamptz;index"`
}
