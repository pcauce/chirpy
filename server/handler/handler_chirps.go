package handler

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pcauce/chirpy/internal/auth"
	"github.com/pcauce/chirpy/internal/config"
	"github.com/pcauce/chirpy/internal/database"
	"github.com/pcauce/chirpy/internal/sqlc"
	"github.com/pcauce/chirpy/server/respond"
	"net/http"
	"sort"
	"time"
)

type Chirp struct {
	ID        uuid.UUID     `json:"id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Body      string        `json:"body"`
	UserID    uuid.NullUUID `json:"user_id"`
}

func CreateChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respond.WithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	userID, err := auth.ValidateJWT(token, config.API.JWTSecret)
	if err != nil {
		respond.WithError(w, http.StatusUnauthorized, "Unauthorized. JWT not valid", err)
		return
	}

	chirpData := map[string]string{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&chirpData)
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, "Couldn't decode JSON", err)
		return
	}

	chirpRecord, err := database.Queries.CreateChirp(r.Context(), sqlc.CreateChirpParams{
		Body:   chirpData["body"],
		UserID: uuid.NullUUID{UUID: userID, Valid: true},
	})
	if err != nil {
		respond.WithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}
	respond.WithJSON(w, http.StatusCreated, Chirp{
		ID:        chirpRecord.ID,
		CreatedAt: chirpRecord.CreatedAt,
		UpdatedAt: chirpRecord.UpdatedAt,
		Body:      chirpRecord.Body,
		UserID:    chirpRecord.UserID,
	})
}

func GetChirps(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Query().Get("author_id") == "" {
	case true:
		GetAllChirps(w, r)
	case false:
		GetChirpsByAuthor(w, r)
	}
}

func GetAllChirps(w http.ResponseWriter, r *http.Request) {
	unformattedChirps, err := database.Queries.GetAllChirps(r.Context())
	if err != nil {
		respond.WithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
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

	sortOrder := r.URL.Query().Get("sort")
	if sortOrder == "" || sortOrder == "asc" {
		respond.WithJSON(w, http.StatusOK, formattedChirps)
		return
	}

	sort.Slice(formattedChirps, func(i, j int) bool {
		return formattedChirps[i].CreatedAt.After(formattedChirps[j].CreatedAt)
	})

	respond.WithJSON(w, http.StatusOK, formattedChirps)
}

func GetChirpsByAuthor(w http.ResponseWriter, r *http.Request) {
	authorID, err := uuid.Parse(r.URL.Query().Get("author_id"))
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, "Couldn't parse author ID", err)
		return
	}

	unformattedChirps, err := database.Queries.GetChirpsByAuthor(r.Context(), uuid.NullUUID{
		UUID:  authorID,
		Valid: true,
	})
	if err != nil {
		respond.WithError(w, http.StatusNotFound, "Couldn't get chirps", err)
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
	respond.WithJSON(w, http.StatusOK, formattedChirps)
}

func GetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirp, err := database.Queries.GetChirpByID(r.Context(), uuid.MustParse(r.PathValue("chirpID")))
	if err != nil {
		respond.WithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	respond.WithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func DeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respond.WithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	userID, err := auth.ValidateJWT(token, config.API.JWTSecret)
	if err != nil {
		respond.WithError(w, http.StatusUnauthorized, "Unauthorized. JWT not valid", err)
		return
	}

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, "Couldn't parse chirp ID", err)
		return
	}

	chirp, err := database.Queries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respond.WithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}
	if chirp.UserID.UUID != userID {
		respond.WithError(w, http.StatusForbidden, "Unauthorized. You can't delete this chirp", err)
	}

	err = database.Queries.DeleteChirp(r.Context(), sqlc.DeleteChirpParams{
		ID: chirpID,
		UserID: uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		},
	})
	if err != nil {
		respond.WithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
