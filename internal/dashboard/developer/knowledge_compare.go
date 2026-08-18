package developer

import "time"

type KnowledgeChangeType string

const (
	KnowledgeChangeAdded   KnowledgeChangeType = "ADDED"
	KnowledgeChangeUpdated KnowledgeChangeType = "UPDATED"
	KnowledgeChangeRemoved KnowledgeChangeType = "REMOVED"
)

type KnowledgeRuleChange struct {
	RuleCode string              `json:"rule_code"`
	Type     KnowledgeChangeType `json:"type"`
	Title    string              `json:"title,omitempty"`
}

type KnowledgeModuleImpact struct {
	Module string `json:"module"`
	Risk   string `json:"risk"`
}

type KnowledgeVersionComparison struct {
	FromVersion   string                  `json:"from_version"`
	ToVersion     string                  `json:"to_version"`
	Added         []KnowledgeRuleChange   `json:"added"`
	Updated       []KnowledgeRuleChange   `json:"updated"`
	Removed       []KnowledgeRuleChange   `json:"removed"`
	Affected      []KnowledgeModuleImpact `json:"affected_modules"`
	Production    string                  `json:"production"`
	PatchRequired bool                    `json:"patch_required"`
}

var FeaturedKnowledgeComparison = KnowledgeVersionComparison{
	FromVersion: "SK-2.4.7",
	ToVersion:   "SK-2.4.8",
	Added: []KnowledgeRuleChange{
		{RuleCode: "API-AUTH-042", Type: KnowledgeChangeAdded},
		{RuleCode: "SESSION-019", Type: KnowledgeChangeAdded},
		{RuleCode: "DEP-031", Type: KnowledgeChangeAdded},
	},
	Updated: []KnowledgeRuleChange{
		{RuleCode: "API-AUTH-018", Type: KnowledgeChangeUpdated},
		{RuleCode: "DATABASE-007", Type: KnowledgeChangeUpdated},
	},
	Removed: []KnowledgeRuleChange{
		{RuleCode: "API-AUTH-003", Type: KnowledgeChangeRemoved},
	},
	Affected: []KnowledgeModuleImpact{
		{Module: "Administrator", Risk: "HIGH"},
		{Module: "Langganan", Risk: "HIGH"},
		{Module: "Pelanggan", Risk: "MEDIUM"},
		{Module: "Technician", Risk: "MEDIUM"},
		{Module: "Sales", Risk: "LOW"},
		{Module: "Customer", Risk: "LOW"},
	},
	Production:    "UNCHANGED",
	PatchRequired: true,
}

type KnowledgePatchPipeline struct {
	KnowledgeVersion  string   `json:"knowledge_version"`
	Stages            []string `json:"stages"`
	ProductionChanged bool     `json:"production_changed"`
}

var DefaultKnowledgePatchPipeline = KnowledgePatchPipeline{
	KnowledgeVersion:  "SK-2.4.8",
	Stages:            []string{"SECURITY_KNOWLEDGE", "AI_ANALYSIS", "FINDING", "RECOMMENDATION", "PATCH_PROPOSAL", "DEVELOPER_PREVIEW", "APPROVAL", "STAGING", "PRODUCTION"},
	ProductionChanged: false,
}

type KnowledgeRollbackRequest struct {
	FromVersion string `json:"from_version"`
	ToVersion   string `json:"to_version"`
	Reason      string `json:"reason"`
	RequestedBy string `json:"requested_by"`
	ApprovedBy  string `json:"approved_by,omitempty"`
	Status      string `json:"status"`
}

type KnowledgeRollbackAudit struct {
	Action      string    `json:"action"`
	FromVersion string    `json:"from_version"`
	ToVersion   string    `json:"to_version"`
	Reason      string    `json:"reason"`
	RequestedBy string    `json:"requested_by"`
	ApprovedBy  string    `json:"approved_by"`
	Timestamp   time.Time `json:"timestamp"`
}

var FeaturedKnowledgeRollbackAudit = KnowledgeRollbackAudit{
	Action: "KNOWLEDGE_ROLLBACK", FromVersion: "SK-2.4.8", ToVersion: "SK-2.4.7", Reason: "Compatibility issue", RequestedBy: "Developer", ApprovedBy: "Developer", Timestamp: time.Date(2026, time.August, 17, 14, 42, 0, 0, time.FixedZone("WITA", 8*60*60)),
}
