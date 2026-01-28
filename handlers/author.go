package handlers

import (
	"art-gallery/models"
	"art-gallery/repository"
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthorHandler struct {
	Repo *repository.GalleryStore
}

func (h *AuthorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(h.Repo.GetAuthors())
	case http.MethodPost:
		h.handlePost(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *AuthorHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	var newA models.Author

	if err := json.NewDecoder(r.Body).Decode(&newA); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.Repo.AddAuthor(newA)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Успешно добавлено")
}
