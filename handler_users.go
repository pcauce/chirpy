package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pcauce/chirpy/internal/auth"
	"github.com/pcauce/chirpy/internal/database"
	"net/http"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	Refresh   string    `json:"refresh_token"`
}

func (cfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {
	userData := map[string]string{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userData)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode JSON", err)
		return
	}

	hashedPassword, err := auth.HashPassword(userData["password"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	createdUser, err := cfg.queries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          userData["email"],
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Email:     createdUser.Email,
	})
}

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	loginData := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginData)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode request data", err)
		return
	}

	user, err := cfg.queries.GetUserByEmail(r.Context(), loginData.Email)
	if err != nil || auth.CheckPasswordHash(user.HashedPassword, loginData.Password) != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	newJwtToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT", err)
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
	}
	err = cfg.queries.StoreRefresh(r.Context(), database.StoreRefreshParams{
		Token:     refreshToken,
		UserID:    uuid.NullUUID{user.ID, true},
		ExpiresAt: time.Now().Add(cfg.tokenDuration["refresh"]),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't store refresh token in database", err)
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     newJwtToken,
		Refresh:   refreshToken,
	})
}
