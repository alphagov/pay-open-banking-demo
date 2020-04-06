package web

import (
	"fmt"
	"net/http"

	"github.com/alphagov/pay-open-banking-demo/database"
	"github.com/labstack/echo/v4"
)

type SelectMethodData struct {
	Payment PaymentData
	Action  string
}

func GetSelectMethodHandler(db *database.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		charge, err := db.GetCharge(c.Param("payment_id"))
		if err != nil {
			return err
		}

		data := SelectMethodData{
			Payment:   NewPaymentData(charge),
			Action:    fmt.Sprintf("/payment/%s/select_bank", charge.ExternalID),
		}
		return c.Render(http.StatusOK, "select_method.html", data)
	}
}
