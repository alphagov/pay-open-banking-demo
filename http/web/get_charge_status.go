package web

import (
	"net/http"
	"log"

	"github.com/alphagov/pay-open-banking-demo/database"
	"github.com/labstack/echo/v4"
)

type GetChargeStatusResponse struct {
	Status string `json:"status"`
}

func GetChargeStatusHandler(db *database.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println("Getting charge status")
		charge, err := db.GetCharge(c.Param("payment_id"))
		if err != nil {
			return err
		}

		response := GetChargeStatusResponse{Status: charge.Status}
		return c.JSON(http.StatusOK, response)
	}
}