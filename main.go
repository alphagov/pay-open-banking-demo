package main

import (
	"fmt"
	"log"
	"os"

	"github.com/stephencdaly/stephens-openbanking-test/database"
	"github.com/stephencdaly/stephens-openbanking-test/http"
)

func Main() error {
	var userName = os.Getenv("DB_USERNAME")
	var password = os.Getenv("DB_PASSWORD")
	var hostname = os.Getenv("DB_HOSTNAME")
	var dbName = os.Getenv("DB_NAME")
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", userName, password, hostname, dbName)
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
