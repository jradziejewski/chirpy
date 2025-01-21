package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/jradziejewski/chirpy/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		os.Exit(1)
	}

	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {
		os.Exit(1)
	}

	dbQueries := database.New(db)

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("."))
	apiCfg := newApiConfig(dbQueries)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))

	mux.HandleFunc("GET /api/healthz", handlerHealth)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	err = server.ListenAndServe()
	if err != nil {
		os.Exit(1)
	}
}
