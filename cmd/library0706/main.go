package main

import (
	"log"

	inmemory "github.com/Dorrrke/library0706/internal/repository/inmemory"
	"github.com/Dorrrke/library0706/internal/server"
)

func main() {
	log.Println("Starting server .....")
	db := inmemory.NewUserStorage()

	srv := server.NewServer(db)

	if err := srv.Start(); err != nil {
		panic(err)
	}
}
