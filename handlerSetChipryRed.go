package main

import (
	"Chirpy/internal/auth"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerSetChipryRed(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string            `json:"event"`
		Data  map[string]string `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "Missing Authorization")
		return
	}
	if apiKey != cfg.polkaKey {
		respondWithError(w, 401, "Error Authorization")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	userIDStr, ok := params.Data["user_id"]
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Missing user_id")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = cfg.database.SetChirpyRed(r.Context(), userID)
	if err != nil {
		respondWithError(w, 404, "Failed to find user")
		return
	}

	w.WriteHeader(204)
}
