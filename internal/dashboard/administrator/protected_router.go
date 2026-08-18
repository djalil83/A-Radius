package administrator

import (
	"net/http"

	"github.com/djalil83/A-Radius/internal/authz"
	"github.com/go-chi/chi/v5"
)

func ProtectedRouter(h *Handler, engine *authz.Engine, audit authz.AuditDecision) http.Handler {
	r := chi.NewRouter()
	read := func(next http.Handler) http.Handler {
		return engine.RequirePermissionHTTP("administrator.dashboard.read", audit, next)
	}
	proposal := func(next http.Handler) http.Handler {
		return engine.RequirePermissionHTTP("administrator.proposal.create", audit, next)
	}
	approve := func(next http.Handler) http.Handler {
		return engine.RequirePermissionHTTP("administrator.proposal.approve", audit, next)
	}

	r.With(read).Get("/api/v1/administrator/modules", h.modules)
	r.With(read).Get("/api/v1/administrator/ai-reports", h.aiReports)
	r.With(proposal).Post("/api/v1/administrator/proposals/preview", h.preview)
	r.With(approve).Post("/api/v1/administrator/proposals/{id}/approve", h.approve)
	r.With(approve).Post("/api/v1/administrator/proposals/{id}/reject", h.reject)
	return r
}
