package postgres

import (
	"context"
	"errors"

	"book-service/internal/book"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type bookRepo struct {
	db *pgxpool.Pool
}

func NewBookRepo(db *pgxpool.Pool) *bookRepo {
	return &bookRepo{db: db}
}

func (r *bookRepo) GetBooks(ctx context.Context) ([]book.Book, error) {
	query := `
	with book_details as (
		select b.id, b.name, b.author_id, jsonb_agg_strict(to_jsonb(g)) as genre from
		(books as b left join book_genres as bg on b.id = bg.book_id)
		left join genres as g on bg.genre_id = g.id group by b.id)

	select bd.id, bd.name, a.id, a.name, a.surname, bd.genre from book_details as bd join authors as a on bd.author_id = a.id
	`

	var books []book.Book

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var b book.Book

		err := rows.Scan(&b.Id, &b.Name, &b.AuthorId, &b.AuthorName, &b.AuthorSurname, &b.Genres)
		if err != nil {
			return nil, err
		}

		books = append(books, b)
	}

	return books, nil
}

func (r *bookRepo) GetBook(ctx context.Context, bookId int) (*book.BookDetails, error) {
	query := `
	with book_details as (
		select b.id, b.name, b.description, b.author_id, jsonb_agg_strict(to_jsonb(g)) as genre
		from (books as b left join book_genres as bg on b.id = bg.book_id)
		left join genres as g on bg.genre_id = g.id where b.id = $1 group by b.id
	)
	
	select bd.id, bd.name, bd.description, a.id as author_id, a.name, a.surname, bd.genre from
	book_details as bd join authors as a on bd.author_id = a.id
	`

	var b book.BookDetails
	err := r.db.QueryRow(ctx, query, bookId).Scan(&b.Id, &b.Name, &b.Description, &b.AuthorId, &b.AuthorName, &b.AuthorSurname, &b.Genres)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, book.ErrNotFound
		}

		return nil, err
	}

	return &b, nil
}

func (r *bookRepo) GetGenres(ctx context.Context) ([]book.GenreDetails, error) {
	query := `select g.id, g.name, count(bg) from genres as g join book_genres as bg on g.id = bg.genre_id group by g.id`

	var genres []book.GenreDetails
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var g book.GenreDetails

		if err := rows.Scan(&g.Id, &g.Name, &g.TotalBooks); err != nil {
			return nil, err
		}

		genres = append(genres, g)
	}

	return genres, nil
}
