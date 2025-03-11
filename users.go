package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/GianniBuoni/chirpy/internal/database"
	"github.com/google/uuid"
)

func (a *apiConfig) handeUsers(w http.ResponseWriter, r *http.Request) {
	// decode req
	decoder := json.NewDecoder(r.Body)
	user := database.CreateUserParams{}
	err := decoder.Decode(&user)
	if err != nil {
		log.Printf("ERROR: could not unmarshal reqest, %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, unexpected)
	}

	// int response
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	newUser, err := a.queries.CreateUser(context.Background(), user)
	if err != nil {
		log.Printf("ERROR: could not create user\n")
		respondWithError(w, http.StatusInternalServerError, unexpected)
	}

	data, err := json.Marshal(newUser)
	if err != nil {
		log.Printf("ERROR: could not marshal new user\n")
		respondWithError(w, http.StatusInternalServerError, unexpected)
	}
	respondWithJSON(w, http.StatusCreated, data)
}
