package pelanggan

import (
	"net/http"
	"strings"

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

	// Outer application router mounts this handler at:
	//
	//	/dashboard/pelanggan/
	//
	// Strip that prefix before dispatching to the role-specific
	// dashboard routes:
	//
	//	/dashboard/pelanggan/          -> /
	//	/dashboard/pelanggan/dashboard -> /dashboard
	//	/dashboard/pelanggan/me        -> /me
	//	/dashboard/pelanggan/services  -> /services
	//
	// This prevents "/" from swallowing all dashboard routes.
	next := http.StripPrefix(
		"/dashboard/pelanggan",
		h.Routes(),
	)

	// Defensive check: StripPrefix only operates on the expected
	// mount path. Requests outside that path must not be dispatched.
	protected := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r == nil ||
			(r.URL.Path != "/dashboard/pelanggan" &&
				!strings.HasPrefix(r.URL.Path, "/dashboard/pelanggan/")) {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})

	return engine.RequirePermissionHTTP(
		"customer.portal.read",
		audit,
		protected,
	)
}
