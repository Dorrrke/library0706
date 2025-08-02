package dbstorage

import (
	"context"
	"errors"
	"log"
	"time"

	domainErrors "github.com/Dorrrke/library0706/internal/domain/errors"
	"github.com/Dorrrke/library0706/internal/domain/models"
	"github.com/google/uuid"
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

func (s *Storage) GetBook(bid string) (models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var book models.Book
	err := s.conn.QueryRow(ctx, "SELECT * FROM books WHERE bid = $1", bid).
		Scan(&book.BookID, &book.Author, &book.Lable, &book.Description, &book.Genre, &book.WritedAt, &book.Count)
	if err != nil {
		return models.Book{}, domainErrors.ErrBookNotFound
	}
	return book, nil
}

func (s *Storage) SaveBook(book models.Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.conn.Exec(ctx, "INSERT INTO books (bid, author, lable, description, genre, writed_at) VALUES ($1, $2, $3, $4, $5, $6)",
		book.BookID, book.Author, book.Lable, book.Description, book.Genre, book.WritedAt,
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	_, err = tx.Prepare(ctx, "save_book",
		"INSERT INTO books (bid, author, lable, description, genre, writed_at) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (lable) DO UPDATE SET count = books.count + 1")
	if err != nil {
		return err
	}

	for _, book := range books {
		_, err = tx.Exec(ctx, "save_book", book.BookID, book.Author, book.Lable, book.Description, book.Genre, book.WritedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *Storage) BorrowBook(bid string, uid string) error {
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

	var book models.Book
	err = tx.QueryRow(ctx, "SELECT * FROM books WHERE bid = $1", bid).
		Scan(&book.BookID, &book.Author, &book.Lable, &book.Description, &book.Genre, &book.WritedAt, &book.Count)
	if err != nil {
		log.Printf("Error get book: %s", err.Error())
		return domainErrors.ErrBookNotFound
	}

	if book.Count == 0 {
		return domainErrors.ErrBooksAreOut
	}

	_, err = tx.Exec(ctx, "UPDATE books SET count = count - 1 WHERE bid = $1", bid)
	if err != nil {
		return err
	}

	tID := uuid.New().String()
	_, err = tx.Exec(ctx, "INSERT INTO borrowed_books (id, book_id, user_id) VALUES ($1, $2, $3)", tID, bid, uid)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *Storage) ReturnBook(bid, uid string) error {
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

	_, err = tx.Exec(ctx, "UPDATE books SET count = count + 1 WHERE bid = $1", bid)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM borrowed_books WHERE book_id = $1 AND user_id = $2", bid, uid)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
