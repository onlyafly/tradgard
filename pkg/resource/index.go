package resource

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/onlyafly/tradgard/pkg/service"
)

// IndexResource represents a page resource
type IndexResource struct {
	PageService *service.PageService
}

// ViewIndex shows the index page
func (r *IndexResource) ViewIndex(c echo.Context) error {
	pageInfos, err := r.PageService.GetRecentlyUpdatedPageInfos(10)
	if err != nil {
		return err
	}

	data := struct {
		Context     echo.Context
		RecentPages []*service.PageInfo
	}{
		c,
		pageInfos,
	}
	return c.Render(http.StatusOK, "index", data)
}
