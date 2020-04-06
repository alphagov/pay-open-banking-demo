package web

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"github.com/alphagov/pay-open-banking-demo/database"
	"github.com/alphagov/pay-open-banking-demo/internal/truelayer"
	"github.com/labstack/echo/v4"
	"github.com/skip2/go-qrcode"
)

type ContinueOnMobileData struct {
	Payment                 PaymentData
	QR                      string
	ContinueOnDesktopAction string
}

func GetContinueToPaymentHandler(db *database.DB, trueLayer *truelayer.TrueLayer) echo.HandlerFunc {
	return func(c echo.Context) error {
		charge, err := db.GetCharge(c.Param("payment_id"))
		if err != nil {
			return err
		}

		var qr []byte
		qr, err = qrcode.Encode(fmt.Sprintf("%s/payment/%s/redirect_to_bank?transferredDevice=true", os.Getenv("APPLICATION_URL"), charge.ExternalID), qrcode.Medium, 512)
		if err != nil {
			return err
		}
		qrBase46 := base64.StdEncoding.EncodeToString(qr)

		return c.Render(http.StatusOK, "continue_to_payment.html", ContinueOnMobileData{
			Payment:                 NewPaymentData(charge),
			QR:                      qrBase46,
			ContinueOnDesktopAction: fmt.Sprintf("/payment/%s/redirect_to_bank", charge.ExternalID)})
	}
}
