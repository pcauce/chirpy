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
	"time"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	Token       string    `json:"token"`
	Refresh     string    `json:"refresh_token"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	userData := map[string]string{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userData)
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, "Couldn't decode JSON", err)
		return
	}

	hashedPassword, err := auth.HashPassword(userData["password"])
	if err != nil {
		respond.WithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	createdUser, err := database.Queries.CreateUser(r.Context(), sqlc.CreateUserParams{
		Email:          userData["email"],
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respond.WithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respond.WithJSON(w, http.StatusCreated, User{
		ID:          createdUser.ID,
		CreatedAt:   createdUser.CreatedAt,
		UpdatedAt:   createdUser.UpdatedAt,
		Email:       createdUser.Email,
		IsChirpyRed: createdUser.IsChirpyRed,
	})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	loginData := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginData)
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, "Couldn't decode request data", err)
		return
	}

	user, err := database.Queries.GetUserByEmail(r.Context(), loginData.Email)
	if err != nil || auth.CheckPasswordHash(user.HashedPassword, loginData.Password) != nil {
		respond.WithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	newJwtToken, err := auth.MakeJWT(user.ID, config.API.JWTSecret, time.Hour)
	if err != nil {
		respond.WithError(w, http.StatusInternalServerError, "Couldn't create JWT", err)
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respond.WithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
	}
	err = database.Queries.StoreRefresh(r.Context(), sqlc.StoreRefreshParams{
		Token:     refreshToken,
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
		ExpiresAt: time.Now().Add(config.API.TokenDuration["refresh"]),
	})
	if err != nil {
		respond.WithError(w, http.StatusInternalServerError, "Couldn't store refresh token in sqlc", err)
	}

	respond.WithJSON(w, http.StatusOK, User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		Token:       newJwtToken,
		Refresh:     refreshToken,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func ChangeUserCredentials(w http.ResponseWriter, r *http.Request) {
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

	credentials := map[string]string{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&credentials)

	rawPassword, ok := credentials["password"]
	if !ok {
		respond.WithError(w, http.StatusBadRequest, "Password missing", err)
		return
	}
	hashPassword, err := auth.HashPassword(rawPassword)
	if err != nil {
		respond.WithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}
	_, err = database.Queries.UpdateUserPassword(r.Context(), sqlc.UpdateUserPasswordParams{
		ID:             userID,
		HashedPassword: hashPassword,
	})
	if err != nil {
		respond.WithError(w, http.StatusInternalServerError, "Couldn't update user password", err)
		return
	}

	email, ok := credentials["email"]
	if !ok {
		respond.WithError(w, http.StatusBadRequest, "Email missing", err)
		return
	}
	user, err := database.Queries.UpdateUserEmail(r.Context(), sqlc.UpdateUserEmailParams{
		ID:    userID,
		Email: email,
	})
	if err != nil {
		respond.WithError(w, http.StatusInternalServerError, "Couldn't update user password", err)
		return
	}

	respond.WithJSON(w, http.StatusOK, User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}
