package inmemory

import (
	"fmt"

	"github.com/Dorrrke/library0706/internal/domain/errors"
	"github.com/Dorrrke/library0706/internal/domain/models"
	"github.com/google/uuid"
)

type Storage struct {
	userDB map[string]models.User
	bookDB map[string]models.Book
}

func NewStorage() *Storage {
	return &Storage{
		userDB: make(map[string]models.User),
		bookDB: make(map[string]models.Book),
	}
}

func (s *Storage) SaveUser(user models.User) error {
	for _, dbUser := range s.userDB {
		if dbUser.Email == user.Email {
			return fmt.Errorf("error in db: %w", errors.ErrUserAlredyExist)
		}
	}

	s.userDB[user.UID] = user
	return nil
}

func (s *Storage) GetUser(user models.UserLogin) (models.User, error) {
	for _, dbUser := range s.userDB {
		if dbUser.Email == user.Email {
			return dbUser, nil
		}
	}
	return models.User{}, errors.ErrIvalidCreds
}

func (s *Storage) GetBooksList() ([]models.Book, error) {
	var booksList []models.Book
	if len(s.bookDB) == 0 {
		return nil, errors.ErrBooksListIsEmpty
	}

	for _, book := range s.bookDB {
		booksList = append(booksList, book)
	}

	return booksList, nil
}

func (s *Storage) SaveBook(book models.Book) {
	for key, b := range s.bookDB {
		if b.Author == book.Author && b.Lable == book.Lable {
			mBook := s.bookDB[key]
			mBook.Count++
			s.bookDB[key] = mBook
			return
		}
	}

	bookID := uuid.New().String()
	book.BookID = bookID
	book.Count = 1

	s.bookDB[bookID] = book
}
