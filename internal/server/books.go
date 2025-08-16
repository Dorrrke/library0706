package server

import (
	"net/http"

	"github.com/Dorrrke/library0706/internal/domain/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (api *LibraryApi) booksList(ctx *gin.Context) {
	books, err := api.db.GetBooksList()
	if err != nil {
		api.log.Error().Err(err).Msg("failed to get books list")
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, books)
}

func (api *LibraryApi) getBook(ctx *gin.Context) {
	id := ctx.Param("bookID")
	api.log.Debug().Msgf("book id: %s", id)
	book, err := api.db.GetBook(id)
	if err != nil {
		api.log.Error().Err(err).Msg("failed to get book")
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, book)
}

func (api *LibraryApi) borrowBook(ctx *gin.Context) {
	bid := ctx.Param("bookID")
	api.log.Debug().Msgf("book id: %s", bid)
	ctxUid, exist := ctx.Get("user_id")
	if !exist {
		api.log.Error().Msg("failed to get user id from context")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	uid := ctxUid.(string)
	api.log.Debug().Msgf("user id: %s", uid)

	err := api.db.BorrowBook(bid, uid)
	if err != nil {
		api.log.Error().Err(err).Msg("failed to borrow book")
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, "book borrowed")
}

func (api *LibraryApi) returnBook(ctx *gin.Context) {
	bid := ctx.Param("bookID")
	api.log.Debug().Msgf("book id: %s", bid)
	ctxUid, exist := ctx.Get("user_id")
	if !exist {
		api.log.Error().Msg("failed to get user id from context")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	uid := ctxUid.(string)
	api.log.Debug().Msgf("user id: %s", uid)

	err := api.db.ReturnBook(bid, uid)
	if err != nil {
		api.log.Error().Err(err).Msg("failed to return book")
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, "book returned")
}

func (api *LibraryApi) newBook(ctx *gin.Context) {
	var book models.Book
	err := ctx.ShouldBindBodyWithJSON(&book)
	if err != nil {
		api.log.Error().Err(err).Msg("failed to unmarshal body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book.BookID = uuid.New().String()
	api.log.Debug().Msgf("book id: %s", book.BookID)

	if err = api.db.SaveBook(book); err != nil {
		api.log.Error().Err(err).Msg("failed to save book")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, book)
}

func (api *LibraryApi) newBooks(ctx *gin.Context) {
	var books []models.Book
	err := ctx.ShouldBindBodyWithJSON(&books)
	if err != nil {
		api.log.Error().Err(err).Msg("failed to unmarshal body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i := range books {
		books[i].BookID = uuid.New().String()
		api.log.Debug().Msgf("book id: %s", books[i].BookID)
	}

	if err = api.db.SaveBooks(books); err != nil {
		api.log.Error().Err(err).Msg("failed to save books")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, "books created")
}
