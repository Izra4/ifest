package repositories

import (
	"IFEST/internals/core/domain"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IUserDocRepository interface {
	Create(userID, docID uuid.UUID) (sql.Result, error)
	FindByUserID(userID uuid.UUID) ([]domain.Docs, error)
	FindByDocID(docID uuid.UUID) ([]domain.User, error)
}

type UserDocRepository struct {
	db *sqlx.DB
}

func NewUserDocRepository(db *sqlx.DB) IUserDocRepository {
	return &UserDocRepository{db: db}
}

func (u UserDocRepository) Create(userID, docID uuid.UUID) (sql.Result, error) {
	query := `INSERT INTO user_doc (user_id, doc_id) VALUES ($1, $2)`

	result, err := u.db.Exec(query, userID, docID)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (u UserDocRepository) FindByUserID(userID uuid.UUID) ([]domain.Docs, error) {
	query := `
        SELECT * FROM user_doc
        INNER JOIN documents ON user_doc.doc_id = documents.doc_id
        WHERE user_doc.user_id = $1
    `

	var documents []domain.Docs
	err := u.db.Select(&documents, query, userID)
	if err != nil {
		return nil, err
	}

	return documents, nil
}

func (u UserDocRepository) FindByDocID(docID uuid.UUID) ([]domain.User, error) {
	query := `
        SELECT * FROM user_doc
        INNER JOIN users ON user_doc.user_id = users.id
        WHERE user_doc.doc_id = $1
    `

	var users []domain.User
	err := u.db.Select(&users, query, docID)
	if err != nil {
		return nil, err
	}
	return users, nil
}
