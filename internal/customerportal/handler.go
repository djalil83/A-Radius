package customerportal

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/djalil83/A-Radius/internal/authz"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// authenticatedUserID obtains the user identity exclusively from
// the authenticated authorization principal.
//
// Do not use client-controlled headers such as X-Actor-ID here.
func authenticatedUserID(r *http.Request) (string, bool) {
	if r == nil {
		return "", false
	}

	principal := authz.PrincipalFromContext(r.Context())
	if principal == nil || principal.UserID == "" {
		return "", false
	}

	return principal.UserID, true
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := authenticatedUserID(r)
	if !ok {
		http.Error(w, "authentication required", http.StatusUnauthorized)
		return
	}

	customer, err := h.service.CustomerProfile(r.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrCustomerNotFound) {
			http.Error(w, "customer not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, customer)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(value)
}

func (h *Handler) Services(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := authenticatedUserID(r)
	if !ok {
		http.Error(w, "authentication required", http.StatusUnauthorized)
		return
	}

	services, err := h.service.CustomerServices(r.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrCustomerNotFound) {
			http.Error(w, "customer not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, services)
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := authenticatedUserID(r)
	if !ok {
		http.Error(w, "authentication required", http.StatusUnauthorized)
		return
	}

	dashboard, err := h.service.CustomerDashboard(r.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrCustomerNotFound) {
			http.Error(w, "customer not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, dashboard)
}
