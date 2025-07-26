package dbstorage

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
)

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
