package main

import (
	"net/http"
	"encoding/json"
	"log"
	"fmt"
)

// All chirps must be no more than 140 characters long

/* Accept a json body that looks like below:
	{
		"body": "This is an opinion I need to share with the world"
	}
*/

type RequestBody struct {
	Body     string  `json:"body"`
}

/* If a chirp is valid, send an appropriate HTTP status code (200) and a json body of this shape:
	{
		"valid": true
	}
*/

type ValidResponse struct {
	Body    string    `json:"cleaned_body"`
}

/* If an error occurs, send an appropriate HTTP status code (400) and a json body of this shape:
	{
		"error": "Something went wrong"
	}
*/

type ErrorResponse struct {
	Error    string  `json:"error"`
}

// function to handle /api/validate_chirp POST requests
func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	reqBody := RequestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		jsonData, err := json.Marshal(ErrorResponse{Error:fmt.Sprintf("Error decoding request json: %s", err)})
		if err != nil {
			log.Printf("Error marshaling error response: %s", err)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(jsonData))
	}

	lenBody := len(reqBody.Body)

	if lenBody > 140 {
		jsonData, err := json.Marshal(ErrorResponse{Error:"Chirp is too long"})
		if err != nil {
			log.Printf("Error marshaling error response: %s", err)
			w.Header().Add("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(500)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte(jsonData))
		return
	}

	newBody := cleanBody(reqBody.Body)

	jsonData, err := json.Marshal(ValidResponse{Body: newBody})
	if err != nil {
		log.Printf("Error marshaling error response: %s", err)
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(jsonData))
}

