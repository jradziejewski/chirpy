package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeSomethingWentWrong(w http.ResponseWriter, errResp any, err error) {
	log.Printf("Error marshalling JSON: %s", err)
	if errDat, err := json.Marshal(errResp); err == nil {
		w.WriteHeader(500)
		w.Write(errDat)
	}
}
