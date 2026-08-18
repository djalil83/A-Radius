package developer

import (
	"encoding/json"
	"net/http"
	"time"
)

// Handler menyediakan endpoint khusus Developer Security Dashboard.
type Handler struct{}

// Routes mendaftarkan endpoint yang selanjutnya wajib dibungkus middleware JWT
// dan permission RBAC pada router aplikasi.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.index)
	mux.HandleFunc("/security/overview", h.securityOverview)
	mux.HandleFunc("/security/features", h.features)
	mux.HandleFunc("/security/scans", h.scans)
	mux.HandleFunc("/security/continuous/policy", h.continuousPolicy)
	mux.HandleFunc("/security/continuous/sources", h.continuousSources)
	mux.HandleFunc("/security/continuous/featured-finding", h.featuredContinuousFinding)
	mux.HandleFunc("/security/continuous/patch-preview", h.featuredPatchPreview)
	mux.HandleFunc("/security/knowledge/policy", h.knowledgePolicy)
	mux.HandleFunc("/security/knowledge/versions", h.knowledgeVersions)
	mux.HandleFunc("/security/knowledge/featured", h.featuredKnowledge)
	mux.HandleFunc("/security/knowledge/new-intelligence", h.newIntelligence)
	mux.HandleFunc("/security/knowledge/compare", h.knowledgeCompare)
	mux.HandleFunc("/security/knowledge/patch-pipeline", h.knowledgePatchPipeline)
	mux.HandleFunc("/security/knowledge/rollback-audit", h.knowledgeRollbackAudit)
	return mux
}

func (h *Handler) index(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"dashboard": "developer-security",
		"status":    "ready",
		"timestamp": time.Now().UTC(),
	})
}

func (h *Handler) securityOverview(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"score": 100,
		"findings": map[string]int{
			"critical": 0, "high": 0, "medium": 0, "low": 0, "informational": 0,
		},
		"scan_status": "not_scanned",
	})
}

func (h *Handler) features(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"items": FeatureContracts})
}

func (h *Handler) scans(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": map[string]string{"code": "METHOD_NOT_ALLOWED", "message": "only POST is supported"}})
		return
	}

	var request struct {
		Type string `json:"type"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10)).Decode(&request); err != nil || request.Type == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": map[string]string{"code": "INVALID_SCAN_REQUEST", "message": "scan type is required"}})
		return
	}
	if request.Type != "full" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": map[string]string{"code": "UNSUPPORTED_SCAN_TYPE", "message": "supported scan type is full"}})
		return
	}

	// Eksekusi scanner nyata akan dibuat sebagai job asynchronous setelah queue,
	// worker, evidence storage, dan approval policy tersedia.
	writeJSON(w, http.StatusAccepted, map[string]any{
		"scan_id": "pending",
		"type":    request.Type,
		"status":  "queued",
	})
}

func (h *Handler) continuousPolicy(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, DefaultContinuousSecurityPolicy)
}

func (h *Handler) continuousSources(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"sources": TrustedKnowledgeSources, "status": "validated"})
}

func (h *Handler) featuredContinuousFinding(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, FeaturedContinuousFinding)
}

func (h *Handler) knowledgePolicy(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, DefaultKnowledgePromotionPolicy)
}

func (h *Handler) knowledgeVersions(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"active_version": FeaturedKnowledgeVersion,
		"counts":         map[string]int{"new": 12, "review": 4, "active": 1, "archived": 18},
		"versions": []KnowledgeVersionRecord{
			FeaturedKnowledgeVersion,
			{Version: "SK-2.4.6", Status: KnowledgeArchived, Findings: 19, Source: "Security AI", DiscoveredAt: time.Date(2026, time.August, 14, 14, 20, 0, 0, time.FixedZone("WITA", 8*60*60)), Environment: KnowledgeEnvironmentIsolated},
			{Version: "SK-2.4.5", Status: KnowledgeArchived, Findings: 17, Source: "Security AI", DiscoveredAt: time.Date(2026, time.August, 11, 14, 20, 0, 0, time.FixedZone("WITA", 8*60*60)), Environment: KnowledgeEnvironmentIsolated},
			{Version: "SK-2.4.4", Status: KnowledgeArchived, Findings: 16, Source: "Security AI", DiscoveredAt: time.Date(2026, time.August, 8, 14, 20, 0, 0, time.FixedZone("WITA", 8*60*60)), Environment: KnowledgeEnvironmentIsolated},
		},
		"environment":        KnowledgeEnvironmentIsolated,
		"production_changed": false,
	})
}

func (h *Handler) featuredKnowledge(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, FeaturedKnowledgeVersion)
}

func (h *Handler) newIntelligence(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, FeaturedNewIntelligence)
}

func (h *Handler) knowledgeCompare(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": map[string]string{"code": "METHOD_NOT_ALLOWED", "message": "only GET is supported"}})
		return
	}
	writeJSON(w, http.StatusOK, FeaturedKnowledgeComparison)
}

func (h *Handler) knowledgePatchPipeline(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": map[string]string{"code": "METHOD_NOT_ALLOWED", "message": "only GET is supported"}})
		return
	}
	writeJSON(w, http.StatusOK, DefaultKnowledgePatchPipeline)
}

func (h *Handler) knowledgeRollbackAudit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": map[string]string{"code": "METHOD_NOT_ALLOWED", "message": "only GET is supported"}})
		return
	}
	writeJSON(w, http.StatusOK, FeaturedKnowledgeRollbackAudit)
}

func (h *Handler) featuredPatchPreview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET, POST")
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": map[string]string{"code": "METHOD_NOT_ALLOWED", "message": "only GET or POST is supported"}})
		return
	}
	if r.Method == http.MethodPost {
		writeJSON(w, http.StatusAccepted, FeaturedPatchPreview)
		return
	}
	writeJSON(w, http.StatusOK, FeaturedPatchPreview)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
