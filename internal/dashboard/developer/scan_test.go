package developer

import (
	"context"
	"testing"
)

func TestScanLifecycle(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryScanRepository()

	job, err := NewScanJob("full")
	if err != nil {
		t.Fatal(err)
	}

	if job.Status != ScanQueued {
		t.Fatalf("expected queued, got %s", job.Status)
	}

	if err := repo.Create(ctx, job); err != nil {
		t.Fatal(err)
	}

	worker := ScanWorker{
		Repository: repo,
		Scanner:    SafeScanner{},
	}

	if err := worker.Run(ctx, job.ID); err != nil {
		t.Fatal(err)
	}

	result, err := repo.Get(ctx, job.ID)
	if err != nil {
		t.Fatal(err)
	}

	if result.Status != ScanCompleted {
		t.Fatalf("expected completed, got %s", result.Status)
	}

	if len(result.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(result.Findings))
	}
}

func TestInvalidScanType(t *testing.T) {
	if _, err := NewScanJob("arbitrary-command"); err != ErrInvalidScan {
		t.Fatalf("expected ErrInvalidScan, got %v", err)
	}
}
