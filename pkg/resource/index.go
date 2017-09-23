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
	pageInfos, err := r.PageService.GetRecentlyUpdatedPageAddressInfos(10)
	if err != nil {
		return err
	}

	data := struct {
		Context     echo.Context
		SiteName    string
		RecentPages []*service.PageAddressInfo
	}{
		c,
		r.PageService.SiteName,
		pageInfos,
	}
	return c.Render(http.StatusOK, "index", data)
}
