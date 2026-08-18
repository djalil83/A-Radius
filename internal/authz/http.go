package authz

import (
	"context"
	"errors"
	"net/http"
)

type contextKey string

const principalContextKey contextKey = "a-radius.authz.principal"

func WithPrincipal(
	ctx context.Context,
	principal *Principal,
) context.Context {
	return context.WithValue(
		ctx,
		principalContextKey,
		principal,
	)
}

func PrincipalFromContext(ctx context.Context) *Principal {
	principal, _ := ctx.Value(principalContextKey).(*Principal)
	return principal
}

type AuditDecision func(
	ctx context.Context,
	principal *Principal,
	permission string,
	allowed bool,
	status int,
	r *http.Request,
)

func (e *Engine) PermissionMiddleware(permission string, audit AuditDecision) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return e.RequirePermissionHTTP(permission, audit, next)
	}
}

func (e *Engine) RequirePermissionHTTP(
	permission string,
	audit AuditDecision,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		principal := PrincipalFromContext(r.Context())

		err := e.RequirePermission(
			r.Context(),
			principal,
			permission,
		)

		if err != nil {
			status := http.StatusInternalServerError

			switch {
			case errors.Is(err, ErrUnauthenticated):
				status = http.StatusUnauthorized
			case errors.Is(err, ErrForbidden):
				status = http.StatusForbidden
			}

			if audit != nil {
				audit(
					r.Context(),
					principal,
					permission,
					false,
					status,
					r,
				)
			}

			http.Error(
				w,
				http.StatusText(status),
				status,
			)
			return
		}

		if audit != nil {
			audit(
				r.Context(),
				principal,
				permission,
				true,
				http.StatusOK,
				r,
			)
		}

		next.ServeHTTP(w, r)
	})
}
