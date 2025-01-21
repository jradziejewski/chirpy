package main

import (
	"net/http"
	"sync/atomic"

	"github.com/jradziejewski/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func newApiConfig(db *database.Queries) *apiConfig {
	cfg := &apiConfig{}
	cfg.fileserverHits.Store(0)
	cfg.db = db
	return cfg
}
