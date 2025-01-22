package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jradziejewski/chirpy/internal/database"
)

type ChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	resp := ChirpResponse{}

	parsedChirpID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, 400, "Provided ChirpID could not be parsed", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), parsedChirpID)
	if err != nil {
		respondWithError(w, 404, "Could not retrieve chirp", err)
		return
	}
	resp.ID = chirp.ID
	resp.CreatedAt = chirp.CreatedAt
	resp.UpdatedAt = chirp.UpdatedAt
	resp.Body = chirp.Body
	resp.UserID = chirp.UserID

	respondWithJson(w, 200, resp)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "Could not retrieve chirps", err)
		return
	}
	var chirpResponses []ChirpResponse
	for _, chirp := range chirps {
		chirpResponses = append(chirpResponses, ChirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJson(w, 200, chirpResponses)
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}
	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	resp := returnVals{}
	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Error decoding JSON", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Body too long", nil)
		return
	} else if len(params.Body) == 0 {
		respondWithError(w, 400, "Empty body", nil)
		return
	}

	cleanBody := replaceProfane(params.Body)
	userID, err := uuid.Parse(params.UserID)
	if err != nil {
		respondWithError(w, 400, "Provided UserID could not be parsed", err)
		return
	}

	now := time.Now().UTC()
	chirpParams := database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Body:      cleanBody,
		UserID:    userID,
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, 500, "Error creating chirp", err)
		return
	}

	resp.ID = chirp.ID
	resp.CreatedAt = chirp.CreatedAt
	resp.UpdatedAt = chirp.UpdatedAt
	resp.Body = chirp.Body
	resp.UserID = chirp.UserID

	respondWithJson(w, 201, resp)
}

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
		respondWithError(w, 500, "Error decoding JSON", err)
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
		respondWithError(w, 500, "Error creating user", err)
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
