package server

import (
	"fmt"

	"github.com/Dorrrke/library0706/internal"
	"github.com/Dorrrke/library0706/internal/domain/models"

	"github.com/gin-gonic/gin"
)

type Repository interface {
	GetUser(email string) (models.User, error)
	SaveUser(user models.User) error
	GetBooksList() ([]models.Book, error)
	GetBook(bid string) (models.Book, error)
	SaveBook(book models.Book) error
	SaveBooks([]models.Book) error
	BorrowBook(bid string, uid string) error
	ReturnBook(bid, uid string) error
	SaveRefreshToken(refreshToken, tokenID, userID string) error
	CheckRefreshToken(tokenID string) (bool, error)
}

type LibraryApi struct {
	db Repository
}

func NewServer(db Repository) *LibraryApi {
	return &LibraryApi{
		db: db,
	}
}

func (s *LibraryApi) Start(cfg internal.Config) error {
	router := gin.Default()
	router.POST("/refresh", s.refreshHandler)
	router.GET("/health", s.healthHandler)
	router.POST("/books")
	books := router.Group("/books")
	{
		books.POST("/create", s.JWTAuthMiddleware(), s.newBook)
		books.POST("/create/batch", s.JWTAuthMiddleware(), s.newBooks)
		books.GET("/list", s.booksList)
		books.GET("/get/:bookID", s.JWTAuthMiddleware(), s.getBook)
		books.PUT("/borrow/:bookID", s.JWTAuthMiddleware(), s.borrowBook)
		books.PUT("/borrow/:bookID/return", s.JWTAuthMiddleware(), s.returnBook)
		books.PUT("/update/:bookID", s.JWTAuthMiddleware())
		books.DELETE("/delete/:bookID", s.JWTAuthMiddleware())
	}
	users := router.Group("/users")
	{
		users.POST("/register", s.register)
		users.POST("/login", s.login)
	}

	return router.Run(fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port)) //0.0.0.0:8080
}

func (s *LibraryApi) healthHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"status": "ok"})
}
