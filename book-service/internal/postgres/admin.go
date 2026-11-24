package postgres

import (
	"context"
	"fmt"
	"strings"

	"book-service/internal/book"

	"github.com/jackc/pgx/v5/pgconn"
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

func (r *adminRepo) UpdateGenre(ctx context.Context, genreId int, name string) error {
	query := "update genres set name = $1 where id = $2"

	res, err := r.db.Exec(ctx, query, name, genreId)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return book.ErrNotFound
	}

	return nil
}

func (r *adminRepo) DeleteGenre(ctx context.Context, genreId int) error {
	query := "delete from genres where id = $1"

	res, err := r.db.Exec(ctx, query, genreId)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return book.ErrNotFound
	}

	return nil
}

func (r *adminRepo) AddBook(ctx context.Context, newBook book.NewBook) error {
	t, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer t.Rollback(ctx)

	query := "insert into books (name, description, author_id) values ($1, $2, $3) returning id"
	var bookId int

	err = t.QueryRow(ctx, query, newBook.Name, newBook.Description, newBook.AuthorId).Scan(&bookId)
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok && err.Code == "23503" {
			return fmt.Errorf("Author %d %w", newBook.AuthorId, book.ErrNotFound)
		}

		return err
	}

	if len(newBook.Genres) > 0 {
		var b strings.Builder
		b.WriteString("insert into book_genres (book_id, genre_id) values ($1, $2)")

		params := []any{bookId, newBook.Genres[0]}

		for i, id := range newBook.Genres[1:] {
			params = append(params, id)
			fmt.Fprintf(&b, ", ($1, $%d)", i+3)
		}

		_, err = t.Exec(ctx, b.String(), params...)
		if err != nil {
			if err, ok := err.(*pgconn.PgError); ok && err.Code == "23503" {
				return fmt.Errorf("Genre %w", book.ErrNotFound)
			}

			return err
		}
	}

	err = t.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *adminRepo) UpdateBook(ctx context.Context, newBook book.BookUpdate) error {
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
	query := updateHelper("books", columns)

	t, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer t.Rollback(ctx)

	res, err := t.Exec(ctx, query, params...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return book.ErrNotFound
	}

	err = t.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *adminRepo) DeleteBook(ctx context.Context, b book.BookDeletedEvent) error {
	t, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer t.Rollback(ctx)

	res, err := t.Exec(ctx, "delete from books where id = $1", b.BookId)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return book.ErrNotFound
	}

	_, err = t.Exec(ctx, "insert into outbox values (default, 'book.deleted', $1)", b)
	if err != nil {
		return err
	}

	err = t.Commit(ctx)
	if err != nil {
		return err
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

func (r *adminRepo) UpdateAuthor(ctx context.Context, authorId int, newAuthor book.NewAuthor) error {
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
	query := updateHelper("authors", columns)

	res, err := r.db.Exec(ctx, query, params...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return book.ErrNotFound
	}

	return nil
}
