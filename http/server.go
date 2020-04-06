package http

import (
	"github.com/alphagov/pay-open-banking-demo/database"
	"github.com/alphagov/pay-open-banking-demo/http/api"
	"github.com/alphagov/pay-open-banking-demo/http/web"
	"github.com/alphagov/pay-open-banking-demo/internal/truelayer"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	DB        *database.DB
	TrueLayer *truelayer.TrueLayer
}

func Start(config Config) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/v1/api/payments", api.CreatePaymentHandler(config.DB))
	e.GET("/v1/api/payments/:payment_id", api.GetPaymentHandler(config.DB))

	web.Routes(e, config.DB, config.TrueLayer)

	e.Logger.Fatal(e.Start(":8080"))
}
