package review

type Review struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	UserId  int    `json:"user_id"`
	BookId  int    `json:"book_id"`
}

type ReviewDetailAdmin struct {
	Id       int    `json:"id"`
	Content  string `json:"content"`
	BookName string `json:"book_name"`
	BookId   int    `json:"book_id"`
	Username string `json:"username"`
	UserId   int    `json:"user_id"`
}

type ReviewDetail struct {
	Id       int    `json:"id"`
	Content  string `json:"content"`
	Username string `json:"username"`
	BookId   int    `json:"book_id"`
	Likes    int    `json:"likes"`
}

type UserReview struct {
	Id       int    `json:"id"`
	Content  string `json:"content"`
	BookName string `json:"book_name"`
	BookId   int    `json:"book_id"`
	Likes    int    `json:"likes"`
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
	BookId int `json:"book_id"`
}

type BookMessage struct {
	BookId int    `json:"book_id"`
	Name   string `json:"name"`
}
