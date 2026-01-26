package main

import (
	"encoding/json"
	"net/http"

	"github.com/tsironi93/WebServer/internal/auth"
)

func (cfg *apiConf) HandlerLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	log := parameters{}
	if err := decoder.Decode(&log); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	userID, err := cfg.db.GetUserIDByEmail(r.Context(), log.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "email doesnt exist in database", err)
		return
	}

	ok, err := auth.CheckPasswordHash(log.Password, userID.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error dehashing password", err)
		return
	}

	if !ok {
		respondWithError(w, http.StatusUnauthorized, "wrong password", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        userID.ID,
		CreatedAt: userID.CreatedAt,
		UpdatedAt: userID.UpdatedAt,
		Email:     userID.Email,
	})
}
