package services

import (
	email2 "IFEST/helpers/email"
	"IFEST/internals/core/domain"
	"IFEST/internals/repositories"
	"fmt"
	"github.com/google/uuid"
	"os"
	"time"
)

type IUserDocService interface {
	Create(userID, docID uuid.UUID, email, name string) (domain.AccessReq, error)
	FindByUserID(userID uuid.UUID) ([]domain.Docs, error)
	FindByDocID(docID uuid.UUID) ([]domain.User, error)
	FindByToken(token string) (domain.AccessReq, error)
	DeleteAccessByToken(token string) error
	DeleteAccessByUserID(userID, docID uuid.UUID) error
	DeleteExpired() error
}

type UserDocService struct {
	userDocRepository repositories.IUserDocRepository
}

func NewUserDocService(userDocRepository repositories.IUserDocRepository) *UserDocService {
	return &UserDocService{
		userDocRepository: userDocRepository,
	}
}

func (u *UserDocService) Create(userID, docID uuid.UUID, email, name string) (domain.AccessReq, error) {
	token := uuid.New().String()
	expiredAt := time.Now().Add(time.Minute * 30)

	req := domain.AccessReq{
		DocID:      docID,
		UserID:     userID,
		Token:      token,
		Expired_at: expiredAt,
	}

	downloadLink := fmt.Sprintf("%s/api/document/download?token=%s", os.Getenv("BASE_URL"), token)
	email2.SendDownloadLink(email, name, downloadLink)
	return u.userDocRepository.Create(&req)
}

func (u *UserDocService) FindByUserID(userID uuid.UUID) ([]domain.Docs, error) {
	return u.userDocRepository.FindByUserID(userID)
}

func (u *UserDocService) FindByDocID(docID uuid.UUID) ([]domain.User, error) {
	return u.userDocRepository.FindByDocID(docID)
}

func (u *UserDocService) FindByToken(token string) (domain.AccessReq, error) {
	return u.userDocRepository.FindByToken(token)
}

func (u *UserDocService) DeleteAccessByToken(token string) error {
	return u.userDocRepository.DeleteAccessByToken(token)
}

func (u *UserDocService) DeleteAccessByUserID(userID, docID uuid.UUID) error {
	return u.userDocRepository.DeleteAccessByUserID(userID, docID)
}

func (u *UserDocService) DeleteExpired() error {
	return u.userDocRepository.DeleteExpired()
}
