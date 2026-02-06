package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConf) HandlerChirpsGetSingle(w http.ResponseWriter, r *http.Request) {
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

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
