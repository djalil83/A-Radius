package developer

import (
	"context"
	"errors"
	"sync"
	"time"

	"crypto/rand"
	"encoding/hex"
)

var (
	ErrScanNotFound = errors.New("scan job not found")
	ErrInvalidScan  = errors.New("invalid scan type")
)

type ScanStatus string

const (
	ScanQueued    ScanStatus = "queued"
	ScanRunning   ScanStatus = "running"
	ScanCompleted ScanStatus = "completed"
	ScanFailed    ScanStatus = "failed"
)

type ScanJob struct {
	ID          string     `json:"scan_id"`
	Type        string     `json:"type"`
	Status      ScanStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Findings    []Finding  `json:"findings,omitempty"`
	Error       string     `json:"error,omitempty"`
}

type ScanRepository interface {
	Create(context.Context, *ScanJob) error
	Get(context.Context, string) (*ScanJob, error)
	Update(context.Context, *ScanJob) error
}

type MemoryScanRepository struct {
	mu   sync.RWMutex
	jobs map[string]*ScanJob
}

func NewMemoryScanRepository() *MemoryScanRepository {
	return &MemoryScanRepository{
		jobs: make(map[string]*ScanJob),
	}
}

func (r *MemoryScanRepository) Create(_ context.Context, job *ScanJob) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if job == nil || job.ID == "" {
		return ErrInvalidScan
	}

	copy := *job
	r.jobs[job.ID] = &copy
	return nil
}

func (r *MemoryScanRepository) Get(_ context.Context, id string) (*ScanJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	job, ok := r.jobs[id]
	if !ok {
		return nil, ErrScanNotFound
	}

	copy := *job
	copy.Findings = append([]Finding(nil), job.Findings...)
	return &copy, nil
}

func (r *MemoryScanRepository) Update(_ context.Context, job *ScanJob) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if job == nil || job.ID == "" {
		return ErrInvalidScan
	}

	if _, ok := r.jobs[job.ID]; !ok {
		return ErrScanNotFound
	}

	copy := *job
	copy.Findings = append([]Finding(nil), job.Findings...)
	r.jobs[job.ID] = &copy
	return nil
}

func NewScanJob(scanType string) (*ScanJob, error) {
	if scanType != "full" {
		return nil, ErrInvalidScan
	}

	return &ScanJob{
		ID:        newScanID(),
		Type:      scanType,
		Status:    ScanQueued,
		CreatedAt: time.Now().UTC(),
		Findings:  []Finding{},
	}, nil
}

func newScanID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "scan-unavailable"
	}
	return "scan-" + hex.EncodeToString(b)
}
