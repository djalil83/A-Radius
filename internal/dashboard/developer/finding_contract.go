package developer

// SecurityFinding merepresentasikan hasil scan yang belum boleh mengubah Production.
type SecurityFinding struct {
	ID             string        `json:"id"`
	Severity       string        `json:"severity"`
	Module         string        `json:"module"`
	Endpoint       string        `json:"endpoint"`
	Problem        string        `json:"problem"`
	Recommendation string        `json:"recommendation"`
	AffectedRoles  []string      `json:"affected_roles"`
	Production     string        `json:"production"`
	Stage          WorkflowStage `json:"stage"`
}

// FeaturedAuthorizationFinding adalah contoh finding yang ditampilkan pada UI.
var FeaturedAuthorizationFinding = SecurityFinding{
	ID:             "SEC-2026-00127",
	Severity:       "HIGH",
	Module:         "Langganan API",
	Endpoint:       "/api/langganan/update",
	Problem:        "Authorization check tidak konsisten.",
	Recommendation: "Tambahkan role-based authorization middleware.",
	AffectedRoles:  []string{"Administrator", "Technician", "Sales", "Reseller"},
	Production:     "UNCHANGED",
	Stage:          StageFinding,
}

// ProposedFixPreview menunjukkan diff usulan tanpa memberi akses deployment kepada AI.
type ProposedFixPreview struct {
	FindingID  string        `json:"finding_id"`
	Before     string        `json:"before"`
	After      string        `json:"after"`
	Stage      WorkflowStage `json:"stage"`
	Production string        `json:"production"`
}

var FeaturedAuthorizationPreview = ProposedFixPreview{
	FindingID:  "SEC-2026-00127",
	Before:     "authorization → partial",
	After:      "authorization → RBAC middleware",
	Stage:      StageDeveloperPreview,
	Production: "UNCHANGED",
}
