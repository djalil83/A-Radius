package subscriptionlifecycle

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const QueueName = "a-radius:subscription:lifecycle"

type Executor interface {
	Execute(ctx context.Context, proposal BulkActionProposal) error
}

type AuditFunc func(ctx context.Context, proposal BulkActionProposal, outcome string, err error)

type Worker struct {
	Redis   *redis.Client
	Execute Executor
	Audit   AuditFunc
}

func (w *Worker) Enqueue(ctx context.Context, proposal BulkActionProposal) error {
	if proposal.ApprovalRequired && proposal.ApprovedBy == nil {
		return fmt.Errorf("proposal %s requires approval", proposal.ID)
	}
	payload, err := json.Marshal(proposal)
	if err != nil {
		return err
	}
	return w.Redis.RPush(ctx, QueueName, payload).Err()
}

func (w *Worker) Run(ctx context.Context) error {
	for {
		result, err := w.Redis.BLPop(ctx, 5*time.Second, QueueName).Result()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			return err
		}
		var proposal BulkActionProposal
		if err := json.Unmarshal([]byte(result[1]), &proposal); err != nil {
			return err
		}
		if proposal.ApprovalRequired && proposal.ApprovedBy == nil {
			if w.Audit != nil {
				w.Audit(ctx, proposal, "REJECTED_NO_APPROVAL", fmt.Errorf("approval required"))
			}
			continue
		}
		err = w.Execute.Execute(ctx, proposal)
		if w.Audit != nil {
			outcome := "SUCCESS"
			if err != nil {
				outcome = "FAILED"
			}
			w.Audit(ctx, proposal, outcome, err)
		}
	}
}
