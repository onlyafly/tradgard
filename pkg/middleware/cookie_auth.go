package middleware

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/onlyafly/tradgard/pkg/service"
)

// CookieBasedAuthentication is middleware for adding an authenticated user to the context
func CookieBasedAuthentication(authService *service.AuthService) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := authService.ExtractDataFromCookie(c); err != nil {
				return fmt.Errorf("Problem extracting data from cookie: %s", err)
			}
			return next(c)
		}
	}
}
