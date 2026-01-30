package handlers

import (
	"art-gallery/repository"
	"encoding/json"
	"net/http"
)

type FactsHandler struct {
	Repo *repository.GalleryStore
}

func (h *FactsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		facts := h.Repo.GetFacts()
		if err := json.NewEncoder(w).Encode(facts); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}
}
