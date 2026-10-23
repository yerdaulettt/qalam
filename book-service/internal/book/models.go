package book

type Genre struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type GenreDetails struct {
	Genre
	TotalBooks int `json:"total_books"`
}

type Book struct {
	Id            int     `json:"id"`
	Name          string  `json:"name"`
	AuthorId      int     `json:"author_id"`
	AuthorName    string  `json:"author_name"`
	AuthorSurname string  `json:"author_surname"`
	Rating        int     `json:"rating"`
	RatingCount   int     `json:"rating_count"`
	Genres        []Genre `json:"genres"`
}

type NewBook struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	AuthorId    int    `json:"author_id"`
	Genres      []int  `json:"genres"`
}

type BookUpdate struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AuthorId    int    `json:"author_id"`
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
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	About   string `json:"about"`
}

type NewAuthor struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	About   string `json:"about"`
}

type BookDeletedMessage struct {
	BookId int `json:"book_id"`
}

type BookMessage struct {
	BookId int    `json:"book_id"`
	Name   string `json:"name"`
}

type Filter struct {
	Column    string
	Operation string
	Value     string
}
