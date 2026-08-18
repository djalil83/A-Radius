package developer

// FindingCategory adalah kategori pemeriksaan yang dapat dijalankan scanner.
type FindingCategory string

const (
	CategorySQLInjection            FindingCategory = "sql_injection"
	CategoryXSS                     FindingCategory = "xss"
	CategoryCSRF                    FindingCategory = "csrf"
	CategoryBrokenAccessControl     FindingCategory = "broken_access_control"
	CategoryAuthenticationWeakness  FindingCategory = "authentication_weakness"
	CategoryInsecureSecrets         FindingCategory = "insecure_secrets"
	CategoryDependencyVulnerability FindingCategory = "dependency_vulnerability"
	CategoryAPIExposure             FindingCategory = "api_exposure"
	CategoryDangerousConfiguration  FindingCategory = "dangerous_configuration"
	CategoryPrivilegeEscalation     FindingCategory = "privilege_escalation"
)

// WorkflowStage mengidentifikasi state yang harus dilalui proposed fix.
type WorkflowStage string

const (
	StageFinding            WorkflowStage = "finding"
	StageRecommendation     WorkflowStage = "ai_recommendation"
	StageProposedFix        WorkflowStage = "proposed_fix"
	StageDeveloperPreview   WorkflowStage = "developer_preview"
	StageRejected           WorkflowStage = "rejected"
	StageApproved           WorkflowStage = "approved"
	StageSecurityTest       WorkflowStage = "security_test"
	StageStagingTest        WorkflowStage = "staging_test"
	StageProductionApproval WorkflowStage = "production_approval"
	StageDeploy             WorkflowStage = "deploy"
	StageHealthCheck        WorkflowStage = "health_check"
	StageSuccess            WorkflowStage = "success"
	StageRollback           WorkflowStage = "rollback"
)

// AIProductionPolicy adalah guard eksplisit untuk memastikan AI hanya bersifat advisory.
type AIProductionPolicy struct {
	AllowDirectProductionAccess bool
	RequiresHumanApproval       bool
	RequiredStages              []WorkflowStage
}

var DefaultAIProductionPolicy = AIProductionPolicy{
	AllowDirectProductionAccess: false,
	RequiresHumanApproval:       true,
	RequiredStages: []WorkflowStage{
		StageDeveloperPreview,
		StageSecurityTest,
		StageStagingTest,
		StageProductionApproval,
		StageHealthCheck,
	},
}
