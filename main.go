package main

import (
	"log"
	"os"

	"github.com/alphagov/pay-open-banking-demo/database"
	"github.com/alphagov/pay-open-banking-demo/http"
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

	http.Start(http.Config{
		DB: db})

	return nil
}

func main() {
	if err := Main(); err != nil {
		log.Fatal(err)
	}
}
