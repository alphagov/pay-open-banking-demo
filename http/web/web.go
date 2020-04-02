package web

import (
	"html/template"
	"io"
	"log"

	"github.com/alphagov/pay-open-banking-demo/database"
	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func Routes(e *echo.Echo, db *database.DB) {
	t := &Template{
		templates: template.Must(template.ParseGlob("http/web/views/*.html")),
	}

	log.Print(t.templates.Name)

	e.Renderer = t

	e.GET("/payment/:payment_id/select_bank", GetSelectProvidersHander(db))
}
