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
	fileserverHits  atomic.Int32
	dbQPtr       *database.Queries
}



func main() {
	const filepathRoot = "."
	const port = "8080"
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("error connecting to postgres DB: %s", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{fileserverHits: atomic.Int32{}, dbQPtr: dbQueries}

	fileserverHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileserverHandler))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	


	srv := &http.Server{
		Addr:     ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

}