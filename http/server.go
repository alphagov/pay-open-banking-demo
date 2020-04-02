package http

import (
	"github.com/alphagov/pay-open-banking-demo/database"
	"github.com/alphagov/pay-open-banking-demo/http/api"
	"github.com/alphagov/pay-open-banking-demo/internal/truelayer"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
)

type Config struct {
	DB *database.DB
}

func Start(config Config) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	truelayerToken := truelayer.GeneratePaymentToken()
	log.Printf("Got token, expires in %d", truelayerToken.ExpiresIn)

	e.POST("/v1/api/payments", api.CreatePaymentHandler(config.DB))

	e.Logger.Fatal(e.Start(":8080"))
}
