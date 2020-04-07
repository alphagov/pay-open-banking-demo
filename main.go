package main

import (
	"log"
	"os"

	"github.com/alphagov/pay-open-banking-demo/database"
	"github.com/alphagov/pay-open-banking-demo/http"
	"github.com/alphagov/pay-open-banking-demo/internal/truelayer"
)

func Main() error {
	connStr := os.Getenv("DATABASE_URL")
	db, err := database.NewDB(connStr)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	if err := db.Init(); err != nil {
		return err
	}

	trueLayer, err := truelayer.NewTruelayer(&truelayer.Config{
		AuthURL:           getEnv("TRUELAYER_AUTH_URL", "https://auth.truelayer-sandbox.com"),
		PayURL:            getEnv("TRUELAYER_PAY_URL", "https://pay-api.truelayer-sandbox.com"),
		ClientID:          os.Getenv("TRUELAYER_CLIENT_ID"),
		ClientSecret:      os.Getenv("TRUELAYER_CLIENT_SECRET"),
		BankAccountNumber: getEnv("BANK_ACCOUNT_NO", "23456789"),
		BankSortCode:      getEnv("BANK_SORT_CODE", "234567")})

	if err != nil {
		return err
	}

	http.Start(http.Config{
		DB:        db,
		TrueLayer: trueLayer})

	return nil
}

func main() {
	if err := Main(); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
