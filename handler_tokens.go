package main

import (
	"github.com/pcauce/chirpy/internal/auth"
	"github.com/pcauce/chirpy/internal/database"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerIssueNewAccess(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	userID, err := cfg.queries.GetUserFromRefresh(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	newAccess, err := auth.MakeJWT(userID.UUID, cfg.jwtSecret, cfg.tokenDuration["access"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT Access Token", err)
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: newAccess,
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerRevokeAccess(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = cfg.queries.RevokeRefresh(r.Context(), database.RevokeRefreshParams{
		Token:     token,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Token doesn't exist", err)
	}

	w.WriteHeader(http.StatusNoContent)
}
