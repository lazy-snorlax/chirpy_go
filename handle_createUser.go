package main

import (
	"context"
	"encoding/json"
	"main/internal/database"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) createUser(resp http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	type returnVals struct {
		Id        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Email     string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Email:     params.Email,
	})
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(resp, http.StatusCreated, returnVals{
		Id:        user.ID.String(),
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
		Email:     params.Email,
	})
}
