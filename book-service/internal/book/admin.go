package book

import (
	"context"
	"encoding/json"
)

type adminRepository interface {
	NewGenre(ctx context.Context, name string) (Genre, error)
	UpdateGenre(ctx context.Context, genreId int, name string) (Genre, error)
	DeleteGenre(ctx context.Context, genreId int) (Genre, error)
	UpdateBook(ctx context.Context, newBook *BookUpdate) (*BookUpdate, error)
	DeleteBook(ctx context.Context, bookId int) error
	AddAuthor(ctx context.Context, newAuthor *NewAuthor) (*Author, error)
	UpdateAuthor(ctx context.Context, authorId int, newAuthor *NewAuthor) (*Author, error)
}

type adminPublisher interface {
	Publish(ctx context.Context, message []byte, queue string) error
}

type AdminService struct {
	repo      adminRepository
	publisher adminPublisher
}

func NewAdminService(r adminRepository, p adminPublisher) *AdminService {
	return &AdminService{repo: r, publisher: p}
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

func (s *AdminService) DeleteBook(ctx context.Context, bookId int) error {
	err := s.repo.DeleteBook(ctx, bookId)
	if err != nil {
		return err
	}

	message, err := json.Marshal(BookDeletedMessage{BookId: bookId})
	if err != nil {
		return err
	}

	err = s.publisher.Publish(ctx, message, "book.deleted")
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

func (s *AdminService) UpdateAuthor(ctx context.Context, authorId int, newAuthor *NewAuthor) (*Author, error) {
	a, err := s.repo.UpdateAuthor(ctx, authorId, newAuthor)
	if err != nil {
		return nil, err
	}

	return a, nil
}
