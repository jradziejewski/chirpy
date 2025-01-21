package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handlerHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

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
	w.WriteHeader(200)
	cfg.fileserverHits.Store(0)
	w.Write([]byte("OK"))
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Error string `json:"error"`
		Valid bool   `json:"valid"`
	}
	resp := returnVals{}
	statusCode := 200

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		errResp := returnVals{
			Valid: false,
			Error: "Something went wrong",
		}
		writeSomethingWentWrong(w, errResp, err)
		return
	}

	if len(params.Body) > 140 {
		statusCode = 400
		resp.Valid = false
		resp.Error = "Chirp is too long"
	} else {
		resp.Valid = true
	}

	dat, err := json.Marshal(resp)
	if err != nil {
		errResp := returnVals{
			Valid: false,
			Error: "Something went wrong",
		}
		writeSomethingWentWrong(w, errResp, err)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(dat)
}
