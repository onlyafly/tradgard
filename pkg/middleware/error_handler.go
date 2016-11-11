package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/color"
)

func respondError(c echo.Context, statusCode int, extraMessage string) error {
	data := struct {
		Context       echo.Context
		StatusCode    int
		StatusMessage string
		ExtraMessage  string
	}{
		c,
		statusCode,
		http.StatusText(statusCode),
		extraMessage,
	}

	return c.Render(statusCode, "error", data)
}

func printError(err error, c echo.Context) {
	colorable := color.New()

	s := fmt.Sprintf(
		"[ERROR] %s %s (%d): %s\n",
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
	statusCode := http.StatusInternalServerError
	extraMessage := ""

	if he, ok := err.(*echo.HTTPError); ok {
		statusCode = he.Code
		extraMessage = he.Message
	}

	if c.Echo().Debug() {
		extraMessage = err.Error()
	}

	if !c.Response().Committed() {
		if c.Request().Method() == "HEAD" { // Echo Issue #608
			c.NoContent(statusCode)
		} else {
			if errDuringResponse := respondError(c, statusCode, extraMessage); errDuringResponse != nil {
				c.String(http.StatusInternalServerError, "How meta! Error rendering error response.")
				printError(errDuringResponse, c)
			}
		}
	}

	printError(err, c)
	//c.Echo().Logger().Error(err)
}
