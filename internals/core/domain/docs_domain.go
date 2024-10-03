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

type UpdateStatusRequest struct {
	Status int `json:"status" validate:"required"`
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
	DocumentNumber string    `json:"-" db:"document_number"`
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

type AccessHistory struct {
	AcessorID    string    `json:"acessor_id"`
	DocID        string    `json:"doc_id"`
	AccessorName string    `json:"accessor_name"`
	Type         string    `json:"type"`
	Number       string    `json:"number"`
	AccessTime   time.Time `json:"access_time"`
}

type AcessDeleteRequest struct {
	DocID     string `json:"doc_id" validate:"required"`
	AcessorID string `json:"accessor_id" validate:"required"`
}

type UnverifiedDocs struct {
	ID     string    `json:"id"`
	Name   string    `json:"name"`
	Type   string    `json:"type"`
	Number string    `json:"number"`
	Date   time.Time `json:"date"`
}
