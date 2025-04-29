package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"
	"sort"
	"github.com/benjaminafoster/chirpy/internal/database"
	"github.com/google/uuid"
)

/* Accepts a JSON body with the following shape
{
	"body": "Hello, world!",
	"user_id": "123e4567-e89b-12d3-a456-426614174000"
}
*/

type ChirpRequest struct {
	Body   string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

/* If successful, return 201 and chirp that matches the following:
	{
	"id": "94b7e44c-3604-42e3-bef7-ebfcc3efff8f",
	"created_at": "2021-01-01T00:00:00Z",
	"updated_at": "2021-01-01T00:00:00Z",
	"body": "Hello, world!",
	"user_id": "123e4567-e89b-12d3-a456-426614174000"
	}
*/

type Chirp struct {
	ID        uuid.UUID    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID    `json:"user_id"`
}

// post one chirp
func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	
	// Decode the json data
	decoder := json.NewDecoder(r.Body)
	reqBody := ChirpRequest{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	
	// Need to check if the user_id exists
	user_id := reqBody.UserID
	_, err = cfg.DbPtr.GetUserById(context.Background(), user_id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "User does not exist in user database", err)
		return
	}

	body := reqBody.Body
	// Need to check if the chirp is valid
	if len(body) > 140 {
		respondWithError(w, http.StatusInternalServerError, "Chirp is too long", fmt.Errorf("Chirp is too long"))
		return
	}

	newBody := cleanBody(body)

	chirpParams := database.CreateChirpParams {
		Body: newBody,
		UserID: user_id,
	}

	chirpDb, err := cfg.DbPtr.CreateChirp(context.Background(), chirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error adding chirp to database", err)
		return
	}

	chirp := Chirp{
		ID: chirpDb.ID,
		CreatedAt: chirpDb.CreatedAt,
		UpdatedAt: chirpDb.UpdatedAt,
		Body: chirpDb.Body,
		UserID: chirpDb.UserID,
	}

	respondWithJSON(w, http.StatusCreated, chirp)

}

// get all chirps (sorted by created_at)
func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirpsDB, err := cfg.DbPtr.GetChirps(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve all chirps", err)
		return
	}

	chirpsSlice := []Chirp{}

	for _, chirp := range chirpsDB {
		chirpsSlice = append(chirpsSlice, Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		})
	}

	sort.Sort(ByDate{chirpsSlice})
	
	respondWithJSON(w, http.StatusOK, chirpsSlice)
}

// get single chirp
func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpString := r.PathValue("chirpID")
	if chirpString == "" {
		respondWithError(w, http.StatusInternalServerError, "No chirp ID provided", fmt.Errorf("no chirp ID provided"))
		return
	}

	chirpID, err := uuid.Parse(chirpString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing chirp string into UUID", err)
		return
	}

	chirpDb, err := cfg.DbPtr.GetChirpByID(context.Background(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	chirp := Chirp{
		ID: chirpDb.ID,
		CreatedAt: chirpDb.CreatedAt,
		UpdatedAt: chirpDb.UpdatedAt,
		Body: chirpDb.Body,
		UserID: chirpDb.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirp)

}


// Sorting by created_at date for chirps
type Chirps []Chirp

func (c Chirps) Len() int {
	return len(c)
}
func (c Chirps) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type ByDate struct {Chirps}

func (s ByDate) Less(i, j int) bool {
	return s.Chirps[i].CreatedAt.Before(s.Chirps[j].CreatedAt)
}


// Auxiliary function to handle cleaning the body of chirps
func cleanBody(str string) string {
	words := strings.Fields(str)
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	for i, word := range words {
		if slices.Contains(profaneWords, strings.ToLower(word)) {
			words[i] = "****"
		}
	}

	newString := strings.Join(words, " ")
	return newString
}