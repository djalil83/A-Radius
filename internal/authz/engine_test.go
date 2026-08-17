package authz

import (
	"context"
	"errors"
	"testing"
)

func TestRequirePermissionNilPrincipal(t *testing.T) {
	engine := NewEngine(nil)

	err := engine.RequirePermission(
		context.Background(),
		nil,
		"users.read",
	)

	if !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("expected ErrUnauthenticated, got %v", err)
	}
}

func TestRequirePermissionEmptyUserID(t *testing.T) {
	engine := NewEngine(nil)

	err := engine.RequirePermission(
		context.Background(),
		&Principal{
			UserID: "",
		},
		"users.read",
	)

	if !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("expected ErrUnauthenticated, got %v", err)
	}
}

func TestRequirePermissionEmptyPermission(t *testing.T) {
	engine := NewEngine(nil)

	err := engine.RequirePermission(
		context.Background(),
		&Principal{
			UserID: "user-id",
		},
		"",
	)

	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestHasPermissionInvalidInputFailsClosed(t *testing.T) {
	engine := NewEngine(nil)

	allowed, err := engine.HasPermission(
		context.Background(),
		"",
		"users.read",
	)

	if err == nil {
		t.Fatal("expected configuration error because engine has no DB")
	}

	if !errors.Is(err, ErrNotConfigured) {
		t.Fatalf("expected ErrNotConfigured, got %v", err)
	}

	if allowed {
		t.Fatal("expected permission to be denied")
	}
}

func TestHasPermissionWithoutDatabaseFailsClosed(t *testing.T) {
	engine := NewEngine(nil)

	allowed, err := engine.HasPermission(
		context.Background(),
		"user-id",
		"users.read",
	)

	if !errors.Is(err, ErrNotConfigured) {
		t.Fatalf("expected ErrNotConfigured, got %v", err)
	}

	if allowed {
		t.Fatal("expected permission to be denied")
	}
}

func TestErrors(t *testing.T) {
	if ErrUnauthenticated == nil {
		t.Fatal("ErrUnauthenticated must exist")
	}

	if ErrForbidden == nil {
		t.Fatal("ErrForbidden must exist")
	}

	if ErrNotConfigured == nil {
		t.Fatal("ErrNotConfigured must exist")
	}
}
