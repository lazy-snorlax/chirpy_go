package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func validateRequest(resp http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
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
		CleanedBody: checkProfane(params.Body),
	})
}

func checkProfane(body string) string {
	profane := []string{"kerfuffle", "sharbert", "fornax"}
	clean := []string{}
	toClean := strings.Split(body, " ")
	for _, word := range toClean {
		badWord := slices.Contains(profane, strings.Trim(strings.ToLower(word), " "))
		if badWord {
			clean = append(clean, "****")
		} else {
			clean = append(clean, word)
		}
	}
	return strings.Join(clean, " ")
}
