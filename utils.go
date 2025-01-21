package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5xx error: %s", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJson(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling json: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Write(dat)
}

func replaceProfane(text string) string {
	words := strings.Fields(text)
	var cleanWords []string

	for _, word := range words {
		lower := strings.ToLower(word)
		if lower == "kerfuffle" {
			word = "****"
		} else if lower == "sharbert" {
			word = "****"
		} else if lower == "fornax" {
			word = "****"
		}

		cleanWords = append(cleanWords, word)
	}

	return strings.Join(cleanWords[:], " ")
}
