package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tsironi93/WebServer/internal/auth"
	"github.com/tsironi93/WebServer/internal/database"
)

type LoginResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *apiConf) HandlerLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	login := parameters{}
	if err := decoder.Decode(&login); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	userID, err := cfg.db.GetUserIDByEmail(r.Context(), login.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid email or password", err)
		return
	}

	ok, err := auth.CheckPasswordHash(login.Password, userID.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error dehashing password", err)
		return
	}

	if !ok {
		respondWithError(w, http.StatusUnauthorized, "invalid email or password", nil)
		return
	}

	expires := time.Hour
	token, err := auth.MakeJWT(userID.ID, cfg.JWTSecret, expires)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create refresh token", err)
		return
	}

	dbRefresh, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: userID.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt create refresh token", err)
		log.Println("CreateRefreshToken error:", err)
		return
	}

	respondWithJSON(w, http.StatusOK, LoginResponse{
		ID:           userID.ID,
		CreatedAt:    userID.CreatedAt,
		UpdatedAt:    userID.UpdatedAt,
		Email:        userID.Email,
		Token:        token,
		RefreshToken: dbRefresh.Token,
	})
}
