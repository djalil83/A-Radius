package pelanggan

import (
	"net/http"

	"github.com/djalil83/A-Radius/internal/customerportal"
)

// Handler adalah adapter role-specific untuk dashboard pelanggan.
// Authentication dan RBAC dilakukan pada router aplikasi.
type Handler struct {
	portal *customerportal.Handler
}

// NewHandler membuat adapter dashboard pelanggan.
func NewHandler(portal *customerportal.Handler) *Handler {
	return &Handler{portal: portal}
}

// Routes menyediakan endpoint dashboard pelanggan.
// Middleware authentication/RBAC harus dipasang oleh caller.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", h.dashboard)
	mux.HandleFunc("/dashboard", h.dashboard)
	mux.HandleFunc("/me", h.me)
	mux.HandleFunc("/services", h.services)

	return mux
}

func (h *Handler) dashboard(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.portal == nil {
		http.Error(
			w,
			"customer dashboard unavailable",
			http.StatusServiceUnavailable,
		)
		return
	}

	h.portal.Dashboard(w, r)
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.portal == nil {
		http.Error(
			w,
			"customer dashboard unavailable",
			http.StatusServiceUnavailable,
		)
		return
	}

	h.portal.Me(w, r)
}

func (h *Handler) services(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.portal == nil {
		http.Error(
			w,
			"customer dashboard unavailable",
			http.StatusServiceUnavailable,
		)
		return
	}

	h.portal.Services(w, r)
}
