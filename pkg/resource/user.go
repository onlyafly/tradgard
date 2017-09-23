package resource

import (
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/onlyafly/tradgard/pkg/service"
)

// UserResource represents a page resource
type UserResource struct {
	SiteName    string
	AuthService *service.AuthService
}

// ViewLogIn shows the Log In page
func (r *UserResource) ViewLogIn(c echo.Context) error {
	loginResult := c.QueryParam("login_result")

	data := struct {
		Context     echo.Context
		SiteName    string
		LoginResult string
	}{
		c,
		r.SiteName,
		loginResult,
	}

	return c.Render(http.StatusOK, "login", data)
}

// ActionLogIn logs in the user
func (r *UserResource) ActionLogIn(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == os.Getenv("ADMIN_USERNAME") && password == os.Getenv("ADMIN_PASSWORD") {
		if err := r.AuthService.StoreUsernameInCookie(c, username); err != nil {
			return err
		}
		return c.Redirect(http.StatusSeeOther, "/")
	}

	return c.Redirect(http.StatusSeeOther, "/user/login?login_result=failed")
}

// ViewLogOut logs out the user
func (r *UserResource) ViewLogOut(c echo.Context) error {
	r.AuthService.ClearCookie(c)
	return c.Redirect(http.StatusSeeOther, "/")
}
