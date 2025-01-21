package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jradziejewski/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
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
	statusCode := 201
	params := parameters{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)
	if err != nil {
		writeSomethingWentWrong(w, nil, err)
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
		writeSomethingWentWrong(w, nil, err)
		return
	}
	resp.ID = user.ID
	resp.Email = user.Email
	resp.CreatedAt = user.CreatedAt
	resp.UpdatedAt = user.UpdatedAt

	dat, err := json.Marshal(resp)
	if err != nil {
		writeSomethingWentWrong(w, nil, err)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(dat)
}

func handlerHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Error       string `json:"error"`
		CleanedBody string `json:"cleaned_body"`
	}
	resp := returnVals{}
	statusCode := 200

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		errResp := returnVals{
			Error: "Something went wrong",
		}
		writeSomethingWentWrong(w, errResp, err)
		return
	}

	if len(params.Body) > 140 {
		statusCode = 400
		resp.Error = "Chirp is too long"
	} else {
		resp.CleanedBody = replaceProfane(params.Body)
	}

	dat, err := json.Marshal(resp)
	if err != nil {
		errResp := returnVals{
			Error: "Something went wrong",
		}
		writeSomethingWentWrong(w, errResp, err)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(dat)
}
