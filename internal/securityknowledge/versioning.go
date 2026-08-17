package securityknowledge

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (r *Repository) Activate(
	ctx context.Context,
	knowledgeKey string,
	version string,
) error {
	if r == nil || r.DB == nil {
		return errors.New("security knowledge repository is not configured")
	}

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin knowledge activation: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	var currentStatus string

	err = tx.QueryRowContext(
		ctx,
		`
SELECT status
FROM apb.security_knowledge_versions
WHERE knowledge_key = $1
  AND version = $2
FOR UPDATE
`,
		knowledgeKey,
		version,
	).Scan(&currentStatus)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}

	if err != nil {
		return fmt.Errorf("lock knowledge version: %w", err)
	}

	if currentStatus == string(StatusRevoked) {
		return fmt.Errorf("%w: revoked version", ErrNotActive)
	}

	_, err = tx.ExecContext(
		ctx,
		`
UPDATE apb.security_knowledge_versions
SET status = 'deprecated'
WHERE knowledge_key = $1
  AND status = 'active'
  AND version <> $2
`,
		knowledgeKey,
		version,
	)

	if err != nil {
		return fmt.Errorf("deprecate previous knowledge: %w", err)
	}

	result, err := tx.ExecContext(
		ctx,
		`
UPDATE apb.security_knowledge_versions
SET status = 'active'
WHERE knowledge_key = $1
  AND version = $2
  AND status = 'draft'
`,
		knowledgeKey,
		version,
	)

	if err != nil {
		return fmt.Errorf("activate knowledge: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check knowledge activation: %w", err)
	}

	if rows != 1 {
		return ErrNotActive
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit knowledge activation: %w", err)
	}

	return nil
}
