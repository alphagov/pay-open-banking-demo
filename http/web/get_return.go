package web

import (
	"log"
	"net/http"

	"github.com/alphagov/pay-open-banking-demo/database"
	"github.com/alphagov/pay-open-banking-demo/internal/truelayer"
	"github.com/labstack/echo/v4"
)

type ReturnData struct {
	Payment PaymentData
}

func GetReturnHandler(db *database.DB, truelayerAccessToken string) echo.HandlerFunc {
	return func(c echo.Context) error {
		providerID := c.QueryParam("payment_id")
		log.Print("Getting payment by provider id: " + providerID)
		charge, err := db.GetChargeByProviderId(providerID)
		if err != nil {
			return err
		}

		response, err := truelayer.GetSinglePaymentInfo(providerID, truelayerAccessToken)
		if err != nil {
			return err
		}

		paymentResult := response.PaymentResult[0]
		log.Printf("Updating charge %s status to %s", charge.ExternalID, paymentResult.Status)
		db.UpdateChargeStatus(charge.ExternalID, paymentResult.Status)

		success := paymentResult.Status != "failed" &&
			paymentResult.Status != "rejected" &&
			paymentResult.Status != "cancelled"

		if success {
			payment := PaymentData{
				ServiceName: "Pay your car tax",
				Description: charge.Description,
				Amount:      charge.Amount,
				Reference:   charge.Reference,
			}
			data := ReturnData{
				Payment: payment,
			}
			return c.Render(http.StatusOK, "success.html", data)
		}

		return c.Render(http.StatusOK, "failed.html", nil)
	}
}
