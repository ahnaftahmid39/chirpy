package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

const lenThresHold = 140

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type reqBody struct {
		Body string `json:"body"`
	}

	type errResBody struct {
		Error string `json:"error"`
	}

	type resBody struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	body := &reqBody{}
	err := decoder.Decode(body)

	if err != nil {
		respondWithError(w, 400, "Unmarshaling failed")
		return
	}

	if len(body.Body) > lenThresHold {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	res := &resBody{
		CleanedBody: replaceProfane(body.Body),
	}

	respondWithJSON(w, 200, res)
}

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
