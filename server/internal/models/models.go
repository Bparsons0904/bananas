package models

import (
	"time"
)

type TestResult struct {
	ID          int        `json:"id" db:"id"`
	Framework   string     `json:"framework" db:"framework"`
	TestType    string     `json:"test_type" db:"test_type"`
	ExecutionMs int        `json:"execution_ms" db:"execution_ms"`
	Success     bool       `json:"success" db:"success"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type Framework struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Type        string    `json:"type" db:"type"` // "backend" or "frontend"
	Description string    `json:"description" db:"description"`
	Enabled     bool      `json:"enabled" db:"enabled"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}