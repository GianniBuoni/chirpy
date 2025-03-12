package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/GianniBuoni/chirpy/internal/auth"
	"github.com/google/uuid"
)

type loginRequest struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type loginResponse struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Id        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (cfg *ApiConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := loginRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("ERROR: could not unmarshal request, %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, unexpected)
		return
	}
	// handle empty fields
	if params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Bad request: experting email\n")
		return
	}
	if params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Bad request: no password set\n")
		return
	}
	user, err := cfg.Queries.GetUser(r.Context(), params.Email)
	if err != nil {
		msg := fmt.Sprintf("Account for %s, not found\n", params.Email)
		respondWithError(w, http.StatusNotFound, msg)
		return
	}
	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	// init response
	data := loginResponse{
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Id:        user.ID,
		Email:     user.Email,
	}
	expiry := tokenDuration
	if params.ExpiresInSeconds > 0 {
		dur := time.Duration(params.ExpiresInSeconds) * time.Second
		expiry = min(tokenDuration, dur)
	}
	data.Token, err = auth.MakeJWT(user.ID, cfg.SignSecret, expiry)
	if err != nil {
		return
	}
	res, err := json.Marshal(data)
	if err != nil {
		log.Printf("ERROR: could not unmarshal request, %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, unexpected)
		return
	}
	respondWithJSON(w, http.StatusOK, res)
}
