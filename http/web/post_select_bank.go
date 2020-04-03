package web

import (
	"net/http"

	"github.com/alphagov/pay-open-banking-demo/database"
	"github.com/alphagov/pay-open-banking-demo/internal/truelayer"
	"github.com/labstack/echo/v4"
)

func PostSelectBankHandler(db *database.DB, truelayerAccessToken string) echo.HandlerFunc {
	return func(c echo.Context) error {
		charge, err := db.GetCharge(c.Param("payment_id"))
		if err != nil {
			return err
		}

		request := truelayer.SinglePaymentRequest{
			Amount:                       charge.Amount,
			Currency:                     "GBP",
			BeneficiaryName:              "GOV.UK Pay",
			BeneficiaryReference:         "GOV.UK PAY DEMO",
			BeneficiarySortCode:          "234567",
			BeneficiaryAccountNumber:     "23456789",
			BeneficiaryRemitterReference: "GOV.UK PAY DEMO",
			RedirectURL:                  "https://console.truelayer-sandbox.com/redirect-page",
			RemitterProviderID:           c.FormValue("select-bank"),
			DirectBankLink:               true,
		}

		response, err := truelayer.CreateSinglePayment(request, truelayerAccessToken)
		if err != nil {
			return err
		}
		paymentResult := response.PaymentResult[0]

		err = db.UpdateChargeWithProviderID(charge.ExternalID, paymentResult.SimpID, "started")
		if err != nil {
			return err
		}

		return c.Redirect(http.StatusSeeOther, paymentResult.AuthURI)
	}
}
