package dbstorage

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Dorrrke/library0706/internal/domain/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Storage) GetBooksList() ([]models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := s.conn.Query(ctx, "SELECT * FROM books")
	if err != nil {
		return nil, err
	}

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.BookID, &book.Author, &book.Lable, &book.Description, &book.Genre, &book.WritedAt, &book.Count)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (s *Storage) SaveBook(book models.Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.conn.Exec(ctx, "INSERT INTO books (bid, author, lable, description, genre, writed_at, count) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		book.BookID, book.Author, book.Lable, book.Description, book.Genre, book.WritedAt, book.Count,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
				_, err := s.conn.Exec(ctx, "UPDATE books SET count = count + 1 WHERE lable = $1", book.Lable)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (s *Storage) SaveBooks(books []models.Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Println("Error rollback: ", err.Error())
		}
	}()

	_, err = tx.Prepare(ctx, "save_book", "INSERT INTO books (bid, author, lable, description, genre, writed_at, count) VALUES ($1, $2, $3, $4, $5, $6, $7)")
	if err != nil {
		return err
	}

	for _, book := range books {
		_, err = tx.Exec(ctx, "save_book", book.BookID, book.Author, book.Lable, book.Description, book.Genre, book.WritedAt, book.Count)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
