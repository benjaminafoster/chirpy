package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/google/uuid"

	//"github.com/benjaminafoster/chirpy/internal/database"
)

/* Accepts a JSON body with the following shape
	{
  		"email": "user@example.com"
	}
*/

type UserRequestBody struct {
	Email           string `json:"email"`
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
	Id              uuid.UUID `json:"uuid"`
	Created_At      time.Time `json:"created_at"`
	Updated_At      time.Time `json:"updated_at"`
	Email           string    `json:"email"`
}


func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// Decode request body

	decoder := json.NewDecoder(r.Body)
	reqBody := UserRequestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Printf("error decoding request body: %s", err)
		jsonData, err := json.Marshal(ErrorResponse{Error:fmt.Sprintf("Error decoding request json: %s", err)})
		if err != nil {
			log.Printf("Error marshaling error response: %s", err)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(jsonData))
		return
	}

	// Create the user in the DB
	user, err := cfg.DbPtr.CreateUser(r.Context(), reqBody.Email)
	if err != nil {
		log.Printf("error adding user to database: %s", err)
		jsonData, err := json.Marshal(ErrorResponse{Error:fmt.Sprintf("error adding user to database: %s", err)})
		if err != nil {
			log.Printf("Error marshaling error response: %s", err)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(jsonData))
		return
	}

	newUser := User{
		Id:            user.ID,
		Created_At:    user.CreatedAt,
		Updated_At:    user.UpdatedAt,
		Email:         user.Email,
	}

	// Write response headers and encode and response body based on success of user creation
	jsonData, err := json.Marshal(newUser)
	if err != nil {
		log.Printf("Error marshaling user into JSON data: %s", err)
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write([]byte(jsonData))
}