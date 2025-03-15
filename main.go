package main

import (
	"Chirpy/internal/database"
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	database       *database.Queries
	platform       string
	secret         string
	polkaKey       string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error opening sql connection")
	}
	dbQueries := database.New(db)
	mux := http.NewServeMux()
	apiConfig := apiConfig{database: dbQueries, platform: os.Getenv("PLATFORM"), secret: os.Getenv("SECRET"), polkaKey: os.Getenv("POLKA_KEY")}
	mux.Handle("/app/", http.StripPrefix("/app/", apiConfig.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("GET /admin/metrics", apiConfig.handlerCountRequest)
	mux.HandleFunc("POST /admin/reset", apiConfig.handlerResetRequest)
	mux.HandleFunc("POST /api/users", apiConfig.handlerNewUser)
	mux.HandleFunc("POST /api/chirps", apiConfig.handlerNewChirps)
	mux.HandleFunc("GET /api/chirps", apiConfig.handlerAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiConfig.handlerGetChirp)
	mux.HandleFunc("POST /api/login", apiConfig.handlerLoginUser)
	mux.HandleFunc("POST /api/refresh", apiConfig.handlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiConfig.handlerRevokeToken)
	mux.HandleFunc("PUT /api/users", apiConfig.handlerUpdateUser)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiConfig.handlerDeleteChirp)
	mux.HandleFunc("POST /api/polka/webhooks", apiConfig.handlerSetChipryRed)
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
