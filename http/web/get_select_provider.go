package web

import (
	"fmt"
	"net/http"

	"github.com/alphagov/pay-open-banking-demo/database"
	"github.com/alphagov/pay-open-banking-demo/internal/truelayer"
	"github.com/labstack/echo/v4"
)

type PaymentData struct {
	ServiceName string
	Description string
	Amount      int
}

type SelectProviderData struct {
	Providers []truelayer.Provider
	Payment   PaymentData
	Action    string
}

func GetSelectProvidersHander(db *database.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		providers := truelayer.GetProviders().Results
		charge, err := db.GetCharge(c.Param("payment_id"))
		if err != nil {
			return err
		}

		payment := PaymentData{
			ServiceName: "Pay your car tax",
			Description: charge.Description,
			Amount:      charge.Amount,
		}
		data := SelectProviderData{
			Providers: providers,
			Payment:   payment,
			Action:    fmt.Sprintf("/v1/payments/%s/select_bank", charge.ExternalId),
		}
		return c.Render(http.StatusOK, "select_provider.html", data)
	}
}
