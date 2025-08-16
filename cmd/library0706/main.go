package main

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/Dorrrke/library0706/internal"
	dbstorage "github.com/Dorrrke/library0706/internal/repository/db-storage"
	"github.com/Dorrrke/library0706/internal/repository/inmemory"
	"github.com/Dorrrke/library0706/internal/server"
	"github.com/Dorrrke/library0706/pkg/logger"
)

func main() {
	cfg := internal.ReadConfig()
	log := logger.Get(cfg.Debug)
	log.Info().Msg("Server starting...")
	log.Debug().Any("config", cfg).Send()

	var repo server.Repository

	log.Debug().Msg("Connect to DB...")
	db, err := dbstorage.NewStorage(cfg.DBDSN)
	if err == nil {
		log.Debug().Msg("Connect to DB success")
		repo = db
		dbstorage.Migrations(cfg.DBDSN, "./migrations")
	} else {
		log.Warn().Err(err).Msg("Connect to DB failed. Using in-memory storage")
		repo = inmemory.NewStorage()
	}

	srv := server.NewServer(repo)

	if err := srv.Start(*cfg); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
