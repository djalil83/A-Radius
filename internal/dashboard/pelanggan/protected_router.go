package pelanggan

import (
	"net/http"

	"github.com/djalil83/A-Radius/internal/authz"
)

// ProtectedRouter memasang authorization/RBAC di depan dashboard pelanggan.
//
// Authentication harus dilakukan oleh caller sebelum handler ini.
// Permission canonical untuk customer portal adalah:
// customer.portal.read
//
// Audit keputusan authorization dilakukan oleh caller.
func ProtectedRouter(
	h *Handler,
	engine *authz.Engine,
	audit authz.AuditDecision,
) http.Handler {
	if h == nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(
				w,
				"customer dashboard unavailable",
				http.StatusServiceUnavailable,
			)
		})
	}

	if engine == nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(
				w,
				"authorization unavailable",
				http.StatusServiceUnavailable,
			)
		})
	}

	return engine.RequirePermissionHTTP(
		"customer.portal.read",
		audit,
		h.Routes(),
	)
}
