package middleware

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

// HTMLTemplateRenderer is an echo.Renderer using Go's HTML template engine
type HTMLTemplateRenderer struct {
	Templates *template.Template
}

// Render implements the echo.Renderer interface
func (r *HTMLTemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return r.Templates.ExecuteTemplate(w, name, data)
}
