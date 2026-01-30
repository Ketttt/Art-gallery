package repository

import (
	"art-gallery/models"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

const fileName = "gallery.json"

type GalleryData struct {
	Paintings []models.Painting `json:"paintings"`
	Authors   []models.Author   `json:"authors"`
	Facts     []models.Facts    `json:"facts"`
}

type GalleryStore struct {
	paintings []models.Painting
	authors   []models.Author
	facts     []models.Facts
}

func NewGalleryStore() *GalleryStore {
	store := &GalleryStore{}

	if err := store.loadFromFile(); err != nil {
		store.authors = []models.Author{
			{ID: "vango", Name: "Vincent van Gogh", BirthYear: 1853, Country: "Netherlands", Portrait: "https://upload.wikimedia.org/wikipedia/commons/thumb/3/38/VanGogh_1887_Selbstbildnis.jpg/960px-VanGogh_1887_Selbstbildnis.jpg"},
			{ID: "da-vinci", Name: "Leonardo da Vinci", BirthYear: 1452, Country: "Italy", Portrait: "http://upload.wikimedia.org/wikipedia/commons/b/ba/Leonardo_self.jpg"},
			{ID: "picasso", Name: "Pablo Picasso", BirthYear: 1881, Country: "Spain", Portrait: "https://upload.wikimedia.org/wikipedia/commons/9/98/Pablo_picasso_1.jpg"},
			{ID: "dali", Name: "Salvador Dalí", BirthYear: 1904, Country: "Spain", Portrait: "https://upload.wikimedia.org/wikipedia/commons/thumb/2/24/Salvador_Dal%C3%AD_1939.jpg/960px-Salvador_Dal%C3%AD_1939.jpg"},
			{ID: "monet", Name: "Claude Monet", BirthYear: 1840, Country: "France", Portrait: "http://upload.wikimedia.org/wikipedia/commons/thumb/3/33/Claude_Monet_1899_Nadar.jpg/500px-Claude_Monet_1899_Nadar.jpg"},
		}

		store.paintings = []models.Painting{
			{ID: "5", Title: "The Weeping Woman", AuthorID: "picasso", Year: 1937, ImageUrl: "https://upload.wikimedia.org/wikipedia/en/1/14/Picasso_The_Weeping_Woman_Tate_identifier_T05010_10.jpg", IsFavorite: false, IsPopular: true},
			{ID: "6", Title: "Mona Lisa", AuthorID: "da-vinci", Year: 1503, ImageUrl: "https://upload.wikimedia.org/wikipedia/commons/thumb/e/ec/Mona_Lisa%2C_by_Leonardo_da_Vinci%2C_from_C2RMF_retouched.jpg/800px-Mona_Lisa%2C_by_Leonardo_da_Vinci%2C_from_C2RMF_retouched.jpg", IsFavorite: false, IsPopular: true},
			{ID: "1", Title: "The Starry Night", AuthorID: "vango", Year: 1889, ImageUrl: "https://upload.wikimedia.org/wikipedia/commons/thumb/e/ea/Van_Gogh_-_Starry_Night_-_Google_Art_Project.jpg/1280px-Van_Gogh_-_Starry_Night_-_Google_Art_Project.jpg", IsFavorite: true, IsPopular: true},
			{ID: "2", Title: "Persistence of Memory", AuthorID: "dali", Year: 1931, ImageUrl: "https://upload.wikimedia.org/wikipedia/en/d/dd/The_Persistence_of_Memory.jpg", IsFavorite: false, IsPopular: true},
			{ID: "3", Title: "Pierrot", AuthorID: "picasso", Year: 1937, ImageUrl: "https://upload.wikimedia.org/wikipedia/en/thumb/7/76/Pablo_Picasso%2C_1918%2C_Pierrot%2C_oil_on_canvas%2C_92.7_x_73_cm%2C_Museum_of_Modern_Art.jpg/330px-Pablo_Picasso%2C_1918%2C_Pierrot%2C_oil_on_canvas%2C_92.7_x_73_cm%2C_Museum_of_Modern_Art.jpg", IsFavorite: false, IsPopular: true},
			{ID: "4", Title: "Water Lilies", AuthorID: "monet", Year: 1919, ImageUrl: "https://upload.wikimedia.org/wikipedia/commons/thumb/c/cb/Claude_Monet_Nympheas_1915_Musee_Marmottan_Paris.jpg/960px-Claude_Monet_Nympheas_1915_Musee_Marmottan_Paris.jpg", IsFavorite: false, IsPopular: true},
		}

		store.facts = []models.Facts{
			{ID: "1", Text: "Van Gogh painted 'The Starry Night' from the window of his asylum room at Saint-Rémy-de-Provence."},
			{ID: "2", Text: "Leonardo da Vinci could write with one hand and draw with the other simultaneously."},
			{ID: "3", Text: "The persistence of memory was inspired by Dalí seeing Camembert cheese melting in the sun."},
		}
		fmt.Printf("Загружено фактов для сохранения: %d\n", len(store.facts))
		store.saveToFile()
	}

	return store
}

func (s *GalleryStore) saveToFile() {
	dataToSave := GalleryData{
		Paintings: s.paintings,
		Authors:   s.authors,
		Facts:     s.facts,
	}

	bytes, _ := json.MarshalIndent(dataToSave, "", " ")
	os.WriteFile(fileName, bytes, 0644)
}

func (s *GalleryStore) loadFromFile() error {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	var data GalleryData
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	s.paintings = data.Paintings
	s.authors = data.Authors
	s.facts = data.Facts
	return nil
}

func (s *GalleryStore) ToggleFavorite(id string) bool {
	for i := range s.paintings {
		if s.paintings[i].ID == id {
			s.paintings[i].IsFavorite = !s.paintings[i].IsFavorite
			s.saveToFile()
			return true
		}
	}
	return false
}

func (s *GalleryStore) DeletePainting(id string) bool {
	for i, p := range s.paintings {
		if p.ID == id {
			s.paintings = append(s.paintings[:i], s.paintings[i+1:]...)
			s.saveToFile()
			return true
		}
	}
	return false
}

func (s *GalleryStore) AddPainting(p models.Painting) error {
	authorExists := false

	for _, a := range s.authors {
		if a.ID == p.AuthorID {
			authorExists = true
			break
		}
	}

	if !authorExists {
		return errors.New("author not found: cannot add painting for non-existent author")
	}

	s.paintings = append(s.paintings, p)
	s.saveToFile()
	return nil
}

func (s *GalleryStore) GetPaintings() []models.PaintingResponse {
	authorMap := make(map[string]models.Author)
	for _, a := range s.authors {
		authorMap[a.ID] = a
	}
	var response []models.PaintingResponse
	for _, p := range s.paintings {
		authors := authorMap[p.AuthorID]

		res := models.PaintingResponse{
			ID:         p.ID,
			Title:      p.Title,
			Author:     authors,
			Year:       p.Year,
			ImageUrl:   p.ImageUrl,
			IsFavorite: p.IsFavorite,
			IsPopular:  p.IsPopular,
		}
		response = append(response, res)
	}
	return response
}

func (s *GalleryStore) GetPaintingsByAuthor(authorID string) []models.PaintingResponse {
	authorMap := make(map[string]models.Author)
	for _, a := range s.authors {
		authorMap[a.ID] = a
	}
	var response []models.PaintingResponse
	for _, p := range s.paintings {
		if p.AuthorID == authorID {
			authorData := authorMap[p.AuthorID]
			res := models.PaintingResponse{
				ID:         p.ID,
				Title:      p.Title,
				Author:     authorData,
				Year:       p.Year,
				ImageUrl:   p.ImageUrl,
				IsFavorite: p.IsFavorite,
			}
			response = append(response, res)
		}
	}
	return response
}

func (s *GalleryStore) GetAuthors() []models.Author {
	return s.authors
}

func (s *GalleryStore) SearchPaintings(query string) []models.PaintingResponse {
	var response []models.PaintingResponse
	query = strings.ToLower(query)

	for _, p := range s.GetPaintings() {
		title := strings.ToLower(p.Title)
		author := strings.ToLower(p.Author.Name)
		if strings.Contains(title, query) || strings.Contains(author, query) {
			response = append(response, p)
		}
	}
	return response
}

func (s *GalleryStore) GetPopularPaintings() []models.PaintingResponse {
	var popular []models.PaintingResponse

	for _, p := range s.GetPaintings() {
		if p.IsPopular {
			popular = append(popular, p)
		}
	}
	return popular
}

func (s *GalleryStore) AddAuthor(a models.Author) error {
	authorExists := false

	for _, auth := range s.authors {
		if auth.ID == a.ID {
			authorExists = true
			break
		}
	}
	if authorExists {
		return errors.New("author already exist")
	}

	if a.Name == "" {
		return errors.New("Please fill in the required name field.")
	}
	s.authors = append(s.authors, a)
	s.saveToFile()
	return nil
}

func (s *GalleryStore) GetFacts() []models.Facts {
	return s.facts
}
