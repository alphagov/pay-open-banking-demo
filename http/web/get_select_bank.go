package web

import (
	"fmt"
	"net/http"

	"github.com/alphagov/pay-open-banking-demo/database"
	"github.com/alphagov/pay-open-banking-demo/internal/truelayer"
	"github.com/labstack/echo/v4"
)

type SelectProviderData struct {
	Providers []truelayer.Provider
	Payment   PaymentData
	Action    string
}

func GetSelectBankHander(db *database.DB, trueLayer *truelayer.TrueLayer) echo.HandlerFunc {
	return func(c echo.Context) error {
		providers := trueLayer.GetProviders().Results
		charge, err := db.GetCharge(c.Param("payment_id"))
		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "select_bank.html", SelectProviderData{
			Providers: providers,
			Payment:   NewPaymentData(charge),
			Action:    fmt.Sprintf("/payment/%s/select_bank", charge.ExternalID)})
	}
}
