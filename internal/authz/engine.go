package authz

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrUnauthenticated = errors.New("authentication required")
	ErrForbidden       = errors.New("permission denied")
	ErrNotConfigured   = errors.New("authorization engine is not configured")
)

type Principal struct {
	UserID   string
	Username string
}

type Querier interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Engine struct {
	db Querier
}

func NewEngine(db Querier) *Engine {
	return &Engine{db: db}
}

// HasPermission evaluates RBAC server-side.
// Authorization is fail-closed.
func (e *Engine) HasPermission(
	ctx context.Context,
	userID string,
	permission string,
) (bool, error) {
	if e == nil || e.db == nil {
		return false, ErrNotConfigured
	}

	if userID == "" || permission == "" {
		return false, nil
	}

	const query = `
SELECT EXISTS (
SELECT 1
FROM apb.user_roles ur
JOIN apb.role_permissions rp
ON rp.role_id = ur.role_id
JOIN apb.permissions p
ON p.id = rp.permission_id
JOIN apb.users u
ON u.id = ur.user_id
WHERE u.id = $1
AND u.status = 'active'
AND p.permission_key = $2
)`

	var allowed bool

	if err := e.db.QueryRowContext(
		ctx,
		query,
		userID,
		permission,
	).Scan(&allowed); err != nil {
		return false, fmt.Errorf("check permission: %w", err)
	}

	return allowed, nil
}

// RequirePermission enforces an application-level permission.
func (e *Engine) RequirePermission(
	ctx context.Context,
	principal *Principal,
	permission string,
) error {
	if principal == nil || principal.UserID == "" {
		return ErrUnauthenticated
	}

	if permission == "" {
		return ErrForbidden
	}

	allowed, err := e.HasPermission(
		ctx,
		principal.UserID,
		permission,
	)
	if err != nil {
		return err
	}

	if !allowed {
		return ErrForbidden
	}

	return nil
}
