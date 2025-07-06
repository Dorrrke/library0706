package dbstorage

import (
	"context"
	"time"

	"github.com/Dorrrke/library0706/internal/domain/models"
	"github.com/jackc/pgx/v5"
)

type Storage struct {
	conn *pgx.Conn
}

func NewStorage() (*Storage, error) {
	conn, err := pgx.Connect(context.Background(), "conn_str") //TODO: добавить реальную строчку подключения
	if err != nil {
		return nil, err
	}
	return &Storage{
		conn: conn,
	}, nil
}

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
		return err
	}

	return nil
}

func (s *Storage) GetUser(email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User

	row := s.conn.QueryRow(ctx, "SELECT * FROM users WHERE email = $1", email)

	err := row.Scan(&user.UID, &user.Name, &user.Age, &user.Email, &user.Pass)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *Storage) SaveUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.conn.Exec(ctx, "INSERT INTO users (uid, name, age, email, pass) VALUES ($1, $2, $3, $4, $5)",
		user.UID, user.Name, user.Age, user.Email, user.Pass,
	)
	if err != nil {
		return err
	}

	return nil
}
