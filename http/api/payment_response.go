package api

import (
	"fmt"
	"os"

	"github.com/alphagov/pay-open-banking-demo/database"
)

type PaymentResponse struct {
	PaymentID   string `json:"payment_id"`
	Reference   string `json:"reference"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	NextURL     string `json:"next_url"`
	Status      string `json:"status"`
}

func NewPaymentResponse(charge database.Charge) *PaymentResponse {
	response := PaymentResponse{
		PaymentID:   charge.ExternalID,
		Reference:   charge.Reference,
		Description: charge.Description,
		Amount:      charge.Amount,
		NextURL:     fmt.Sprintf("%spayment/%s/select_method", os.Getenv("APPLICATION_URL"), charge.ExternalID),
		Status:      charge.Status,
	}
	return &response
}
