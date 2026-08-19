package migrations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Runner struct {
	DB  *sql.DB
	Dir string
}

func NewRunner(db *sql.DB, dir string) *Runner {
	return &Runner{
		DB:  db,
		Dir: dir,
	}
}

func (r *Runner) Run(ctx context.Context) error {
	if r == nil || r.DB == nil {
		return errors.New("migration runner database is not configured")
	}

	if strings.TrimSpace(r.Dir) == "" {
		return errors.New("migration directory is not configured")
	}

	files, err := migrationFiles(r.Dir)
	if err != nil {
		return err
	}

	if err := ensureMigrationTable(ctx, r.DB); err != nil {
		return err
	}

	applied, err := appliedVersions(ctx, r.DB)
	if err != nil {
		return err
	}

	for _, file := range files {
		if _, ok := applied[file.Version]; ok {
			continue
		}

		if err := r.apply(ctx, file); err != nil {
			return fmt.Errorf(
				"apply migration %s: %w",
				file.Version,
				err,
			)
		}
	}

	return nil
}

type migrationFile struct {
	Version string
	Path    string
}

func migrationFiles(dir string) ([]migrationFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf(
			"read migration directory: %w",
			err,
		)
	}

	files := make([]migrationFile, 0)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		if !strings.HasSuffix(name, ".sql") {
			continue
		}

		version := strings.TrimSuffix(name, ".sql")

		if strings.TrimSpace(version) == "" {
			return nil, fmt.Errorf(
				"invalid migration filename %q",
				name,
			)
		}

		files = append(files, migrationFile{
			Version: version,
			Path:    filepath.Join(dir, name),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Version < files[j].Version
	})

	seen := make(map[string]struct{}, len(files))

	for _, file := range files {
		if _, exists := seen[file.Version]; exists {
			return nil, fmt.Errorf(
				"duplicate migration version %q",
				file.Version,
			)
		}

		seen[file.Version] = struct{}{}
	}

	return files, nil
}

func ensureMigrationTable(
	ctx context.Context,
	db *sql.DB,
) error {
	const query = `
CREATE SCHEMA IF NOT EXISTS apb;

CREATE TABLE IF NOT EXISTS apb.schema_migrations (
version VARCHAR(100) PRIMARY KEY,
applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
`

	if _, err := db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf(
			"ensure schema migrations table: %w",
			err,
		)
	}

	return nil
}

func appliedVersions(
	ctx context.Context,
	db *sql.DB,
) (map[string]struct{}, error) {
	rows, err := db.QueryContext(
		ctx,
		`SELECT version FROM apb.schema_migrations`,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"query applied migrations: %w",
			err,
		)
	}
	defer rows.Close()

	result := make(map[string]struct{})

	for rows.Next() {
		var version string

		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf(
				"scan migration version: %w",
				err,
			)
		}

		result[version] = struct{}{}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate migration versions: %w",
			err,
		)
	}

	return result, nil
}

func (r *Runner) apply(
	ctx context.Context,
	file migrationFile,
) error {
	data, err := os.ReadFile(file.Path)
	if err != nil {
		return fmt.Errorf(
			"read migration file: %w",
			err,
		)
	}

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf(
			"begin migration transaction: %w",
			err,
		)
	}

	committed := false

	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	if _, err := tx.ExecContext(ctx, string(data)); err != nil {
		return fmt.Errorf(
			"execute migration SQL: %w",
			err,
		)
	}

	if _, err := tx.ExecContext(
		ctx,
		`
INSERT INTO apb.schema_migrations(version)
VALUES ($1)
ON CONFLICT (version) DO NOTHING
`,
		file.Version,
	); err != nil {
		return fmt.Errorf(
			"record migration %s: %w",
			file.Version,
			err,
		)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf(
			"commit migration %s: %w",
			file.Version,
			err,
		)
	}

	committed = true

	return nil
}
