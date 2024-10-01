package repositories

import (
	"IFEST/internals/core/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IUserDocRepository interface {
	Create(accesReq *domain.AccessReq) (domain.AccessReq, error)
	FindByUserID(userID uuid.UUID) ([]domain.Docs, error)
	FindByDocID(docID uuid.UUID) ([]domain.User, error)
}

type UserDocRepository struct {
	db *sqlx.DB
}

func NewUserDocRepository(db *sqlx.DB) IUserDocRepository {
	return &UserDocRepository{db: db}
}

func (u UserDocRepository) Create(accesReq *domain.AccessReq) (domain.AccessReq, error) {
	query := `
		INSERT INTO user_doc_access(user_id, doc_id) 
		VALUES (:user_id, :doc_id)
		RETURNING user_id, doc_id
	`
	result, err := u.db.NamedQuery(query, accesReq)
	if err != nil {
		return domain.AccessReq{}, err
	}

	if result.Next() {
		err = result.StructScan(accesReq)
	}

	return *accesReq, nil
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
