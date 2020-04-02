package api

import (
	"net/http"

	"github.com/alphagov/pay-open-banking-demo/database"
	"github.com/labstack/echo/v4"
)

func GetPaymentHandler(db *database.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		charge, err := db.GetCharge(c.Param("payment_id"))
		if err != nil {
			return err
		}

		response := NewPaymentResponse(charge)
		return c.JSON(http.StatusOK, response)
	}
}