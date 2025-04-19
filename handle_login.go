package main

import (
	"context"
	"encoding/json"
	"main/internal/auth"
	"net/http"
)

func (cfg *apiConfig) login(resp http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(context.Background(), params.Email)
	if err != nil {
		respondWithError(resp, http.StatusNotFound, "Incorrect email or password", err)
		return
	}

	if err = auth.CheckPasswordHash(user.HashedPassword, params.Password); err != nil {
		respondWithError(resp, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	respondWithJSON(resp, http.StatusOK, User{
		Id:        user.ID.String(),
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
		Email:     params.Email,
	})
}
