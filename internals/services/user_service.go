package services

import (
	"IFEST/helpers"
	"IFEST/internals/core/domain"
	"IFEST/internals/repositories"
	"errors"
	"github.com/google/uuid"
	"time"
)

type IUserService interface {
	Create(user *domain.UserRequest) (domain.User, error)
	GetByEmail(email string) (domain.User, error)
	GetByID(id string) (domain.User, error)
}

type UserService struct {
	userRepository repositories.IUserRepository
}

func NewUserService(userRepository repositories.IUserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (u *UserService) Create(user *domain.UserRequest) (domain.User, error) {
	ID := uuid.New()

	_, err := u.userRepository.GetByEmail(user.Email)
	if err == nil {
		return domain.User{}, errors.New("user already exist")
	}

	hashedPass, err := helpers.HashPassword(user.Password)
	if err != nil {
		return domain.User{}, err
	}

	userDomain := domain.User{
		ID:           ID,
		Name:         user.Name,
		Email:        user.Email,
		Password:     hashedPass,
		IsGoogleAuth: false,
		CreatedAt:    time.Now(),
	}

	newUser, err := u.userRepository.Create(&userDomain)
	if err != nil {
		return domain.User{}, err
	}

	return newUser, nil
}

func (u *UserService) GetByEmail(email string) (domain.User, error) {
	return u.userRepository.GetByEmail(email)
}

func (u *UserService) GetByID(id string) (domain.User, error) {
	return u.userRepository.GetByID(id)
}
