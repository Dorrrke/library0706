package main

import (
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

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

	log.Println("Try connect to DB: ", cfg.DBDSN)
	db, err := dbstorage.NewStorage(cfg.DBDSN)
	if err == nil {
		log.Println("Connect to DB successfully")
		repo = db
		dbstorage.Migrations(cfg.DBDSN, "./migrations")
	} else {
		log.Printf("Connect to DB failed: %s\nUsing in-memory storage\n", err.Error())
		repo = inmemory.NewStorage()
	}

	srv := server.NewServer(repo)

	if err := srv.Start(*cfg); err != nil {
		panic(err)
	}
}
