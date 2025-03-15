package main

import (
	"Chirpy/internal/auth"
	"net/http"
)

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 500, "Failed to get refresh token")
		return
	}

	err = cfg.database.UpdateToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, 500, "Failed to revoke token")
		return
	}

	w.WriteHeader(204)
}
