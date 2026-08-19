package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/djalil83/A-Radius/internal/authz"
)

type userContextKey struct{}

func WithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(
		ctx,
		userContextKey{},
		user,
	)
}

func UserFromContext(ctx context.Context) *User {
	user, ok := ctx.Value(userContextKey{}).(User)
	if !ok {
		return nil
	}

	return &user
}

func (h *Handler) RequireSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h == nil || h.Service == nil {
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
			return
		}

		cookie, err := r.Cookie(SessionCookieName)
		if err != nil || cookie.Value == "" {
			http.Error(
				w,
				http.StatusText(http.StatusUnauthorized),
				http.StatusUnauthorized,
			)
			return
		}

		user, err := h.Service.Authenticate(
			r.Context(),
			cookie.Value,
		)
		if err != nil {
			status := http.StatusUnauthorized

			if errors.Is(err, ErrUserDisabled) {
				status = http.StatusForbidden
			}

			http.Error(
				w,
				http.StatusText(status),
				status,
			)
			return
		}

		ctx := WithUser(r.Context(), user)

		principal := &authz.Principal{
			UserID:   user.ID,
			Username: user.Username,
		}

		ctx = authz.WithPrincipal(ctx, principal)

		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)
	})
}
