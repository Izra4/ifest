package repositories

import (
	"IFEST/internals/core/domain"
	"github.com/jmoiron/sqlx"
)

type IDocsRepository interface {
	Upload(docs domain.Docs) (domain.Docs, error)
	FindByID(id string) (domain.Docs, error)
	FindByUserID(id string) ([]domain.Docs, error)
}

type DocsRepository struct {
	db *sqlx.DB
}

func NewDocsRepository(db *sqlx.DB) IDocsRepository {
	return &DocsRepository{
		db: db,
	}
}

func (dr *DocsRepository) Upload(docs domain.Docs) (domain.Docs, error) {

	query := `
		INSERT INTO documents (id,name,user_id,type,status,number,created_at,updated_at) 
		VALUES (:id,:name,:user_id,:type,:status,:number,:created_at,:updated_at)
		RETURNING id,user_id,type,status,created_at,updated_at
	`

	data, err := dr.db.NamedQuery(query, &docs)
	if err != nil {
		return domain.Docs{}, err
	}

	if data.Next() {
		err = data.StructScan(&docs)
	}

	return docs, nil
}

func (dr *DocsRepository) FindByID(id string) (domain.Docs, error) {
	query := `
		SELECT * FROM documents WHERE id = $1
	`
	var data domain.Docs
	err := dr.db.Get(&data, query, id)
	return data, err
}

func (dr *DocsRepository) FindByUserID(id string) ([]domain.Docs, error) {
	query := `
		SELECT * FROM documents WHERE user_id = $1
	`
	var data []domain.Docs
	err := dr.db.Select(&data, query, id)
	return data, err
}
