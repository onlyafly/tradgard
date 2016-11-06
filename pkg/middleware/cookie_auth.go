package middleware

import (
	"github.com/labstack/echo"
	"github.com/onlyafly/tradgard/pkg/service"
)

// CookieBasedAuthentication is middleware for adding an authenticated user to the context
func CookieBasedAuthentication(authService *service.AuthService) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := authService.ExtractDataFromCookie(c); err != nil {
				return err
			}
			return next(c)
		}
	}
}
