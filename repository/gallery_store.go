package repository

import (
	"art-gallery/models"
	"encoding/json"
	"errors"
	"os"
	"strings"
)

const fileName = "gallery.json"

// Упаковка. Нужна только в момент чтения/записи файла. В свифт это была бы вспомогательная структура для Codable.
type GalleryData struct {
	Paintings []models.Painting `json:"paintings"`
	Authors   []models.Author   `json:"authors"`
}

// Живая память. Это твой сейф, который стоит в офисе. В нем лежат «папки» с картинами и авторами, пока сервер включен.
type GalleryStore struct {
	paintings []models.Painting
	authors   []models.Author
}

// Запуск. Когда сервер включается, эта функция проверяет: «Есть ли данные на диске?». Если нет — создает начальный набор.
func NewGalleryStore() *GalleryStore {
	store := &GalleryStore{}

	if err := store.loadFromFile(); err != nil {
		store.authors = []models.Author{
			{ID: "vango", Name: "Винсент Ван Гог", BirthYear: 185, Country: "Нидерланды", Portrait: "https://upload.wikimedia.org/wikipedia/commons/thumb/3/38/VanGogh_1887_Selbstbildnis.jpg/960px-VanGogh_1887_Selbstbildnis.jpg"},
			{ID: "da-vinci", Name: "Леонардо да Винчи", BirthYear: 1452, Country: "Италия", Portrait: "http://upload.wikimedia.org/wikipedia/commons/b/ba/Leonardo_self.jpg"},
		}

		store.paintings = []models.Painting{
			{ID: "1", Title: "Звездная ночь", AuthorID: "vango", Year: 1889, IsFavorite: true},
			{ID: "6", Title: "Мона Лиза", AuthorID: "da-vinci", Year: 1503, IsFavorite: false},
		}
		store.saveToFile()
	}

	return store
}

// Кнопка «Сохранить». Каждый раз, когда мы что-то меняем (лайк, удаление, добавление),
// мы вызываем её, чтобы изменения не пропали после выключения сервера.
func (s *GalleryStore) saveToFile() {
	dataToSave := GalleryData{
		Paintings: s.paintings,
		Authors:   s.authors,
	}

	bytes, _ := json.MarshalIndent(dataToSave, "", " ")
	os.WriteFile(fileName, bytes, 0644)
}

// Чтение при старте. Перекладывает данные из скучного текстового файла в быстрые слайсы (массивы) в оперативной памяти.
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
		painting := strings.ToLower(p.Title)
		if strings.Contains(painting, query) {
			response = append(response, p)
		}
	}
	return response
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
