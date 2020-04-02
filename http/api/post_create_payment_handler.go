package api

import (
	"log"
	"net/http"

	"github.com/alphagov/pay-open-banking-demo/database"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type createPaymentRequest struct {
	Reference   string `json:"reference"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	ReturnURL   string `json:"return_url"`
}

func CreatePaymentHandler(db *database.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var request createPaymentRequest
		err := c.Bind(&request)
		if err != nil {
			return err
		}

		paymentId := uuid.New().String()
		charge := database.Charge{
			Reference:   request.Reference,
			Description: request.Description,
			Amount:      request.Amount,
			ReturnUrl:   request.ReturnURL,
			ExternalId:  paymentId,
			Status:      "created",
		}

		err = db.Ping()
		if err != nil {
			log.Print("Failed to ping db")
			return err
		}

		err = db.InsertCharge(charge)
		if err != nil {
			log.Print("Failed to insert charge")
			return err
		}

		response := NewPaymentResponse(charge)
		return c.JSON(http.StatusAccepted, response)
	}
}
