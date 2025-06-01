package main

import (
	"net/http"
	"os"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("PLATFORM") != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err := cfg.queries.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't reset database", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
