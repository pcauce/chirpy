package handler

import (
	"github.com/pcauce/chirpy/internal/auth"
	"github.com/pcauce/chirpy/internal/config"
	"github.com/pcauce/chirpy/internal/database"
	"github.com/pcauce/chirpy/internal/sqlc"
	"github.com/pcauce/chirpy/server/respond"
	"net/http"
	"time"
)

func IssueNewAccessToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respond.WithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	userID, err := database.Queries.GetUserFromRefresh(r.Context(), token)
	if err != nil {
		respond.WithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	newAccess, err := auth.MakeJWT(userID.UUID, config.API.JWTSecret, config.API.TokenDuration["access"])
	if err != nil {
		respond.WithError(w, http.StatusInternalServerError, "Couldn't create JWT Access Token", err)
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: newAccess,
	}
	respond.WithJSON(w, http.StatusOK, response)
}

func RevokeAccessToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = database.Queries.RevokeRefresh(r.Context(), sqlc.RevokeRefreshParams{
		Token:     token,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		respond.WithError(w, http.StatusInternalServerError, "Token doesn't exist", err)
	}

	w.WriteHeader(http.StatusNoContent)
}
