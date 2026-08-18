package developer

import "time"

// KnowledgeSource identifies a trusted external security source.
type KnowledgeSource string

const (
	SourceCISAKEV        KnowledgeSource = "cisa_kev"
	SourceOWASPAPI       KnowledgeSource = "owasp_api_security"
	SourceCSAF           KnowledgeSource = "oasis_csaf"
	SourceDependency     KnowledgeSource = "dependency_advisory"
	SourceVendorAdvisory KnowledgeSource = "vendor_advisory"
)

// KnowledgeItem is immutable evidence ingested from an external source.
type KnowledgeItem struct {
	ID            string          `json:"id"`
	Source        KnowledgeSource `json:"source"`
	SourceURL     string          `json:"source_url"`
	AdvisoryID    string          `json:"advisory_id,omitempty"`
	Title         string          `json:"title"`
	ContentHash   string          `json:"content_hash"`
	RetrievedAt   time.Time       `json:"retrieved_at"`
	PublishedAt   *time.Time      `json:"published_at,omitempty"`
	ParserVersion string          `json:"parser_version"`
	Validation    string          `json:"validation"`
	Confidence    string          `json:"confidence"`
}

// SecurityAnalysis links validated knowledge with an A-RADIUS target.
type SecurityAnalysis struct {
	ID           string   `json:"id"`
	Target       string   `json:"target"`
	Components   []string `json:"components"`
	KnowledgeIDs []string `json:"knowledge_ids"`
	Status       string   `json:"status"`
}

// ContinuousFinding is advisory output; it never grants deployment authority.
type ContinuousFinding struct {
	ID             string   `json:"id"`
	Severity       string   `json:"severity"`
	Category       string   `json:"category"`
	Title          string   `json:"title"`
	Module         string   `json:"module"`
	Evidence       []string `json:"evidence"`
	Recommendation string   `json:"recommendation"`
	Production     string   `json:"production"`
	Status         string   `json:"status"`
}

// PatchPreview is a proposed change that must be reviewed by a human.
type PatchPreview struct {
	ID             string `json:"id"`
	FindingID      string `json:"finding_id"`
	Before         string `json:"before"`
	After          string `json:"after"`
	Environment    string `json:"environment"`
	Production     string `json:"production"`
	ApprovalStatus string `json:"approval_status"`
}

// ContinuousSecurityPolicy is the fail-closed boundary around production.
type ContinuousSecurityPolicy struct {
	AIProductionAccess  bool     `json:"ai_production_access"`
	HumanApproval       bool     `json:"human_approval"`
	RequiredStages      []string `json:"required_stages"`
	AllowedEnvironments []string `json:"allowed_environments"`
}

var DefaultContinuousSecurityPolicy = ContinuousSecurityPolicy{
	AIProductionAccess: false,
	HumanApproval:      true,
	RequiredStages: []string{
		"developer_approval", "automated_security_test", "staging_test", "production_review", "health_check",
	},
	AllowedEnvironments: []string{"sandbox", "staging"},
}

var TrustedKnowledgeSources = []KnowledgeSource{
	SourceCISAKEV, SourceOWASPAPI, SourceCSAF, SourceDependency, SourceVendorAdvisory,
}

var FeaturedContinuousFinding = ContinuousFinding{
	ID: "SEC-2026-00127", Severity: "HIGH", Category: "broken_access_control",
	Title:  "Authorization bypass pada API langganan",
	Module: "Langganan API", Evidence: []string{"authorization policy tidak konsisten pada /api/langganan/update"},
	Recommendation: "Tambahkan RBAC middleware dan uji function-level authorization.", Production: "UNCHANGED", Status: "open",
}

var FeaturedPatchPreview = PatchPreview{
	ID: "FIX-2026-00421", FindingID: "SEC-2026-00127", Before: "authorization -> partial",
	After: "authorization -> RBAC middleware", Environment: "developer_preview", Production: "UNCHANGED", ApprovalStatus: "pending",
}
