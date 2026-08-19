package subscriptionproduction

import (
	"net/http"

	"github.com/djalil83/A-Radius/internal/authz"
	"github.com/go-chi/chi/v5"
)

func ProtectedRouter(h *Handler, engine *authz.Engine, audit authz.AuditDecision) http.Handler {
	r := chi.NewRouter()
	read := func(next http.Handler) http.Handler {
		return engine.RequirePermissionHTTP("subscription:read", audit, next)
	}
	preview := func(next http.Handler) http.Handler {
		return engine.RequirePermissionHTTP("subscription:preview", audit, next)
	}

	r.With(read).Get("/api/v1/subscription-production/policy", h.policy)
	r.With(read).Get("/api/v1/subscription-production/integrations", h.integrations)
	r.With(read).Get("/api/v1/subscription-production/readiness", h.readiness)
	r.With(preview).Post("/api/v1/subscription-production/preview", h.preview)
	return r
}
