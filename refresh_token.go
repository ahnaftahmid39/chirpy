package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/ahnaftahmid39/chirpy/internal/auth"
	"github.com/ahnaftahmid39/chirpy/internal/database"
)

func (cfg *apiConfig) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	currentRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	tokenRecord, err := cfg.db.GetRefreshTokenByToken(r.Context(), currentRefreshToken)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	if tokenRecord.RevokedAt.Valid {
		respondWithError(w, 401, fmt.Sprintf("refresh token revoked at %v, need to login again", tokenRecord.RevokedAt.Time))
		return
	}

	if tokenRecord.ExpiresAt.Before(time.Now()) {
		respondWithError(w, 401, "refresh token has expired. need login again")
		return
	}

	newAccessToken, err := auth.MakeJWT(tokenRecord.UserID, cfg.jwtSecret, time.Hour*1)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	res := struct {
		Token string `json:"token"`
	}{Token: newAccessToken}
	respondWithJSON(w, 200, res)
}

func (cfg *apiConfig) handleRevokeToken(w http.ResponseWriter, r *http.Request) {
	currentRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 400, err.Error())
	}

	_, err = cfg.db.RevokeRefreshTokenByToken(r.Context(), database.RevokeRefreshTokenByTokenParams{
		Token: currentRefreshToken,
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: time.Now(),
	})

	if err != nil {
		respondWithError(w, 400, err.Error())
	}

	w.WriteHeader(204)
}
