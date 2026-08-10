package book

import "encoding/json"

type Genre struct {
	Id   int    `json:"genre_id"`
	Name string `json:"genre_name"`
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
