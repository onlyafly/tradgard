package server

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	echoMiddleware "github.com/labstack/echo/middleware"
	"github.com/onlyafly/tradgard/pkg/middleware"
	"github.com/onlyafly/tradgard/pkg/resource"
	"github.com/onlyafly/tradgard/pkg/service"
	"github.com/russross/blackfriday"
)

// Config is the config for starting the server
type Config struct {
	Port     string
	Database *sqlx.DB
}

// Start the web server
func Start(config Config) {
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
		Format: "${method} ${uri} (status=${status})\n",
	}))

	r := &middleware.HTMLTemplateRenderer{
		Templates: template.Must(template.ParseGlob("etc/views/*.html")),
	}
	e.SetRenderer(r)

	// Services

	pageService := &service.PageService{
		DB: config.Database,
	}

	// Resources

	pageResource := &resource.PageResource{
		PageService: pageService,
	}

	// Routes

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/test", func(c echo.Context) error {
		markdownContent := "**Hi**"
		output := blackfriday.MarkdownCommon([]byte(markdownContent))
		return c.HTML(http.StatusOK, string(output))
	})

	e.GET("/page/:id", pageResource.View)
	e.GET("/page/:id/edit", pageResource.ViewEdit)
	e.POST("/page/:id/save", pageResource.PostSave)

	e.GET("/name/:name", func(c echo.Context) error {
		name := c.Param("name")
		markdownContent := fmt.Sprintf("Hi, **%s**!", name)
		output := blackfriday.MarkdownCommon([]byte(markdownContent))

		data := struct {
			DudeName string
			Content  template.HTML
		}{
			name,
			template.HTML(string(output)), // convert the string to HTML so that html/templates knows it can be trusted
		}

		return c.Render(http.StatusOK, "hello", data)
	})

	e.Static("/", "static")

	fmt.Println("[INFO] Server starting on port " + config.Port + "!")

	if err := e.Run(standard.New(":" + config.Port)); err != nil {
		fmt.Fprintln(os.Stderr, "[FATAL] Error starting server", err)
	}
}
