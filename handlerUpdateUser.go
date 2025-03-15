package main

import (
	"Chirpy/internal/auth"
	"Chirpy/internal/database"
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 401, "Something went wrong")
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

	password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 401, "Failed to hash new password")
		return
	}

	err = cfg.database.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             user_id,
		Email:          params.Email,
		HashedPassword: password,
	})
	if err != nil {
		respondWithError(w, 401, "Failed to update user info")
		return
	}

	respondWithJSON(w, 200, map[string]string{"email": params.Email})
}
