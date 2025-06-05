package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

type PolkaWebhook struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	webhook := PolkaWebhook{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&webhook)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode JSON", err)
		return
	}

	userID, err := uuid.Parse(webhook.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse user ID", err)
		return
	}

	switch webhook.Event {
	case "user.upgraded":
		err = cfg.queries.UpgradeChirpyRed(r.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Couldn't upgrade user", err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	default:
		w.WriteHeader(http.StatusNoContent)
		return
	}
}
