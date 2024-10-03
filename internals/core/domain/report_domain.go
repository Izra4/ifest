package domain

import (
	"github.com/google/uuid"
	"time"
)

type Report struct {
	ID         uuid.UUID `db:"id" json:"id"`
	UserID     uuid.UUID `db:"user_id" json:"user_id"`
	ReportText string    `db:"report_text" json:"report_text" validate:"required"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type ReportCreateRequest struct {
	UserID     uuid.UUID `json:"user_id" validate:"required"`
	ReportText string    `json:"report_text" validate:"required"`
}

type ReportUpdateRequest struct {
	ReportText string `json:"report_text" validate:"required"`
}

type ReportAccessInfo struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	ReportText string    `json:"report_text"`
	CreatedAt  time.Time `json:"created_at"`
}
