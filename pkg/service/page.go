package service

import (
	"database/sql"
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// PageService is a page service
type PageService struct {
	DB *sqlx.DB
}

// PageModel represents a page in the DB
type PageModel struct {
	ID      int64     `db:"id"`
	Name    string    `db:"name"`
	Content string    `db:"content"`
	Created time.Time `db:"created"`
	Updated time.Time `db:"updated"`
}

// PageInfo represents a page info
type PageInfo struct {
	Name        string
	EscapedName string
}

// GenerateHTML generates the HTML content for the page
func (s *PageService) GenerateHTML(p *PageModel) string {
	transformedContent := transformWikiLinks(p.Content)

	// See blackfriday's Markdown rendering: https://github.com/russross/blackfriday
	unsafeHTMLContent := blackfriday.MarkdownCommon([]byte(transformedContent))

	// See how bluemonday prevents XSS here: https://github.com/microcosm-cc/bluemonday
	safeHTMLContent := bluemonday.UGCPolicy().SanitizeBytes(unsafeHTMLContent)

	return string(safeHTMLContent)
}

// transformWikiLinks finds all wiki links:
//     {like this}
// and turns them into standard Markdown links:
//     [like this](like%20this)
func transformWikiLinks(s string) string {
	re := regexp.MustCompile(`{([^{]+)}`)
	matches := re.FindAllStringSubmatchIndex(s, -1)

	current := 0
	output := ""
	for _, match := range matches {
		start := match[0]
		nameStart := match[2]
		nameEnd := match[3]
		name := s[nameStart:nameEnd]
		link := fmt.Sprintf("[%s](%s)", name, url.QueryEscape(name))
		output = output + s[current:start] + link
		current = match[1]
	}
	output = output + s[current:]
	return output
}

// GetRecentlyUpdatedPageInfos get recently updated page names
func (s *PageService) GetRecentlyUpdatedPageInfos(limit int) ([]*PageInfo, error) {
	stmt := `SELECT name
	         FROM pages
	         ORDER BY updated DESC
	         LIMIT $1`
	pages := []*PageModel{}
	err := s.DB.Select(&pages, stmt, limit)
	if err != nil {
		return nil, err
	}

	infos := make([]*PageInfo, len(pages))
	for i, page := range pages {
		infos[i] = &PageInfo{
			Name:        page.Name,
			EscapedName: url.QueryEscape(page.Name),
		}
	}

	return infos, nil
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
	           name = :name,
			   updated = now()
	         WHERE id = :id`
	_, err := s.DB.NamedExec(stmt, p)
	return err
}

// Create a page in the DB
func (s *PageService) Create(p *PageModel) error {
	stmt := `INSERT INTO pages
	           (name, content)
	         VALUES
	           (:name, :content)`
	_, err := s.DB.NamedExec(stmt, p)
	return err
}
