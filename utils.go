package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func writeSomethingWentWrong(w http.ResponseWriter, errResp any, err error) {
	log.Printf("Error marshalling JSON: %s", err)
	if errDat, err := json.Marshal(errResp); err == nil {
		w.WriteHeader(500)
		w.Write(errDat)
	}
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
