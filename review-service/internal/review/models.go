package review

type Review struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	BookId  int    `json:"book_id"`
}

type Message struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
