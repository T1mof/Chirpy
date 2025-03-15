package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerAllChirps(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("author_id")
	sortParam := r.URL.Query().Get("sort")

	chirps, err := cfg.database.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if sortParam == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}

	if s != "" {
		s, err := uuid.Parse(s)
		if err != nil {
			respondWithError(w, 401, "Error author_id format")
			return
		}
		authroIDChirps := make([]Chirp, 0)
		for i := range chirps {
			if chirps[i].UserID == s {
				authroIDChirps = append(authroIDChirps, Chirp{
					ID:        chirps[i].ID,
					CreatedAt: chirps[i].CreatedAt,
					UpdatedAt: chirps[i].UpdatedAt,
					Body:      chirps[i].Body,
					UserID:    chirps[i].UserID,
				})
			}
		}
		respondWithJSON(w, 200, authroIDChirps)
		return
	}

	allChirps := make([]Chirp, 0)
	for i := range chirps {
		allChirps = append(allChirps, Chirp{
			ID:        chirps[i].ID,
			CreatedAt: chirps[i].CreatedAt,
			UpdatedAt: chirps[i].UpdatedAt,
			Body:      chirps[i].Body,
			UserID:    chirps[i].UserID,
		})
	}
	respondWithJSON(w, 200, allChirps)
}
