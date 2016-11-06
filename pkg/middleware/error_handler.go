package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/color"
)

func printError(err error, c echo.Context) {
	colorable := color.New()

	s := fmt.Sprintf(
		"%s %s (status=%d): [ERROR] %s\n",
		c.Request().Method(),
		c.Request().URL().Path(),
		c.Response().Status(),
		err.Error(),
	)

	coloredText := colorable.Red(s)
	os.Stdout.Write([]byte(coloredText))
}

// CustomHTTPErrorHandler handles errors
func CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	msg := http.StatusText(code)
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	}
	if c.Echo().Debug() {
		msg = err.Error()
	}
	if !c.Response().Committed() {
		if c.Request().Method() == "HEAD" { // Echo Issue #608
			c.NoContent(code)
		} else {
			c.String(code, msg)
		}
	}

	printError(err, c)
	//c.Echo().Logger().Error(err)
}
