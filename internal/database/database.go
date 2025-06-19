package database

import (
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/pcauce/chirpy/internal/sqlc"
	"log"
	"os"
)

var queries *sqlc.Queries

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	queries = sqlc.New(db)
}

func Queries() *sqlc.Queries {
	return queries
}
