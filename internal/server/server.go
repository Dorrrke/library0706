package server

import (
	inmemory "github.com/Dorrrke/library0706/internal/repository/inmemory"

	"github.com/gin-gonic/gin"
)

type LibraryApi struct {
	db *inmemory.UserStrage
}

func NewServer(db *inmemory.UserStrage) *LibraryApi {
	return &LibraryApi{
		db: db,
	}
}

func (s *LibraryApi) Start() error {
	router := gin.Default()
	router.POST("/books")
	task := router.Group("/books")
	{
		task.POST("/save")
		task.PUT("/:id")
		task.DELETE("/:id")
		task.GET("/:id")
	}
	users := router.Group("/users")
	{
		users.POST("/register")
		users.POST("/login")
	}

	return router.Run(":8080")
}
