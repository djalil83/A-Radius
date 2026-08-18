package auth

import "net/http"

// RegisterRoutes registers authentication endpoints.
func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
) {
	if mux == nil || handler == nil {
		return
	}

	mux.HandleFunc(
		"/api/v1/auth/login",
		handler.Login,
	)

	mux.HandleFunc(
		"/api/v1/auth/logout",
		handler.Logout,
	)
}
