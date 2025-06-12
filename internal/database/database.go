package database

import (
	"database/sql"
	"github.com/pcauce/chirpy/internal/sqlc"
	"log"
	"os"
)

var connection = func() *sql.DB {
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	return db
}()

var Queries = sqlc.New(connection)
