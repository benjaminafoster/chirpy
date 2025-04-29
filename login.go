package main

import (
	"context"
	"encoding/json"
	"net/http"
	"github.com/benjaminafoster/chirpy/internal/auth"
	"log"
)

/* accepts a request body with the following shape and saves in UserRequestBody type (found in users.go)
{
	"password": "04234",
	"email": "lane@example.com"
}
*/

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	reqBody := UserRequestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request body", err)
		return		
	}

	// Look up if user exists in database (by email)
	userDb, err := cfg.DbPtr.GetUserByEmail(context.Background(), reqBody.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Email doesn't appear in users database", err)
		return
	}

	// check password against stored hash. reject if not (with 401 Unauthorized), accept if yes
	log.Printf("Checking password against user with email: %s", reqBody.Email)
	req_pwd := reqBody.Password
	stored_pwd := userDb.HashedPassword
	err = auth.CheckPasswordHash(stored_pwd, req_pwd)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Password doesn't match stored hash -- unauthorized", err)
		return
	}
	
	// Return user data (User type in users.go) with 200 OK
	user := User{
		Id: userDb.ID,
		Created_At: userDb.CreatedAt,
		Updated_At: userDb.UpdatedAt,
		Email: userDb.Email,
	}

	respondWithJSON(w, http.StatusOK, user)
}