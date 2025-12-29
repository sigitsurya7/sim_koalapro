package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"koalbot_api/internal/model"
)

type MasterPenggunaRepository struct {
	db *sql.DB
}

func NewMasterPenggunaRepository(db *sql.DB) *MasterPenggunaRepository {
	return &MasterPenggunaRepository{db: db}
}

func (r *MasterPenggunaRepository) Create(ctx context.Context, idPengguna int64, telegram *string, jenis string, active bool) (model.MasterPengguna, error) {
	var result model.MasterPengguna
	var telegramVal sql.NullString
	if telegram != nil {
		telegramVal = sql.NullString{String: *telegram, Valid: true}
	}

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO m_pengguna (id_pengguna, telegram, jenis, active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, uuid, id_pengguna, telegram, jenis, active, created_at
	`, idPengguna, telegramVal, jenis, active).Scan(
		&result.ID,
		&result.UUID,
		&result.IDPengguna,
		&telegramVal,
		&result.Jenis,
		&result.Active,
		&result.CreatedAt,
	)
	if err != nil {
		return model.MasterPengguna{}, err
	}
	if telegramVal.Valid {
		val := telegramVal.String
		result.Telegram = &val
	}

	return result, nil
}

func (r *MasterPenggunaRepository) GetByID(ctx context.Context, id int64) (model.MasterPengguna, error) {
	var result model.MasterPengguna
	var telegram sql.NullString
	var updatedAt sql.NullTime
	var deletedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, `
		SELECT id, uuid, id_pengguna, telegram, jenis, active, created_at, updated_at, deleted_at
		FROM m_pengguna
		WHERE id = $1
	`, id).Scan(
		&result.ID,
		&result.UUID,
		&result.IDPengguna,
		&telegram,
		&result.Jenis,
		&result.Active,
		&result.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		return model.MasterPengguna{}, err
	}

	if telegram.Valid {
		val := telegram.String
		result.Telegram = &val
	}
	if updatedAt.Valid {
		val := updatedAt.Time
		result.UpdatedAt = &val
	}
	if deletedAt.Valid {
		val := deletedAt.Time
		result.DeletedAt = &val
	}

	return result, nil
}

func (r *MasterPenggunaRepository) GetByIDPengguna(ctx context.Context, idPengguna int64) (model.MasterPengguna, error) {
	var result model.MasterPengguna
	var telegram sql.NullString
	var updatedAt sql.NullTime
	var deletedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, `
		SELECT id, uuid, id_pengguna, telegram, jenis, active, created_at, updated_at, deleted_at
		FROM m_pengguna
		WHERE id_pengguna = $1
	`, idPengguna).Scan(
		&result.ID,
		&result.UUID,
		&result.IDPengguna,
		&telegram,
		&result.Jenis,
		&result.Active,
		&result.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		return model.MasterPengguna{}, err
	}

	if telegram.Valid {
		val := telegram.String
		result.Telegram = &val
	}
	if updatedAt.Valid {
		val := updatedAt.Time
		result.UpdatedAt = &val
	}
	if deletedAt.Valid {
		val := deletedAt.Time
		result.DeletedAt = &val
	}

	return result, nil
}

func (r *MasterPenggunaRepository) List(ctx context.Context, search string, jenis string, limit, offset int) ([]model.MasterPengguna, int, error) {
	conditions := []string{
		"deleted_at IS NULL",
		"($1 = '' OR telegram ILIKE '%' || $1 || '%' OR CAST(id_pengguna AS TEXT) ILIKE '%' || $1 || '%')",
	}
	args := []any{search}
	argPos := 2

	if jenis != "" {
		conditions = append(conditions, fmt.Sprintf("jenis = $%d", argPos))
		args = append(args, jenis)
		argPos++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM m_pengguna WHERE %s", strings.Join(conditions, " AND "))
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	listQuery := fmt.Sprintf(`
		SELECT id, uuid, id_pengguna, telegram, jenis, active, created_at, updated_at
		FROM m_pengguna
		WHERE %s
		ORDER BY id DESC
		LIMIT $%d OFFSET $%d
	`, strings.Join(conditions, " AND "), argPos, argPos+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]model.MasterPengguna, 0)
	for rows.Next() {
		var item model.MasterPengguna
		var telegram sql.NullString
		var updatedAt sql.NullTime
		if err := rows.Scan(
			&item.ID,
			&item.UUID,
			&item.IDPengguna,
			&telegram,
			&item.Jenis,
			&item.Active,
			&item.CreatedAt,
			&updatedAt,
		); err != nil {
			return nil, 0, err
		}
		if telegram.Valid {
			val := telegram.String
			item.Telegram = &val
		}
		if updatedAt.Valid {
			val := updatedAt.Time
			item.UpdatedAt = &val
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *MasterPenggunaRepository) Update(ctx context.Context, id int64, idPengguna *int64, telegram *string, jenis *string, active *bool) error {
	setClauses := make([]string, 0, 4)
	args := make([]any, 0, 6)
	argPos := 1

	if idPengguna != nil {
		setClauses = append(setClauses, fmt.Sprintf("id_pengguna = $%d", argPos))
		args = append(args, *idPengguna)
		argPos++
	}
	if telegram != nil {
		setClauses = append(setClauses, fmt.Sprintf("telegram = $%d", argPos))
		args = append(args, *telegram)
		argPos++
	}
	if jenis != nil {
		setClauses = append(setClauses, fmt.Sprintf("jenis = $%d", argPos))
		args = append(args, *jenis)
		argPos++
	}
	if active != nil {
		setClauses = append(setClauses, fmt.Sprintf("active = $%d", argPos))
		args = append(args, *active)
		argPos++
	}

	setClauses = append(setClauses, "updated_at = NOW()")

	args = append(args, id)
	query := fmt.Sprintf("UPDATE m_pengguna SET %s WHERE id = $%d", strings.Join(setClauses, ", "), argPos)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *MasterPenggunaRepository) SoftDelete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE m_pengguna
		SET deleted_at = NOW(), updated_at = NOW(), active = FALSE
		WHERE id = $1
	`, id)
	return err
}

func (r *MasterPenggunaRepository) CountByActive(ctx context.Context) (int, int, error) {
	var activeCount int
	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM m_pengguna
		WHERE deleted_at IS NULL AND active = TRUE
	`).Scan(&activeCount); err != nil {
		return 0, 0, err
	}

	var inactiveCount int
	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM m_pengguna
		WHERE deleted_at IS NULL AND active = FALSE
	`).Scan(&inactiveCount); err != nil {
		return 0, 0, err
	}

	return activeCount, inactiveCount, nil
}
