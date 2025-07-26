package dbstorage

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
)

type Storage struct {
	conn *pgx.Conn
}

func NewStorage(connStr string) (*Storage, error) {
	conn, err := pgx.Connect(context.Background(), connStr) //TODO: добавить реальную строчку подключения
	if err != nil {
		return nil, err
	}
	return &Storage{
		conn: conn,
	}, nil
}

func Migrations(dbDsn string, migratePath string) error {
	mPath := fmt.Sprintf("file://%s", migratePath)
	m, err := migrate.New(mPath, dbDsn)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Database is up to date")
		}
		return err
	}

	log.Println("Database migrated successfully")
	return nil
}
