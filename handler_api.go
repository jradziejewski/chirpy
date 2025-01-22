package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jradziejewski/chirpy/internal/auth"
	"github.com/jradziejewski/chirpy/internal/database"
)

// Chirps

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
	var userID interface{}
	authorID := r.URL.Query().Get("author_id")
	if authorID != "" {
		parsed, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, 500, "Could not parse author_id", err)
			return
		}
		userID = parsed
	}
	chirps, err := cfg.db.GetChirps(r.Context(), userID)
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
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}

	type parameters struct {
		Body string `json:"body"`
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
	err = decoder.Decode(&params)
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

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}

	chirpID := r.PathValue("chirpID")

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
	if chirp.UserID != userID {
		respondWithError(w, 403, "Forbidden", nil)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), parsedChirpID)
	if err != nil {
		respondWithError(w, 500, "Could not delete chirp", err)
		return
	}

	w.WriteHeader(204)
}

// Users

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerCredentialsChange(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := parameters{}
	resp := UserResponse{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Error decoding JSON", err)
		return
	}

	reqToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}

	userID, err := auth.ValidateJWT(reqToken, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "An error occurred while hashing password", err)
		return
	}

	updateParams := database.UpdateEmailAndPasswordParams{
		Email: params.Email,
		HashedPassword: sql.NullString{
			String: hashedPassword,
			Valid:  true,
		},
		ID: userID,
	}

	user, err := cfg.db.UpdateEmailAndPassword(r.Context(), updateParams)
	if err != nil {
		respondWithError(w, 500, "An error occurred while updating credentials", err)
		return
	}
	resp.ID = user.ID
	resp.CreatedAt = user.CreatedAt
	resp.UpdatedAt = user.UpdatedAt
	resp.Email = user.Email
	resp.IsChirpyRed = user.IsChirpyRed.Bool

	respondWithJson(w, 200, resp)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		UserResponse
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	resp := response{}
	params := parameters{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Error decoding JSON", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, "Wrong credentials", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword.String)
	if err != nil {
		respondWithError(w, 401, "Wrong credentials", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, 500, "Error generating JWT", err)
		return
	}

	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, "Error generating refresh token", err)
		return
	}

	now := time.Now().UTC()
	expiresAt := now.Add(time.Duration(60) * 24 * time.Hour)
	refreshTokenParams := database.CreateRefreshTokenParams{
		Token:     refreshTokenString,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: expiresAt,
		UserID:    user.ID,
	}

	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), refreshTokenParams)
	if err != nil {
		respondWithError(w, 500, "Error saving refresh token", err)
		return
	}

	resp.ID = user.ID
	resp.Email = user.Email
	resp.CreatedAt = user.CreatedAt
	resp.UpdatedAt = user.UpdatedAt
	resp.IsChirpyRed = user.IsChirpyRed.Bool
	resp.Token = token
	resp.RefreshToken = refreshToken.Token

	respondWithJson(w, 200, resp)
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	resp := UserResponse{}
	params := parameters{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Error decoding JSON", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Error hashing password", err)
		return
	}

	now := time.Now().UTC()
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Email:     params.Email,
		HashedPassword: sql.NullString{
			String: hashedPassword,
			Valid:  true,
		},
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
	resp.IsChirpyRed = user.IsChirpyRed.Bool

	respondWithJson(w, 201, resp)
}

// Webhooks

func (cfg *apiConfig) HandlerUpdateIsChirpyRed(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if apiKey != cfg.polkaKey || err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 204, "Error decoding JSON", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	updateParams := database.UpdateIsChirpyRedParams{
		IsChirpyRed: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
		ID: params.Data.UserID,
	}

	_, err = cfg.db.UpdateIsChirpyRed(r.Context(), updateParams)
	if err != nil {
		respondWithError(w, 404, "User not found", err)
		return
	}
	w.WriteHeader(204)
}

// Other

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	resp := response{}

	reqToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), reqToken)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, 500, "An error occurred", err)
		return
	}

	resp.Token = accessToken
	respondWithJson(w, 200, resp)
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	reqToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}

	now := time.Now().UTC()
	params := database.RevokeTokenParams{
		RevokedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
		Token: reqToken,
	}

	err = cfg.db.RevokeToken(r.Context(), params)
	if err != nil {
		respondWithError(w, 500, "Could not revoke token", err)
		return
	}

	w.WriteHeader(204)
}

func handlerHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
