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
	mux.HandleFunc("GET /api/healthz", handleReadiness)

	mux.HandleFunc("POST /admin/reset", cfg.handleReset)
	mux.HandleFunc("GET /admin/metrics/", cfg.handleMetrics)
	// mux.HandleFunc("GET /api/metrics", cfg.handleNumberOfRequests)

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

func (cfg *apiConfig) handleReset(resp http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	resp.Header().Add("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(200)
	io.WriteString(resp, "Hits Reset")
}

func (cfg *apiConfig) handleMetrics(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "text/html; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
	`, cfg.fileserverHits.Load())))
}
