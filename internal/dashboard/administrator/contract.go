package administrator

import "time"

type Module string

const (
	ModuleVoucher Module = "VOUCHER"
	ModuleMitra   Module = "MITRA"
	ModuleBilling Module = "BILLING"
	ModulePayment Module = "PAYMENT"
	ModuleNMS     Module = "NMS"
	ModuleTeknisi Module = "TEKNISI"
	ModuleFinance Module = "FINANCE"
	ModuleSistem  Module = "SISTEM"
)

type RiskLevel string

const (
	RiskInfo     RiskLevel = "INFO"
	RiskLow      RiskLevel = "LOW"
	RiskMedium   RiskLevel = "MEDIUM"
	RiskHigh     RiskLevel = "HIGH"
	RiskCritical RiskLevel = "CRITICAL"
)

type ProposalStatus string

const (
	StatusPreview         ProposalStatus = "PREVIEW"
	StatusPendingApproval ProposalStatus = "PENDING_APPROVAL"
	StatusApproved        ProposalStatus = "APPROVED"
	StatusRejected        ProposalStatus = "REJECTED"
	StatusQueued          ProposalStatus = "QUEUED"
	StatusRunning         ProposalStatus = "RUNNING"
	StatusSucceeded       ProposalStatus = "SUCCEEDED"
	StatusFailed          ProposalStatus = "FAILED"
	StatusRolledBack      ProposalStatus = "ROLLED_BACK"
)

type ModuleDefinition struct {
	Key              Module `json:"key"`
	Label            string `json:"label"`
	Description      string `json:"description"`
	Permission       string `json:"permission"`
	ApprovalRequired bool   `json:"approval_required"`
}

var ModuleCatalog = []ModuleDefinition{
	{ModuleVoucher, "Voucher", "Pembuatan, distribusi, aktivasi, dan pencabutan voucher.", "administrator.voucher.manage", true},
	{ModuleMitra, "Mitra", "Pengelolaan reseller/mitra dan komisi.", "administrator.mitra.manage", true},
	{ModuleBilling, "Billing", "Invoice, jatuh tempo, denda, dan isolir terkontrol.", "administrator.billing.manage", true},
	{ModulePayment, "Payment", "Rekonsiliasi pembayaran dan webhook gateway.", "administrator.payment.manage", true},
	{ModuleNMS, "NMS", "Monitoring perangkat, router, link, dan alarm jaringan.", "administrator.nms.manage", true},
	{ModuleTeknisi, "Teknisi", "Penugasan teknisi, work order, dan validasi pekerjaan.", "administrator.teknisi.manage", true},
	{ModuleFinance, "Finance", "Pendapatan, biaya, komisi, dan laporan keuangan cabang.", "administrator.finance.view", false},
	{ModuleSistem, "Sistem", "Konfigurasi cabang, integrasi, lisensi, dan health check.", "administrator.system.manage", true},
}

type ActionProposal struct {
	ID                string         `json:"id"`
	BranchID          string         `json:"branch_id"`
	Module            Module         `json:"module"`
	Action            string         `json:"action"`
	TargetType        string         `json:"target_type"`
	TargetIDs         []string       `json:"target_ids"`
	BeforeState       map[string]any `json:"before_state"`
	ProposedState     map[string]any `json:"proposed_state"`
	RiskLevel         RiskLevel      `json:"risk_level"`
	Reason            string         `json:"reason"`
	Status            ProposalStatus `json:"status"`
	RequestedBy       string         `json:"requested_by"`
	ApprovedBy        *string        `json:"approved_by,omitempty"`
	WorkerID          *string        `json:"worker_id,omitempty"`
	ErrorMessage      *string        `json:"error_message,omitempty"`
	ProductionChanged bool           `json:"production_changed"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

type AIReport struct {
	ID                string    `json:"id"`
	BranchID          string    `json:"branch_id"`
	Module            Module    `json:"module"`
	Title             string    `json:"title"`
	Severity          RiskLevel `json:"severity"`
	Finding           string    `json:"finding"`
	Recommendation    string    `json:"recommendation"`
	Impact            []string  `json:"impact"`
	Status            string    `json:"status"`
	ProposalID        *string   `json:"proposal_id,omitempty"`
	ProductionChanged bool      `json:"production_changed"`
	CreatedAt         time.Time `json:"created_at"`
}

func IsProductionSafe(p ActionProposal) bool { return !p.ProductionChanged }
