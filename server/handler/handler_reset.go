package handler

import (
	"github.com/pcauce/chirpy/internal/database"
	"github.com/pcauce/chirpy/server/respond"
	"net/http"
	"os"
)

func ResetDatabase(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("PLATFORM") != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err := database.Queries.DeleteAllUsers(r.Context())
	if err != nil {
		respond.WithError(w, http.StatusInternalServerError, "Couldn't reset sqlc", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
