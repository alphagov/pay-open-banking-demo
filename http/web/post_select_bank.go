package web

import (
	"fmt"
	"net/http"

	"github.com/alphagov/pay-open-banking-demo/database"
	"github.com/alphagov/pay-open-banking-demo/internal/truelayer"
	"github.com/labstack/echo/v4"
)

func PostSelectBankHandler(db *database.DB, trueLayer *truelayer.TrueLayer) echo.HandlerFunc {
	return func(c echo.Context) error {
		charge, err := db.GetCharge(c.Param("payment_id"))
		if err != nil {
			return err
		}

		// TODO detect whether we're on mobile or desktop if on mobile, create payment
		// with truelayer and redirect to auth_uri immediately

		// redirectURL := os.Getenv("APPLICATION_URL") + "/return"
		// paymentResult, err := CreateTrueLayerPayment(trueLayer, charge, c.FormValue("select-bank"), redirectURL)
		// if err != nil {
		// 	return err
		// }

		// return c.Redirect(http.StatusSeeOther, paymentResult.AuthURI)

		err = db.UpdateChargeBank(charge.ExternalID, c.FormValue("select-bank"))
		if err != nil {
			return err
		}

		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/payment/%s/continue_to_payment", charge.ExternalID))
	}
}
