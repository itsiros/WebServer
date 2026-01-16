package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type userData struct {
	Id        uuid.UUID `json:"id"`
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
		respondWithError(w, http.StatusUnauthorized, "Wrong email format", fmt.Errorf("Not valid email"))
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), create.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	resp := &userData{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	dat, err := json.Marshal(resp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Marshal failed", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(dat)
}
