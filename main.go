package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/app/", http.StripPrefix("/app", fileServer))
	mux.HandleFunc("/healthz", handler)
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}
