package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const fileRoot = "."

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(fileRoot))))
	mux.HandleFunc("/healthz", handleReadiness)

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
