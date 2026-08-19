package pelanggan

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/djalil83/A-Radius/internal/customerportal"
	"github.com/djalil83/A-Radius/web"
)

// Handler adalah adapter role-specific untuk dashboard pelanggan.
type Handler struct {
	portal *customerportal.Handler
}

// NewHandler membuat adapter dashboard pelanggan.
func NewHandler(portal *customerportal.Handler) *Handler {
	return &Handler{portal: portal}
}

// Routes menyediakan UI dashboard sekaligus endpoint internal dashboard.
//
// Middleware authentication/RBAC tetap dipasang oleh caller.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", h.dashboardUI)
	mux.HandleFunc("/dashboard", h.dashboard)
	mux.HandleFunc("/me", h.me)
	mux.HandleFunc("/services", h.services)

	mux.HandleFunc("/dashboard.js", h.dashboardJS)

	return mux
}

func (h *Handler) dashboardUI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data, err := fs.ReadFile(web.Assets, "dashboards/pelanggan/index.html")
	if err != nil {
		http.Error(
			w,
			"customer dashboard unavailable",
			http.StatusServiceUnavailable,
		)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")

	_, _ = w.Write(data)
}

func (h *Handler) dashboardJS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	if !strings.HasSuffix(r.URL.Path, "/dashboard.js") {
		http.NotFound(w, r)
		return
	}

	data, err := fs.ReadFile(web.Assets, "dashboards/pelanggan/dashboard.js")
	if err != nil {
		http.Error(
			w,
			"customer dashboard asset unavailable",
			http.StatusServiceUnavailable,
		)
		return
	}

	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")

	_, _ = w.Write(data)
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
