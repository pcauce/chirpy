package database

import (
	"database/sql"
	"github.com/pcauce/chirpy/internal/sqlc"
	"log"
	"os"
)

var Queries *sqlc.Queries

func Init() {
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	Queries = sqlc.New(db)
}
