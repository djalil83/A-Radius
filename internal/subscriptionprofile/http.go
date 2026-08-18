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

type PermissionMiddleware func(permission string, next http.Handler) http.Handler

func Router(h *Handler, protect PermissionMiddleware) http.Handler {
	r := chi.NewRouter()
	route := func(method, path, permission string, handler http.HandlerFunc) {
		if protect == nil {
			r.Method(method, path, handler)
			return
		}
		r.With(func(next http.Handler) http.Handler { return protect(permission, next) }).Method(method, path, handler)
	}
	route(http.MethodGet, "/api/v1/subscription-profiles", "subscription_profiles.read", h.list)
	route(http.MethodPost, "/api/v1/subscription-profiles", "subscription_profiles.create", h.create)
	route(http.MethodGet, "/api/v1/subscription-profiles/{id}", "subscription_profiles.read", h.get)
	route(http.MethodPatch, "/api/v1/subscription-profiles/{id}", "subscription_profiles.update", h.update)
	route(http.MethodDelete, "/api/v1/subscription-profiles/{id}", "subscription_profiles.archive", h.archive)
	route(http.MethodGet, "/api/v1/subscription-profiles/{id}/revisions", "subscription_profiles.read_history", h.revisions)
	return r
}

func identityFromRequest(w http.ResponseWriter, r *http.Request) (tenantID, actorID string, ok bool) {
	principal := authz.PrincipalFromContext(r.Context())
	if principal == nil || principal.UserID == "" || principal.TenantID == "" {
		writeError(w, http.StatusUnauthorized, "UNAUTHENTICATED", "a valid JWT identity is required")
		return "", "", false
	}
	return principal.TenantID, principal.UserID, true
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	tenant, _, ok := identityFromRequest(w, r)
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
	tenant, _, ok := identityFromRequest(w, r)
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
	tenant, actor, ok := identityFromRequest(w, r)
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
	tenant, actor, ok := identityFromRequest(w, r)
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
	tenant, actor, ok := identityFromRequest(w, r)
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
	tenant, _, ok := identityFromRequest(w, r)
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
