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

// GetByID returns a page
func (r *PageResource) GetByID(c echo.Context) error {
	idString := c.Param("id")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return err
	}

	p, err := r.PageService.Get(id)
	if err != nil {
		return err
	} else if p == nil {
		return echo.NewHTTPError(http.StatusNotFound, "page not found")
	}

	htmlContent := blackfriday.MarkdownCommon([]byte(p.Content))

	data := struct {
		DudeName string
		Content  template.HTML
	}{
		idString,
		template.HTML(string(htmlContent)), // convert the string to HTML so that html/templates knows it can be trusted
	}

	return c.Render(http.StatusOK, "hello", data)
}
