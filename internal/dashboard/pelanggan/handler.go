package pelanggan

import (
	"net/http"
)

// Handler menyediakan endpoint khusus untuk Pelanggan Dashboard.
type Handler struct{}

// Routes mendaftarkan route dengan middleware role="pelanggan" pada router induk.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.index)
	return mux
}

func (h *Handler) index(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	_, _ = w.Write([]byte(`{"error":{"code":"NOT_IMPLEMENTED","message":"dashboard endpoint belum diimplementasikan"}}`))
}
