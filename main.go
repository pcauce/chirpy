package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/pcauce/chirpy/internal/database"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

import _ "github.com/lib/pq"

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	const (
		port = "8080"
	)
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		queries:        dbQueries,
		platform:       os.Getenv("PLATFORM"),
		jwtSecret:      os.Getenv("JWT_SECRET"),
		tokenDuration: map[string]time.Duration{
			"access":  time.Hour,
			"refresh": time.Hour * 24 * 60,
		},
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidate)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUserCreate)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	mux.HandleFunc("POST /api/login", apiCfg.handlerUserLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerIssueNewAccess)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeAccess)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

type apiConfig struct {
	fileserverHits atomic.Int32
	queries        *database.Queries
	platform       string
	jwtSecret      string
	tokenDuration  map[string]time.Duration
}
