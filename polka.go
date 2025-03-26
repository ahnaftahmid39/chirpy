package main

import (
	"encoding/json"
	"net/http"

	"github.com/ahnaftahmid39/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlePolkaWebhookEvent(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, 401, "invalid api key provided")
		return
	}

	decoder := json.NewDecoder(r.Body)
	body := &struct {
		Event string `json:"event"`
		Data  struct {
			UserId string `json:"user_id"`
		} `json:"data"`
	}{}
	err = decoder.Decode(body)
	if err != nil {
		respondWithError(w, 400, "could not parse request body")
		return
	}

	if body.Event == "user.upgraded" {
		userId, err := uuid.Parse(body.Data.UserId)
		if err != nil {
			respondWithError(w, 400, err.Error())
			return
		}
		_, err = cfg.db.MakeUserChirpyRedById(r.Context(), userId)
		if err != nil {
			respondWithError(w, 404, err.Error())
			return
		}
	}

	respondWithJSON(w, 204, nil)
}
