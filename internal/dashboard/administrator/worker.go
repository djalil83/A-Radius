package administrator

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const QueueName = "a-radius:administrator:approved-actions"

type Executor interface {
	Execute(context.Context, ActionProposal) error
}
type AuditFunc func(context.Context, ActionProposal, string, error)

type Worker struct {
	Redis   *redis.Client
	Execute Executor
	Audit   AuditFunc
}

func (w *Worker) EnqueueApproved(ctx context.Context, proposal ActionProposal) error {
	if proposal.Status != StatusApproved {
		return fmt.Errorf("proposal %s is not approved", proposal.ID)
	}
	if proposal.ProductionChanged {
		return fmt.Errorf("proposal %s cannot change production directly", proposal.ID)
	}
	if w.Redis == nil {
		return fmt.Errorf("redis worker is not configured")
	}
	payload, err := json.Marshal(proposal)
	if err != nil {
		return err
	}
	return w.Redis.RPush(ctx, QueueName, payload).Err()
}

func (w *Worker) Run(ctx context.Context) error {
	if w.Redis == nil || w.Execute == nil {
		return fmt.Errorf("worker dependencies are not configured")
	}
	for {
		item, err := w.Redis.BLPop(ctx, 5*time.Second, QueueName).Result()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			return err
		}
		var proposal ActionProposal
		if err := json.Unmarshal([]byte(item[1]), &proposal); err != nil {
			return err
		}
		if proposal.Status != StatusApproved || proposal.ProductionChanged {
			if w.Audit != nil {
				w.Audit(ctx, proposal, "REJECTED_WORKER_GUARD", fmt.Errorf("approval or production guard failed"))
			}
			continue
		}
		err = w.Execute.Execute(ctx, proposal)
		if w.Audit != nil {
			outcome := "SUCCEEDED"
			if err != nil {
				outcome = "FAILED"
			}
			w.Audit(ctx, proposal, outcome, err)
		}
	}
}
