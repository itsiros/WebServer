package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConf) HandlerSingleChirp(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		ID uuid.UUID `json:"id"`
	}

	uuidString := r.PathValue("/api/chirps/")
	uuid, err := uuid.Parse(uuidString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed parsing uuid", err)
		return
	}

	chirp, err := cfg.db.GetSingleChirp(r.Context(), uuid)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not in the database", err)
		return
	}

	dat, err := json.Marshal(chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Marshal failed", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(dat)
}
