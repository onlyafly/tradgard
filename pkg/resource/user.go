package resource

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/onlyafly/tradgard/pkg/service"
)

// UserResource represents a page resource
type UserResource struct {
	AuthService *service.AuthService
}

// ViewLogIn shows the Log In page
func (r *UserResource) ViewLogIn(c echo.Context) error {
	data := struct {
		Context echo.Context
	}{
		c,
	}

	return c.Render(http.StatusOK, "login", data)
}

// PostLogInDo logs in the user
func (r *UserResource) PostLogInDo(c echo.Context) error {
	username := c.FormValue("username")
	if err := r.AuthService.StoreUsernameInCookie(c, username); err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, "/")
}

// GetLogOutDo logs out the user
func (r *UserResource) GetLogOutDo(c echo.Context) error {
	r.AuthService.ClearCookie(c)
	return c.Redirect(http.StatusSeeOther, "/")
}
