package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ahnaftahmid39/chirpy/internal/auth"
	"github.com/ahnaftahmid39/chirpy/internal/database"
	"github.com/google/uuid"
)

type userInfo struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := &reqBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(params)
	if err != nil {
		respondWithError(w, 400, "could not parse request body")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("hashing password failed, %v\n", err))
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("could not create user, %v\n", err))
		return
	}

	respondWithJSON(w, 201, userInfo{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := &struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err := decoder.Decode(params)
	if err != nil {
		respondWithError(w, 400, "could not parse request body")
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("email does not exist, full error: %v\n", err))
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "email and password does not match\n")
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour*1)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintln(err))
	}

	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, fmt.Sprintln(err))
	}

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refresh_token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("error saving refresh token to database, full error: %v\n", err))
	}

	respondWithJSON(w, 200, userInfo{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		IsChirpyRed:  user.IsChirpyRed,
		Token:        token,
		RefreshToken: refresh_token,
	})
}

func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := &struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err = decoder.Decode(params)
	if err != nil {
		respondWithError(w, 400, "could not parse request body")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	updatedUser, err := cfg.db.UpdateUserById(r.Context(), database.UpdateUserByIdParams{
		ID:             userId,
		Email:          params.Email,
		HashedPassword: hashedPassword,
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	respondWithJSON(w, 200, userInfo{
		ID:          updatedUser.ID,
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
		Email:       updatedUser.Email,
		IsChirpyRed: updatedUser.IsChirpyRed,
	})
}
