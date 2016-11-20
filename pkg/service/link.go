package service

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// LinkService is a page service
type LinkService struct {
	DB *sqlx.DB
}

// LinkModel represents a link in the DB
type LinkModel struct {
	ID         int64     `db:"id"`
	FromPageID int64     `db:"from_page_id"`
	ToPageID   int64     `db:"to_page_id"`
	Created    time.Time `db:"created"`
	Updated    time.Time `db:"updated"`
}

// AddLink adds a link to the DB
func (s *LinkService) AddLink(fromPageID int64, toPageID int64) error {
	lm := &LinkModel{
		FromPageID: fromPageID,
		ToPageID:   toPageID,
	}
	stmt := `INSERT INTO links
	           (from_page_id, to_page_id)
	         VALUES
	           (:from_page_id, :to_page_id)`
	_, err := s.DB.NamedExec(stmt, lm)
	return err
}

// DeleteByFromPageID removes links from the DB
func (s *LinkService) DeleteByFromPageID(fromPageID int64) error {
	stmt := `DELETE FROM links
	         WHERE from_page_id = $1`
	_, err := s.DB.Exec(stmt, fromPageID)
	return err
}
