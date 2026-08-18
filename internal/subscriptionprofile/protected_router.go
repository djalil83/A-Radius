package subscriptionprofile

import (
	"net/http"

	"github.com/djalil83/A-Radius/internal/authz"
	"github.com/go-chi/chi/v5"
)

// ProtectedRouter mendaftarkan endpoint Subscription Profile dengan permission
// yang eksplisit. Middleware autentikasi JWT harus dipasang oleh caller di luar
// router ini agar Principal sudah tersedia sebelum pengecekan permission.
func ProtectedRouter(h *Handler, engine *authz.Engine, audit authz.AuditDecision) http.Handler {
	r := chi.NewRouter()

	read := func(next http.Handler) http.Handler {
		return engine.RequirePermissionHTTP("profile:read", audit, next)
	}
	write := func(next http.Handler) http.Handler {
		return engine.RequirePermissionHTTP("profile:write", audit, next)
	}

	r.With(read).Get("/api/v1/subscription-profiles", h.list)
	r.With(write).Post("/api/v1/subscription-profiles", h.create)
	r.With(read).Get("/api/v1/subscription-profiles/{id}", h.get)
	r.With(write).Patch("/api/v1/subscription-profiles/{id}", h.update)
	r.With(write).Delete("/api/v1/subscription-profiles/{id}", h.archive)
	r.With(read).Get("/api/v1/subscription-profiles/{id}/revisions", h.revisions)

	return r
}
