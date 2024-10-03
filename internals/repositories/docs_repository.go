package repositories

import (
	"IFEST/internals/core/domain"
	"github.com/jmoiron/sqlx"
	"log"
)

type IDocsRepository interface {
	Upload(docs domain.Docs) (domain.Docs, error)
	FindByID(id string) (domain.DocumentAccessInfo, error)
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

func (dr *DocsRepository) FindByID(id string) (domain.DocumentAccessInfo, error) {
	query := `
	SELECT 
		d.id AS document_id,
		d.user_id AS user_id,
		d.name AS document_name,
		d.number AS document_number,
		d.type AS document_type,
		d.status AS document_status,
		COUNT(uda.user_id) AS access_count,
		COALESCE(STRING_AGG(u_access.email, ', '), '') AS access_emails
	FROM 
		documents d
	JOIN 
		users u_owner ON d.user_id = u_owner.id
	LEFT JOIN 
		user_doc_access uda ON d.id = uda.doc_id
	LEFT JOIN 
		users u_access ON uda.user_id = u_access.id 
	WHERE 
		d.id = $1
	GROUP BY 
		d.id, u_owner.email;

	`
	var data domain.DocumentAccessInfo
	err := dr.db.Get(&data, query, id)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (dr *DocsRepository) FindByUserID(id string) ([]domain.Docs, error) {
	log.Println(id)
	query := `
		SELECT documents.*, COUNT(user_doc_access.user_id) AS user_count 
		FROM documents
		LEFT JOIN user_doc_access ON documents.id = user_doc_access.doc_id
		WHERE documents.user_id = $1
		GROUP BY documents.id
	`
	var data []domain.Docs
	err := dr.db.Select(&data, query, id)
	return data, err
}
