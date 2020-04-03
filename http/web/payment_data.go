package web

import (
	"fmt"

	"github.com/alphagov/pay-open-banking-demo/database"
)

type PaymentData struct {
	ServiceName   string
	Description   string
	AmountInPence string
	Reference     string
	ReturnURL     string
}

func NewPaymentData(charge database.Charge) PaymentData {
	return PaymentData{
		ServiceName:   "Buy a cake",
		Description:   charge.Description,
		AmountInPence: fmt.Sprintf("%.2f", float32(charge.Amount)/100.0),
		Reference:     charge.Reference,
		ReturnURL: charge.ReturnURL,
	}
}
