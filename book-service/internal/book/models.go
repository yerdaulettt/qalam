package book

import "encoding/json"

type Genre struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type GenreDetails struct {
	Genre
	TotalBooks int `json:"total_books"`
}

type GenreSlice []Genre

func (g *GenreSlice) Scan(src any) error {
	var data []byte

	switch v := src.(type) {
	case []byte:
		data = v
	}

	return json.Unmarshal(data, g)
}

type Book struct {
	Id            int     `json:"id"`
	Name          string  `json:"name"`
	AuthorId      int     `json:"author_id"`
	AuthorName    string  `json:"author_name"`
	AuthorSurname string  `json:"author_surname"`
	Genres        []Genre `json:"genres"`
}

type BookDetails struct {
	Id            int     `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	AuthorId      int     `json:"author_id"`
	AuthorName    string  `json:"author_name"`
	AuthorSurname string  `json:"author_surname"`
	Genres        []Genre `json:"genres"`
}

type Author struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	About      string `json:"about"`
	TotalBooks int    `json:"total_books"`
}
