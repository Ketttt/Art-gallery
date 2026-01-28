package models

type Author struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	BirthYear int    `json:"birth_year"`
	Country   string `json:"country"`
	Portrait  string `json:"portrait"`
}

type Painting struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	AuthorID   string `json:"author_id"`
	Year       int    `json:"year"`
	ImageUrl   string `json:"image_url"`
	IsFavorite bool   `json:"isFavorite"`
}

type PaintingResponse struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Author     Author `json:"author"`
	Year       int    `json:"year"`
	ImageUrl   string `json:"image_url"`
	IsFavorite bool   `json:"isFavorite"`
}
