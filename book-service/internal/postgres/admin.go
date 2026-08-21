package postgres

import (
	"context"
	"errors"
	"log"

	"book-service/internal/book"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type adminRepo struct {
	db *pgxpool.Pool
}

func NewAdminRepo(db *pgxpool.Pool) *adminRepo {
	return &adminRepo{db: db}
}

func (r *adminRepo) NewGenre(ctx context.Context, name string) (book.Genre, error) {
	query := "insert into genres (name) values ($1) returning id, name"

	var g book.Genre
	if err := r.db.QueryRow(ctx, query, name).Scan(&g.Id, &g.Name); err != nil {
		return book.Genre{}, err
	}

	return g, nil
}

func (r *adminRepo) UpdateGenre(ctx context.Context, genreId int, name string) (book.Genre, error) {
	query := "update genres set name = $1 where id = $2 returning id, name"

	var g book.Genre
	if err := r.db.QueryRow(ctx, query, name, genreId).Scan(&g.Id, &g.Name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return book.Genre{}, book.ErrNotFound
		}

		return book.Genre{}, err
	}

	return g, nil
}

func (r *adminRepo) DeleteGenre(ctx context.Context, genreId int) (book.Genre, error) {
	query := "delete from genres where id = $1 returning id, name"

	var g book.Genre
	if err := r.db.QueryRow(ctx, query, genreId).Scan(&g.Id, &g.Name); err != nil {
		return book.Genre{}, err
	}

	return g, nil
}

func (r *adminRepo) UpdateBook(ctx context.Context, newBook *book.BookUpdate) (*book.BookUpdate, error) {
	columns := []string{}
	params := []any{}

	if newBook.Name != "" {
		columns = append(columns, "name")
		params = append(params, newBook.Name)
	}

	if newBook.Description != "" {
		columns = append(columns, "description")
		params = append(params, newBook.Description)
	}

	if newBook.AuthorId != 0 {
		columns = append(columns, "author_id")
		params = append(params, newBook.AuthorId)
	}

	params = append(params, newBook.Id)

	query := updateHelper("books", columns, []string{"id, name, description, author_id"})
	log.Println(params...)
	var b book.BookUpdate

	err := r.db.QueryRow(ctx, query, params...).Scan(&b.Id, &b.Name, &b.Description, &b.AuthorId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, book.ErrNotFound
		}

		return nil, err
	}

	return &b, nil
}

func (r *adminRepo) DeleteBook(ctx context.Context, bookId int) error {
	res, err := r.db.Exec(ctx, "delete from books where id = $1", bookId)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return book.ErrNotFound
	}

	return nil
}

func (r *adminRepo) AddAuthor(ctx context.Context, newAuthor *book.NewAuthor) (*book.Author, error) {
	query := "insert into authors (name, surname, about) values ($1, $2, $3) returning id, name, surname, about"

	var a book.Author
	err := r.db.QueryRow(ctx, query, newAuthor.Name, newAuthor.Surname, newAuthor.About).Scan(&a.Id, &a.Name, &a.Surname, &a.About)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (r *adminRepo) UpdateAuthor(ctx context.Context, authorId int, newAuthor *book.NewAuthor) (*book.Author, error) {
	columns := []string{}
	params := []any{}

	if newAuthor.Name != "" {
		columns = append(columns, "name")
		params = append(params, newAuthor.Name)
	}

	if newAuthor.Surname != "" {
		columns = append(columns, "surname")
		params = append(params, newAuthor.Surname)
	}

	if newAuthor.About != "" {
		columns = append(columns, "about")
		params = append(params, newAuthor.About)
	}

	params = append(params, authorId)

	query := updateHelper("authors", columns, []string{"id, name, surname, about"})
	var a book.Author

	if err := r.db.QueryRow(ctx, query, params...).Scan(&a.Id, &a.Name, &a.Surname, &a.About); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, book.ErrNotFound
		}

		return nil, err
	}

	return &a, nil
}
