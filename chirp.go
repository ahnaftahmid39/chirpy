package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/ahnaftahmid39/chirpy/internal/auth"
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
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	body := &reqBody{}
	err = decoder.Decode(body)

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
		UserID: userId,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error creating chirp, full error: %v\n", err))
		return
	}

	respondWithJSON(w, 201, chirp)
}

func (cfg *apiConfig) handleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	authorIdStr := r.URL.Query().Get("author_id")
	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "asc"
	}

	var chirps []database.Chirp
	var err error
	if authorIdStr != "" {
		authorId, parseErr := uuid.Parse(authorIdStr)
		if parseErr != nil {
			respondWithError(w, 400, fmt.Sprintf("error parsing authorId, full error: %v\n", err))
			return
		}
		chirps, err = cfg.db.GetChirpsByAuthorId(r.Context(), authorId)
	} else {
		chirps, err = cfg.db.GetAllChiprs(r.Context())
	}

	if sort == "desc" {
		slices.SortFunc(chirps, func(a database.Chirp, b database.Chirp) int {
			return b.CreatedAt.Compare(a.CreatedAt)
		})
	}
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
		respondWithError(w, 404, fmt.Sprintf("error getting chirp by id, full error: %v\n", err))
		return
	}

	respondWithJSON(w, 200, chirp)
}

func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {
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
	chirpIdStr := r.PathValue("chirpId")
	chirpId, err := uuid.Parse(chirpIdStr)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	chirp, err := cfg.db.GetChirpById(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, 404, err.Error())
		return
	}

	if chirp.UserID != userId {
		respondWithError(w, 403, "seems like you are trying to access someone else's chirp")
		return
	}

	err = cfg.db.DeleteChripById(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	respondWithJSON(w, 204, nil)
}
