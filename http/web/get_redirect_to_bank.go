package web

import (
	"net/http"
	"os"

	"github.com/alphagov/pay-open-banking-demo/database"
	"github.com/alphagov/pay-open-banking-demo/internal/truelayer"
	"github.com/labstack/echo/v4"
)

func GetRedirectToBankHandler(db *database.DB, trueLayer *truelayer.TrueLayer) echo.HandlerFunc {
	return func(c echo.Context) error {
		charge, err := db.GetCharge(c.Param("payment_id"))
		if err != nil {
			return err
		}

		transferredDevice := c.QueryParam("transferredDevice")
		var redirectURL string
		if transferredDevice == "true" {
			redirectURL = os.Getenv("APPLICATION_URL") + "/return/back_to_desktop"
		} else {
			redirectURL = os.Getenv("APPLICATION_URL") + "/return"
		}

		paymentResult, err := CreateTrueLayerPayment(db, trueLayer, charge, charge.Bank.String, redirectURL)
		if err != nil {
			return err
		}

		return c.Redirect(http.StatusSeeOther, paymentResult.AuthURI)
	}
}
