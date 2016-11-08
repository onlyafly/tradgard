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
		Templates: template.Must(template.ParseGlob("etc/views/*.html")),
	}
	e.SetRenderer(r)

	// Services

	pageService := &service.PageService{
		DB: config.Database,
	}
	authService := &service.AuthService{
		CookieName:   cookieName,
		SecureCookie: extension.NewSecureCookie(),
	}

	// Service-Dependent Middleware

	e.Use(middleware.CookieBasedAuthentication(authService))

	// Resources

	pageResource := &resource.PageResource{
		PageService: pageService,
		AuthService: authService,
	}
	userResource := &resource.UserResource{
		AuthService: authService,
	}

	// Routes

	e.GET("/", func(c echo.Context) error {
		data := struct {
			Context echo.Context
		}{
			c,
		}
		return c.Render(http.StatusOK, "home", data)
	})

	e.GET("/page/:name", pageResource.ViewByName)
	e.GET("/page/:name/edit", pageResource.ViewEditByName)

	e.GET("/page/id/:id", pageResource.ViewByID)
	e.GET("/page/id/:id/edit", pageResource.ViewEditByID)
	e.POST("/page/id/:id/actions/update", pageResource.ActionUpdateByID)
	e.POST("/page/actions/create", pageResource.ActionCreate)

	e.GET("/login", userResource.ViewLogIn)
	e.POST("/login/do", userResource.PostLogInDo)
	e.GET("/logout", userResource.GetLogOutDo)

	// Static

	e.Static("/", "static")

	// Start the server

	fmt.Println("[INFO] Server starting on port " + config.Port + "!")

	if err := e.Run(standard.New(":" + config.Port)); err != nil {
		fmt.Fprintln(os.Stderr, "[FATAL] Error starting server", err)
	}
}
