package main

import (
	"net/http"
)

func (cfg *apiConf) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting chirps from database", err)
		return
	}

	resp := make([]Chirp, len(chirps))
	for i, c := range chirps {
		resp[i] = Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		}
	}

	respondWithJSON(w, http.StatusOK, resp)
}
