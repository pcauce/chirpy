package main

import (
	"github.com/pcauce/chirpy/internal/config"
	"github.com/pcauce/chirpy/server/handler"
	"log"
	"net/http"
)

import _ "github.com/lib/pq"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /admin/reset", handler.ResetDatabase)
	mux.HandleFunc("POST /api/users", handler.CreateUser)
	mux.HandleFunc("PUT /api/users", handler.ChangeUserCredentials)
	mux.HandleFunc("POST /api/login", handler.LoginUser)
	mux.HandleFunc("POST /api/refresh", handler.IssueNewAccessToken)
	mux.HandleFunc("POST /api/revoke", handler.RevokeAccessToken)
	mux.HandleFunc("POST /api/chirps", handler.CreateChirp)
	mux.HandleFunc("POST /api/validate_chirp", handler.ValidateChirp)
	mux.HandleFunc("GET /api/chirps", handler.GetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", handler.GetChirpByID)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", handler.DeleteChirp)
	mux.HandleFunc("POST /api/polka/webhooks", handler.PolkaWebhooks)

	server := http.Server{
		Addr:    ":" + config.Port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
