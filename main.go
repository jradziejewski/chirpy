package main

import (
	"database/sql"
	"fmt"
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

	// Users
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	// Chirps
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)

	// Admin
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	fmt.Println("Chirpy server started!")
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("Chirpy server started!")
		os.Exit(1)
	}
}
