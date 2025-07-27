package errors

import "errors"

var (
	ErrUserAlredyExist = errors.New("user alredy exist")
	ErrIvalidCreds     = errors.New("ivalid email or pass")
)

var (
	ErrBooksListIsEmpty = errors.New("books list is empty")
	ErrBookNotFound     = errors.New("book not found")
	ErrBooksAreOut      = errors.New("all these books are sorted out")
)
