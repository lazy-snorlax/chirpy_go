package main

import (
	"context"
	"encoding/json"
	"main/internal/auth"
	"main/internal/database"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) createChirp(resp http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(resp, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(resp, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
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
		UserID:    userID,
	})
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}
	respondWithJSON(resp, http.StatusCreated, Chirp{
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
