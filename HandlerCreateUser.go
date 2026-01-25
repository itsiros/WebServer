package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConf) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type createUser struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	create := createUser{}
	if err := decoder.Decode(&create); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	if !strings.Contains(create.Email, "@") || !strings.Contains(create.Email, ".") {
		respondWithError(w, http.StatusBadRequest, "Wrong email format", fmt.Errorf("Not valid email"))
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), create.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}
