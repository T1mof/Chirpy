package main

import "net/http"

func (cfg *apiConfig) handlerResetRequest(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Acsess denied")
		return
	}

	cfg.fileserverHits.Store(0)
	err := cfg.database.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, 500, "Error while delete users")
		return
	}

	respondWithJSON(w, 200, map[string]string{"status": "Ok"})
}
