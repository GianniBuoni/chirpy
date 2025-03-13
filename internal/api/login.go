package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GianniBuoni/chirpy/internal/auth"
	"github.com/GianniBuoni/chirpy/internal/database"
	"github.com/google/uuid"
)

type loginRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type loginResponse struct {
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Id           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *ApiConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := loginRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "decoder.Decode", err)
		return
	}
	// handle empty fields
	if params.Email == "" {
		respondWithInfoError(w, r.Pattern, http.StatusBadRequest, "expecting email")
		return
	}
	if params.Password == "" {
		respondWithInfoError(w, r.Pattern, http.StatusBadRequest, "no password set")
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
		respondWithInfoError(w, r.Pattern, http.StatusUnauthorized)
		return
	}

	// init response
	data := loginResponse{
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Id:        user.ID,
		Email:     user.Email,
	}
	data.Token, err = auth.MakeJWT(user.ID, cfg.SignSecret, JWTDuration)
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "MakeJWT", err)
		return
	}
	data.RefreshToken, err = auth.MakeRefreshToken()
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "MakeRefreshToken", err)
		return
	}
	// store RefreshToken in database
	_, err = cfg.Queries.CreateRefreshToken(
		r.Context(),
		database.CreateRefreshTokenParams{
			Token:     data.RefreshToken,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(RefreshDuration),
		},
	)

	res, err := json.Marshal(data)
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "json.Marshal", err)
		return
	}
	respondWithJSON(w, http.StatusOK, res)
}
