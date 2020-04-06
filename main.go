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
		AuthURL:      os.Getenv("TRUELAYER_AUTH_URL"),
		PayURL:       os.Getenv("TRUELAYER_PAY_URL"),
		ClientID:     os.Getenv("TRUELAYER_CLIENT_ID"),
		ClientSecret: os.Getenv("TRUELAYER_CLIENT_SECRET")})
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
