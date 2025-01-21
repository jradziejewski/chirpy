package main

import (
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/jradziejewski/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func newApiConfig(db *database.Queries) *apiConfig {
	godotenv.Load()
	platform := os.Getenv("PLATFORM")
	cfg := &apiConfig{}
	cfg.fileserverHits.Store(0)
	cfg.db = db
	cfg.platform = platform
	return cfg
}
