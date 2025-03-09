package main

import (
	"Chirpy/internal/database"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
}

func (cfg *apiConfig) handlerNewChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	params.Body = badWordReplace(params.Body)
	chirp, err := cfg.database.CreateChirp(r.Context(), database.CreateChirpParams{Body: params.Body, UserID: params.UserID})
	if err != nil {
		respondWithError(w, 400, "Error while create chirp")
		return
	}

	c := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	respondWithJSON(w, 201, c)
}

func badWordReplace(s string) string {
	a := strings.Split(s, " ")
	b := map[string]bool{"kerfuffle": true, "sharbert": true, "fornax": true}
	for i := range a {
		if b[strings.ToLower(a[i])] {
			a[i] = "****"
		}
	}
	return strings.Join(a, " ")
}
