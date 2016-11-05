package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/russross/blackfriday"
)

const (
	defaultPort = "5000"
)

func Start() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/test", func(c echo.Context) error {
		input := "**Hi**"
		output := blackfriday.MarkdownCommon([]byte(input))
		return c.HTML(http.StatusOK, string(output))
	})

	port := getEnvOr("PORT", defaultPort)
	fmt.Println("Starting on port " + port)
	e.Run(standard.New(":" + port))
}

func getEnvOr(envVar, fallback string) string {
	if result := os.Getenv(envVar); result != "" {
		return result
	}
	return fallback
}
