package administrator

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/djalil83/A-Radius/internal/authz"
	"github.com/go-chi/chi/v5"
)

type Handler struct{ DB *sql.DB }

func NewHandler(db *sql.DB) *Handler { return &Handler{DB: db} }

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/api/v1/administrator/modules", h.modules)
	r.Get("/api/v1/administrator/ai-reports", h.aiReports)
	r.Post("/api/v1/administrator/proposals/preview", h.preview)
	r.Post("/api/v1/administrator/proposals/{id}/approve", h.approve)
	r.Post("/api/v1/administrator/proposals/{id}/reject", h.reject)
	return r
}

func (h *Handler) modules(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"modules": ModuleCatalog, "production_changed": false})
}

func (h *Handler) aiReports(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		writeJSON(w, http.StatusServiceUnavailable, errorBody("DATABASE_NOT_CONFIGURED"))
		return
	}
	rows, err := h.DB.QueryContext(r.Context(), `SELECT id, branch_id, module, title, severity, finding, recommendation, impact, status, proposal_id, production_changed, created_at FROM apb.admin_ai_reports WHERE branch_id = $1 ORDER BY created_at DESC LIMIT 100`, queryBranch(r))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorBody("REPORT_QUERY_FAILED"))
		return
	}
	defer rows.Close()
	reports := make([]AIReport, 0)
	for rows.Next() {
		var report AIReport
		var impact []byte
		var proposal sql.NullString
		if err := rows.Scan(&report.ID, &report.BranchID, &report.Module, &report.Title, &report.Severity, &report.Finding, &report.Recommendation, &impact, &report.Status, &proposal, &report.ProductionChanged, &report.CreatedAt); err != nil {
			writeJSON(w, http.StatusInternalServerError, errorBody("REPORT_SCAN_FAILED"))
			return
		}
		_ = json.Unmarshal(impact, &report.Impact)
		if proposal.Valid {
			report.ProposalID = &proposal.String
		}
		reports = append(reports, report)
	}
	writeJSON(w, http.StatusOK, map[string]any{"reports": reports, "production_changed": false})
}

type previewRequest struct {
	BranchID      string         `json:"branch_id"`
	Module        Module         `json:"module"`
	Action        string         `json:"action"`
	TargetType    string         `json:"target_type"`
	TargetIDs     []string       `json:"target_ids"`
	BeforeState   map[string]any `json:"before_state"`
	ProposedState map[string]any `json:"proposed_state"`
	RiskLevel     RiskLevel      `json:"risk_level"`
	Reason        string         `json:"reason"`
}

func (h *Handler) preview(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		writeJSON(w, http.StatusServiceUnavailable, errorBody("DATABASE_NOT_CONFIGURED"))
		return
	}
	principal := authz.PrincipalFromContext(r.Context())
	if principal == nil {
		writeJSON(w, http.StatusUnauthorized, errorBody("UNAUTHENTICATED"))
		return
	}
	var req previewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody("INVALID_JSON"))
		return
	}
	if !validModule(req.Module) || strings.TrimSpace(req.Action) == "" || strings.TrimSpace(req.Reason) == "" || strings.TrimSpace(req.BranchID) == "" {
		writeJSON(w, http.StatusBadRequest, errorBody("INVALID_PROPOSAL"))
		return
	}
	id := newID()
	now := time.Now().UTC()
	targetIDs, _ := json.Marshal(req.TargetIDs)
	before, _ := json.Marshal(defaultMap(req.BeforeState))
	proposed, _ := json.Marshal(defaultMap(req.ProposedState))
	_, err := h.DB.ExecContext(r.Context(), `INSERT INTO apb.admin_action_proposals (id,branch_id,module,action,target_type,target_ids,before_state,proposed_state,risk_level,reason,status,requested_by,production_changed,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,'PENDING_APPROVAL',$11,false,$12,$12)`, id, req.BranchID, req.Module, req.Action, req.TargetType, targetIDs, before, proposed, req.RiskLevel, req.Reason, principal.UserID, now)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorBody("PROPOSAL_CREATE_FAILED"))
		return
	}
	h.recordAudit(r, id, req.BranchID, principal.UserID, req.Module, req.Action, req.TargetType, req.TargetIDs, "PENDING_APPROVAL")
	writeJSON(w, http.StatusCreated, map[string]any{"proposal_id": id, "status": StatusPendingApproval, "production_changed": false, "approval_required": true})
}

func (h *Handler) approve(w http.ResponseWriter, r *http.Request) { h.decide(w, r, StatusApproved) }
func (h *Handler) reject(w http.ResponseWriter, r *http.Request)  { h.decide(w, r, StatusRejected) }

func (h *Handler) decide(w http.ResponseWriter, r *http.Request, status ProposalStatus) {
	if h.DB == nil {
		writeJSON(w, http.StatusServiceUnavailable, errorBody("DATABASE_NOT_CONFIGURED"))
		return
	}
	principal := authz.PrincipalFromContext(r.Context())
	if principal == nil {
		writeJSON(w, http.StatusUnauthorized, errorBody("UNAUTHENTICATED"))
		return
	}
	id := chi.URLParam(r, "id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, errorBody("INVALID_PROPOSAL_ID"))
		return
	}
	result, err := h.DB.ExecContext(r.Context(), `UPDATE apb.admin_action_proposals SET status=$1, approved_by=$2, updated_at=now(), production_changed=false WHERE id=$3 AND status='PENDING_APPROVAL' AND requested_by <> $2`, status, principal.UserID, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorBody("PROPOSAL_DECISION_FAILED"))
		return
	}
	count, _ := result.RowsAffected()
	if count != 1 {
		writeJSON(w, http.StatusConflict, errorBody("PROPOSAL_NOT_PENDING_OR_SELF_APPROVAL"))
		return
	}
	var branchID string
	var module Module
	var action, targetType string
	var targetIDs []byte
	if err := h.DB.QueryRowContext(r.Context(), `SELECT branch_id, module, action, target_type, target_ids FROM apb.admin_action_proposals WHERE id=$1`, id).Scan(&branchID, &module, &action, &targetType, &targetIDs); err == nil {
		var ids []string
		_ = json.Unmarshal(targetIDs, &ids)
		h.recordAudit(r, id, branchID, principal.UserID, module, action, targetType, ids, string(status))
	}
	writeJSON(w, http.StatusOK, map[string]any{"proposal_id": id, "status": status, "production_changed": false, "execution": "NOT_EXECUTED"})
}

func (h *Handler) recordAudit(r *http.Request, proposalID, branchID, actorID string, module Module, action, targetType string, targetIDs []string, result string) {
	ids, _ := json.Marshal(targetIDs)
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = newID()
	}
	_, _ = h.DB.ExecContext(r.Context(), `INSERT INTO apb.admin_module_audit (proposal_id,branch_id,actor_id,actor_role,module,action,target_type,target_ids,result,request_id,user_agent) VALUES ($1,$2,$3,'administrator',$4,$5,$6,$7,$8,$9,$10)`, proposalID, branchID, actorID, module, action, targetType, ids, result, requestID, r.UserAgent())
}

func validModule(m Module) bool {
	for _, item := range ModuleCatalog {
		if item.Key == m {
			return true
		}
	}
	return false
}
func defaultMap(v map[string]any) map[string]any {
	if v == nil {
		return map[string]any{}
	}
	return v
}
func queryBranch(r *http.Request) string {
	if v := r.URL.Query().Get("branch_id"); v != "" {
		return v
	}
	if p := authz.PrincipalFromContext(r.Context()); p != nil && p.TenantID != "" {
		return p.TenantID
	}
	return "0"
}
func errorBody(code string) map[string]any {
	return map[string]any{"error": map[string]string{"code": code}}
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(errors.New("secure id generation failed"))
	}
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}
