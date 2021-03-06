package resource

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"

	"github.com/labstack/echo"
	"github.com/onlyafly/tradgard/pkg/service"
)

type pageViewTemplateContext struct {
	PageID        int64
	SiteName      string
	HeaderTitle   string
	PageName      string
	PageContent   template.HTML
	EditPagePath  string
	IncomingLinks []*service.PageAddressInfo
	OutgoingLinks []*service.PageAddressInfo
	Exists        bool
	Context       echo.Context
}

// PageResource represents a page resource
type PageResource struct {
	PageService *service.PageService
	AuthService *service.AuthService
	LinkService *service.LinkService
}

// ViewByName shows a page given its name
func (r *PageResource) ViewByName(c echo.Context) error {
	unescapedName, err := url.QueryUnescape(c.Param("name"))
	if err != nil {
		return nil
	}

	p, err := r.PageService.GetByName(unescapedName)
	if err != nil {
		return err
	} else if p == nil {
		data := pageViewTemplateContext{
			PageID:       -1,
			SiteName:     r.PageService.SiteName,
			HeaderTitle:  unescapedName + " | " + r.PageService.SiteName,
			PageName:     unescapedName,
			PageContent:  "",
			EditPagePath: generateEditPagePath(&service.PageModel{Name: unescapedName}),
			Exists:       false,
			Context:      c,
		}
		return c.Render(http.StatusOK, "page_view", data)
		//return echo.NewHTTPError(http.StatusNotFound, "No custom page with that name found!")
	}

	generatedHTML := r.PageService.GenerateHTML(p)

	incomingLinks, err := r.PageService.GetLinkedPageAddressInfosByToPageID(p.ID)
	if err != nil {
		return err
	}

	outgoingLinks, err := r.PageService.GetLinkedPageAddressInfosByFromPageID(p.ID)
	if err != nil {
		return err
	}

	data := pageViewTemplateContext{
		PageID:        p.ID,
		SiteName:      r.PageService.SiteName,
		HeaderTitle:   p.Name + " | " + r.PageService.SiteName,
		PageName:      p.Name,
		PageContent:   template.HTML(generatedHTML), // convert the string to HTML so that html/templates knows it can be trusted
		EditPagePath:  generateEditPagePath(p),
		IncomingLinks: incomingLinks,
		OutgoingLinks: outgoingLinks,
		Exists:        true,
		Context:       c,
	}
	c.Set("showEditButton", true)
	return c.Render(http.StatusOK, "page_view", data)
}

// ViewEditByName shows the editor for a page
func (r *PageResource) ViewEditByName(c echo.Context) error {
	if !r.AuthService.IsAuthenticated(c) {
		return echo.NewHTTPError(http.StatusUnauthorized, "Need to be logged in to edit pages")
	}

	unescapedName, err := url.QueryUnescape(c.Param("name"))
	if err != nil {
		return nil
	}

	data := struct {
		PageID       int64
		SiteName     string
		HeaderTitle  string
		PageName     string
		PageContent  string
		SavePagePath string
		PagePath     string
		PageExists   bool
		Context      echo.Context
	}{
		Context:     c,
		SiteName:    r.PageService.SiteName,
		HeaderTitle: unescapedName + " | " + r.PageService.SiteName,
	}

	p, err := r.PageService.GetByName(unescapedName)
	if err != nil {
		return err
	}

	if p != nil {
		data.PageID = p.ID
		data.PageName = p.Name
		data.PageContent = p.Content
		data.SavePagePath = generateUpdatePagePath(p)
		data.PagePath = generateViewPagePath(p)
		data.PageExists = true
	} else {
		data.PageID = -1
		data.PageName = unescapedName
		data.PageContent = ""
		data.SavePagePath = generateCreatePagePath()
		data.PagePath = generateViewPagePath(&service.PageModel{Name: unescapedName})
		data.PageExists = false
	}

	return c.Render(http.StatusOK, "page_edit", data)
}

// ActionUpdateByID updates a page
func (r *PageResource) ActionUpdateByID(c echo.Context) error {
	p, err := r.fetchPageFromIDString(c.Param("id"))
	if err != nil {
		return err
	}

	p.Name = c.FormValue("page_name")
	p.Content = c.FormValue("page_content")

	if err := r.PageService.Update(p); err != nil {
		return err
	}
	if err := r.PageService.RegeneratePageLinks(p); err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, generateViewPagePath(p))
}

// ActionCreate creates a page
func (r *PageResource) ActionCreate(c echo.Context) error {
	pCreate := &service.PageModel{
		Name:    c.FormValue("page_name"),
		Content: c.FormValue("page_content"),
	}

	if err := r.PageService.Create(pCreate); err != nil {
		return err
	}
	if err := r.PageService.RegeneratePageLinks(pCreate); err != nil {
		return err
	}

	p, err := r.PageService.GetByName(pCreate.Name)
	if err != nil {
		return err
	} else if p == nil {
		return fmt.Errorf("Problem finding newly created page: '%s'", pCreate.Name)
	}

	return c.Redirect(http.StatusSeeOther, generateViewPagePath(p))
}

func generateViewPagePath(p *service.PageModel) string {
	return fmt.Sprintf("/%s", url.QueryEscape(p.Name))
}

func generateEditPagePath(p *service.PageModel) string {
	return fmt.Sprintf("/%s/edit", url.QueryEscape(p.Name))
}

func generateUpdatePagePath(p *service.PageModel) string {
	return fmt.Sprintf("/actions/update_page/id/%d", p.ID)
}

func generateCreatePagePath() string {
	return fmt.Sprintf("/actions/create_page")
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
