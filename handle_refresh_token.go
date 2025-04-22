package main

import (
	"main/internal/auth"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerRefresh(resp http.ResponseWriter, req *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(resp, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(resp, http.StatusUnauthorized, "Couldn't get user for refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		respondWithError(resp, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	respondWithJSON(resp, http.StatusOK, response{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRevoke(resp http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(resp, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	_, err = cfg.db.RevokeRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, "Couldn't revoke session", err)
		return
	}

	resp.WriteHeader(http.StatusNoContent)
}
