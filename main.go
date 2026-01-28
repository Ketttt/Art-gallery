package main

import (
	"art-gallery/handlers"
	"art-gallery/repository"
	"fmt"
	"net/http"
)

func main() {
	store := repository.NewGalleryStore()
	galleryH := &handlers.GalleryHandler{Repo: store}
	authorH := &handlers.AuthorHandler{Repo: store}
	http.Handle("/gallery", galleryH)
	http.Handle("/authors", authorH)
	fmt.Println("Сервер ArtGallery запущен на http://localhost:8080/gallery")
	fmt.Println("Сервер ArtGallery запущен на http://localhost:8080/authors")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска:", err)
	}
}
