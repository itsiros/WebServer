package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConf) HandlerChirpsGetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	isDesc := false

	sortOrder := r.URL.Query().Get("sort")
	if sortOrder == "desc" {
		isDesc = true
	}

	var authorID uuid.UUID = uuid.Nil
	if s := r.URL.Query().Get("author_id"); s != "" {
		id, err := uuid.Parse(s)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "couldn't parse author_id", err)
			return
		}
		authorID = id
	}

	chirps, err := cfg.db.GetChirps(ctx, authorID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting chirps", err)
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

	if isDesc {
		sort.Slice(resp, func(i, j int) bool {
			return resp[j].CreatedAt.Before(resp[i].CreatedAt)
		})
	}
	respondWithJSON(w, http.StatusOK, resp)
}
