package server

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	echoMiddleware "github.com/labstack/echo/middleware"
	"github.com/onlyafly/tradgard/pkg/extension"
	"github.com/onlyafly/tradgard/pkg/middleware"
	"github.com/onlyafly/tradgard/pkg/resource"
	"github.com/onlyafly/tradgard/pkg/service"
)

const cookieName = "tradgard-cookie"

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

	/*
		e.Use(echoMiddleware.BasicAuthWithConfig(echoMiddleware.BasicAuthConfig{
			Skipper: func(c echo.Context) bool {
				if
			},
			Validator: func(username, password string) bool {
				if username == "joe" && password == "secret" {
					return true
				}
				return false
			},
		}))
	*/

	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
		Format: "${method} ${uri} (status=${status})\n",
	}))

	r := &middleware.HTMLTemplateRenderer{
		Templates: template.Must(parseHTMLTemplatesRecursively("etc/views")),
	}
	e.SetRenderer(r)

	// Services

	linkService := &service.LinkService{
		DB: config.Database,
	}
	pageService := &service.PageService{
		DB:          config.Database,
		LinkService: linkService,
	}
	authService := &service.AuthService{
		CookieName:   cookieName,
		SecureCookie: extension.NewSecureCookie(),
	}

	// Service-Dependent Middleware

	e.Use(middleware.CookieBasedAuthentication(authService))

	// Resources

	indexResource := &resource.IndexResource{
		PageService: pageService,
	}
	pageResource := &resource.PageResource{
		PageService: pageService,
		AuthService: authService,
	}
	userResource := &resource.UserResource{
		AuthService: authService,
	}

	// Routes

	e.GET("/", indexResource.ViewIndex)

	e.GET("/:name", pageResource.ViewByName)
	e.GET("/:name/edit", pageResource.ViewEditByName)

	e.POST("/actions/update_page/id/:id", pageResource.ActionUpdateByID)
	e.POST("/actions/create_page", pageResource.ActionCreate)

	e.GET("/user/login", userResource.ViewLogIn)
	e.GET("/user/logout", userResource.ViewLogOut)

	e.POST("/actions/login", userResource.ActionLogIn)

	// Static

	e.Static("/", "static")

	// Start the server

	fmt.Println("[INFO] Server starting on port " + config.Port + "!")

	if err := e.Run(standard.New(":" + config.Port)); err != nil {
		fmt.Fprintln(os.Stderr, "[FATAL] Error starting server", err)
	}
}

func parseHTMLTemplatesRecursively(basePath string) (*template.Template, error) {
	templ := template.New("")
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			_, err = templ.ParseFiles(path)
			if err != nil {
				return err // This is a return from the Walk function
			}
		}

		return err // This is a return from the Walk function
	})

	return templ, err
}
