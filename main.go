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
	apiConfig := apiConfig{database: dbQueries, platform: os.Getenv("PLATFORM")}
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
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
