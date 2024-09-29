package domain

import (
	"github.com/google/uuid"
	"time"
)

type Docs struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	Type      string    `json:"type" db:"type"`
	Status    int       `json:"status" db:"status"`
	Number    string    `json:"number" db:"number"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type DocsUpload struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Name   string    `json:"name" validate:"required"`
	Number string    `json:"number" validate:"required"`
	Type   string    `json:"type" validate:"required"`
}
