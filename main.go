package main

import (
	"database/sql"
	"io"
	"log"
	"main/internal/database"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
	polkaSecret    string
}

func main() {
	const port = "8080"
	const fileRoot = "."

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	dbCon, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbCon)

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set in environment file")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set in environment file")
	}
	polkaSecret := os.Getenv("POLKA_KEY")
	if polkaSecret == "" {
		log.Fatal("JWT_SECRET must be set in environment file")
	}

	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		jwtSecret:      jwtSecret,
		polkaSecret:    polkaSecret,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(fileRoot)))))
	mux.HandleFunc("GET /api/healthz", handleReadiness)

	mux.HandleFunc("POST /api/polka/webhooks", cfg.userChirpyRed)

	mux.HandleFunc("POST /api/login", cfg.login)
	mux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevoke)

	mux.HandleFunc("POST /api/users", cfg.createUser)
	mux.HandleFunc("PUT /api/users", cfg.updateUser)

	mux.HandleFunc("POST /api/chirps", cfg.createChirp)
	mux.HandleFunc("GET /api/chirps", cfg.getAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.deleteChirp)

	mux.HandleFunc("POST /admin/reset", cfg.handleReset)
	mux.HandleFunc("GET /admin/metrics/", cfg.handleMetrics)

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
