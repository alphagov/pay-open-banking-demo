package web

import (
	"log"
	"net/http"

	"github.com/alphagov/pay-open-banking-demo/database"
	"github.com/alphagov/pay-open-banking-demo/internal/truelayer"
	"github.com/labstack/echo/v4"
)

type PaymentCompleteData struct {
	Payment           PaymentData
	ContinueOnDesktop bool
}

func PaymentComplete(c echo.Context, db *database.DB, trueLayer *truelayer.TrueLayer,
	charge database.Charge, continueOnDesktop bool) error {
	response, err := trueLayer.GetSinglePaymentInfo(charge.ProviderID.String)
	if err != nil {
		return err
	}

	paymentResult := response.PaymentResult[0]
	log.Printf("Updating charge %s status to %s", charge.ExternalID, paymentResult.Status)
	db.UpdateChargeStatus(charge.ExternalID, paymentResult.Status)

	data := PaymentCompleteData{
		Payment:           NewPaymentData(charge),
		ContinueOnDesktop: continueOnDesktop}

	if paymentResult.Status == "failed" || paymentResult.Status == "rejected" {
		return c.Render(http.StatusOK, "failed.html", data)
	}
	if paymentResult.Status == "cancelled" {
		return c.Render(http.StatusOK, "cancelled.html", data)
	}
	return c.Render(http.StatusOK, "success.html", data)
}
