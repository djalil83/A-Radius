package securityknowledge

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
)

type Loader struct {
	Repository *Repository
}

func NewLoader(repository *Repository) *Loader {
	return &Loader{
		Repository: repository,
	}
}

func (l *Loader) LoadManifestFile(
	ctx context.Context,
	path string,
	source string,
) (*KnowledgeVersion, error) {
	if l == nil || l.Repository == nil {
		return nil, errors.New("security knowledge loader is not configured")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read knowledge manifest: %w", err)
	}

	manifest, err := ValidateManifest(data)
	if err != nil {
		return nil, fmt.Errorf("validate knowledge manifest: %w", err)
	}

	sum := sha256.Sum256(data)
	contentHash := fmt.Sprintf("%x", sum)

	knowledge, err := l.Repository.CreateDraft(
		ctx,
		manifest.KnowledgeKey,
		manifest.Version,
		contentHash,
		source,
	)
	if err != nil {
		return nil, fmt.Errorf("create knowledge draft: %w", err)
	}

	return knowledge, nil
}
