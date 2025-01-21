package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jradziejewski/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}
	resp := returnVals{}
	params := parameters{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, err.Error(), err)
		return
	}

	now := time.Now().UTC()
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Email:     params.Email,
	}

	user, err := cfg.db.CreateUser(r.Context(), userParams)
	if err != nil {
		respondWithError(w, 500, err.Error(), err)
		return
	}
	resp.ID = user.ID
	resp.Email = user.Email
	resp.CreatedAt = user.CreatedAt
	resp.UpdatedAt = user.UpdatedAt

	respondWithJson(w, 201, resp)
}

func handlerHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Error       string `json:"error"`
		CleanedBody string `json:"cleaned_body"`
	}
	resp := returnVals{}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	if len(params.Body) > 140 {
		resp.Error = "Chirp is too long"
	} else {
		resp.CleanedBody = replaceProfane(params.Body)
	}

	respondWithJson(w, 200, resp)
}
