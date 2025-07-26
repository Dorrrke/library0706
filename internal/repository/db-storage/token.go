package dbstorage

import (
	"context"
	"time"
)

func (s *Storage) SaveRefreshToken(refreshToken, tokenID, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.conn.Exec(ctx, "INSERT INTO tokens (tid, uid, token) VALUES ($1, $2, $3)", tokenID, userID, refreshToken)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) CheckRefreshToken(tokenID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var tokenValid bool
	err := s.conn.QueryRow(ctx, "SELECT valid FROM tokens WHERE tid = $1", tokenID).Scan(&tokenValid)
	if err != nil {
		return false, err
	}

	return tokenValid, nil
}
