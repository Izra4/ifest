package services

import (
	"IFEST/internals/core/domain"
	"IFEST/internals/repositories"
	"github.com/google/uuid"
)

type IUserDocService interface {
	Create(userID, docID uuid.UUID) (domain.AccessReq, error)
	FindByUserID(userID uuid.UUID) ([]domain.Docs, error)
	FindByDocID(docID uuid.UUID) ([]domain.User, error)
}

type UserDocService struct {
	userDocRepository repositories.IUserDocRepository
}

func NewUserDocService(userDocRepository repositories.IUserDocRepository) *UserDocService {
	return &UserDocService{
		userDocRepository: userDocRepository,
	}
}

func (u *UserDocService) Create(userID, docID uuid.UUID) (domain.AccessReq, error) {

	req := domain.AccessReq{
		DocID:  docID,
		UserID: userID,
	}

	return u.userDocRepository.Create(&req)
}

func (u *UserDocService) FindByUserID(userID uuid.UUID) ([]domain.Docs, error) {
	return u.userDocRepository.FindByUserID(userID)
}

func (u *UserDocService) FindByDocID(docID uuid.UUID) ([]domain.User, error) {
	return u.userDocRepository.FindByDocID(docID)
}
