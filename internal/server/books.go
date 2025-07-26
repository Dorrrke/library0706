package server

import (
	"log"
	"net/http"

	"github.com/Dorrrke/library0706/internal/domain/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (api *LibraryApi) booksList(ctx *gin.Context) {
	books, err := api.db.GetBooksList()
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, books)
}

func (api *LibraryApi) newBook(ctx *gin.Context) {
	var book models.Book
	err := ctx.ShouldBindBodyWithJSON(&book)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book.BookID = uuid.New().String()
	api.db.SaveBook(book)

	ctx.JSON(http.StatusCreated, book)
}

func (api *LibraryApi) newBooks(ctx *gin.Context) {
	var books []models.Book
	err := ctx.ShouldBindBodyWithJSON(&books)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i := range books {
		books[i].BookID = uuid.New().String()
	}

	if err = api.db.SaveBooks(books); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, "books created")
}
