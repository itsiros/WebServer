package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/tsironi93/WebServer/internal/auth"
)

func (cfg *apiConf) HandlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no bearer header", err)
		return
	}

	user, err := auth.ValidateJWT(bearer, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid JWT", err)
		return
	}

	uuidString := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(uuidString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed parsing uuid", err)
		return
	}

	chirp, err := cfg.db.GetSingleChirp(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not in the database", err)
		return
	}

	if chirp.UserID != user {
		respondWithError(w, http.StatusForbidden, "Forbiden", err)
		return
	}

	if err := cfg.db.DeleteSingleChirp(r.Context(), chirp.ID); err != nil {
		respondWithError(w, http.StatusNotFound, "Error deleting chirp", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
