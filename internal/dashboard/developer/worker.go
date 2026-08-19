package developer

import (
	"context"
	"fmt"
	"time"
)

type Scanner interface {
	Scan(context.Context, ScanJob) ([]Finding, error)
}

type SafeScanner struct{}

func (SafeScanner) Scan(_ context.Context, job ScanJob) ([]Finding, error) {
	if job.Type != "full" {
		return nil, ErrInvalidScan
	}

	return []Finding{
		{
			ID:          fmt.Sprintf("%s-001", job.ID),
			RuleID:      "DEV-BASELINE-001",
			Title:       "Baseline security scan completed",
			Severity:    SeverityInfo,
			Description: "Initial developer security baseline completed.",
			Evidence:    "scanner=safe-baseline",
		},
	}, nil
}

type ScanWorker struct {
	Repository ScanRepository
	Scanner    Scanner
}

func (w *ScanWorker) Run(ctx context.Context, id string) error {
	job, err := w.Repository.Get(ctx, id)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	job.Status = ScanRunning
	job.StartedAt = &now

	if err := w.Repository.Update(ctx, job); err != nil {
		return err
	}

	findings, err := w.Scanner.Scan(ctx, *job)
	if err != nil {
		job.Status = ScanFailed
		job.Error = err.Error()
		_ = w.Repository.Update(ctx, job)
		return err
	}

	completed := time.Now().UTC()
	job.Status = ScanCompleted
	job.CompletedAt = &completed
	job.Findings = findings

	return w.Repository.Update(ctx, job)
}
