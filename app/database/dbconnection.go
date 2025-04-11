package database

import (
	"github.com/jmoiron/sqlx"
	"log"
)

const pgConnectionString = "host=localhost port=5432 user=postgres password=francium dbname=BookMySpot sslmode=disable"

func ConnecttoDB() (*sqlx.DB, error) {

	db, err := sqlx.Open("postgres", pgConnectionString)
	if err != nil {
		log.Fatalf("Error connecting to postgresSQL: %v", err)
		return nil, err
	}
	return db, nil
}
