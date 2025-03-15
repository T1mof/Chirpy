package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
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

	chirp, err := cfg.database.GetOneChirp(r.Context(), uuidValue)
	if err != nil {
		respondWithError(w, 404, "Chirp dont find")
		return
	}
	c := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	respondWithJSON(w, 200, c)
}
