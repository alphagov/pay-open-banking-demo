package web

import (
	"html/template"
	"io"

	"github.com/alphagov/pay-open-banking-demo/database"
	"github.com/alphagov/pay-open-banking-demo/internal/truelayer"
	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func Routes(e *echo.Echo, db *database.DB, trueLayer *truelayer.TrueLayer) {
	t := &Template{
		templates: template.Must(template.ParseGlob("http/web/views/*.html")),
	}

	e.Static("/", "public")
	e.Renderer = t

	e.GET("/payment/:payment_id/select_method", GetSelectMethodHandler(db))
	e.GET("/payment/:payment_id/select_bank", GetSelectBankHander(db, trueLayer))
	e.POST("/payment/:payment_id/select_bank", PostSelectBankHandler(db, trueLayer))
	e.GET("/payment/:payment_id/continue_to_payment", GetContinueOnMobileHandler(db, trueLayer))
	e.GET("/return", GetReturnHandler(db, trueLayer))
}
