package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"main/internal/auth"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) userChirpyRed(resp http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId uuid.UUID `json:"user_id"`
		}
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(resp, http.StatusUnauthorized, "Unauthorized apiKey", err)
		return
	}

	if apiKey != cfg.polkaSecret {
		respondWithError(resp, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	if params.Event != "user.upgraded" {
		resp.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpgradeUserToChirpyRed(req.Context(), params.Data.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(resp, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		respondWithError(resp, http.StatusInternalServerError, "error updating user", err)
		return
	}
	resp.WriteHeader(http.StatusNoContent)
}
