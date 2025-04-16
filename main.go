package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const port = "8080"
	const fileRoot = "."

	cfg := apiConfig{}

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(fileRoot)))))
	mux.HandleFunc("/healthz", handleReadiness)
	mux.HandleFunc("/metrics", cfg.handleNumberOfRequests)
	mux.HandleFunc("/reset", cfg.handleReset)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}

func handleReadiness(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(200)
	io.WriteString(resp, "OK")
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleNumberOfRequests(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(200)
	respStr := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	io.WriteString(resp, respStr)
}

func (cfg *apiConfig) handleReset(resp http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	resp.Header().Add("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(200)
	io.WriteString(resp, "Hits Reset")
}
