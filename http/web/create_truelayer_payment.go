package web

import (
	"github.com/alphagov/pay-open-banking-demo/database"
	"github.com/alphagov/pay-open-banking-demo/internal/truelayer"
)

func CreateTrueLayerPayment(db *database.DB, trueLayer *truelayer.TrueLayer, charge database.Charge, bank string, redirectURL string) (truelayer.PaymentResult, error) {
	request := &truelayer.SinglePaymentRequest{
		Amount:                       charge.Amount,
		Currency:                     "GBP",
		BeneficiaryName:              "GOV.UK Pay Cake Service",
		BeneficiaryReference:         "GOV.UK PAY DEMO",
		BeneficiaryRemitterReference: "GOV.UK PAY DEMO",
		RedirectURL:                  redirectURL,
		RemitterProviderID:           bank,
		DirectBankLink:               true,
	}

	paymentResult := truelayer.PaymentResult{}
	response, err := trueLayer.CreateSinglePayment(request)
	if err != nil {
		return paymentResult, err
	}
	paymentResult = response.PaymentResult[0]

	err = db.UpdateChargeWithProviderID(charge.ExternalID, paymentResult.SimpID, "started")
	if err != nil {
		return paymentResult, err
	}

	return paymentResult, nil
}
