package securityknowledge

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CreateDraft(
	ctx context.Context,
	knowledgeKey string,
	version string,
	contentHash string,
	source string,
) (*KnowledgeVersion, error) {
	if r == nil || r.DB == nil {
		return nil, errors.New("security knowledge repository is not configured")
	}

	if !knowledgeKeyPattern.MatchString(knowledgeKey) {
		return nil, ErrInvalidKnowledgeKey
	}

	if !versionPattern.MatchString(version) {
		return nil, ErrInvalidVersion
	}

	if !ValidateHash(contentHash) {
		return nil, ErrInvalidHash
	}

	const query = `
INSERT INTO apb.security_knowledge_versions
(knowledge_key, version, content_hash, source, status)
VALUES
($1, $2, $3, $4, 'draft')
RETURNING
id,
knowledge_key,
version,
content_hash,
source,
status,
learned_at,
created_at
`

	var k KnowledgeVersion

	err := r.DB.QueryRowContext(
		ctx,
		query,
		knowledgeKey,
		version,
		contentHash,
		source,
	).Scan(
		&k.ID,
		&k.KnowledgeKey,
		&k.Version,
		&k.ContentHash,
		&k.Source,
		&k.Status,
		&k.LearnedAt,
		&k.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("create knowledge draft: %w", err)
	}

	return &k, nil
}

func (r *Repository) GetActive(
	ctx context.Context,
	knowledgeKey string,
) (*KnowledgeVersion, error) {
	if r == nil || r.DB == nil {
		return nil, errors.New("security knowledge repository is not configured")
	}

	const query = `
SELECT
id,
knowledge_key,
version,
content_hash,
source,
status,
learned_at,
created_at
FROM apb.security_knowledge_versions
WHERE knowledge_key = $1
  AND status = 'active'
ORDER BY created_at DESC
LIMIT 1
`

	var k KnowledgeVersion

	err := r.DB.QueryRowContext(
		ctx,
		query,
		knowledgeKey,
	).Scan(
		&k.ID,
		&k.KnowledgeKey,
		&k.Version,
		&k.ContentHash,
		&k.Source,
		&k.Status,
		&k.LearnedAt,
		&k.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get active knowledge: %w", err)
	}

	return &k, nil
}
