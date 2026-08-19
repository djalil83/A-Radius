package customerportal

import "net/http"

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
) {
	mux.HandleFunc("/api/v1/customer/dashboard", handler.Dashboard)
	mux.HandleFunc("/api/v1/customer/me", handler.Me)
	mux.HandleFunc("/api/v1/customer/services", handler.Services)
}
