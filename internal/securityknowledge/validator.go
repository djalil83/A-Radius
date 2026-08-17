package securityknowledge

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
)

var (
	knowledgeKeyPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,149}$`)
	versionPattern      = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
)

type Manifest struct {
	KnowledgeKey string         `json:"knowledge_key"`
	Version      string         `json:"version"`
	Status       string         `json:"status"`
	Scope        []string       `json:"scope"`
	Rules        map[string]any `json:"rules"`
}

func ValidateManifest(data []byte) (*Manifest, error) {
	var manifest Manifest

	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}

	if !knowledgeKeyPattern.MatchString(manifest.KnowledgeKey) {
		return nil, ErrInvalidKnowledgeKey
	}

	if !versionPattern.MatchString(manifest.Version) {
		return nil, ErrInvalidVersion
	}

	if manifest.Status == "" {
		return nil, ErrInvalidStatus
	}

	if len(manifest.Scope) == 0 {
		return nil, fmt.Errorf("knowledge scope must not be empty")
	}

	if manifest.Rules == nil {
		return nil, fmt.Errorf("knowledge rules must not be empty")
	}

	return &manifest, nil
}

func ValidateHash(hash string) bool {
	if len(hash) != 64 {
		return false
	}

	_, err := hex.DecodeString(hash)
	return err == nil
}
