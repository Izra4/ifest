package services

import (
	"IFEST/internals/core/domain"
	"IFEST/internals/repositories"
	"github.com/google/uuid"
)

type IReportsService interface {
	CreateReport(report domain.ReportCreateRequest) (domain.Report, error)
	GetReportByID(id uuid.UUID) (domain.ReportAccessInfo, error)
	GetReports() ([]domain.Report, error)
	UpdateReport(id uuid.UUID, report domain.ReportUpdateRequest) (domain.Report, error)
	DeleteReport(id uuid.UUID) error
}

type ReportsService struct {
	repo repositories.IReportsRepository
}

func NewReportsService(repo repositories.IReportsRepository) IReportsService {
	return &ReportsService{
		repo: repo,
	}
}

func (rs *ReportsService) CreateReport(report domain.ReportCreateRequest) (domain.Report, error) {
	return rs.repo.CreateReport(report)
}

func (rs *ReportsService) GetReportByID(id uuid.UUID) (domain.ReportAccessInfo, error) {
	return rs.repo.GetReportByID(id)
}

func (rs *ReportsService) GetReports() ([]domain.Report, error) {
	return rs.repo.GetReports()
}

func (rs *ReportsService) UpdateReport(id uuid.UUID, report domain.ReportUpdateRequest) (domain.Report, error) {
	return rs.repo.UpdateReport(id, report)
}

func (rs *ReportsService) DeleteReport(id uuid.UUID) error {
	return rs.repo.DeleteReport(id)
}
