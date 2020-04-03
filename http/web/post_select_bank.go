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
			BeneficiarySortCode:          "234567",
			BeneficiaryAccountNumber:     "23456789",
			BeneficiaryRemitterReference: "GOV.UK PAY DEMO",
			RedirectURL:                  "https://www.payments.service.gov.uk/",
		}

		response, err := truelayer.CreateSinglePayment(request, truelayerAccessToken)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, response)
		// return c.Redirect(http.StatusSeeOther, truelayerRes.PaymentResult[0].RedirectURI)
	}
}
