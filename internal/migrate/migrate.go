package migrate

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

func Run(db *sql.DB, path string) error {
	paths := candidatePaths(path)
	var lastErr error

	for _, p := range paths {
		if p == "" {
			continue
		}
		sqlBytes, err := os.ReadFile(p)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				lastErr = err
				continue
			}
			return err
		}
		if len(sqlBytes) == 0 {
			return nil
		}
		_, err = db.Exec(string(sqlBytes))
		return err
	}

	if lastErr != nil {
		return fmt.Errorf("migration file not found: %w", lastErr)
	}
	return nil
}

func candidatePaths(path string) []string {
	if path != "" {
		return []string{path, "/db/schema.sql", "db/schema.sql"}
	}
	return []string{"/db/schema.sql", "db/schema.sql"}
}
