package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"os"
	"database/sql"
	"fmt"
	"github.com/benjaminafoster/chirpy/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)


type apiConfig struct {
	FileserverHits  atomic.Int32
	DbPtr       *database.Queries
	Platform    string
}



func main() {
	const filepathRoot = "."
	const port = "8080"
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("Platform env variable must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("error connecting to postgres DB: %s", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{FileserverHits: atomic.Int32{}, DbPtr: dbQueries, Platform: platform}

	fileserverHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileserverHandler))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	
	


	srv := &http.Server{
		Addr:     ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

}