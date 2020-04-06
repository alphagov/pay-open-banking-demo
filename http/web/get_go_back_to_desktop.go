package web

import (
	"log"
	"net/http"

	"github.com/alphagov/pay-open-banking-demo/database"
	"github.com/alphagov/pay-open-banking-demo/internal/truelayer"
	"github.com/labstack/echo/v4"
)

type BackToDesktopData struct {
	Payment PaymentData
}

func GetGoBackToDesktopHandler(db *database.DB, trueLayer *truelayer.TrueLayer) echo.HandlerFunc {
	return func(c echo.Context) error {
		providerID := c.QueryParam("payment_id")
		log.Print("Getting payment by provider id: " + providerID)
		charge, err := db.GetChargeByProviderId(providerID)
		if err != nil {
			return err
		}

		response, err := trueLayer.GetSinglePaymentInfo(providerID)
		if err != nil {
			return err
		}

		paymentResult := response.PaymentResult[0]
		log.Printf("Updating charge %s status to %s", charge.ExternalID, paymentResult.Status)
		db.UpdateChargeStatus(charge.ExternalID, paymentResult.Status)

		return c.Render(http.StatusOK, "go_back_to_desktop.html",
			BackToDesktopData{Payment: NewPaymentData(charge)})
	}
}
