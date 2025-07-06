package main

import (
	"log"

	"github.com/Dorrrke/library0706/internal"
	dbstorage "github.com/Dorrrke/library0706/internal/repository/db-storage"
	"github.com/Dorrrke/library0706/internal/repository/inmemory"
	"github.com/Dorrrke/library0706/internal/server"
)

func main() {
	cfg := internal.ReadConfig()

	log.Printf("\nServer addr: %s\nServer port: %d\n\n", cfg.Addr, cfg.Port)

	log.Println("Starting server .....")

	var repo server.Repository

	db, err := dbstorage.NewStorage()
	if err == nil {
		repo = db
	} else {
		repo = inmemory.NewStorage()
	}

	srv := server.NewServer(repo)

	if err := srv.Start(*cfg); err != nil {
		panic(err)
	}
}
