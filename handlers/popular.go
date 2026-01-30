package handlers

import (
	"art-gallery/repository"
	"encoding/json"
	"net/http"
)

type PopularHandler struct {
	Repo *repository.GalleryStore
}

func (h *PopularHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		h.handlerGet(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *PopularHandler) handlerGet(w http.ResponseWriter, r *http.Request) {
	popularPaintings := h.Repo.GetPopularPaintings()

	err := json.NewEncoder(w).Encode(popularPaintings)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
