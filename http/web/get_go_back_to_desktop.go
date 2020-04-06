package web

import (
	"log"

	"github.com/alphagov/pay-open-banking-demo/database"
	"github.com/alphagov/pay-open-banking-demo/internal/truelayer"
	"github.com/labstack/echo/v4"
)

func GetGoBackToDesktopHandler(db *database.DB, trueLayer *truelayer.TrueLayer) echo.HandlerFunc {
	return func(c echo.Context) error {
		providerID := c.QueryParam("payment_id")
		log.Print("Getting payment by provider id: " + providerID)
		charge, err := db.GetChargeByProviderId(providerID)
		if err != nil {
			return err
		}

		return PaymentComplete(c, db, trueLayer, charge, true)
	}
}
