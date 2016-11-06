package resource

import (
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
}

// View shows a page
func (r *PageResource) View(c echo.Context) error {
	p, err := r.fetchPageFromParam(c, "id")
	if err != nil {
		return err
	}

	htmlContent := blackfriday.MarkdownCommon([]byte(p.Content))

	data := struct {
		DudeName string
		Content  template.HTML
	}{
		"dude",
		template.HTML(string(htmlContent)), // convert the string to HTML so that html/templates knows it can be trusted
	}

	return c.Render(http.StatusOK, "page_view", data)
}

// ViewEdit shows the editor for a page
func (r *PageResource) ViewEdit(c echo.Context) error {
	p, err := r.fetchPageFromParam(c, "id")
	if err != nil {
		return err
	}

	data := struct {
		MarkdownContent string
	}{
		p.Content,
	}

	return c.Render(http.StatusOK, "page_edit", data)
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
