package developer

// SecurityLayer membedakan kontrol preventif, detektif, dan respons.
type SecurityLayer string

const (
	LayerPrevention SecurityLayer = "prevention"
	LayerDetection  SecurityLayer = "detection"
	LayerResponse   SecurityLayer = "response"
)

// SecurityControl adalah katalog kontrol yang dapat dipantau lintas aplikasi.
type SecurityControl struct {
	Key              string        `json:"key"`
	Name             string        `json:"name"`
	Layer            SecurityLayer `json:"layer"`
	OwnerRoles       []string      `json:"owner_roles"`
	Status           string        `json:"status"`
	RequiresApproval bool          `json:"requires_approval"`
}

var CrossApplicationControls = []SecurityControl{
	{Key: "rbac", Name: "RBAC", Layer: LayerPrevention, OwnerRoles: []string{"developer", "administrator"}, Status: "active"},
	{Key: "mfa", Name: "MFA", Layer: LayerPrevention, OwnerRoles: []string{"developer", "administrator"}, Status: "required"},
	{Key: "session-security", Name: "Session security", Layer: LayerPrevention, OwnerRoles: []string{"developer", "administrator"}, Status: "monitored"},
	{Key: "rate-limiting", Name: "Rate limiting", Layer: LayerPrevention, OwnerRoles: []string{"developer"}, Status: "active"},
	{Key: "input-validation", Name: "Input validation", Layer: LayerPrevention, OwnerRoles: []string{"developer"}, Status: "enforced"},
	{Key: "secret-management", Name: "Secret management", Layer: LayerPrevention, OwnerRoles: []string{"developer"}, Status: "vault-backed"},
	{Key: "api-authorization", Name: "API authorization", Layer: LayerPrevention, OwnerRoles: []string{"developer", "administrator"}, Status: "enforced"},
	{Key: "ai-security-checker", Name: "AI Security Checker", Layer: LayerDetection, OwnerRoles: []string{"developer"}, Status: "ready"},
	{Key: "anomaly-detection", Name: "Anomaly detection", Layer: LayerDetection, OwnerRoles: []string{"developer"}, Status: "monitoring"},
	{Key: "login-api-monitoring", Name: "Login/API monitoring", Layer: LayerDetection, OwnerRoles: []string{"developer", "administrator"}, Status: "monitoring"},
	{Key: "change-monitoring", Name: "File, database, permission changes", Layer: LayerDetection, OwnerRoles: []string{"developer", "administrator"}, Status: "audited"},
	{Key: "alert", Name: "Alert", Layer: LayerResponse, OwnerRoles: []string{"developer", "administrator"}, Status: "proposed"},
	{Key: "quarantine-lock", Name: "Quarantine / lock", Layer: LayerResponse, OwnerRoles: []string{"developer", "administrator"}, Status: "approval_required", RequiresApproval: true},
	{Key: "rollback", Name: "Rollback", Layer: LayerResponse, OwnerRoles: []string{"developer", "administrator"}, Status: "approval_required", RequiresApproval: true},
	{Key: "disable-credential", Name: "Disable credential", Layer: LayerResponse, OwnerRoles: []string{"developer", "administrator"}, Status: "approval_required", RequiresApproval: true},
	{Key: "incident-report", Name: "Incident report", Layer: LayerResponse, OwnerRoles: []string{"developer", "administrator"}, Status: "available"},
	{Key: "audit-trail", Name: "Audit trail", Layer: LayerResponse, OwnerRoles: []string{"developer", "administrator"}, Status: "append_only"},
}

// DisruptiveActionPolicy mencegah AI atau service scanner mengeksekusi tindakan
// yang dapat mengganggu pelanggan tanpa human approval yang terverifikasi.
type DisruptiveActionPolicy struct {
	Action               string   `json:"action"`
	AllowedProposers     []string `json:"allowed_proposers"`
	RequiredApprovers    []string `json:"required_approvers"`
	RequiresChangeTicket bool     `json:"requires_change_ticket"`
	RequiresAuditEvent   bool     `json:"requires_audit_event"`
}

var DisruptiveActions = []DisruptiveActionPolicy{
	{Action: "block_account", AllowedProposers: []string{"developer", "administrator"}, RequiredApprovers: []string{"developer", "administrator"}, RequiresChangeTicket: true, RequiresAuditEvent: true},
	{Action: "disconnect_api", AllowedProposers: []string{"developer", "administrator"}, RequiredApprovers: []string{"developer", "administrator"}, RequiresChangeTicket: true, RequiresAuditEvent: true},
	{Action: "change_firewall", AllowedProposers: []string{"developer", "administrator"}, RequiredApprovers: []string{"developer", "administrator"}, RequiresChangeTicket: true, RequiresAuditEvent: true},
	{Action: "delete_credential", AllowedProposers: []string{"developer", "administrator"}, RequiredApprovers: []string{"developer", "administrator"}, RequiresChangeTicket: true, RequiresAuditEvent: true},
	{Action: "deploy", AllowedProposers: []string{"developer", "administrator"}, RequiredApprovers: []string{"developer", "administrator"}, RequiresChangeTicket: true, RequiresAuditEvent: true},
}
