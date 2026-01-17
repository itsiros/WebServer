package main

import "net/http"

func (cfg *apiConf) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting chirps from database", err)
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
