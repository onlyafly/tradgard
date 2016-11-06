package service

import (
	"fmt"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/labstack/echo"
)

// AuthService assists with authentication and authorization
type AuthService struct {
	CookieName   string
	SecureCookie *securecookie.SecureCookie
}

// StoreUsernameInCookie stores the username in the cookie
func (s *AuthService) StoreUsernameInCookie(c echo.Context, username string) error {
	cookie := new(echo.Cookie)
	cookie.SetName(s.CookieName)

	decodedCookie := map[string]string{
		"username": username,
	}

	encodedCookieString, err := s.encodeCookie(decodedCookie)
	if err != nil {
		return err
	}

	cookie.SetValue(encodedCookieString)
	cookie.SetExpires(time.Now().Add(24 * time.Hour))
	cookie.SetPath("/")

	c.SetCookie(cookie)

	fmt.Println("[DEBUG] Storing cookie:", decodedCookie)

	return nil
}

// ClearCookie erases the cookie
func (s *AuthService) ClearCookie(c echo.Context) error {
	cookie := new(echo.Cookie)
	cookie.SetName(s.CookieName)
	cookie.SetValue("")
	cookie.SetExpires(time.Now())
	cookie.SetPath("/")
	c.SetCookie(cookie)

	fmt.Println("[DEBUG] Clearing cookie")

	return nil
}

// ExtractDataFromCookie puts data from the cookie into the context
func (s *AuthService) ExtractDataFromCookie(c echo.Context) error {
	cookie, err := c.Cookie(s.CookieName)
	if err != nil {
		fmt.Println("[DEBUG] No cookie session")
	} else {
		encodedCookieString := cookie.Value()
		decodedCookie, err := s.decodeCookie(encodedCookieString)
		if err != nil {
			return err
		}

		fmt.Println("[DEBUG] Data found in cookie session:", decodedCookie)
		c.Set("username", decodedCookie["username"])
	}
	return nil
}

// IsAuthenticated returns true if there is an authenticated user in the context
func (s *AuthService) IsAuthenticated(c echo.Context) bool {
	username := c.Get("username")
	if username != nil && username != "" {
		return true
	}
	return false
}

func (s *AuthService) encodeCookie(value map[string]string) (string, error) {
	return s.SecureCookie.Encode(s.CookieName, value)
}

func (s *AuthService) decodeCookie(encodedCookieString string) (map[string]string, error) {
	decodedCookie := make(map[string]string)
	if err := s.SecureCookie.Decode(s.CookieName, encodedCookieString, &decodedCookie); err != nil {
		return nil, err
	}
	return decodedCookie, nil
}
