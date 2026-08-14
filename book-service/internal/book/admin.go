package book

import (
	"context"
)

type adminRepository interface {
	NewGenre(ctx context.Context, name string) (Genre, error)
	UpdateGenre(ctx context.Context, genreId int, name string) (Genre, error)
	DeleteGenre(ctx context.Context, genreId int) (Genre, error)
	UpdateBook(ctx context.Context, newBook *BookUpdate) (*BookUpdate, error)
	AddAuthor(ctx context.Context, newAuthor *NewAuthor) (*Author, error)
	UpdateAuthor(ctx context.Context, authorId int, newAuthor *NewAuthor) (*Author, error)
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

func (s *AdminService) UpdateGenre(ctx context.Context, genreId int, name string) (Genre, error) {
	g, err := s.repo.UpdateGenre(ctx, genreId, name)
	if err != nil {
		return Genre{}, err
	}

	return g, nil
}

func (s *AdminService) DeleteGenre(ctx context.Context, genreId int) (Genre, error) {
	g, err := s.repo.DeleteGenre(ctx, genreId)
	if err != nil {
		return Genre{}, err
	}

	return g, nil
}

func (s *AdminService) UpdateBook(ctx context.Context, newBook *BookUpdate) (*BookUpdate, error) {
	b, err := s.repo.UpdateBook(ctx, newBook)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *AdminService) AddAuthor(ctx context.Context, newAuthor *NewAuthor) (*Author, error) {
	a, err := s.repo.AddAuthor(ctx, newAuthor)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (s *AdminService) UpdateAuthor(ctx context.Context, authorId int, newAuthor *NewAuthor) (*Author, error) {
	a, err := s.repo.UpdateAuthor(ctx, authorId, newAuthor)
	if err != nil {
		return nil, err
	}

	return a, nil
}
