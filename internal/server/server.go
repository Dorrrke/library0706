package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Dorrrke/library0706/internal"
	domainErrors "github.com/Dorrrke/library0706/internal/domain/errors"
	"github.com/Dorrrke/library0706/internal/domain/models"
	inmemory "github.com/Dorrrke/library0706/internal/repository/inmemory"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

type LibraryApi struct {
	db *inmemory.Storage
}

func NewServer(db *inmemory.Storage) *LibraryApi {
	return &LibraryApi{
		db: db,
	}
}

func (s *LibraryApi) Start(cfg internal.Config) error {
	router := gin.Default()
	router.POST("/books")
	books := router.Group("/books")
	{
		books.POST("/create", s.newBook)
		books.GET("/list", s.booksList)
		books.GET("/get/:bookID")
		books.PUT("/update/:bookID")
		books.DELETE("/delete/:bookID")
	}
	users := router.Group("/users")
	{
		users.POST("/register", s.register)
		users.POST("/login", s.login)
	}

	return router.Run(fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port)) //0.0.0.0:8080
}

func (api *LibraryApi) register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		log.Println("failed to bind body: ", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valid := validator.New()

	err := valid.Struct(user)
	if err != nil {
		log.Println("Failed validation: ", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.UID = uuid.New().String()
	log.Println(user.UID)

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Pass), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Pass = string(hash)

	err = api.db.SaveUser(user)
	if err != nil {
		if errors.Is(err, domainErrors.ErrUserAlredyExist) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.String(200, "OK")
}

func (api *LibraryApi) login(ctx *gin.Context) {
	var user models.UserLogin
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		log.Println("failed to bind body: ", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valid := validator.New()

	err := valid.Struct(user)
	if err != nil {
		log.Println("Failed validation: ", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbUser, err := api.db.GetUser(user)
	if err != nil {
		if errors.Is(err, domainErrors.ErrIvalidCreds) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Pass), []byte(user.Pass))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dbUser)
}

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

	api.db.SaveBook(book)

	ctx.JSON(http.StatusCreated, book)
}
