package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/benjaminafoster/chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/benjaminafoster/chirpy/internal/auth"
)

/* Accepts a JSON body with the following shape
	{
  		"email": "user@example.com",
		"password": "example_password"
	}
*/

type UserRequestBody struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
}

/* Returns 201 Created if user is successfully created
	{
		"id": "50746277-23c6-4d85-a890-564c0044c2fb",
		"created_at": "2021-07-07T00:00:00Z",
		"updated_at": "2021-07-07T00:00:00Z",
		"email": "user@example.com"
	}
*/
type User struct {
	Id              uuid.UUID `json:"id"`
	Created_At      time.Time `json:"created_at"`
	Updated_At      time.Time `json:"updated_at"`
	Email           string    `json:"email"`
}

// post one user
func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// Decode request body

	decoder := json.NewDecoder(r.Body)
	reqBody := UserRequestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode user request", err)
		return
	}

	hashed_pwd, err := auth.HashPassword(reqBody.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	// Create user parameters
	params := database.CreateUserParams{
		Email: reqBody.Email,
		HashedPassword: hashed_pwd,
	}

	// Create the user in the DB
	user, err := cfg.DbPtr.CreateUser(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user in users database", err)
		return
	}

	newUser := User{
		Id:            user.ID,
		Created_At:    user.CreatedAt,
		Updated_At:    user.UpdatedAt,
		Email:         user.Email,
	}

	respondWithJSON(w, http.StatusCreated, newUser)
}