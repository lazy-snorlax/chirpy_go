package main

import (
	"context"
	"net/http"
)

func (cfg *apiConfig) handleReset(resp http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		resp.WriteHeader(http.StatusForbidden)
		resp.Write([]byte("Reset is only allowed in dev environment."))
		return
	}
	err := cfg.db.DeleteUsers(context.Background())
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, "Couldn't reset users", err)
		return
	}
	cfg.fileserverHits.Store(0)
	resp.Header().Add("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("Hits reset to 0 and database reset to initial state"))
}
