package handler

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pcauce/chirpy/internal/auth"
	"github.com/pcauce/chirpy/internal/config"
	"github.com/pcauce/chirpy/internal/database"
	"github.com/pcauce/chirpy/server/respond"
	"net/http"
)

type PolkaWebhook struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func PolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respond.WithError(w, http.StatusUnauthorized, "Bad Request", err)
		return
	}

	if apiKey != config.APIConfig().PolkaKey {
		respond.WithError(w, http.StatusUnauthorized, "Unauthorized", err)
	}

	webhook := PolkaWebhook{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&webhook)
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, "Couldn't decode JSON", err)
		return
	}

	userID, err := uuid.Parse(webhook.Data.UserID)
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, "Couldn't parse user ID", err)
		return
	}

	switch webhook.Event {
	case "user.upgraded":
		err = database.Queries().UpgradeChirpyRed(r.Context(), userID)
		if err != nil {
			respond.WithError(w, http.StatusNotFound, "Couldn't upgrade user", err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	default:
		w.WriteHeader(http.StatusNoContent)
		return
	}
}
