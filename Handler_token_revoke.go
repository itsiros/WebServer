package main

import (
	"net/http"

	"github.com/tsironi93/WebServer/internal/auth"
)

func (cfg *apiConf) HandlerTokenRevoke(w http.ResponseWriter, r *http.Request) {
	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldnt get token", err)
		return
	}

	if err := cfg.db.Revoke(r.Context(), bearer); err != nil {
		respondWithError(w, http.StatusUnauthorized, "token not int he database", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
