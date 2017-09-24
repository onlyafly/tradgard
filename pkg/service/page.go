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
	DB          *sqlx.DB
	LinkService *LinkService
	SiteName    string
}

// PageModel represents a page in the DB
type PageModel struct {
	ID      int64     `db:"id"`
	Name    string    `db:"name"`
	Content string    `db:"content"`
	Created time.Time `db:"created"`
	Updated time.Time `db:"updated"`
}

// PageAddressInfo represents a page info
type PageAddressInfo struct {
	Name        string
	EscapedName string
}

// GenerateHTML generates the HTML content for the page
func (s *PageService) GenerateHTML(p *PageModel) string {
	transformedContent := transformWikiLinks(p.Content)

	// See blackfriday's Markdown rendering: https://github.com/russross/blackfriday

	unsafeHTMLContent := blackfriday.MarkdownCommon(
		[]byte(transformedContent),
	)
	/* TODO: after upgrading to 2.0.0
	unsafeHTMLContent := blackfriday.Run(
		[]byte(transformedContent),
		blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.HardLineBreak|blackfriday.AutoHeadingIDs|blackfriday.Autolink),
	)
	*/

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

// GetLinkedPageAddressInfosByToPageID gets from DB
func (s *PageService) GetLinkedPageAddressInfosByToPageID(toPageID int64) ([]*PageAddressInfo, error) {
	stmt := `SELECT name
	         FROM pages
	         WHERE id IN (
				 SELECT from_page_id
				 FROM links
				 WHERE to_page_id = $1
			 )`
	pages := []*PageModel{}
	err := s.DB.Select(&pages, stmt, toPageID)
	if err != nil {
		return nil, err
	}

	return generatePageAddressInfosFromPages(pages), nil
}

// GetLinkedPageAddressInfosByFromPageID gets from DB
func (s *PageService) GetLinkedPageAddressInfosByFromPageID(fromPageID int64) ([]*PageAddressInfo, error) {
	stmt := `SELECT name
	         FROM pages
	         WHERE id IN (
				 SELECT to_page_id
				 FROM links
				 WHERE from_page_id = $1
			 )`
	pages := []*PageModel{}
	err := s.DB.Select(&pages, stmt, fromPageID)
	if err != nil {
		return nil, err
	}

	return generatePageAddressInfosFromPages(pages), nil
}

// GetRecentlyUpdatedPageAddressInfos get recently updated page names
func (s *PageService) GetRecentlyUpdatedPageAddressInfos(limit int) ([]*PageAddressInfo, error) {
	stmt := `SELECT name
	         FROM pages
	         ORDER BY updated DESC
	         LIMIT $1`
	pages := []*PageModel{}
	err := s.DB.Select(&pages, stmt, limit)
	if err != nil {
		return nil, err
	}

	return generatePageAddressInfosFromPages(pages), nil
}

func generatePageAddressInfosFromPages(pages []*PageModel) []*PageAddressInfo {
	infos := make([]*PageAddressInfo, len(pages))
	for i, page := range pages {
		infos[i] = &PageAddressInfo{
			Name:        page.Name,
			EscapedName: url.QueryEscape(page.Name),
		}
	}
	return infos
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

// RegeneratePageLinks regenerates all links for this page
func (s *PageService) RegeneratePageLinks(p *PageModel) error {
	// Remove all existing links from this page
	if err := s.LinkService.DeleteByFromPageID(p.ID); err != nil {
		return err
	}

	// Create the new links
	linkNames := findLinkNamesInContent(p.Content)

	for _, name := range linkNames {
		pother, err := s.GetByName(name)
		if err != nil {
			return err
		}

		if pother == nil {
			// Page does not yet exist, so we should create it so that we can link to it
			pCreate := &PageModel{
				Name:    name,
				Content: "",
			}
			if err = s.Create(pCreate); err != nil {
				return err
			}
			pother, err = s.GetByName(name)
			if err != nil {
				return err
			}
		}

		if err := s.LinkService.AddLink(p.ID, pother.ID); err != nil {
			return err
		}
	}

	return nil
}

// findLinkNamesInContent finds all wiki links {like this} and returns the names of the links
func findLinkNamesInContent(content string) []string {
	re := regexp.MustCompile(`{([^{]+)}`)
	matches := re.FindAllStringSubmatchIndex(content, -1)

	links := []string{}
	for _, match := range matches {
		nameStart := match[2]
		nameEnd := match[3]
		name := content[nameStart:nameEnd]
		links = append(links, name)
	}
	return links
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
