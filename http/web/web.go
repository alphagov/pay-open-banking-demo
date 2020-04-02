package web

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func Routes(e *echo.Echo) {
	t := &Template{
		templates: template.Must(template.ParseGlob("http/web/views/*.html")),
	}

	e.Renderer = t

	e.GET("/payment/:payment_id", GetSelectProviders())
}
