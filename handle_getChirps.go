package main

import (
	"context"
	"net/http"
)

type Chirp struct {
	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Body      string `json:"body"`
	UserId    string `json:"user_id"`
}

func (cfg *apiConfig) getAllChirps(resp http.ResponseWriter, req *http.Request) {
	allChirps, err := cfg.db.GetAllChirps(context.Background())
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	var chirps []Chirp
	for _, chirp := range allChirps {
		dat := Chirp{
			Id:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt.String(),
			UpdatedAt: chirp.UpdatedAt.String(),
			Body:      chirp.Body,
			UserId:    chirp.UserID.String(),
		}
		chirps = append(chirps, dat)
	}

	respondWithJSON(resp, http.StatusOK, chirps)
}
