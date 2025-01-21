package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("."))
	apiCfg := newApiConfig()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("GET /api/healthz", handlerHealth)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}
