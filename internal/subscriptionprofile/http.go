package subscriptionprofile

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/djalil83/A-Radius/internal/authz"
	"github.com/go-chi/chi/v5"
)

type Handler struct{ repo *Repository }

func NewHandler(repo *Repository) *Handler { return &Handler{repo: repo} }

func Router(h *Handler) http.Handler {
	r := chi.NewRouter()
	r.Get("/api/v1/subscription-profiles", h.list)
	r.Post("/api/v1/subscription-profiles", h.create)
	r.Get("/api/v1/subscription-profiles/{id}", h.get)
	r.Patch("/api/v1/subscription-profiles/{id}", h.update)
	r.Delete("/api/v1/subscription-profiles/{id}", h.archive)
	r.Get("/api/v1/subscription-profiles/{id}/revisions", h.revisions)
	return r
}

// identityHeaders memakai Principal terverifikasi pada router produksi.
// Fallback header dipertahankan hanya untuk Router legacy dan test lama; jangan
// expose Router legacy tanpa authn.Middleware di lingkungan produksi.
func identityHeaders(w http.ResponseWriter, r *http.Request) (tenantID, actorID string, ok bool) {
	if principal := authz.PrincipalFromContext(r.Context()); principal != nil {
		if strings.TrimSpace(principal.UserID) == "" || strings.TrimSpace(principal.TenantID) == "" {
			writeError(w, http.StatusUnauthorized, "UNAUTHENTICATED", "verified principal is incomplete")
			return "", "", false
		}
		return principal.TenantID, principal.UserID, true
	}

	// Compatibility adapter only. ProtectedRouter never relies on this path.
	tenantID = strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
	actorID = strings.TrimSpace(r.Header.Get("X-Actor-ID"))
	if tenantID == "" || actorID == "" {
		writeError(w, http.StatusUnauthorized, "UNAUTHENTICATED", "verified principal is required")
		return "", "", false
	}
	return tenantID, actorID, true
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	tenant, _, ok := identityHeaders(w, r)
	if !ok {
		return
	}
	limit, offset := boundedInt(r.URL.Query().Get("limit"), 50, 1, 100), boundedInt(r.URL.Query().Get("offset"), 0, 0, 1000000)
	items, err := h.repo.List(r.Context(), tenant, r.URL.Query().Get("q"), strings.ToUpper(r.URL.Query().Get("service_type")), strings.ToUpper(r.URL.Query().Get("status")), limit, offset)
	if err != nil {
		writeError(w, 500, "INTERNAL_ERROR", "failed to list subscription profiles")
		return
	}
	writeJSON(w, http.StatusOK, ListResult{Items: items, Limit: limit, Offset: offset})
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	tenant, _, ok := identityHeaders(w, r)
	if !ok {
		return
	}
	p, err := h.repo.Get(r.Context(), tenant, chi.URLParam(r, "id"))
	if err != nil {
		mapError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	tenant, actor, ok := identityHeaders(w, r)
	if !ok {
		return
	}
	var req CreateRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	req, err := NormalizeCreate(req)
	if err != nil {
		writeError(w, 400, "VALIDATION_ERROR", err.Error())
		return
	}
	p, err := h.repo.Create(r.Context(), tenant, actor, req)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			writeError(w, 409, "DUPLICATE_NAME", "a profile with this name already exists")
		} else {
			writeError(w, 500, "INTERNAL_ERROR", "failed to create subscription profile")
		}
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	tenant, actor, ok := identityHeaders(w, r)
	if !ok {
		return
	}
	var req UpdateRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	req, err := NormalizeUpdate(req)
	if err != nil {
		writeError(w, 400, "VALIDATION_ERROR", err.Error())
		return
	}
	p, err := h.repo.Update(r.Context(), tenant, chi.URLParam(r, "id"), actor, req)
	if err != nil {
		mapError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *Handler) archive(w http.ResponseWriter, r *http.Request) {
	tenant, actor, ok := identityHeaders(w, r)
	if !ok {
		return
	}
	version, err := strconv.ParseInt(r.URL.Query().Get("version"), 10, 64)
	if err != nil || version < 1 {
		writeError(w, 400, "VALIDATION_ERROR", "version query parameter must be a positive integer")
		return
	}
	if err := h.repo.Archive(r.Context(), tenant, chi.URLParam(r, "id"), actor, version); err != nil {
		mapError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) revisions(w http.ResponseWriter, r *http.Request) {
	tenant, _, ok := identityHeaders(w, r)
	if !ok {
		return
	}
	limit := boundedInt(r.URL.Query().Get("limit"), 50, 1, 100)
	items, err := h.repo.Revisions(r.Context(), tenant, chi.URLParam(r, "id"), limit)
	if err != nil {
		writeError(w, 500, "INTERNAL_ERROR", "failed to list profile revisions")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "limit": limit})
}

func mapError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		writeError(w, 404, "NOT_FOUND", "subscription profile not found")
	case errors.Is(err, ErrConflict):
		writeError(w, 409, "VERSION_CONFLICT", "profile was changed by another request; reload before updating")
	case errors.Is(err, ErrValidation):
		writeError(w, 400, "VALIDATION_ERROR", err.Error())
	default:
		writeError(w, 500, "INTERNAL_ERROR", "internal server error")
	}
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	if r.Body == nil {
		writeError(w, 400, "INVALID_JSON", "request body is required")
		return false
	}
	defer r.Body.Close()
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		writeError(w, 400, "INVALID_JSON", "request body is invalid")
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}
func boundedInt(raw string, fallback, min, max int) int {
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < min {
		return fallback
	}
	if n > max {
		return max
	}
	return n
}
