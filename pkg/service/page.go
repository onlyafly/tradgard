package service

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// PageService is a page service
type PageService struct {
	DB *sqlx.DB
}

// PageModel represents a page in the DB
type PageModel struct {
	ID      int64  `db:"id"`
	Name    string `db:"name"`
	Content string `db:"content"`
}

// Get one page
func (s *PageService) Get(key int64) (*PageModel, error) {
	stmt := `SELECT *
           FROM pages
           WHERE id=$1
           LIMIT 1`
	p := &PageModel{}
	err := s.DB.Get(p, stmt, key)
	switch err {
	case nil:
		return p, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

// GetByName gets one page by name
func (s *PageService) GetByName(name string) (*PageModel, error) {
	stmt := `SELECT *
           FROM pages
           WHERE name=$1
           LIMIT 1`
	p := &PageModel{}
	err := s.DB.Get(p, stmt, name)
	switch err {
	case nil:
		return p, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

// Update a page in the DB
func (s *PageService) Update(p *PageModel) error {
	stmt := `UPDATE pages
	         SET
					   content = :content,
						 name = :name
           WHERE id = :id`
	_, err := s.DB.NamedExec(stmt, p)
	return err
}
