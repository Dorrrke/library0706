package main

import (
	"log"

	"github.com/Dorrrke/library0706/internal"
	inmemory "github.com/Dorrrke/library0706/internal/repository/inmemory"
	"github.com/Dorrrke/library0706/internal/server"
)

func main() {
	cfg := internal.ReadConfig()

	log.Printf("\nServer addr: %s\nServer port: %d\n\n", cfg.Addr, cfg.Port)

	log.Println("Starting server .....")
	db := inmemory.NewStorage()

	srv := server.NewServer(db)

	if err := srv.Start(*cfg); err != nil {
		panic(err)
	}
}
