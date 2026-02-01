package main

import (
	"net/http"
	"time"

	"github.com/tsironi93/WebServer/internal/auth"
)

type RefreshResp struct {
	Token string `json:"token"`
}

func (cfg *apiConf) HandlerTokenRefresh(w http.ResponseWriter, r *http.Request) {

	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldnt get token", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), bearer)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token not in the database", err)
		return
	}

	expires := time.Hour
	token, err := auth.MakeJWT(user, cfg.JWTSecret, expires)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, RefreshResp{
		Token: token,
	})
}
