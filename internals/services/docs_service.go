package services

import (
	"IFEST/internals/core/domain"
	"IFEST/internals/repositories"
	"github.com/google/uuid"
	"time"
)

type IDocsService interface {
	Upload(docs domain.DocsUpload) (domain.Docs, error)
	FindByID(id string) (domain.DocumentAccessInfo, error)
	FindByUserID(id string) ([]domain.Docs, error)
	UpdateStatus(id uuid.UUID, status int) error
	GetAllDocsByStatus(status int) ([]domain.Docs, error)
}

type DocsService struct {
	docsRepository repositories.IDocsRepository
}

func NewDocsService(docsRepository repositories.IDocsRepository) *DocsService {
	return &DocsService{
		docsRepository: docsRepository,
	}
}

func (d *DocsService) Upload(docs domain.DocsUpload) (domain.Docs, error) {
	ID := uuid.New()
	doc := domain.Docs{
		ID:        ID,
		UserID:    docs.UserID,
		Name:      docs.Name,
		Type:      docs.Type,
		Status:    0,
		Number:    docs.Number,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := d.docsRepository.Upload(doc)
	if err != nil {
		return domain.Docs{}, err
	}

	return result, nil
}

func (d *DocsService) FindByID(id string) (domain.DocumentAccessInfo, error) {
	return d.docsRepository.FindByID(id)
}

func (d *DocsService) FindByUserID(id string) ([]domain.Docs, error) {
	return d.docsRepository.FindByUserID(id)
}

func (s *DocsService) UpdateStatus(id uuid.UUID, status int) error {
	return s.docsRepository.UpdateStatus(id, status)
}

func (s *DocsService) GetAllDocsByStatus(status int) ([]domain.Docs, error) {
	return s.docsRepository.GetAllDocsByStatus(status)
}
