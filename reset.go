package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)

	if cfg.platform != "dev" {
		respondWithError(w, 403, "can only perform reset in dev mode")
		return
	}

	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("could not delete all users, %v\n", err))
		return
	}

	respondWithJSON(w, 200, map[string]string{"message": "Hits reset to 0 and deleted all users"})
}
