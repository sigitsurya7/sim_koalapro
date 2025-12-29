package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"koalbot_api/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User
	var deletedAt sql.NullTime
	var lastSeen sql.NullTime

	err := r.db.QueryRowContext(ctx, `
		SELECT id, uid, username, password, role, active, deleted_at, last_seen
		FROM m_user
		WHERE username = $1
		LIMIT 1
	`, username).Scan(
		&user.ID,
		&user.UID,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.Active,
		&deletedAt,
		&lastSeen,
	)
	if err != nil {
		return model.User{}, err
	}
	if deletedAt.Valid {
		t := deletedAt.Time
		user.DeletedAt = &t
	}
	if lastSeen.Valid {
		t := lastSeen.Time
		user.LastSeen = &t
	}

	return user, nil
}

func (r *UserRepository) ListUsers(ctx context.Context, search string, limit, offset int) ([]model.User, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM m_user
		WHERE deleted_at IS NULL
			AND ($1 = '' OR username ILIKE '%' || $1 || '%')
	`, search).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, uid, username, role, active, created_at, created_by, updated_at, updated_by, last_seen
		FROM m_user
		WHERE deleted_at IS NULL
			AND ($1 = '' OR username ILIKE '%' || $1 || '%')
		ORDER BY id DESC
		LIMIT $2 OFFSET $3
	`, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]model.User, 0)
	for rows.Next() {
		var user model.User
		var createdBy sql.NullString
		var updatedAt sql.NullTime
		var updatedBy sql.NullString
		var lastSeen sql.NullTime

		if err := rows.Scan(
			&user.ID,
			&user.UID,
			&user.Username,
			&user.Role,
			&user.Active,
			&user.CreatedAt,
			&createdBy,
			&updatedAt,
			&updatedBy,
			&lastSeen,
		); err != nil {
			return nil, 0, err
		}

		if createdBy.Valid {
			val := createdBy.String
			user.CreatedBy = &val
		}
		if updatedAt.Valid {
			val := updatedAt.Time
			user.UpdatedAt = &val
		}
		if updatedBy.Valid {
			val := updatedBy.String
			user.UpdatedBy = &val
		}
		if lastSeen.Valid {
			val := lastSeen.Time
			user.LastSeen = &val
		}

		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) UpdateLastSeen(ctx context.Context, uid string, lastSeen time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE m_user
		SET last_seen = $2, updated_at = NOW()
		WHERE uid = $1
	`, uid, lastSeen)
	return err
}

func (r *UserRepository) CreateUser(ctx context.Context, username, passwordHash, role, createdBy string) (string, string, error) {
	var uid string
	var storedRole string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO m_user (username, password, role, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING uid, role
	`, username, passwordHash, role, createdBy).Scan(&uid, &storedRole)
	if err != nil {
		return "", "", err
	}
	return uid, storedRole, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, uid string, req UpdateUserRequest) error {
	setClauses := make([]string, 0, 5)
	args := make([]any, 0, 6)
	argPos := 1

	if req.Username != nil {
		setClauses = append(setClauses, fmt.Sprintf("username = $%d", argPos))
		args = append(args, *req.Username)
		argPos++
	}
	if req.Password != nil {
		setClauses = append(setClauses, fmt.Sprintf("password = $%d", argPos))
		args = append(args, *req.Password)
		argPos++
	}
	if req.Role != nil {
		setClauses = append(setClauses, fmt.Sprintf("role = $%d", argPos))
		args = append(args, *req.Role)
		argPos++
	}
	if req.Active != nil {
		setClauses = append(setClauses, fmt.Sprintf("active = $%d", argPos))
		args = append(args, *req.Active)
		argPos++
	}
	if req.UpdatedBy != nil {
		setClauses = append(setClauses, fmt.Sprintf("updated_by = $%d", argPos))
		args = append(args, *req.UpdatedBy)
		argPos++
	}

	setClauses = append(setClauses, "updated_at = NOW()")

	args = append(args, uid)
	query := fmt.Sprintf("UPDATE m_user SET %s WHERE uid = $%d", strings.Join(setClauses, ", "), argPos)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *UserRepository) SoftDeleteUser(ctx context.Context, uid, deletedBy string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE m_user
		SET deleted_at = NOW(), deleted_by = $2, active = FALSE, updated_at = NOW(), updated_by = $2
		WHERE uid = $1
	`, uid, deletedBy)
	return err
}
