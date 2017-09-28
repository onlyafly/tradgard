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
	recents, err := r.PageService.GetRecentlyUpdatedPageAddressInfos(10)
	if err != nil {
		return err
	}

	orphans, err := r.PageService.GetOrphanedPageAddressInfos(10)
	if err != nil {
		return err
	}

	data := struct {
		Context     echo.Context
		SiteName    string
		HeaderTitle string
		RecentPages []*service.PageAddressInfo
		OrphanPages []*service.PageAddressInfo
	}{
		c,
		r.PageService.SiteName,
		r.PageService.SiteName,
		recents,
		orphans,
	}
	return c.Render(http.StatusOK, "index", data)
}
