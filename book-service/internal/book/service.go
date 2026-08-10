package book

import "context"

type BookRepository interface {
	GetBooks(ctx context.Context) ([]Book, error)
}

type BookService struct {
	repo BookRepository
}

func NewBookService(r BookRepository) *BookService {
	return &BookService{repo: r}
}

func (s *BookService) GetBooks(ctx context.Context) ([]Book, error) {
	books, err := s.repo.GetBooks(ctx)
	if err != nil {
		return nil, err
	}

	return books, nil
}
