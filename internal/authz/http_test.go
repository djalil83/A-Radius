package authz

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequirePermissionHTTPUnauthenticated(t *testing.T) {
	engine := NewEngine(nil)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler must not be called")
	})

	handler := engine.RequirePermissionHTTP(
		"users.read",
		nil,
		next,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/users",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected HTTP 401, got %d",
			rec.Code,
		)
	}
}

func TestRequirePermissionHTTPAuthorizationEngineError(t *testing.T) {
	engine := NewEngine(nil)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler must not be called")
	})

	handler := engine.RequirePermissionHTTP(
		"users.read",
		nil,
		next,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/users",
		nil,
	)

	req = req.WithContext(
		WithPrincipal(
			context.Background(),
			&Principal{
				UserID: "user-1",
			},
		),
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected HTTP 500, got %d",
			rec.Code,
		)
	}
}
