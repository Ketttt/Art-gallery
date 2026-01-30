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
	popularH := &handlers.PopularHandler{Repo: store}
	factsH := &handlers.FactsHandler{Repo: store}
	http.Handle("/gallery", galleryH)
	http.Handle("/authors", authorH)
	http.Handle("/popular", popularH)
	http.Handle("/facts", factsH)
	fmt.Println("Сервер ArtGallery запущен на http://localhost:8080/gallery")
	fmt.Println("Сервер ArtGallery запущен на http://localhost:8080/authors")
	fmt.Println("Сервер ArtGallery запущен на http://localhost:8080/popular")
	fmt.Println("Сервер ArtGallery запущен на http://localhost:8080/facts")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска:", err)
	}
}
