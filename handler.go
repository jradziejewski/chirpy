package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	htmlTemplate := `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`
	hits := cfg.fileserverHits.Load()
	responseText := fmt.Sprintf(htmlTemplate, hits)
	w.Write([]byte(responseText))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if cfg.platform != "dev" {
		w.WriteHeader(403)
		w.Write([]byte("Reset only available on dev platform"))
		return
	}
	w.WriteHeader(200)
	cfg.fileserverHits.Store(0)

	// Delete users
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error deleting users"))
		return
	}

	// Delete chirps
	err = cfg.db.DeleteChirps(r.Context())
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error deleting chirps"))
		return
	}

	w.Write([]byte("OK"))
}
