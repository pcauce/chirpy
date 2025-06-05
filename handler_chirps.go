package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pcauce/chirpy/internal/auth"
	"github.com/pcauce/chirpy/internal/database"
	"net/http"
	"time"
)

type Chirp struct {
	ID        uuid.UUID     `json:"id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Body      string        `json:"body"`
	UserID    uuid.NullUUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirpCreate(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized. JWT not valid", err)
		return
	}

	chirpData := map[string]string{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&chirpData)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode JSON", err)
		return
	}

	chirpRecord, err := cfg.queries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   chirpData["body"],
		UserID: uuid.NullUUID{UUID: userID, Valid: true},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirpRecord.ID,
		CreatedAt: chirpRecord.CreatedAt,
		UpdatedAt: chirpRecord.UpdatedAt,
		Body:      chirpRecord.Body,
		UserID:    chirpRecord.UserID,
	})
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	unformattedChirps, err := cfg.queries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	var formattedChirps []Chirp
	for _, chirp := range unformattedChirps {
		formattedChirps = append(formattedChirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, formattedChirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirp, err := cfg.queries.GetChirp(r.Context(), uuid.MustParse(r.PathValue("chirpID")))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized. JWT not valid", err)
		return
	}

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse chirp ID", err)
		return
	}

	chirp, err := cfg.queries.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}
	if chirp.UserID.UUID != userID {
		respondWithError(w, http.StatusForbidden, "Unauthorized. You can't delete this chirp", err)
	}

	err = cfg.queries.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID: chirpID,
		UserID: uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
