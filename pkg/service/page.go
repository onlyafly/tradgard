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
	Content string `db:"content"`
}

// Get retrieves one page
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
