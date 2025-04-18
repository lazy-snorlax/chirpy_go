package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(resp http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(resp, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(resp http.ResponseWriter, code int, payload interface{}) {
	resp.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		resp.WriteHeader(500)
		return
	}
	resp.WriteHeader(code)
	resp.Write(dat)
}
