package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

const SessionCookieName = "a_radius_session"

type AuthService interface {
	Login(
		ctx context.Context,
		username string,
		password string,
		ipAddress string,
		userAgent string,
	) (User, string, time.Time, error)

	Logout(
		ctx context.Context,
		rawToken string,
	) error

	Authenticate(
		ctx context.Context,
		rawToken string,
	) (User, error)
}

type Handler struct {
	Service AuthService
}

func NewHandler(service AuthService) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.Service == nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}

	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
		)
		return
	}

	user, token, expiresAt, err := h.Service.Login(
		r.Context(),
		req.Username,
		req.Password,
		clientIP(r),
		r.UserAgent(),
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			http.Error(
				w,
				http.StatusText(http.StatusUnauthorized),
				http.StatusUnauthorized,
			)

		case errors.Is(err, ErrUserDisabled):
			http.Error(
				w,
				http.StatusText(http.StatusForbidden),
				http.StatusForbidden,
			)

		default:
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
		}

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   int(time.Until(expiresAt).Seconds()),
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})

	writeJSON(
		w,
		http.StatusOK,
		LoginResponse{
			UserID:    user.ID,
			Username:  user.Username,
			ExpiresAt: expiresAt,
		},
	)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.Service == nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}

	cookie, err := r.Cookie(SessionCookieName)
	if err == nil && cookie.Value != "" {
		if err := h.Service.Logout(
			r.Context(),
			cookie.Value,
		); err != nil {
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(1, 0),
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "logged out",
	})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())

	if user == nil {
		http.Error(
			w,
			http.StatusText(http.StatusUnauthorized),
			http.StatusUnauthorized,
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		map[string]any{
			"user_id":  user.ID,
			"username": user.Username,
			"status":   user.Status,
		},
	)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(value)
}

func clientIP(r *http.Request) string {
	if r == nil {
		return ""
	}

	host := r.RemoteAddr

	if i := strings.LastIndex(host, ":"); i >= 0 {
		host = strings.Trim(host[:i], "[]")
	}

	return host
}
