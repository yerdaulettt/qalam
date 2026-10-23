package postgres

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"book-service/internal/book"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context, dbUrl string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func filterHelper(query string, filters []book.Filter, columns map[string]string) (string, []any) {
	var b strings.Builder
	b.WriteString(query)
	b.WriteString(" where ")

	var res []string
	var params []any

	var c string
	for i, f := range filters {
		c = columns[f.Column]
		if c == "" {
			c = f.Column
		}

		res = append(res, fmt.Sprintf("%s %s $%d", c, f.Operation, i+1))
		params = append(params, f.Value)
	}

	b.WriteString(strings.Join(res, " and "))

	return b.String(), params
}

func updateHelper(table string, columns, columnsReturn []string) string {
	var b strings.Builder

	if len(columns) == 1 {
		fmt.Fprintf(&b, "update %s set %s = $1 where id = $2 returning %s", table, columns[0], strings.Join(columnsReturn, ", "))
		return b.String()
	}

	fmt.Fprintf(&b, "update %s set (%s) = (", table, strings.Join(columns, ", "))

	b.WriteString("$1")
	for i := range columns[1:] {
		b.WriteString(", $")
		b.WriteString(strconv.Itoa(i + 2))
	}

	fmt.Fprintf(&b, ") where id = $%d returning %s", len(columns)+1, strings.Join(columnsReturn, ", "))

	log.Println(b.String())

	return b.String()
}
