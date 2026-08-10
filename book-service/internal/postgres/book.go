package postgres

import (
	"context"

	"book-service/internal/book"

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
		select b.id, b.name, b.author_id, json_agg(json_build_object('genre_id', g.id, 'genre_name', g.name)) as genre from
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
