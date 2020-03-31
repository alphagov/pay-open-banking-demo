package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type createPaymentRequest struct {
	Reference   string
	Description string
	Amount      int
	ReturnURL   string `json:"return_url"`
}

type createPaymentResponse struct {
	PaymentID   string `json:"payment_id"`
	Reference   string `json:"reference"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	NextURL     string `json:"next_url"`
}

func newCreatePaymentResponse(PaymentID string,
	Reference string,
	Description string,
	Amount int,
	NextURL string) *createPaymentResponse {
	response := createPaymentResponse{PaymentID: PaymentID,
		Reference:   Reference,
		Description: Description,
		Amount:      Amount,
		NextURL:     NextURL}
	return &response
}

func createPaymentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var req createPaymentRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var paymentID = uuid.New().String()
		var response = newCreatePaymentResponse(paymentID,
			req.Reference,
			req.Description,
			req.Amount,
			fmt.Sprintf("https://localhost:8080/secure/payment/%s", paymentID))

		marshalled, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "%+v", string(marshalled))

	default:
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}

func main() {
	http.HandleFunc("/v1/api/payments", createPaymentHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
