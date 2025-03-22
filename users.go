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

type userWithoutPassword struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
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
	fmt.Println(params)

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("could not create user, %v\n", err))
		return
	}

	respondWithJSON(w, 201, userWithoutPassword{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
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

	respondWithJSON(w, 200, userWithoutPassword{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

}
