package book

import (
	"context"
)

type BookRepository interface {
	GetBooks(ctx context.Context) ([]Book, error)
	GetBook(ctx context.Context, bookId int) (*BookDetails, error)
	GetGenres(ctx context.Context) ([]GenreDetails, error)
	AddRating(ctx context.Context, bookId, userId, rating int) error
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

func (s *BookService) GetBook(ctx context.Context, bookId int) (*BookDetails, error) {
	b, err := s.repo.GetBook(ctx, bookId)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *BookService) GetGenres(ctx context.Context) ([]GenreDetails, error) {
	g, err := s.repo.GetGenres(ctx)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (s *BookService) AddRating(ctx context.Context, bookId, userId, rating int) error {
	err := s.repo.AddRating(ctx, bookId, userId, rating)
	if err != nil {
		return err
	}

	return nil
}
