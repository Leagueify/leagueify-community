package web

import (
	"embed"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

var (
	//go:embed all:*
	web   embed.FS
	webFS = echo.MustSubFS(web, "")
)

// setup template rendering
type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func Routes(e *echo.Echo) {
	// template rendering
	e.Renderer = &Template{
		templates: template.Must(template.ParseFS(webFS, "templates/*.html")),
	}
	e.StaticFS("", webFS)
	// route initialization
	API(e)
}
