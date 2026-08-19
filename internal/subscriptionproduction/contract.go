package subscriptionproduction

import "time"

type Stage string

const (
	StageDeveloper  Stage = "DEVELOPER"
	StagePreview    Stage = "PREVIEW"
	StageApproval   Stage = "APPROVAL"
	StageProduction Stage = "PRODUCTION"
)

type IntegrationDomain string

const (
	DomainCustomer IntegrationDomain = "CUSTOMER"
	DomainService  IntegrationDomain = "SERVICE"
	DomainBilling  IntegrationDomain = "BILLING"
	DomainRadius   IntegrationDomain = "RADIUS"
	DomainMikrotik IntegrationDomain = "MIKROTIK"
	DomainPayment  IntegrationDomain = "PAYMENT"
	DomainFinance  IntegrationDomain = "FINANCE"
)

type SubscriptionStatus string

const (
	StatusDraft           SubscriptionStatus = "DRAFT"
	StatusPreview         SubscriptionStatus = "PREVIEW"
	StatusPendingApproval SubscriptionStatus = "PENDING_APPROVAL"
	StatusActive          SubscriptionStatus = "ACTIVE"
	StatusInactive        SubscriptionStatus = "INACTIVE"
	StatusIsolated        SubscriptionStatus = "ISOLATED"
	StatusSuspended       SubscriptionStatus = "SUSPENDED"
)

type IntegrationBinding struct {
	Domain        IntegrationDomain `json:"domain"`
	Required      bool              `json:"required"`
	Readiness     string            `json:"readiness"`
	FailureMode   string            `json:"failure_mode"`
	SourceOfTruth string            `json:"source_of_truth"`
}

var ProductionIntegrationBindings = []IntegrationBinding{
	{Domain: DomainCustomer, Required: true, Readiness: "customer_id_validated", FailureMode: "block_activation", SourceOfTruth: "customer_service"},
	{Domain: DomainService, Required: true, Readiness: "service_profile_validated", FailureMode: "block_activation", SourceOfTruth: "subscription_service"},
	{Domain: DomainBilling, Required: true, Readiness: "billing_schedule_validated", FailureMode: "block_activation", SourceOfTruth: "billing_service"},
	{Domain: DomainRadius, Required: true, Readiness: "radius_mapping_ready", FailureMode: "quarantine_change", SourceOfTruth: "radius_controller"},
	{Domain: DomainMikrotik, Required: true, Readiness: "router_mapping_ready", FailureMode: "quarantine_change", SourceOfTruth: "network_controller"},
	{Domain: DomainPayment, Required: false, Readiness: "payment_method_available", FailureMode: "mark_payment_pending", SourceOfTruth: "payment_service"},
	{Domain: DomainFinance, Required: true, Readiness: "ledger_mapping_ready", FailureMode: "block_invoice_posting", SourceOfTruth: "finance_service"},
}

type SubscriptionChange struct {
	ID                string               `json:"id"`
	SubscriptionID    string               `json:"subscription_id"`
	Action            string               `json:"action"`
	RequestedBy       string               `json:"requested_by"`
	Stage             Stage                `json:"stage"`
	Status            string               `json:"status"`
	ExpectedVersion   int64                `json:"expected_version"`
	Preview           map[string]any       `json:"preview"`
	Integrations      []IntegrationBinding `json:"integrations"`
	ProductionChanged bool                 `json:"production_changed"`
	CreatedAt         time.Time            `json:"created_at"`
}

type ProductionPolicy struct {
	AIProductionAccess       bool     `json:"ai_production_access"`
	RequiresPreview          bool     `json:"requires_preview"`
	RequiresApproval         bool     `json:"requires_approval"`
	RequiresIntegrationCheck bool     `json:"requires_integration_check"`
	RequiredApproverRoles    []string `json:"required_approver_roles"`
	ProductionActions        []string `json:"production_actions"`
}

var DefaultProductionPolicy = ProductionPolicy{
	AIProductionAccess:       false,
	RequiresPreview:          true,
	RequiresApproval:         true,
	RequiresIntegrationCheck: true,
	RequiredApproverRoles:    []string{"administrator", "developer"},
	ProductionActions:        []string{"ACTIVATE", "SET_INACTIVE", "SET_ISOLATED", "CHANGE_ROUTER", "CHANGE_RADIUS", "CHANGE_BILLING", "POST_INVOICE", "POST_LEDGER"},
}

type ActivationReadiness struct {
	SubscriptionID    string               `json:"subscription_id"`
	Ready             bool                 `json:"ready"`
	Blockers          []string             `json:"blockers"`
	Bindings          []IntegrationBinding `json:"bindings"`
	ProductionChanged bool                 `json:"production_changed"`
}
