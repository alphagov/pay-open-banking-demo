package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
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
	var userName = os.Getenv("PAY_OPEN_BANKING_DEMO_USERNAME")
	var password = os.Getenv("PAY_OPEN_BANKING_DEMO_PASSWORD")
	var hostname = os.Getenv("PAY_OPEN_BANKING_DEMO_HOSTNAME")
	var dbName = os.Getenv("PAY_OPEN_BANKING_DEMO_DB_NAME")
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", userName, password, hostname, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	driver, err := postgres.WithInstance(db, new(postgres.Config))
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres",
		driver)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
	http.HandleFunc("/v1/api/payments", createPaymentHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
