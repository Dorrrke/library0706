package dbstorage

import (
	"context"
	"time"

	"github.com/Dorrrke/library0706/internal/domain/models"
)

func (s *Storage) GetUser(email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User

	row := s.conn.QueryRow(ctx, "SELECT * FROM users WHERE email = $1", email)

	err := row.Scan(&user.UID, &user.Name, &user.Age, &user.Email, &user.Pass)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *Storage) SaveUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.conn.Exec(ctx, "INSERT INTO users (uid, name, age, email, pass) VALUES ($1, $2, $3, $4, $5)",
		user.UID, user.Name, user.Age, user.Email, user.Pass,
	)
	if err != nil {
		return err
	}

	return nil
}
