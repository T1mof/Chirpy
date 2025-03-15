package main

import (
	"Chirpy/internal/auth"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		respondWithError(w, 404, "chirp ID is missing")
		return
	}

	uuidValue, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, 404, "Error parsing UUID:")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Error token auth")
		return
	}

	user_id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	chirp, err := cfg.database.GetOneChirp(r.Context(), uuidValue)
	if err != nil {
		respondWithError(w, 404, "Failed to find chirps")
		return
	}

	if chirp.UserID != user_id {
		respondWithError(w, 403, "Failed to accsess")
		return
	}

	err = cfg.database.DeleteChirps(r.Context(), uuidValue)
	if err != nil {
		respondWithError(w, 401, "Failed to delete")
		return
	}

	w.WriteHeader(204)
}
