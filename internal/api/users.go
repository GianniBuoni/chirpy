package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/GianniBuoni/chirpy/internal/auth"
	"github.com/GianniBuoni/chirpy/internal/database"
	"github.com/google/uuid"
)

func (a *ApiConfig) HandeUsers(w http.ResponseWriter, r *http.Request) {
	// decode req
	decoder := json.NewDecoder(r.Body)
	user := database.CreateUserParams{}
	err := decoder.Decode(&user)
	if err != nil {
		log.Printf("ERROR: could not unmarshal request, %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, unexpected)
		return
	}
	// handle empty fields
	if user.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Bad request: experting email\n")
		return
	}
	if user.HashedPassword == "" {
		respondWithError(w, http.StatusBadRequest, "Bad request: no password set\n")
		return
	}
	// init response
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.HashedPassword, err = auth.HashPassword(user.HashedPassword)
	if err != nil {
		log.Printf(
			"ERROR: could not create new user with password. Check auth.\n",
		)
		respondWithError(w, http.StatusInternalServerError, unexpected)
		return
	}
	newUser, err := a.Queries.CreateUser(r.Context(), user)
	if err != nil {
		log.Printf("ERROR: could not create user\n")
		respondWithError(w, http.StatusInternalServerError, unexpected)
		return
	}
	data, err := json.Marshal(newUser)
	if err != nil {
		log.Printf("ERROR: could not marshal new user\n")
		respondWithError(w, http.StatusInternalServerError, unexpected)
		return
	}
	respondWithJSON(w, http.StatusCreated, data)
}
