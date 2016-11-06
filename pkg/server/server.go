package server

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	echoMiddleware "github.com/labstack/echo/middleware"
	"github.com/onlyafly/tradgard/pkg/middleware"
	"github.com/russross/blackfriday"
)

const (
	defaultPort = "5000"
)

// Start the web server
func Start() {
	e := echo.New()

	// SetLogLevel sets the log level for the logger. Default value 5 (OFF). Possible values:
	// 0 (DEBUG)
	// 1 (INFO)
	// 2 (WARN)
	// 3 (ERROR)
	// 4 (FATAL)
	// 5 (OFF)
	e.SetLogLevel(0)
	e.SetDebug(false)
	e.SetHTTPErrorHandler(middleware.CustomHTTPErrorHandler)

	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${uri} (status=${status})\n",
	}))

	r := &middleware.HTMLTemplateRenderer{
		Templates: template.Must(template.ParseGlob("etc/views/*.html")),
	}
	e.SetRenderer(r)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/test", func(c echo.Context) error {
		input := "**Hi**"
		output := blackfriday.MarkdownCommon([]byte(input))
		return c.HTML(http.StatusOK, string(output))
	})

	e.GET("/page/:name", func(c echo.Context) error {
		name := c.Param("name")
		input := fmt.Sprintf("Hi, **%s**!", name)
		output := blackfriday.MarkdownCommon([]byte(input))

		data := struct {
			DudeName string
			Content  template.HTML
		}{
			name,
			template.HTML(string(output)), // convert the string to HTML so that html/templates knows it can be trusted
		}

		return c.Render(http.StatusOK, "hello", data)
	})

	e.Static("/static", "static")

	port := getEnvOr("PORT", defaultPort)
	fmt.Println("Tradgard starting on port " + port + "!")
	e.Run(standard.New(":" + port))
}

func getEnvOr(envVar, fallback string) string {
	if result := os.Getenv(envVar); result != "" {
		return result
	}
	return fallback
}
