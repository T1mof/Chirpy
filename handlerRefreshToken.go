package main

import (
	"Chirpy/internal/auth"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 500, "Failed to refresh token auth")
		return
	}

	tokenData, err := cfg.database.GetOneRefreshToken(r.Context(), refreshToken)
	if err != nil || tokenData.ExpiresAt.Time.Before(time.Now()) || tokenData.RevokedAt.Valid {
		respondWithError(w, 401, "Failed to get info from database")
		return
	}

	JWTtoken, err := auth.MakeJWT(tokenData.UserID, cfg.secret, time.Duration(3600)*time.Second)
	if err != nil {
		respondWithError(w, 401, "Failed to create JWT token")
		return
	}

	respondWithJSON(w, 200, map[string]string{"token": JWTtoken})
}
