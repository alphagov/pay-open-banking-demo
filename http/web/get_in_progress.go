package web

import (
	"fmt"
	"net/http"

	"github.com/alphagov/pay-open-banking-demo/database"
	"github.com/labstack/echo/v4"
)

type InProgressData struct {
	Payment            PaymentData
	GetChargeStatusURL string
	CompleteURL        string
}

func GetInProgressHandler(db *database.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		charge, err := db.GetCharge(c.Param("payment_id"))
		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "in_progress.html", InProgressData{
			Payment:            NewPaymentData(charge),
			GetChargeStatusURL: fmt.Sprintf("/payment/%s/status", charge.ExternalID),
			CompleteURL:        fmt.Sprintf("/payment/%s/complete", charge.ExternalID)})
	}
}
