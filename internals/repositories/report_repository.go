package repositories

import (
	"IFEST/internals/core/domain"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type IReportsRepository interface {
	CreateReport(report domain.ReportCreateRequest) (domain.Report, error)
	GetReportByID(id uuid.UUID) (domain.ReportAccessInfo, error)
	GetReports() ([]domain.Report, error)
	UpdateReport(id uuid.UUID, report domain.ReportUpdateRequest) (domain.Report, error)
	DeleteReport(id uuid.UUID) error
}

type ReportsRepository struct {
	db *sqlx.DB
}

func NewReportsRepository(db *sqlx.DB) IReportsRepository {
	return &ReportsRepository{
		db: db,
	}
}

func (rr *ReportsRepository) CreateReport(report domain.ReportCreateRequest) (domain.Report, error) {
	query := `
        INSERT INTO reports (id, user_id, report_text, created_at)
        VALUES (:id, :user_id, :report_text, CURRENT_TIMESTAMP)
        RETURNING id, user_id, report_text, created_at
    `

	reportID := uuid.New()
	reportData := domain.Report{
		ID:         reportID,
		UserID:     report.UserID,
		ReportText: report.ReportText,
		CreatedAt:  time.Now(),
	}

	rows, err := rr.db.NamedQuery(query, &reportData)
	if err != nil {
		return domain.Report{}, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&reportData)
		if err != nil {
			return domain.Report{}, err
		}
	}

	return reportData, nil
}

func (rr *ReportsRepository) GetReportByID(id uuid.UUID) (domain.ReportAccessInfo, error) {
	query := `
        SELECT 
            id,
            user_id,
            report_text,
            created_at
        FROM 
            reports
        WHERE 
            id = $1
    `

	var report domain.ReportAccessInfo
	err := rr.db.Get(&report, query, id)
	if err != nil {
		return report, err
	}

	return report, nil
}

func (rr *ReportsRepository) GetReports() ([]domain.Report, error) {
	query := `
        SELECT 
           *
        FROM 
            reports
        ORDER BY created_at DESC
    `

	var reports []domain.Report
	err := rr.db.Select(&reports, query)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

func (rr *ReportsRepository) UpdateReport(id uuid.UUID, report domain.ReportUpdateRequest) (domain.Report, error) {
	query := `
        UPDATE reports
        SET 
            report_text = :report_text
        WHERE id = :id
        RETURNING id, user_id, report_text, created_at
    `

	reportData := domain.Report{
		ID:         id,
		ReportText: report.ReportText,
	}

	rows, err := rr.db.NamedQuery(query, &reportData)
	if err != nil {
		return domain.Report{}, err
	}
	defer rows.Close()

	var updatedReport domain.Report
	if rows.Next() {
		err = rows.StructScan(&updatedReport)
		if err != nil {
			return domain.Report{}, err
		}
	} else {
		return domain.Report{}, fmt.Errorf("no report found with id %s", id)
	}

	return updatedReport, nil
}

func (rr *ReportsRepository) DeleteReport(id uuid.UUID) error {
	query := `
        DELETE FROM reports
        WHERE id = $1
    `
	result, err := rr.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no report found with id %s", id)
	}

	return nil
}
