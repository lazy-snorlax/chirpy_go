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

type Chirp struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    string    `json:"user_id"`
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
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID.String(),
	}
	respondWithJSON(resp, http.StatusOK, dat)
}

func (cfg *apiConfig) getAllChirps(resp http.ResponseWriter, req *http.Request) {
	filter := req.URL.Query().Get("author_id")
	sort := req.URL.Query().Get("sort")

	var allChirps []database.Chirp
	var err error
	if filter == "" {
		allChirps, err = cfg.db.GetAllChirps(context.Background())
		if err != nil {
			respondWithError(resp, http.StatusInternalServerError, "Couldn't get chirps", err)
			return
		}
	} else {
		allChirps, err = cfg.db.GetAllChirpsByUserId(req.Context(), uuid.MustParse(filter))
		if err != nil {
			respondWithError(resp, http.StatusInternalServerError, "Couldn't get chirps", err)
			return
		}
	}

	var chirps []Chirp
	for _, chirp := range allChirps {
		dat := Chirp{
			Id:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID.String(),
		}
		chirps = append(chirps, dat)
	}

	if sort == "desc" {
		slices.SortFunc(chirps, func(a, b Chirp) int {
			return b.CreatedAt.Compare(a.CreatedAt)
		})
	}

	respondWithJSON(resp, http.StatusOK, chirps)
}

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
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
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

func (cfg *apiConfig) deleteChirp(resp http.ResponseWriter, req *http.Request) {
	chirpID := req.PathValue("chirpID")
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(resp, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(resp, http.StatusForbidden, "Couldn't validate JWT", err)
		return
	}

	chirp, err := cfg.db.GetChirp(context.Background(), uuid.MustParse(chirpID))
	if err != nil {
		respondWithError(resp, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	if userID != chirp.UserID {
		respondWithError(resp, http.StatusForbidden, "Unauthorized to delete chirp", err)
		return
	}

	err = cfg.db.DeleteChirpById(req.Context(), uuid.MustParse(chirpID))
	if err != nil {
		respondWithError(resp, http.StatusNotFound, "Couldn't delete chirp", err)
		return
	}

	resp.Header().Add("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(http.StatusNoContent)
}
