package repository

import (
	"context"
	"database/sql"
	"time"
)

type TokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) SaveRefreshToken(ctx context.Context, userUID, token string, lastSeen, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO t_user_token (user_uid, token, last_seen, expires_at)
		VALUES ($1, $2, $3, $4)
	`, userUID, token, lastSeen, expiresAt)
	return err
}
