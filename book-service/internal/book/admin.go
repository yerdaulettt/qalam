package book

import (
	"context"
)

type adminRepository interface {
	NewGenre(ctx context.Context, name string) (Genre, error)
	UpdateGenre(ctx context.Context, genreId int, name string) error
	DeleteGenre(ctx context.Context, genreId int) error
	AddBook(ctx context.Context, newBook NewBook) error
	UpdateBook(ctx context.Context, newBook BookUpdate) error
	DeleteBook(ctx context.Context, b BookDeletedEvent) error
	AddAuthor(ctx context.Context, newAuthor *NewAuthor) (*Author, error)
	UpdateAuthor(ctx context.Context, authorId int, newAuthor NewAuthor) error
}

type AdminService struct {
	repo adminRepository
}

func NewAdminService(r adminRepository) *AdminService {
	return &AdminService{repo: r}
}

func (s *AdminService) NewGenre(ctx context.Context, name string) (Genre, error) {
	g, err := s.repo.NewGenre(ctx, name)
	if err != nil {
		return Genre{}, err
	}

	return g, nil
}

func (s *AdminService) UpdateGenre(ctx context.Context, genreId int, name string) error {
	if err := s.repo.UpdateGenre(ctx, genreId, name); err != nil {
		return err
	}

	return nil
}

func (s *AdminService) DeleteGenre(ctx context.Context, genreId int) error {
	if err := s.repo.DeleteGenre(ctx, genreId); err != nil {
		return err
	}

	return nil
}

func (s *AdminService) AddBook(ctx context.Context, newBook NewBook) error {
	err := s.repo.AddBook(ctx, newBook)
	if err != nil {
		return err
	}

	return nil
}

func (s *AdminService) UpdateBook(ctx context.Context, newBook BookUpdate) error {
	err := s.repo.UpdateBook(ctx, newBook)
	if err != nil {
		return err
	}

	return nil
}

func (s *AdminService) DeleteBook(ctx context.Context, bookId int) error {
	err := s.repo.DeleteBook(ctx, BookDeletedEvent{BookId: bookId})
	if err != nil {
		return err
	}

	return nil
}

func (s *AdminService) AddAuthor(ctx context.Context, newAuthor *NewAuthor) (*Author, error) {
	a, err := s.repo.AddAuthor(ctx, newAuthor)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (s *AdminService) UpdateAuthor(ctx context.Context, authorId int, newAuthor NewAuthor) error {
	err := s.repo.UpdateAuthor(ctx, authorId, newAuthor)
	if err != nil {
		return err
	}

	return nil
}
