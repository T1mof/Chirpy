package main

import "net/http"

func (cfg *apiConfig) handlerAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.database.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
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
