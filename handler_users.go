package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {
	email := map[string]string{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode JSON", err)
		return
	}

	dbUser, err := cfg.queries.CreateUser(r.Context(), email["email"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	})
}
