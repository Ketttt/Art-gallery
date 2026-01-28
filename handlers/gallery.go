package handlers

import (
	"art-gallery/models"
	"art-gallery/repository"
	"encoding/json"
	"fmt"
	"net/http"
)

type GalleryHandler struct {
	Repo *repository.GalleryStore
}

func (h *GalleryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		h.handlerGet(w, r)
	case http.MethodPost:
		h.handlePost(w, r)
	case http.MethodDelete:
		h.handleDelete(w, r)
	case http.MethodPatch:
		h.handlePatch(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *GalleryHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	var newP models.Painting

	if err := json.NewDecoder(r.Body).Decode(&newP); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.Repo.AddPainting(newP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Успешно добавлено")
}

func (h *GalleryHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if h.Repo.DeletePainting(id) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Удалено")
	} else {
		http.Error(w, "Не найдено", http.StatusNotFound)
	}
}

func (h *GalleryHandler) handlePatch(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if h.Repo.ToggleFavorite(id) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Статус изменен")
	} else {
		http.Error(w, "Не найдено", http.StatusNotFound)
	}
}

func (h *GalleryHandler) handlerGet(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author")
	searchQuery := r.URL.Query().Get("search")

	var result interface{}
	if authorID != "" {
		result = h.Repo.GetPaintingsByAuthor(authorID)
	} else if searchQuery != "" {
		result = h.Repo.SearchPaintings(searchQuery)
	} else {
		result = h.Repo.GetPaintings()
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Ошибка кодирования", http.StatusInternalServerError)
		return
	}
}
