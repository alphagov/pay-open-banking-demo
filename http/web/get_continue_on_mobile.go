package web

import (
	"encoding/base64"
	"fmt"
	"net/http"

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

func GetContinueOnMobileHandler(db *database.DB, trueLayer *truelayer.TrueLayer) echo.HandlerFunc {
	return func(c echo.Context) error {
		charge, err := db.GetCharge(c.Param("payment_id"))
		if err != nil {
			return err
		}

		var qr []byte
		qr, err = qrcode.Encode(fmt.Sprintf("/payment/%s/mobile_redirect", charge.ExternalID), qrcode.Medium, 256)
		if err != nil {
			return err
		}
		qrBase46 := base64.StdEncoding.EncodeToString(qr)

		return c.Render(http.StatusOK, "continue_to_payment.html", ContinueOnMobileData{
			Payment:                 NewPaymentData(charge),
			QR:                      qrBase46,
			ContinueOnDesktopAction: fmt.Sprintf("/payment/%s/continue_to_payment", charge.ExternalID)})
	}
}
