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

// View shows a page
func (r *PageResource) View(c echo.Context) error {
	p, err := r.fetchPageFromParam(c, "id")
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

// ViewEdit shows the editor for a page
func (r *PageResource) ViewEdit(c echo.Context) error {
	if !r.AuthService.IsAuthenticated(c) {
		return c.String(http.StatusUnauthorized, "not authorized")
	}

	p, err := r.fetchPageFromParam(c, "id")
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

// PostSave shows the editor for a page
func (r *PageResource) PostSave(c echo.Context) error {
	p, err := r.fetchPageFromParam(c, "id")
	if err != nil {
		return err
	}

	p.Name = c.FormValue("page_name")
	p.Content = c.FormValue("page_content")

	if err := r.PageService.Update(p); err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/page/%d", p.ID))
}

func (r *PageResource) fetchPageFromParam(c echo.Context, idParam string) (*service.PageModel, error) {
	idString := c.Param(idParam)
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return nil, err
	}

	p, err := r.PageService.Get(id)
	if err != nil {
		return nil, err
	} else if p == nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "page not found")
	}

	return p, nil
}
