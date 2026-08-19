package subscriptionproduction

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/policy", h.policy)
	mux.HandleFunc("/integrations", h.integrations)
	mux.HandleFunc("/readiness", h.readiness)
	mux.HandleFunc("/preview", h.preview)
	return mux
}

func (h *Handler) policy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	writeJSON(w, http.StatusOK, DefaultProductionPolicy)
}

func (h *Handler) integrations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": ProductionIntegrationBindings, "production_changed": false})
}

func (h *Handler) readiness(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	id := strings.TrimSpace(r.URL.Query().Get("subscription_id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "SUBSCRIPTION_ID_REQUIRED", "subscription_id is required")
		return
	}
	writeJSON(w, http.StatusOK, ActivationReadiness{SubscriptionID: id, Ready: false, Blockers: []string{"preview_required", "approval_required", "integration_checks_pending"}, Bindings: ProductionIntegrationBindings, ProductionChanged: false})
}

func (h *Handler) preview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}
	var request struct {
		SubscriptionID  string `json:"subscription_id"`
		Action          string `json:"action"`
		ExpectedVersion int64  `json:"expected_version"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10)).Decode(&request); err != nil || strings.TrimSpace(request.SubscriptionID) == "" || strings.TrimSpace(request.Action) == "" {
		writeError(w, http.StatusBadRequest, "INVALID_PREVIEW_REQUEST", "subscription_id and action are required")
		return
	}
	writeJSON(w, http.StatusAccepted, SubscriptionChange{SubscriptionID: request.SubscriptionID, Action: request.Action, RequestedBy: "principal", Stage: StagePreview, Status: "PREVIEW_CREATED", ExpectedVersion: request.ExpectedVersion, Integrations: ProductionIntegrationBindings, ProductionChanged: false})
}

func methodNotAllowed(w http.ResponseWriter, allow string) {
	w.Header().Set("Allow", allow)
	writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method is not supported")
}
func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
