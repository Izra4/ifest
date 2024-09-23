package repositories

import (
	"IFEST/internals/core/domain"
	"github.com/jmoiron/sqlx"
)

type IUserRepository interface {
	Create(user *domain.User) (domain.User, error)
	GetByEmail(email string) (domain.User, error)
}

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) IUserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) Create(user *domain.User) (domain.User, error) {
	query := `
		INSERT INTO users (id,name,email,password,isgoogleauth) 
		VALUES (:id,:name,:email,:password,:isgoogleauth)
		RETURNING id,name,email,password,isgoogleauth
	`

	data, err := ur.db.NamedQuery(query, &user)
	if err != nil {
		return domain.User{}, err
	}

	if data.Next() {
		err = data.StructScan(&user)
	}

	return *user, nil
}

func (ur *UserRepository) GetByEmail(email string) (domain.User, error) {
	var user domain.User
	query := "SELECT * FROM users WHERE email = $1"
	err := ur.db.Get(&user, query, email)
	return user, err
}
