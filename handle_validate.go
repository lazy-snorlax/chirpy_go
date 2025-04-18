package main

import (
	"encoding/json"
	"net/http"
)

func validateRequest(resp http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(resp, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respondWithJSON(resp, http.StatusOK, returnVals{
		Valid: true,
	})
}
