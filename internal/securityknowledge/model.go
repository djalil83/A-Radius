package securityknowledge

import (
	"errors"
	"time"
)

var (
	ErrInvalidKnowledgeKey = errors.New("invalid knowledge key")
	ErrInvalidVersion      = errors.New("invalid knowledge version")
	ErrInvalidHash         = errors.New("invalid content hash")
	ErrInvalidStatus       = errors.New("invalid knowledge status")
	ErrNotFound            = errors.New("knowledge version not found")
	ErrAlreadyExists       = errors.New("knowledge version already exists")
	ErrNotActive           = errors.New("knowledge version is not active")
)

type Status string

const (
	StatusDraft      Status = "draft"
	StatusActive     Status = "active"
	StatusDeprecated Status = "deprecated"
	StatusRevoked    Status = "revoked"
)

type KnowledgeVersion struct {
	ID           string
	KnowledgeKey string
	Version      string
	ContentHash  string
	Source       string
	Status       Status
	LearnedAt    time.Time
	CreatedAt    time.Time
}

func (s Status) Valid() bool {
	switch s {
	case StatusDraft, StatusActive, StatusDeprecated, StatusRevoked:
		return true
	default:
		return false
	}
}
