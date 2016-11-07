package resource

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/onlyafly/tradgard/pkg/service"
	"github.com/russross/blackfriday"
)

// PageResource represents a page resource
type PageResource struct {
	PageService *service.PageService
	AuthService *service.AuthService
}

// ViewByName shows a page given its name
func (r *PageResource) ViewByName(c echo.Context) error {
	p, err := r.fetchPageFromName(c.Param("name"))
	if err != nil {
		return err
	}

	htmlContent := blackfriday.MarkdownCommon([]byte(p.Content))

	data := struct {
		PageID      int64
		PageName    string
		PageContent template.HTML
		Context     echo.Context
	}{
		p.ID,
		p.Name,
		template.HTML(string(htmlContent)), // convert the string to HTML so that html/templates knows it can be trusted
		c,
	}

	return c.Render(http.StatusOK, "page_view", data)
}

// ViewEditByName shows the editor for a page
func (r *PageResource) ViewEditByName(c echo.Context) error {
	if !r.AuthService.IsAuthenticated(c) {
		return c.String(http.StatusUnauthorized, "not authorized")
	}

	p, err := r.fetchPageFromName(c.Param("name"))
	if err != nil {
		return err
	}

	data := struct {
		PageID      int64
		PageName    string
		PageContent string
		Context     echo.Context
	}{
		p.ID,
		p.Name,
		p.Content,
		c,
	}

	return c.Render(http.StatusOK, "page_edit", data)
}

// ViewByID shows a page
func (r *PageResource) ViewByID(c echo.Context) error {
	p, err := r.fetchPageFromIDString(c.Param("id"))
	if err != nil {
		return err
	}

	htmlContent := blackfriday.MarkdownCommon([]byte(p.Content))

	data := struct {
		PageID      int64
		PageName    string
		PageContent template.HTML
		Context     echo.Context
	}{
		p.ID,
		p.Name,
		template.HTML(string(htmlContent)), // convert the string to HTML so that html/templates knows it can be trusted
		c,
	}

	return c.Render(http.StatusOK, "page_view", data)
}

// ViewEditByID shows the editor for a page
func (r *PageResource) ViewEditByID(c echo.Context) error {
	if !r.AuthService.IsAuthenticated(c) {
		return c.String(http.StatusUnauthorized, "not authorized")
	}

	p, err := r.fetchPageFromIDString(c.Param("id"))
	if err != nil {
		return err
	}

	data := struct {
		PageID      int64
		PageName    string
		PageContent string
		Context     echo.Context
	}{
		p.ID,
		p.Name,
		p.Content,
		c,
	}

	return c.Render(http.StatusOK, "page_edit", data)
}

// PostSaveByID shows the editor for a page
func (r *PageResource) PostSaveByID(c echo.Context) error {
	p, err := r.fetchPageFromIDString(c.Param("id"))
	if err != nil {
		return err
	}

	p.Name = c.FormValue("page_name")
	p.Content = c.FormValue("page_content")

	if err := r.PageService.Update(p); err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/page/%s", p.Name))
}

func (r *PageResource) fetchPageFromIDString(idString string) (*service.PageModel, error) {
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Uh oh! That's not a valid ID for a page!")
	}

	p, err := r.PageService.Get(id)
	if err != nil {
		return nil, err
	} else if p == nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Uh oh! There's no page with that ID!")
	}

	return p, nil
}

func (r *PageResource) fetchPageFromName(name string) (*service.PageModel, error) {
	p, err := r.PageService.GetByName(name)
	if err != nil {
		return nil, err
	} else if p == nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Uh oh! There's no page with that name!")
	}

	return p, nil
}
