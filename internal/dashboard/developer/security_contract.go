package developer

// Permission adalah kapabilitas backend yang dapat diberikan melalui RBAC.
type Permission string

const (
	PermissionSystemRead         Permission = "system:read"
	PermissionSystemWrite        Permission = "system:write"
	PermissionSecurityRead       Permission = "security:read"
	PermissionSecurityScan       Permission = "security:scan"
	PermissionThreatRead         Permission = "threat:read"
	PermissionCredentialRead     Permission = "credential:read"
	PermissionCredentialRotate   Permission = "credential:rotate"
	PermissionDependencyRead     Permission = "dependency:read"
	PermissionDependencyFix      Permission = "dependency:fix"
	PermissionAPIRead            Permission = "api:read"
	PermissionDatabaseRead       Permission = "database:read"
	PermissionAuditRead          Permission = "audit:read"
	PermissionCodeRead           Permission = "code:read"
	PermissionCodeWrite          Permission = "code:write"
	PermissionPreviewRead        Permission = "preview:read"
	PermissionApprovalRead       Permission = "approval:read"
	PermissionApprovalDecide     Permission = "approval:decide"
	PermissionDeploymentRead     Permission = "deployment:read"
	PermissionDeploymentRun      Permission = "deployment:run"
	PermissionDeploymentRollback Permission = "deployment:rollback"
	PermissionAIKnowledgeRead    Permission = "ai:knowledge:read"
	PermissionAIResearch         Permission = "ai:research"
	PermissionPatchPreview       Permission = "patch:preview"
	PermissionAIApproval         Permission = "ai:approval"
)

// FeatureContract menghubungkan menu UI dengan permission backend.
type FeatureContract struct {
	Key              string
	Label            string
	Permission       Permission
	Risk             string
	RequiresApproval bool
}

var FeatureContracts = []FeatureContract{
	{Key: "ai-dashboard", Label: "AI Dashboard", Permission: PermissionSystemRead, Risk: "low"},
	{Key: "ai-knowledge", Label: "AI Knowledge", Permission: PermissionAIKnowledgeRead, Risk: "medium"},
	{Key: "learning-history", Label: "Learning History", Permission: PermissionAuditRead, Risk: "medium"},
	{Key: "new-security-intelligence", Label: "New Security Intelligence", Permission: PermissionSecurityRead, Risk: "high"},
	{Key: "application-analysis", Label: "Application Analysis", Permission: PermissionAIResearch, Risk: "high"},
	{Key: "ai-recommendations", Label: "AI Recommendations", Permission: PermissionAIResearch, Risk: "medium"},
	{Key: "security-score", Label: "Security Score", Permission: PermissionSecurityRead, Risk: "low"},
	{Key: "ai-security-checker", Label: "AI Security Checker", Permission: PermissionSecurityScan, Risk: "high", RequiresApproval: true},
	{Key: "threat-detection", Label: "Threat Detection", Permission: PermissionThreatRead, Risk: "high"},
	{Key: "vulnerability-scanner", Label: "Vulnerability Scanner", Permission: PermissionSecurityScan, Risk: "high", RequiresApproval: true},
	{Key: "dependency-scanner", Label: "Dependency Scanner", Permission: PermissionDependencyRead, Risk: "high"},
	{Key: "api-security", Label: "API Security", Permission: PermissionAPIRead, Risk: "high"},
	{Key: "database-security", Label: "Database Security", Permission: PermissionDatabaseRead, Risk: "critical"},
	{Key: "security-audit", Label: "Security Audit", Permission: PermissionAuditRead, Risk: "high"},
	{Key: "security-knowledge", Label: "Security Knowledge", Permission: PermissionAIKnowledgeRead, Risk: "medium"},
	{Key: "cve-intelligence", Label: "CVE Intelligence", Permission: PermissionAIKnowledgeRead, Risk: "high"},
	{Key: "owasp-rules", Label: "OWASP Rules", Permission: PermissionAIKnowledgeRead, Risk: "medium"},
	{Key: "dependency-intelligence", Label: "Dependency Intelligence", Permission: PermissionDependencyRead, Risk: "high"},
	{Key: "framework-updates", Label: "Framework Updates", Permission: PermissionAIKnowledgeRead, Risk: "medium"},
	{Key: "infrastructure-security", Label: "Infrastructure Security", Permission: PermissionSecurityRead, Risk: "critical"},
	{Key: "knowledge-version", Label: "Knowledge Version", Permission: PermissionAIKnowledgeRead, Risk: "medium"},
	{Key: "ai-research", Label: "AI Research", Permission: PermissionAIResearch, Risk: "high"},
	{Key: "generate-patch-preview", Label: "Generate Patch Preview", Permission: PermissionPatchPreview, Risk: "high", RequiresApproval: true},
	{Key: "preview", Label: "Preview", Permission: PermissionPreviewRead, Risk: "medium"},
	{Key: "security-fix-preview", Label: "Security Fix Preview", Permission: PermissionPatchPreview, Risk: "high", RequiresApproval: true},
	{Key: "database-migration-preview", Label: "Database Migration Preview", Permission: PermissionPreviewRead, Risk: "critical", RequiresApproval: true},
	{Key: "approval", Label: "Approval", Permission: PermissionApprovalRead, Risk: "critical"},
	{Key: "security-fix-approval", Label: "Security Fix", Permission: PermissionAIApproval, Risk: "critical", RequiresApproval: true},
	{Key: "production-deployment", Label: "Production Deployment", Permission: PermissionDeploymentRun, Risk: "critical", RequiresApproval: true},
	{Key: "audit-trail", Label: "Audit Trail", Permission: PermissionAuditRead, Risk: "high"},
	{Key: "developer-approval-audit", Label: "Developer Approval", Permission: PermissionAuditRead, Risk: "critical"},
	{Key: "rollback-audit", Label: "Rollback", Permission: PermissionAuditRead, Risk: "critical"},
}
