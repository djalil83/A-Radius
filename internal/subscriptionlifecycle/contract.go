package subscriptionlifecycle

import "time"

type Status string

const (
	StatusActive       Status = "ACTIVE"
	StatusWarning      Status = "WARNING"
	StatusIsolated     Status = "ISOLATED"
	StatusReactivating Status = "REACTIVATING"
	StatusInactive     Status = "INACTIVE"
)

type Action string

const (
	ActionDelete        Action = "DELETE"
	ActionSetInactive   Action = "SET_INACTIVE"
	ActionSetIsolated   Action = "SET_ISOLATED"
	ActionChangeRouter  Action = "CHANGE_ROUTER"
	ActionChangeProfile Action = "CHANGE_PROFILE"
	ActionChangeBilling Action = "CHANGE_BILLING"
)

type EventType string

const (
	EventInvoiceDue       EventType = "INVOICE_DUE"
	EventIsolationDue     EventType = "ISOLATION_DUE"
	EventPaymentReceived  EventType = "PAYMENT_RECEIVED"
	EventServiceActivated EventType = "SERVICE_ACTIVATED"
	EventProfileChanged   EventType = "PROFILE_CHANGED"
)

type LifecycleEvent struct {
	ID             string         `json:"id"`
	SubscriptionID string         `json:"subscription_id"`
	Type           EventType      `json:"type"`
	Source         string         `json:"source"`
	OccurredAt     time.Time      `json:"occurred_at"`
	IdempotencyKey string         `json:"idempotency_key"`
	Payload        map[string]any `json:"payload"`
}

type BulkActionProposal struct {
	ID               string            `json:"id"`
	Action           Action            `json:"action"`
	TargetCount      int               `json:"target_count"`
	TargetFilter     map[string]string `json:"target_filter"`
	From             string            `json:"from"`
	To               string            `json:"to"`
	Risk             string            `json:"risk"`
	Status           string            `json:"status"`
	AIRecommendation bool              `json:"ai_recommendation"`
	PreviewOnly      bool              `json:"preview_only"`
	ApprovalRequired bool              `json:"approval_required"`
	RequestedBy      string            `json:"requested_by"`
	ApprovedBy       *string           `json:"approved_by,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
}

type AutomationDecision struct {
	CurrentStatus  Status   `json:"current_status"`
	NextStatus     Status   `json:"next_status"`
	Reason         string   `json:"reason"`
	Effects        []string `json:"effects"`
	RequiresWorker bool     `json:"requires_worker"`
}

var ApprovalRequiredActions = map[Action]bool{
	ActionDelete: true, ActionSetInactive: true, ActionSetIsolated: true,
	ActionChangeRouter: true, ActionChangeProfile: true, ActionChangeBilling: true,
}

func Decide(event LifecycleEvent, status Status) AutomationDecision {
	switch event.Type {
	case EventInvoiceDue:
		return AutomationDecision{CurrentStatus: status, NextStatus: StatusWarning, Reason: "invoice is due", Effects: []string{"billing.warning", "customer.notification"}}
	case EventIsolationDue:
		return AutomationDecision{CurrentStatus: status, NextStatus: StatusIsolated, Reason: "isolation threshold reached", Effects: []string{"radius.block", "mikrotik.block", "service.block"}, RequiresWorker: true}
	case EventPaymentReceived:
		return AutomationDecision{CurrentStatus: status, NextStatus: StatusReactivating, Reason: "payment received", Effects: []string{"payment.reconcile", "radius.reactivate", "mikrotik.reactivate"}, RequiresWorker: true}
	case EventServiceActivated:
		return AutomationDecision{CurrentStatus: status, NextStatus: StatusActive, Reason: "service activated", Effects: []string{"service.on", "radius.enable", "customer_app.refresh"}}
	default:
		return AutomationDecision{CurrentStatus: status, NextStatus: status, Reason: "no automatic transition", Effects: []string{}}
	}
}
