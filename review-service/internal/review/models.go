package review

type Review struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	UserId  int    `json:"user_id"`
	BookId  int    `json:"book_id"`
}

type ReviewReq struct {
	Content string `json:"content"`
	UserId  int
	BookId  int
}

type ReviewUpdate struct {
	Content  string `json:"content"`
	ReviewId int
	UserId   int
}

type Message struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
