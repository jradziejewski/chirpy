package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("."))
	apiCfg := newApiConfig()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("/healthz", handlerHealth)
	mux.HandleFunc("/metrics", apiCfg.handlerHits)
	mux.HandleFunc("/reset", apiCfg.handlerReset)
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}
