package main

import (
	"context"
	"encoding/json"
	"main/internal/database"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) createChirp(resp http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body   string `json:"body"`
		UserId string `json:"user_id"`
	}
	type returnVals struct {
		Id        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Body      string `json:"body"`
		UserId    string `json:"user_id"`
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

	chirp, err := cfg.db.CreateChirp(context.Background(), database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Body:      checkProfane(params.Body),
		UserID:    uuid.MustParse(params.UserId),
	})
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}
	respondWithJSON(resp, http.StatusCreated, returnVals{
		Id:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt.String(),
		UpdatedAt: chirp.UpdatedAt.String(),
		Body:      chirp.Body,
		UserId:    chirp.UserID.String(),
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
