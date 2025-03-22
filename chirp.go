package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ahnaftahmid39/chirpy/internal/database"
	"github.com/google/uuid"
)

const LENGTH_THRESHOLD = 140

func replaceProfane(s string) (cleaned string) {
	split := strings.Split(s, " ")
	cleanedSplit := make([]string, 0, len(split))
	for _, word := range split {
		switch strings.ToLower(word) {
		case "kerfuffle", "sharbert", "fornax":
			cleanedSplit = append(cleanedSplit, "****")
		default:
			cleanedSplit = append(cleanedSplit, word)
		}
	}

	return strings.Join(cleanedSplit, " ")
}

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	body := &reqBody{}
	err := decoder.Decode(body)

	if err != nil {
		respondWithError(w, 400, "Unmarshaling failed")
		return
	}

	if len(body.Body) > LENGTH_THRESHOLD {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   replaceProfane(body.Body),
		UserID: body.UserId,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error creating chirp, full error: %v\n", err))
		return
	}

	respondWithJSON(w, 201, chirp)
}

func (cfg *apiConfig) handleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChiprs(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error getting all chirps, full error: %v\n", err))
		return
	}

	respondWithJSON(w, 200, chirps)
}

func (cfg *apiConfig) handleGetChirpById(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpId")
	chirpUUID, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid uuid given, full error: %v\n", err))
		return
	}

	chirp, err := cfg.db.GetChirpById(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error getting chirp by id, full error: %v\n", err))
		return
	}

	respondWithJSON(w, 200, chirp)
}
