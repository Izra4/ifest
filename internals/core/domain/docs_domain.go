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
	UserCount int       `json:"user_count" db:"user_count"`
}

type DocsUpload struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Name   string    `json:"name" validate:"required"`
	Number string    `json:"number" validate:"required"`
	Type   string    `json:"type" validate:"required"`
}

type DocumentAccessInfo struct {
	DocumentID     uuid.UUID `json:"id" db:"document_id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	DocumentName   string    `json:"-" db:"document_name"`
	DocumentType   string    `json:"type" db:"document_type"`
	DocumentStatus int       `json:"status" db:"document_status"`
	AccessCount    int       `json:"user_count" db:"access_count"`
	AccessEmails   string    `json:"-" db:"access_emails"`
	FixedEmails    []string  `json:"emails"`
}

type AccessReq struct {
	DocID      uuid.UUID `json:"doc_id" validate:"required" db:"doc_id"`
	UserID     uuid.UUID `json:"user_id" validate:"required" db:"user_id"`
	Token      string    `json:"token" validate:"required" db:"token"`
	Expired_at time.Time `json:"expired_at" validate:"required" db:"expired_at"`
}
