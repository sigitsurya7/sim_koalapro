package seed

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Users(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if os.Getenv("SEED_ENABLED") != "true" {
		return nil
	}

	var count int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM m_user").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	adminUsername, err := getRequiredEnv("SEED_ADMIN_USERNAME")
	if err != nil {
		return err
	}
	adminPassword, err := getRequiredEnv("SEED_ADMIN_PASSWORD")
	if err != nil {
		return err
	}
	viewerUsername, err := getRequiredEnv("SEED_VIEWER_USERNAME")
	if err != nil {
		return err
	}
	viewerPassword, err := getRequiredEnv("SEED_VIEWER_PASSWORD")
	if err != nil {
		return err
	}

	adminHash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	viewerHash, err := bcrypt.GenerateFromPassword([]byte(viewerPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, `
		INSERT INTO m_user (username, password, role, created_by)
		VALUES
			($1, $2, 'admin', 'seed'),
			($3, $4, 'viewer', 'seed')
	`, adminUsername, string(adminHash), viewerUsername, string(viewerHash))
	return err
}

func getRequiredEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("%s is required when SEED_ENABLED=true", key)
	}
	return val, nil
}
