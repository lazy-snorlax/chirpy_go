package main

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type Chirp struct {
	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Body      string `json:"body"`
	UserId    string `json:"user_id"`
}

func (cfg *apiConfig) getChirp(resp http.ResponseWriter, req *http.Request) {
	chirpID := req.PathValue("chirpID")
	chirp, err := cfg.db.GetChirp(context.Background(), uuid.MustParse(chirpID))
	if err != nil {
		respondWithError(resp, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}
	dat := Chirp{
		Id:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt.String(),
		UpdatedAt: chirp.UpdatedAt.String(),
		Body:      chirp.Body,
		UserId:    chirp.UserID.String(),
	}
	respondWithJSON(resp, http.StatusOK, dat)
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
