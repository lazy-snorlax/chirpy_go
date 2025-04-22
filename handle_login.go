package main

import (
	"context"
	"encoding/json"
	"main/internal/auth"
	"main/internal/database"
	"net/http"
	"time"
)

func (cfg *apiConfig) login(resp http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, "Couldn't create refresh token", err)
	}
	_, err = cfg.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, "Couldn't save refresh token", err)
		return
	}

	respondWithJSON(resp, http.StatusOK, response{
		User: User{
			Id:          user.ID.String(),
			CreatedAt:   user.CreatedAt.String(),
			UpdatedAt:   user.UpdatedAt.String(),
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}
